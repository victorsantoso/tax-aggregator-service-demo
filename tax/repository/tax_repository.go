package repository

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"tax-aggregator-service-demo/tax/domain"
	"tax-aggregator-service-demo/tax/entity"
)

type taxRepository struct {
	sourceConn  *sql.DB
	serviceConn *sql.DB
}

func NewTaxRepository(sourceConn, serviceConn *sql.DB) domain.TaxRepository {
	return &taxRepository{
		sourceConn:  sourceConn,
		serviceConn: serviceConn,
	}
}

// get deposit rp total amount query from source database.
const getDepositRpTotalAmount = `
	SELECT
		DATE FORMAT(CONVERT_TZ(FROM_UNIXTIME(success_time), 'UTC', 'Asia/Jakarta'), '%e') AS day_of_month,
		SUM(rp) AS total_rp,
		SUM(amount) AS total_amount,
		SUM(subsidi_fee) AS total_subsidi_fee
	FROM
		deposit_rp
	WHERE
		success_time >= ?
	AND
		success_time < ?
	GROUP BY
		day_of_month
	ORDER BY
		day_of_month
	ASC
`

func (tr *taxRepository) GetDepositRpTotalAmount(startDate, endDate int64) ([]entity.DepositRpTotalAmount, error) {
	sourceConn := tr.sourceConn
	depositRpTotalAmount := []entity.DepositRpTotalAmount{}
	query := func(query string, db *sql.DB) (*sql.Rows, error) {
		return db.Query(query, startDate, endDate)
	}
	r, err := query(getDepositRpTotalAmount, sourceConn)
	if err != nil {
		log.Println("[TaxRepository.GetDepositRpTotalAmount]:: server getting deposit_rp_total_amount from source database.")
		return nil, err
	}
	depositRpTotalAmountPerDay := &entity.DepositRpTotalAmount{}
	for r.Next() {
		if err := r.Scan(
			&depositRpTotalAmountPerDay.DayOfMonth,
			&depositRpTotalAmountPerDay.TotalRp,
			&depositRpTotalAmountPerDay.TotalAmount,
			&depositRpTotalAmountPerDay.TotalSubsidiFee,
		); err != nil {
			log.Println("[TaxRepository.GetDepositRpTotalAmount]:: error scanning deposit_rp_total_amount from source database.")
			return nil, err
		}
		depositRpTotalAmount = append(depositRpTotalAmount, *depositRpTotalAmountPerDay)
	}
	r.Close()
	return depositRpTotalAmount, nil
}

// get total withdraw rp query from source database.
const getTotalWithdrawRp = `
	SELECT
		DATE_FORMAT(CONVERT_TZ(FROM_UNIXTIME(success_time), 'UTC', 'Asia/Jakarta'), '%e') AS day_of_month,
		SUM(rp) AS total_rp
	FROM
		withdraw_rp
	WHERE
		success_time >= ?
	AND
		success_time < ?
	AND
		type != 'coupon'
	GROUP BY
		day_of_month
	ORDER BY
		day_of_month
	ASC
`

func (tr *taxRepository) GetTotalWithdrawRp(startDate, endDate int64) ([]entity.TotalWithdrawRp, error) {
	sourceConn := tr.sourceConn
	totalWithdrawRp := []entity.TotalWithdrawRp{}
	query := func(query string, db *sql.DB) (*sql.Rows, error) {
		return db.Query(query, startDate, endDate)
	}
	r, err := query(getTotalWithdrawRp, sourceConn)
	if err != nil {
		log.Println("[TaxRepository.GetTotalWithdrawRp]:: error on getting withdraw_rp_total_amount from source database.")
		return nil, err
	}
	totalWithdrawRpPerDay := &entity.TotalWithdrawRp{}
	for r.Next() {
		if err := r.Scan(
			&totalWithdrawRpPerDay.DayOfMonth,
			&totalWithdrawRpPerDay.TotalRp,
		); err != nil {
			log.Println("[TaxRepository.GetTotalWithdrawRp]:: error on scanning withdraw_rp_total_amount from source database.")
		}
		totalWithdrawRp = append(totalWithdrawRp, *totalWithdrawRpPerDay)
	}
	r.Close()
	return totalWithdrawRp, nil
}

// get fees query from source database.
const getFees = `
	SELECT
		DATE_FORMAT(CONVERT_TZ(FROM_UNIXTIME(waktu_transaksi), 'UTC', 'Asia/Jakarta'), '%e') AS day_of_month,
		SUM(fee) AS total_fee,
		SUM(upline_bonus) AS total_upline_bonus,
		SUM(remain) AS
	FROM
		fees
	WHERE
		waktu_transaksi >= ? AND waktu_transaksi < ?
	AND
		type NOT IN('deposit', 'tax')
	AND
		upline_id != 1
	GROUP BY
		day_of_month
	ORDER BY
		day_of_month
	ASC
`

func (tr *taxRepository) GetFees(startDate, endDate int64) ([]entity.TotalFee, error) {
	sourceConn := tr.sourceConn
	totalFees := []entity.TotalFee{}
	query := func(query string, db *sql.DB) (*sql.Rows, error) {
		return db.Query(query, startDate, endDate)
	}
	r, err := query(getFees, sourceConn)
	if err != nil {
		log.Println("[TaxRepository.GetFees]:: error getting total_fee from source database.")
		return nil, err
	}
	totalFeesPerDay := &entity.TotalFee{}
	for r.Next() {
		if err := r.Scan(
			&totalFeesPerDay.DayOfMonth,
			&totalFeesPerDay.TotalFee,
			&totalFeesPerDay.TotalUplineBonus,
			&totalFeesPerDay.TotalRemain,
		); err != nil {
			log.Println("[TaxRepository.GetFees]:: error scanning fees from source database.")
			return nil, err
		}
		totalFees = append(totalFees, *totalFeesPerDay)
	}
	r.Close()
	return totalFees, nil
}

// get old fees query from source database.
const getOldFees = `
	SELECT
		DATE_FORMAT(CONVERT_TZ(FROM_UNIXTIME(waktu_transaksi), 'UTC', 'Asia/Jakarta'), '%e') AS day_of_month,
		SUM(fee) AS total_fee,
		SUM(upline_bonus) AS total_upline_bonus,
		SUM(remain) AS total_remain
	FROM
		fees_old
	WHERE
		waktu_transaksi >= ? AND waktu_transaksi < ?
	AND
		type NOT IN('deposit', 'tax')
	AND
		upline_id != 1
	GROUP BY
		day_of_month
	ORDER BY
		day_of_month
	ASC
`

func (tr *taxRepository) GetOldFees(startDate, endDate int64) ([]entity.TotalFee, error) {
	sourceConn := tr.sourceConn
	totalFees := []entity.TotalFee{}
	query := func(query string, db *sql.DB) (*sql.Rows, error) {
		return db.Query(query, startDate, endDate)
	}
	r, err := query(getOldFees, sourceConn)
	if err != nil {
		log.Println("[TaxRepository.GetOldFees]:: error getting total_fee from source database.")
		return nil, err
	}
	totalFeesPerDay := &entity.TotalFee{}
	for r.Next() {
		if err := r.Scan(
			&totalFeesPerDay.DayOfMonth,
			&totalFeesPerDay.TotalFee,
			&totalFeesPerDay.TotalUplineBonus,
			&totalFeesPerDay.TotalRemain,
		); err != nil {
			log.Println("[TaxRepository.GetOldFees]:: error scanning total_fee from source database.")
			return nil, err
		}
		totalFees = append(totalFees, *totalFeesPerDay)
	}
	r.Close()
	return totalFees, nil
}

// get counter fees query from source database.
const getCounterFees = `
	SELECT
		DATE_FORMAT(CONVERT_TZ(FROM_UNIXTIME(success_time), 'UTC', 'Asia/Jakarta'), '%e') AS day_of_month,
		SUM(fee) AS total_fee
	FROM
		counter_buy_btc
	WHERE
		status = 'success'
	AND
		success_time >= ? AND success_time < ?
	GROUP BY
		day_of_month
	ORDER BY
		day_of_month
	ASC
`

func (tr *taxRepository) GetCounterFees(startDate, endDate int64) ([]entity.CounterFee, error) {
	sourceConn := tr.sourceConn
	counterFees := []entity.CounterFee{}
	query := func(query string, db *sql.DB) (*sql.Rows, error) {
		return db.Query(query, startDate, endDate)
	}
	r, err := query(getCounterFees, sourceConn)
	if err != nil {
		log.Println("[TaxRepository.GetCounterFees]:: error getting counter_fee from source database.")
		return nil, err
	}
	couterFeesPerDay := &entity.CounterFee{}
	for r.Next() {
		if err := r.Scan(
			&couterFeesPerDay.DayOfMonth,
			&couterFeesPerDay.TotalFee,
		); err != nil {
			log.Println("[TaxRepository.GetCounterFees]:: error scanning counter_fee from source database.")
			return nil, err
		}
		counterFees = append(counterFees, *couterFeesPerDay)
	}
	r.Close()
	return counterFees, nil
}

func (tr *taxRepository) GetFeesPerDay(startDate, endDate int64) (*entity.TotalFee, error) {
	sourceConn := tr.sourceConn
	totalFees := new(entity.TotalFee)
	query := func(query string, db *sql.DB) (*sql.Rows, error) {
		return db.Query(query, startDate, endDate)
	}
	r, err := query(getFees, sourceConn)
	if err != nil {
		log.Println("[TaxRepository.GetFeesPerDay]:: error getting total_fee per day from source database.")
		return nil, err
	}
	for r.Next() {
		if err := r.Scan(
			&totalFees.DayOfMonth,
			&totalFees.TotalFee,
			&totalFees.TotalUplineBonus,
			&totalFees.TotalRemain,
		); err != nil {
			log.Println("[TaxRepository.GetFeesPerDay]:: error scanning total_fee per day from source database.")
			return nil, err
		}
	}
	r.Close()
	return totalFees, nil
}

func (tr *taxRepository) GetOldFeesPerDay(startDate, endDate int64) (*entity.TotalFee, error) {
	sourceConn := tr.sourceConn
	totalFees := new(entity.TotalFee)
	query := func(query string, db *sql.DB) (*sql.Rows, error) {
		return db.Query(query, startDate, endDate)
	}
	r, err := query(getOldFees, sourceConn)
	if err != nil {
		log.Println("[TaxRepository.GetOldFeesPerDay]:: error getting total_fee per day from source database.")
		return nil, err
	}
	for r.Next() {
		if err := r.Scan(
			&totalFees.DayOfMonth,
			&totalFees.TotalFee,
			&totalFees.TotalUplineBonus,
			&totalFees.TotalRemain,
		); err != nil {
			log.Println("[TaxRepository.GetOldFeesPerDay]:: error scanning total_fee per day from source database.")
			return nil, err
		}
	}
	r.Close()
	return totalFees, nil
}

// get tax transactions query from service database.
const getTaxTransactions = `
	SELECT
		DATE_PART('day', TO_TIMESTAMP(t.transaction_date)) AS day_of_month,
		t.deposit_rp,
		t.withdraw_rp,
		t.fee,
		t.upline_bonus,
		t.remain,
		t.ppn
	FROM
		tax_transaction AS t
	WHERE
		t.transaction_date >= $1 AND t.transaction_date < $2
`

func (tr *taxRepository) GetTaxTransactions(startDate, endDate int64) ([]entity.TaxTransactionSummary, error) {
	serviceConn := tr.serviceConn
	taxTransactionSummaries := []entity.TaxTransactionSummary{}
	query := func(query string, db *sql.DB) (*sql.Rows, error) {
		return db.Query(query, startDate, endDate)
	}
	r, err := query(getTaxTransactions, serviceConn)
	if err != nil {
		log.Println("[TaxRepository.GetTaxTransactions]:: error getting tax_transactions from service database.")
		return nil, err
	}
	taxTransactionPerDay := &entity.TaxTransactionSummary{}
	for r.Next() {
		if err := r.Scan(
			&taxTransactionPerDay.DayOfMonth,
			&taxTransactionPerDay.DepositRp,
			&taxTransactionPerDay.WithdrawRp,
			&taxTransactionPerDay.Fee,
			&taxTransactionPerDay.UplineBonus,
			&taxTransactionPerDay.Remain,
			&taxTransactionPerDay.Ppn,
		); err != nil {
			log.Println("[TaxRepository.GetTaxTransactions]:: error scanning tax_transactions from service database.")
			return nil, err
		}
		taxTransactionSummaries = append(taxTransactionSummaries, *taxTransactionPerDay)
	}
	r.Close()
	return taxTransactionSummaries, nil
}

// insert tax transaction query from service database.
const insertTaxTransaction = `
	INSERT INTO
		tax_transaction(transaction_date, deposit_rp, withdraw_rp, fee, upline_bonus, remain, ppn)
	VALUES
`

func (tr *taxRepository) InsertTaxTransactions(transactionDate int64, taxTransactions []entity.TaxTransaction) error {
	serviceConn := tr.serviceConn
	tx, err := serviceConn.Begin()
	if err != nil {
		log.Println("[TaxRepository.InsertTaxTransaction]:: error begin database transaction in service database.")
		tx.Rollback()
		return err
	}
	var inserts []string
	var args []interface{}
	var begin int64 = 1
	for _, v := range taxTransactions {
		queryStr := fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d)", begin, begin+1, begin+2, begin+3, begin+4, begin+5, begin+6)
		begin += 7
		inserts = append(inserts, queryStr)
		args = append(args, (v.TransactionDate), v.DepositRp, v.WithdrawRp, v.Fee, v.UplineBonus, v.Remain, v.Ppn)
	}
	queryVals := strings.Join(inserts, ",")
	query := insertTaxTransaction + queryVals
	res, err := tx.Exec(query, args...)
	if err != nil {
		log.Println("[TaxRepository.InsertTaxTransactions]:: error insert tax_transaction.")
		tx.Rollback()
	}
	rows, err := res.RowsAffected()
	if err != nil {
		log.Println("[TaxRepository.InsertTaxTransactions]:: error when finding rows.")
		return err
	}
	tx.Commit()
	log.Printf("[TaxRepository.InsertTaxTransactions]:: created %d tax_transactions simultaneously.\n", rows)
	return nil
}
