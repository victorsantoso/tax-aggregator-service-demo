package domain

import (
	"tax-aggregator-service-demo/tax/entity"

	"github.com/labstack/echo/v4"
)

// tax handler interface
type TaxHandler interface {
	Routes(route *echo.Echo)
	GetTax(ctx echo.Context) error
}

// tax configuration from monolith application, this config can be moved into service config like config.json
type TaxConfig struct {
	TimeStartPpn    int64 `query:"time_start_ppn"`
	TimeStartPpnNew int64 `query:"time_start_ppn_new"`
	TarifPpn        int64 `query:"tarif_ppn"`
	TarifPpnNew     int64 `query:"tarif_ppn_new"`
}

type Response struct {
	Data    any    `json:"data"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

type TaxDate struct {
	StartDate    int64
	EndDate      int64
	StartDay     int
	AmountOfDays int
}

type TaxSourceDate struct {
	StartDate    int64
	EndDate      int64
	StartDay     int
	AmountOfDays int
}

// tax usecase interface contract for business layer level in tax usecase
type TaxUsecase interface {
	GetTax(taxDate *TaxDate) (*TaxResponse, error)
	FetchSourceTax(taxSourceDate *TaxSourceDate) (*TaxResponse, error)
}

// tax response for tax_usecase from business layer in tax usecase
type TaxResponse struct {
	Summary          []TaxSummary `json:"summary"`
	TotalRevenue     int64        `json:"total_revenue"`
	TotalBankFee     int64        `json:"total_bank_fee"`
	TotalUplineBonus int64        `json:"total_upline_bonus"`
	TotalRemain      int64        `json:"total_remain"`
	TotalPpn         int64        `json:"total_ppn"`
}

// tax summary for tax bounded context
type TaxSummary struct {
	DepositRp   int64 `json:"deposit_rp"`
	WithdrawRp  int64 `json:"withdraw_rp"`
	Fee         int64 `json:"fee"`
	UplineBonus int64 `json:"upline_bonus"`
	Remain      int64 `json:"remain"`
	Ppn         int64 `json:"ppn"`
	DayOfMonth  int   `json:"day_of_month"`
}

// aggregate fee for tax bounded context
type AggregateFee struct {
	TotalFee         int64 `json:"total_fee"`
	TotalUplineBonus int64 `json:"total_upline_bonus"`
	TotalRemain      int64 `json:"total_remain"`
	DayOfMonth       int   `json:"day_of_month"`
}

// tax repository interface contract for repository layer
type TaxRepository interface {
	GetDepositRpTotalAmount(startDate, endDate int64) ([]entity.DepositRpTotalAmount, error)
	GetTotalWithdrawRp(startDate, calculationDate int64) ([]entity.TotalWithdrawRp, error)
	GetFees(startDate, endDate int64) ([]entity.TotalFee, error)
	GetOldFees(startDate, endDate int64) ([]entity.TotalFee, error)
	GetCounterFees(startDate, endDate int64) ([]entity.CounterFee, error)
	GetFeesPerDay(startTime, endTime int64) (*entity.TotalFee, error)
	GetOldFeesPerDay(startTime, endTime int64) (*entity.TotalFee, error)

	GetTaxTransactions(startDate, endDate int64) ([]entity.TaxTransactionSummary, error)
	InsertTaxTransactions(transactionDate int64, taxTransactions []entity.TaxTransaction) error
}
