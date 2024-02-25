package mailer

import (
	"time"

	"github.com/leekchan/accounting"
	"github.com/shopspring/decimal"
)

func formatMoney(amount decimal.Decimal) string {
	ac := accounting.Accounting{Symbol: "$", Precision: 2}
	return ac.FormatMoney(amount)
}

func monthMapper(month int) string {
	names := []string{"January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"}
	if month >= 0 && month < 12 {
		return names[month]
	}
	return ""
}

func getYear() int {
	return time.Now().Year()
}
