package utils

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

type testCase struct {
	name     string
	dividend decimal.Decimal
	divisor  decimal.Decimal
	expected decimal.Decimal
}

func TestSafeDiv(t *testing.T) {
	tests := []testCase{
		{name: "Non-zero divisor", dividend: decimal.NewFromFloat(10.0), divisor: decimal.NewFromFloat(2.0), expected: decimal.New(50000000000000000, -16)},
		{name: "Zero divisor", dividend: decimal.NewFromFloat(10.0), divisor: decimal.Zero, expected: decimal.Zero},
		{name: "Negative dividend, positive divisor", dividend: decimal.NewFromFloat(-122.3232), divisor: decimal.NewFromFloat(3.0), expected: decimal.New(-407744000000000000, -16)},
		{name: "Positive dividend, negative divisor", dividend: decimal.NewFromFloat(315.765), divisor: decimal.NewFromFloat(-2.0), expected: decimal.New(-1578825000000000000, -16)},
		{name: "Negative dividend, negative divisor", dividend: decimal.NewFromFloat(-542.432), divisor: decimal.NewFromFloat(-2.0), expected: decimal.New(2712160000000000000, -16)},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := SafeDiv(tc.dividend, tc.divisor)
			assert.Equal(t, tc.expected, result)
		})
	}
}
