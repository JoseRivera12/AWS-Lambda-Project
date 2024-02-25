package transactions

import (
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"

	fileutils "transaction-mailer/internal/fileutils"
)

type testCase struct {
	name     string
	records  []*fileutils.Record
	expected *TransactionsData
}

func TestNewTransactionsData(t *testing.T) {
	currentYear := time.Now().Year()

	testCases := []testCase{
		{
			name:    "No Transactions",
			records: []*fileutils.Record{},
			expected: &TransactionsData{
				TotalBalance:         decimal.Zero,
				AvgDebitAmount:       decimal.Zero,
				AvgCreditAmount:      decimal.Zero,
				TransactionsPerMonth: make([]int, 12),
			},
		},
		{
			name: "Mixed transactions",
			records: []*fileutils.Record{
				{Date: fileutils.DateTime{Time: time.Date(currentYear, 1, 1, 0, 0, 0, 0, time.Local)}, Transaction: decimal.NewFromInt(100)},
				{Date: fileutils.DateTime{Time: time.Date(currentYear, 2, 2, 0, 0, 0, 0, time.Local)}, Transaction: decimal.NewFromInt(-50)},
				{Date: fileutils.DateTime{Time: time.Date(currentYear, 3, 3, 0, 0, 0, 0, time.Local)}, Transaction: decimal.NewFromInt(200)},
				{Date: fileutils.DateTime{Time: time.Date(currentYear, 12, 3, 0, 0, 0, 0, time.Local)}, Transaction: decimal.NewFromFloat(-20.5)},
				{Date: fileutils.DateTime{Time: time.Date(currentYear, 3, 3, 0, 0, 0, 0, time.Local)}, Transaction: decimal.NewFromFloat(-10)},
			},
			expected: &TransactionsData{
				TotalBalance:         decimal.New(2195, -1),
				AvgDebitAmount:       decimal.New(-268333333333333333, -16),
				AvgCreditAmount:      decimal.New(1500000000000000000, -16),
				TransactionsPerMonth: []int{1, 1, 2, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			},
		},
		{
			name: "Only credit transactions",
			records: []*fileutils.Record{
				{Date: fileutils.DateTime{Time: time.Date(currentYear, 1, 1, 0, 0, 0, 0, time.Local)}, Transaction: decimal.NewFromInt(100)},
				{Date: fileutils.DateTime{Time: time.Date(currentYear, 2, 2, 0, 0, 0, 0, time.Local)}, Transaction: decimal.NewFromInt(50)},
				{Date: fileutils.DateTime{Time: time.Date(currentYear, 3, 3, 0, 0, 0, 0, time.Local)}, Transaction: decimal.NewFromInt(200)},
				{Date: fileutils.DateTime{Time: time.Date(currentYear, 12, 3, 0, 0, 0, 0, time.Local)}, Transaction: decimal.NewFromFloat(20.5)},
				{Date: fileutils.DateTime{Time: time.Date(currentYear, 3, 3, 0, 0, 0, 0, time.Local)}, Transaction: decimal.NewFromFloat(+10)},
			},
			expected: &TransactionsData{
				TotalBalance:         decimal.New(3805, -1),
				AvgDebitAmount:       decimal.Zero,
				AvgCreditAmount:      decimal.New(761000000000000000, -16),
				TransactionsPerMonth: []int{1, 1, 2, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			},
		},
		{
			name: "Only debit transactions",
			records: []*fileutils.Record{
				{Date: fileutils.DateTime{Time: time.Date(currentYear, 1, 1, 0, 0, 0, 0, time.Local)}, Transaction: decimal.NewFromInt(-190)},
				{Date: fileutils.DateTime{Time: time.Date(currentYear, 2, 2, 0, 0, 0, 0, time.Local)}, Transaction: decimal.NewFromInt(-52)},
				{Date: fileutils.DateTime{Time: time.Date(currentYear, 3, 3, 0, 0, 0, 0, time.Local)}, Transaction: decimal.NewFromInt(-43)},
				{Date: fileutils.DateTime{Time: time.Date(currentYear, 12, 3, 0, 0, 0, 0, time.Local)}, Transaction: decimal.NewFromFloat(-20.5)},
				{Date: fileutils.DateTime{Time: time.Date(currentYear, 3, 3, 0, 0, 0, 0, time.Local)}, Transaction: decimal.NewFromFloat(-10)},
			},
			expected: &TransactionsData{
				TotalBalance:         decimal.New(-3155, -1),
				AvgDebitAmount:       decimal.New(-631000000000000000, -16),
				AvgCreditAmount:      decimal.Zero,
				TransactionsPerMonth: []int{1, 1, 2, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			data := NewTransactionsData(tc.records)
			assert.Equal(t, tc.expected, data)
		})
	}
}
