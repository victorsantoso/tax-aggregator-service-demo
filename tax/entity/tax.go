package entity

import "database/sql"

// Value Objects and Entities that will be mapped into service_database

type CounterFee struct {
	DayOfMonth sql.NullInt64 `json:"day_of_month"`
	TotalFee   sql.NullInt64 `json:"total_fee"`
}

type DepositRpTotalAmount struct {
	DayOfMonth      sql.NullInt64 `json:"day_of_month"`
	TotalRp         sql.NullInt64 `json:"total_rp"`
	TotalAmount     sql.NullInt64 `json:"total_amount"`
	TotalSubsidiFee sql.NullInt64 `json:"total_subsidi_fee"`
}

type TotalFee struct {
	DayOfMonth       sql.NullInt64 `json:"day_of_month"`
	TotalFee         sql.NullInt64 `json:"total_fee"`
	TotalUplineBonus sql.NullInt64 `json:"total_upline_bonus"`
	TotalRemain      sql.NullInt64 `json:"total_remain"`
}

type TotalWithdrawRp struct {
	DayOfMonth sql.NullInt64 `json:"day_of_month"`
	TotalRp    sql.NullInt64 `json:"total_rp"`
}

type TaxTransaction struct {
	TransactionDate int64 `json:"transaction_date"`
	DepositRp       int64 `json:"deposit_rp"`
	WithdrawRp      int64 `json:"withdraw_rp"`
	Fee             int64 `json:"fee"`
	UplineBonus     int64 `json:"upline_bonus"`
	Remain          int64 `json:"remain"`
	Ppn             int64 `json:"ppn"`
}

type TaxTransactionSummary struct {
	DayOfMonth  int64 `json:"day_of_month"`
	DepositRp   int64 `json:"deposit_rp"`
	WithdrawRp  int64 `json:"withdraw_rp"`
	Fee         int64 `json:"fee"`
	UplineBonus int64 `json:"upline_bonus"`
	Remain      int64 `json:"remain"`
	Ppn         int64 `json:"ppn"`
}
