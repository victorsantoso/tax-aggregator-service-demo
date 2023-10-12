package usecase

import (
	"math"
	"sync"
	"tax-aggregator-service-demo/tax"
	"tax-aggregator-service-demo/tax/domain"
	"tax-aggregator-service-demo/tax/entity"
	"time"
)

type taxUsecase struct {
	taxRepository domain.TaxRepository
	taxConfig     *domain.TaxConfig
}

func NewTaxUsecase(taxRepository domain.TaxRepository, taxConfig *domain.TaxConfig) domain.TaxUsecase {
	return &taxUsecase{
		taxRepository: taxRepository,
		taxConfig:     taxConfig,
	}
}

func (tu *taxUsecase) GetTax(taxDate *domain.TaxDate) (*domain.TaxResponse, error) {
	taxResponse := &domain.TaxResponse{}
	summaries := []domain.TaxSummary{}
	aggregateFees := []domain.AggregateFee{}
	bankFee := 0
	beginDate := tax.RoundDay(taxDate.StartDate)
	endDate := beginDate + int64(taxDate.AmountOfDays*86400)

	if beginDate < 1393632000 { // if the start_time is on February 2014 special case
		for i := 1; i <= int(taxDate.AmountOfDays); i++ {
			summary := new(domain.TaxSummary)
			aggregateFee := new(domain.AggregateFee)
			summary.DayOfMonth = i + 14
			aggregateFee.DayOfMonth = i + 14
			summaries = append(summaries, *summary)
			aggregateFees = append(aggregateFees, *aggregateFee)
		}
	} else { // normal case
		for i := 1; i <= int(taxDate.AmountOfDays); i++ {
			summary := new(domain.TaxSummary)
			aggregateFee := new(domain.AggregateFee)
			summary.DayOfMonth = i
			aggregateFee.DayOfMonth = i
			summaries = append(summaries, *summary)
			aggregateFees = append(aggregateFees, *aggregateFee)
		}
	}
	dayToBeQueried := 1
	taxTransactionValid := true
	taxTransactionSummaries, err := tu.taxRepository.GetTaxTransactions(beginDate, endDate)
	if err != nil {
		return nil, err
	}
	if len(taxTransactionSummaries) == 0 {
		taxTransactionValid = false
	} else {
		if len(taxTransactionSummaries) < int(taxDate.AmountOfDays) {
			for _, tts := range taxTransactionSummaries {
				if int64(dayToBeQueried) < tts.DayOfMonth {
					dayToBeQueried = int(tts.DayOfMonth)
				}
			}
			dayToBeQueried += 1
			taxTransactionValid = false
		}
	}
	for _, serviceTax := range taxTransactionSummaries {
		summaries[serviceTax.DayOfMonth-1].DepositRp = serviceTax.DepositRp
		summaries[serviceTax.DayOfMonth-1].WithdrawRp = serviceTax.WithdrawRp
		summaries[serviceTax.DayOfMonth-1].Fee = serviceTax.Fee
		summaries[serviceTax.DayOfMonth-1].UplineBonus = serviceTax.UplineBonus
		summaries[serviceTax.DayOfMonth-1].Remain = serviceTax.Remain
		summaries[serviceTax.DayOfMonth-1].Ppn = serviceTax.Ppn
		taxResponse.TotalRevenue += serviceTax.Fee
		taxResponse.TotalBankFee += int64(bankFee)
		taxResponse.TotalUplineBonus += serviceTax.UplineBonus
		taxResponse.TotalRemain += serviceTax.Remain
		taxResponse.TotalPpn += serviceTax.Ppn
	}
	if taxTransactionValid { // only valid if data is fully existed on service database.
		return taxResponse, nil
	} else { // if the data is not fully available in service database, query to source database & save to service database.
		continueDate := (beginDate + (int64(dayToBeQueried-1) * 86400))
		taxResponseFromSource, err := tu.FetchSourceTax(&domain.TaxSourceDate{
			StartDate:    continueDate,
			EndDate:      endDate,
			StartDay:     dayToBeQueried,
			AmountOfDays: taxDate.AmountOfDays - (dayToBeQueried - 1),
		})
		if err != nil {
			return nil, err
		}

		taxTransactions := []entity.TaxTransaction{}
		for _, trfs := range taxResponseFromSource.Summary {
			if continueDate < 1393632000 {
				dayToBeQueried = 15
			}
			summaries[trfs.DayOfMonth-1].DepositRp = trfs.DepositRp
			summaries[trfs.DayOfMonth-1].WithdrawRp = trfs.WithdrawRp
			summaries[trfs.DayOfMonth-1].Fee = trfs.Fee
			summaries[trfs.DayOfMonth-1].UplineBonus = trfs.UplineBonus
			summaries[trfs.DayOfMonth-1].Remain = trfs.Remain
			summaries[trfs.DayOfMonth-1].Ppn = trfs.Ppn

			taxTransaction := new(entity.TaxTransaction)
			taxTransaction.TransactionDate = beginDate + (int64(trfs.DayOfMonth-1) * 86400)
			taxTransaction.DepositRp = trfs.DepositRp
			taxTransaction.WithdrawRp = trfs.WithdrawRp
			taxTransaction.Fee = trfs.Fee
			taxTransaction.UplineBonus = trfs.UplineBonus
			taxTransaction.Remain = trfs.Remain
			taxTransaction.Ppn = trfs.Ppn

			if (time.Now().Year() == time.Unix(beginDate, 0).Year()) && (time.Now().Month() == time.Unix(beginDate, 0).Month()) && (trfs.DayOfMonth >= time.Now().Day()) {
				continue
			} else {
				taxTransactions = append(taxTransactions, *taxTransaction)
			}
		}
		taxResponse.TotalRevenue += taxResponseFromSource.TotalRevenue
		taxResponse.TotalBankFee += taxResponseFromSource.TotalBankFee
		taxResponse.TotalUplineBonus += taxResponseFromSource.TotalUplineBonus
		taxResponse.TotalRemain += taxResponseFromSource.TotalRemain
		taxResponse.TotalPpn += taxResponseFromSource.TotalPpn
		if len(taxTransactions) > 0 {
			if err := tu.taxRepository.InsertTaxTransactions(continueDate, taxTransactions); err != nil {
				return nil, err
			}
		}
	}

	taxResponse.Summary = summaries
	return taxResponse, nil
}

func (tu *taxUsecase) FetchSourceTax(taxSourceDate *domain.TaxSourceDate) (*domain.TaxResponse, error) {
	taxResponse := &domain.TaxResponse{}
	var summaries []domain.TaxSummary
	var aggregateFees []domain.AggregateFee
	bankFee := 0
	if taxSourceDate.StartDate < 1393632000 { // if startDate is on February 2014 special calculation from 15 February (Special case)
		taxSourceDate.StartDay = 15
		for i := 1; i <= int(taxSourceDate.AmountOfDays); i++ {
			summary := new(domain.TaxSummary)
			aggregateFee := new(domain.AggregateFee)
			aggregateFee.DayOfMonth = int(taxSourceDate.StartDay) + (i - 1)
			summary.DayOfMonth = int(taxSourceDate.StartDay) + (i - 1)
			summaries = append(summaries, *summary)
			aggregateFees = append(aggregateFees, *aggregateFee)
		}
	} else {
		for i := 1; i <= int(taxSourceDate.AmountOfDays); i++ {
			summary := new(domain.TaxSummary)
			aggregateFee := new(domain.AggregateFee)
			aggregateFee.DayOfMonth = int(taxSourceDate.StartDay) + (i - 1)
			summary.DayOfMonth = int(taxSourceDate.StartDay) + (i - 1)
			summaries = append(summaries, *summary)
			aggregateFees = append(aggregateFees, *aggregateFee)
		}
	}

	depositRpTotalAmount, err := tu.taxRepository.GetDepositRpTotalAmount(taxSourceDate.StartDate, taxSourceDate.EndDate)
	if err != nil {
		return nil, err
	}
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func([]entity.DepositRpTotalAmount) {
		for _, depositRp := range depositRpTotalAmount {
			if !depositRp.DayOfMonth.Valid {
				continue
			}
			summaries[int(depositRp.DayOfMonth.Int64)-taxSourceDate.StartDay].DepositRp = depositRp.TotalAmount.Int64
		}
		wg.Done()
	}(depositRpTotalAmount)

	withdrawRpTotalAmount, err := tu.taxRepository.GetTotalWithdrawRp(taxSourceDate.StartDate, taxSourceDate.EndDate)
	if err != nil {
		return nil, err
	}
	wg.Add(1)
	go func([]entity.TotalWithdrawRp) {
		for _, withdrawRp := range withdrawRpTotalAmount {
			if !withdrawRp.DayOfMonth.Valid {
				continue
			}
			summaries[int(withdrawRp.DayOfMonth.Int64)-taxSourceDate.StartDay].DepositRp = withdrawRp.TotalRp.Int64
		}
	}(withdrawRpTotalAmount)

	if (taxSourceDate.StartDate <= 1662742800) && (taxSourceDate.EndDate >= 1662829199) {
		oldFees, err := tu.taxRepository.GetOldFees(taxSourceDate.StartDate, 1662742800)
		if err != nil {
			return nil, err
		}
		wg.Add(1)
		go func([]entity.TotalFee) {
			for _, oldFee := range oldFees {
				if !oldFee.DayOfMonth.Valid {
					continue
				}
				aggregateFees[int(oldFee.DayOfMonth.Int64)-taxSourceDate.StartDay].TotalFee = oldFee.TotalFee.Int64
				aggregateFees[int(oldFee.DayOfMonth.Int64)-taxSourceDate.StartDay].TotalRemain = oldFee.TotalRemain.Int64
				aggregateFees[int(oldFee.DayOfMonth.Int64)-taxSourceDate.StartDay].TotalUplineBonus = oldFee.TotalUplineBonus.Int64
			}
			wg.Done()
		}(oldFees)

		migrationNewFees, err := tu.taxRepository.GetFeesPerDay(1662742800, 1662829200)
		if err != nil {
			return nil, err
		}

		migrationOldFees, err := tu.taxRepository.GetOldFeesPerDay(1662742800, 1662829200)
		if err != nil {
			return nil, err
		}

		migrationFees := &domain.AggregateFee{
			DayOfMonth:       int(migrationNewFees.DayOfMonth.Int64),
			TotalFee:         (migrationNewFees.TotalFee.Int64 + migrationOldFees.TotalFee.Int64),
			TotalUplineBonus: (migrationNewFees.TotalUplineBonus.Int64 + migrationOldFees.TotalUplineBonus.Int64),
			TotalRemain:      (migrationNewFees.TotalRemain.Int64 + migrationOldFees.TotalRemain.Int64),
		}
		aggregateFees[migrationFees.DayOfMonth-taxSourceDate.StartDay].TotalFee = migrationFees.TotalFee
		aggregateFees[migrationFees.DayOfMonth-taxSourceDate.StartDay].TotalRemain = migrationFees.TotalRemain
		aggregateFees[migrationFees.DayOfMonth-taxSourceDate.StartDay].TotalUplineBonus = migrationFees.TotalUplineBonus

		newFees, err := tu.taxRepository.GetFees(1662829200, taxSourceDate.EndDate) // fees calculation with fees after migration date
		if err != nil {
			return nil, err
		}
		wg.Add(1)
		go func([]entity.TotalFee) {
			for _, newFee := range newFees {
				if !newFee.DayOfMonth.Valid {
					continue
				}
				aggregateFees[int(newFee.DayOfMonth.Int64)-taxSourceDate.StartDay].TotalFee = newFee.TotalFee.Int64
				aggregateFees[int(newFee.DayOfMonth.Int64)-taxSourceDate.StartDay].TotalUplineBonus = newFee.TotalUplineBonus.Int64
				aggregateFees[int(newFee.DayOfMonth.Int64)-taxSourceDate.StartDay].TotalRemain = newFee.TotalRemain.Int64
			}
			wg.Done()
		}(newFees)
	}
	wg.Wait()

	counterFees, err := tu.taxRepository.GetCounterFees(taxSourceDate.StartDate, taxSourceDate.EndDate)
	if err != nil {
		return nil, err
	}
	for _, counterFee := range counterFees {
		if !counterFee.DayOfMonth.Valid {
			continue
		}
		aggregateFees[int(counterFee.DayOfMonth.Int64)-taxSourceDate.StartDay].TotalFee += counterFee.TotalFee.Int64
		aggregateFees[int(counterFee.DayOfMonth.Int64)-taxSourceDate.StartDay].TotalRemain += counterFee.TotalFee.Int64 - int64(bankFee)
	}
	var ppn int64
	for _, aggregateFee := range aggregateFees {
		if taxSourceDate.StartDate+(int64(aggregateFee.DayOfMonth-taxSourceDate.StartDay)*86400) >= tu.taxConfig.TimeStartPpn {
			if (taxSourceDate.StartDate + (int64(aggregateFee.DayOfMonth-int(taxSourceDate.StartDay)) * 86400)) < tu.taxConfig.TimeStartPpnNew {
				ppn = int64(math.Ceil(float64(aggregateFee.TotalFee*tu.taxConfig.TarifPpn) / float64(100+tu.taxConfig.TarifPpn)))
			} else {
				ppn = int64(math.Ceil(float64(aggregateFee.TotalFee*tu.taxConfig.TarifPpnNew) / float64(100+tu.taxConfig.TarifPpnNew)))
			}
		}
		aggregateFee.TotalFee -= ppn
		aggregateFee.TotalRemain -= ppn
		summaries[aggregateFee.DayOfMonth-taxSourceDate.StartDay].Ppn = ppn
		summaries[aggregateFee.DayOfMonth-taxSourceDate.StartDay].Fee = aggregateFee.TotalFee
		summaries[aggregateFee.DayOfMonth-taxSourceDate.StartDay].UplineBonus = aggregateFee.TotalUplineBonus
		summaries[aggregateFee.DayOfMonth-taxSourceDate.StartDay].Remain = aggregateFee.TotalRemain

		taxResponse.TotalRevenue += aggregateFee.TotalFee
		taxResponse.TotalBankFee += int64(bankFee)
		taxResponse.TotalUplineBonus += aggregateFee.TotalUplineBonus
		taxResponse.TotalRemain += aggregateFee.TotalRemain
		taxResponse.TotalPpn += ppn
	}
	taxResponse.Summary = summaries
	return taxResponse, nil
}
