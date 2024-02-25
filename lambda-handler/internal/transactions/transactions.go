package transactions

import (
	"time"
	"transaction-mailer/internal/fileutils"
	"transaction-mailer/internal/utils"

	"github.com/shopspring/decimal"
)

type TransactionsData struct {
	TotalBalance         decimal.Decimal
	AvgDebitAmount       decimal.Decimal
	AvgCreditAmount      decimal.Decimal
	TransactionsPerMonth []int
}

func NewTransactionsData(records []*fileutils.Record) *TransactionsData {
	transactionsPerMonth := make([]int, 12)
	data := &TransactionsData{TransactionsPerMonth: transactionsPerMonth}
	data.processTransactions(records)
	return data
}

func (td *TransactionsData) processTransactions(transactions []*fileutils.Record) {
	currentYear := time.Now().Year()
	totalBalance := decimal.Zero
	debitBalance, creditBalance := decimal.Zero, decimal.Zero
	var debitTransactions, creditTransactions int64

	for _, record := range transactions {
		amount := record.Transaction
		if record.Date.Year() != currentYear {
			continue // Skip transactions from different years
		}

		totalBalance = totalBalance.Add(amount)
		month := int(record.Date.Month())
		td.TransactionsPerMonth[month-1]++

		if amount.LessThan(decimal.Zero) {
			debitBalance = debitBalance.Add(amount)
			debitTransactions++
		} else {
			creditBalance = creditBalance.Add(amount)
			creditTransactions++
		}
	}

	td.TotalBalance = totalBalance
	td.AvgDebitAmount = utils.SafeDiv(debitBalance, decimal.NewFromInt(debitTransactions))
	td.AvgCreditAmount = utils.SafeDiv(creditBalance, decimal.NewFromInt(creditTransactions))
}
