package repository

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestTaxRepository_GetDepositRpTotalAmount(t *testing.T) {
	type args struct {
		startDate int64
		endDate   int64
	}
	tests := []struct {
		name         string
		args         args
		testFunction func(t *testing.T, tt args)
	}{
		{
			name: "test get deposit total amount with 0 data",
			args: args{
				startDate: 1680321600,
				endDate:   1682852400,
			},
			testFunction: func(t *testing.T, tt args) {
				sourceConn, sourceMock, err := sqlmock.New()
				assert.NoError(t, err)
				serviceConn, _, err := sqlmock.New()
				assert.NoError(t, err)
				rows := sqlmock.NewRows([]string{"day_of_month", "total_rp", "total_amount", "total_subsidi_fee"})
				sourceMock.ExpectQuery(regexp.QuoteMeta(getDepositRpTotalAmount)).WillReturnRows(rows)
				taxRepository := NewTaxRepository(sourceConn, serviceConn)
				depositRpTotalAmount, err := taxRepository.GetDepositRpTotalAmount(tt.startDate, tt.endDate)
				assert.NoError(t, err)
				assert.Empty(t, depositRpTotalAmount)
			},
		},
		{
			name: "test get deposit rp total amount success",
			args: args{
				startDate: 1680321600,
				endDate:   1682852400,
			},
			testFunction: func(t *testing.T, tt args) {
				sourceConn, sourceMock, err := sqlmock.New()
				assert.NoError(t, err)
				serviceConn, _, err := sqlmock.New()
				assert.NoError(t, err)
				rows := sqlmock.NewRows([]string{"day_of_month", "total_rp", "total_amount", "total_subsidi_fee"})
				rows.AddRow("1", "3403357", "3403357", "10")
				rows.AddRow("2", "3403358", "3403357", "20")
				sourceMock.ExpectQuery(regexp.QuoteMeta(getDepositRpTotalAmount)).WillReturnRows(rows)
				taxRepository := NewTaxRepository(sourceConn, serviceConn)
				depositRpTotalAmount, err := taxRepository.GetDepositRpTotalAmount(tt.startDate, tt.endDate)
				assert.NoError(t, err)
				assert.NotNil(t, depositRpTotalAmount)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.testFunction(t, tt.args)
		})
	}
}

func TestTaxRepository_GetTotalWithdrawRp(t *testing.T) {
	type args struct {
		startDate int64
		endDate   int64
	}
	tests := []struct {
		name         string
		args         args
		testFunction func(t *testing.T, tt args)
	}{
		{
			name: "test get total withdraw rp with 0 data",
			args: args{
				startDate: 1680321600,
				endDate:   1682852400,
			},
			testFunction: func(t *testing.T, tt args) {
				sourceConn, sourceMock, err := sqlmock.New()
				assert.NoError(t, err)
				serviceConn, _, err := sqlmock.New()
				assert.NoError(t, err)
				rows := sqlmock.NewRows([]string{"day_of_month", "total_rp"})
				sourceMock.ExpectQuery(regexp.QuoteMeta(getTotalWithdrawRp)).WillReturnRows(rows)
				taxRepository := NewTaxRepository(sourceConn, serviceConn)
				totalWithdrawRp, err := taxRepository.GetTotalWithdrawRp(tt.startDate, tt.endDate)
				assert.NoError(t, err)
				assert.Empty(t, totalWithdrawRp)
			},
		},
		{
			name: "test get total withdraw rp success",
			args: args{
				startDate: 1680321600,
				endDate:   1682852400,
			},
			testFunction: func(t *testing.T, tt args) {
				sourceConn, sourceMock, err := sqlmock.New()
				assert.NoError(t, err)
				serviceConn, _, err := sqlmock.New()
				assert.NoError(t, err)
				rows := sqlmock.NewRows([]string{"day_of_month", "total_rp"})
				rows.AddRow("1", "3403357")
				sourceMock.ExpectQuery(regexp.QuoteMeta(getTotalWithdrawRp)).WillReturnRows(rows)
				taxRepository := NewTaxRepository(sourceConn, serviceConn)
				totalWithdrawRp, err := taxRepository.GetTotalWithdrawRp(tt.startDate, tt.endDate)
				assert.NoError(t, err)
				assert.NotNil(t, totalWithdrawRp)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.testFunction(t, tt.args)
		})
	}
}

func TestTaxRepository_GetFees(t *testing.T) {
	type args struct {
		startDate int64
		endDate   int64
	}
	tests := []struct {
		name         string
		args         args
		testFunction func(t *testing.T, tt args)
	}{
		{
			name: "test get fees with 0 data",
			args: args{
				startDate: 1680321600,
				endDate:   1682852400,
			},
			testFunction: func(t *testing.T, tt args) {
				sourceConn, sourceMock, err := sqlmock.New()
				assert.NoError(t, err)
				serviceConn, _, err := sqlmock.New()
				assert.NoError(t, err)
				rows := sqlmock.NewRows([]string{"day_of_month", "total_fee", "total_upline_bonus", "total_remain"})
				sourceMock.ExpectQuery(regexp.QuoteMeta(getFees)).WillReturnRows(rows)
				taxRepository := NewTaxRepository(sourceConn, serviceConn)
				fees, err := taxRepository.GetFees(tt.startDate, tt.endDate)
				assert.NoError(t, err)
				assert.Empty(t, fees)
			},
		},
		{
			name: "test get fees success",
			args: args{
				startDate: 1680321600,
				endDate:   1682852400,
			},
			testFunction: func(t *testing.T, tt args) {
				sourceConn, sourceMock, err := sqlmock.New()
				assert.NoError(t, err)
				serviceConn, _, err := sqlmock.New()
				assert.NoError(t, err)
				rows := sqlmock.NewRows([]string{"day_of_month", "total_fee", "total_upline_bonus", "total_remain"})
				rows.AddRow("1", "10000", "20000", "20000")
				rows.AddRow("2", "20000", "40000", "40000")
				sourceMock.ExpectQuery(regexp.QuoteMeta(getFees)).WillReturnRows(rows)
				taxRepository := NewTaxRepository(sourceConn, serviceConn)
				fees, err := taxRepository.GetFees(tt.startDate, tt.endDate)
				assert.NoError(t, err)
				assert.NotNil(t, fees)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.testFunction(t, tt.args)
		})
	}
}

func TestTaxRepository_GetOldFees(t *testing.T) {
	type args struct {
		startDate int64
		endDate   int64
	}
	tests := []struct {
		name         string
		args         args
		testFunction func(t *testing.T, tt args)
	}{
		{
			name: "test get old fees with 0 data",
			args: args{
				startDate: 1680321600,
				endDate:   1682852400,
			},
			testFunction: func(t *testing.T, tt args) {
				sourceConn, sourceMock, err := sqlmock.New()
				assert.NoError(t, err)
				serviceConn, _, err := sqlmock.New()
				assert.NoError(t, err)
				rows := sqlmock.NewRows([]string{"day_of_month", "total_fee", "total_upline_bonus", "total_remain"})
				sourceMock.ExpectQuery(regexp.QuoteMeta(getOldFees)).WillReturnRows(rows)
				taxRepository := NewTaxRepository(sourceConn, serviceConn)
				oldFees, err := taxRepository.GetOldFees(tt.startDate, tt.endDate)
				assert.NoError(t, err)
				assert.Empty(t, oldFees)
			},
		},
		{
			name: "test get fees success",
			args: args{
				startDate: 1680321600,
				endDate:   1682852400,
			},
			testFunction: func(t *testing.T, tt args) {
				sourceConn, sourceMock, err := sqlmock.New()
				assert.NoError(t, err)
				serviceConn, _, err := sqlmock.New()
				assert.NoError(t, err)
				rows := sqlmock.NewRows([]string{"day_of_month", "total_fee", "total_upline_bonus", "total_remain"})
				rows.AddRow("1", "10000", "20000", "20000")
				rows.AddRow("2", "20000", "40000", "40000")
				sourceMock.ExpectQuery(regexp.QuoteMeta(getFees)).WillReturnRows(rows)
				taxRepository := NewTaxRepository(sourceConn, serviceConn)
				oldFees, err := taxRepository.GetFees(tt.startDate, tt.endDate)
				assert.NoError(t, err)
				assert.NotNil(t, oldFees)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.testFunction(t, tt.args)
		})
	}
}

func TestTaxRepository_GetCounterFees(t *testing.T) {
	type args struct {
		startDate int64
		endDate   int64
	}
	tests := []struct {
		name         string
		args         args
		testFunction func(t *testing.T, tt args)
	}{
		{
			name: "test get counter fees with 0 data",
			args: args{
				startDate: 1680321600,
				endDate:   1682852400,
			},
			testFunction: func(t *testing.T, tt args) {
				sourceConn, sourceMock, err := sqlmock.New()
				assert.NoError(t, err)
				serviceConn, _, err := sqlmock.New()
				assert.NoError(t, err)
				rows := sqlmock.NewRows([]string{"day_of_month", "total_fee"})
				sourceMock.ExpectQuery(regexp.QuoteMeta(getCounterFees)).WillReturnRows(rows)
				taxRepository := NewTaxRepository(sourceConn, serviceConn)
				counterFees, err := taxRepository.GetCounterFees(tt.startDate, tt.endDate)
				assert.NoError(t, err)
				assert.Empty(t, counterFees)
			},
		},
		{
			name: "test get counter fees success",
			args: args{
				startDate: 1680321600,
				endDate:   1682852400,
			},
			testFunction: func(t *testing.T, tt args) {
				sourceConn, sourceMock, err := sqlmock.New()
				assert.NoError(t, err)
				serviceConn, _, err := sqlmock.New()
				assert.NoError(t, err)
				rows := sqlmock.NewRows([]string{"day_of_month", "total_fee"})
				rows.AddRow("1", "10000")
				rows.AddRow("2", "20000")
				sourceMock.ExpectQuery(regexp.QuoteMeta(getCounterFees)).WillReturnRows(rows)
				taxRepository := NewTaxRepository(sourceConn, serviceConn)
				counterFees, err := taxRepository.GetCounterFees(tt.startDate, tt.endDate)
				assert.NoError(t, err)
				assert.NotNil(t, counterFees)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.testFunction(t, tt.args)
		})
	}
}

func TestTaxRepository_GetFeesPerDay(t *testing.T) {
	type args struct {
		startDate int64
		endDate   int64
	}
	tests := []struct {
		name         string
		args         args
		testFunction func(t *testing.T, tt args)
	}{
		{
			name: "test get fees per day with 0 data",
			args: args{
				startDate: 1680321600,
				endDate:   1682852400,
			},
			testFunction: func(t *testing.T, tt args) {
				sourceConn, sourceMock, err := sqlmock.New()
				assert.NoError(t, err)
				serviceConn, _, err := sqlmock.New()
				assert.NoError(t, err)
				rows := sqlmock.NewRows([]string{"day_of_month", "total_fee", "total_upline_bonus", "total_remain"})
				sourceMock.ExpectQuery(regexp.QuoteMeta(getFees)).WillReturnRows(rows)
				taxRepository := NewTaxRepository(sourceConn, serviceConn)
				fees, err := taxRepository.GetFeesPerDay(tt.startDate, tt.endDate)
				assert.NoError(t, err)
				assert.Empty(t, fees)
			},
		},
		{
			name: "test get fees per day success",
			args: args{
				startDate: 1680321600,
				endDate:   1682852400,
			},
			testFunction: func(t *testing.T, tt args) {
				sourceConn, sourceMock, err := sqlmock.New()
				assert.NoError(t, err)
				serviceConn, _, err := sqlmock.New()
				assert.NoError(t, err)
				rows := sqlmock.NewRows([]string{"day_of_month", "total_fee", "total_upline_bonus", "total_remain"})
				rows.AddRow("1", "100000", "20000", "20000")
				sourceMock.ExpectQuery(regexp.QuoteMeta(getFees)).WillReturnRows(rows)
				taxRepository := NewTaxRepository(sourceConn, serviceConn)
				fees, err := taxRepository.GetFeesPerDay(tt.startDate, tt.endDate)
				assert.NoError(t, err)
				assert.NotNil(t, fees)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.testFunction(t, tt.args)
		})
	}
}

func TestTaxRepository_GetOldFeesPerDay(t *testing.T) {
	type args struct {
		startDate int64
		endDate   int64
	}
	tests := []struct {
		name         string
		args         args
		testFunction func(t *testing.T, tt args)
	}{
		{
			name: "test get old fees per day with 0 data",
			args: args{
				startDate: 1680321600,
				endDate:   1682852400,
			},
			testFunction: func(t *testing.T, tt args) {
				sourceConn, sourceMock, err := sqlmock.New()
				assert.NoError(t, err)
				serviceConn, _, err := sqlmock.New()
				assert.NoError(t, err)
				rows := sqlmock.NewRows([]string{"day_of_month", "total_fee", "total_upline_bonus", "total_remain"})
				sourceMock.ExpectQuery(regexp.QuoteMeta(getOldFees)).WillReturnRows(rows)
				taxRepository := NewTaxRepository(sourceConn, serviceConn)
				oldFees, err := taxRepository.GetOldFees(tt.startDate, tt.endDate)
				assert.NoError(t, err)
				assert.Empty(t, oldFees)
			},
		},
		{
			name: "test get old fees per day success",
			args: args{
				startDate: 1680321600,
				endDate:   1682852400,
			},
			testFunction: func(t *testing.T, tt args) {
				sourceConn, sourceMock, err := sqlmock.New()
				assert.NoError(t, err)
				serviceConn, _, err := sqlmock.New()
				assert.NoError(t, err)
				rows := sqlmock.NewRows([]string{"day_of_month", "total_fee", "total_upline_bonus", "total_remain"})
				rows.AddRow("1", "100000", "20000", "20000")
				sourceMock.ExpectQuery(regexp.QuoteMeta(getOldFees)).WillReturnRows(rows)
				taxRepository := NewTaxRepository(sourceConn, serviceConn)
				oldFees, err := taxRepository.GetOldFeesPerDay(tt.startDate, tt.endDate)
				assert.NoError(t, err)
				assert.NotNil(t, oldFees)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.testFunction(t, tt.args)
		})
	}
}

func TestTaxRepository_GetTaxTransaction(t *testing.T) {
	type args struct {
		startDate int64
		endDate   int64
	}
	tests := []struct {
		name         string
		args         args
		testFunction func(t *testing.T, tt args)
	}{
		{
			name: "test get tax transaction from service database with 0 data",
			args: args{
				startDate: 1680321600,
				endDate:   1682852400,
			},
			testFunction: func(t *testing.T, tt args) {
				sourceConn, _, err := sqlmock.New()
				assert.NoError(t, err)
				serviceConn, serviceMock, err := sqlmock.New()
				assert.NoError(t, err)
				rows := sqlmock.NewRows([]string{"day_of_month", "deposit_rp", "withdraw_rp", "fee", "upline_bonus", "remain", "ppn"})
				serviceMock.ExpectQuery(regexp.QuoteMeta(getTaxTransactions)).WillReturnRows(rows)
				taxRepository := NewTaxRepository(sourceConn, serviceConn)
				taxTransactions, err := taxRepository.GetTaxTransactions(tt.startDate, tt.endDate)
				assert.NoError(t, err)
				assert.Empty(t, taxTransactions)
			},
		},
		{
			name: "test get tax transaction success",
			args: args{
				startDate: 1680321600,
				endDate:   1682852400,
			},
			testFunction: func(t *testing.T, tt args) {
				sourceConn, _, err := sqlmock.New()
				assert.NoError(t, err)
				serviceConn, serviceMock, err := sqlmock.New()
				assert.NoError(t, err)
				rows := sqlmock.NewRows([]string{"day_of_month", "deposit_rp", "withdraw_rp", "fee", "upline_bonus", "remain", "ppn"})
				rows.AddRow("1", "1000000000", "500000000", "300000000", "30000", "30000", "200000")
				serviceMock.ExpectQuery(regexp.QuoteMeta(getTaxTransactions)).WillReturnRows(rows)
				taxRepository := NewTaxRepository(sourceConn, serviceConn)
				taxTransactions, err := taxRepository.GetTaxTransactions(tt.startDate, tt.endDate)
				assert.NoError(t, err)
				assert.NotNil(t, taxTransactions)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.testFunction(t, tt.args)
		})
	}
}
