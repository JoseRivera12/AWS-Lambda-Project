package utils

import (
	"github.com/shopspring/decimal"
)

func SafeDiv(dividend, divisor decimal.Decimal) decimal.Decimal {
	if divisor.IsZero() {
		return decimal.Zero
	}
	return dividend.Div(divisor)
}
