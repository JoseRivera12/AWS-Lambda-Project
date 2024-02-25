package fileutils

import (
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

type testCase struct {
	name        string
	data        []byte
	expected    []*Record
	expectedErr error
}

func TestGetCSVRows(t *testing.T) {
	testCases := []testCase{
		{
			name: "Valid data",
			data: []byte(`Id,Date,Transaction,Currency,Description
				1,2024-01-01T00:00:00,+1234.56,USD,Salary`,
			),
			expected: []*Record{
				{
					Date:        DateTime{Time: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
					Currency:    "USD",
					Description: "Salary",
					Transaction: decimal.NewFromFloat(1234.56),
				},
			},
		},
		{
			name: "Invalid date format",
			data: []byte(`Id,Date,Transaction,Currency,Description
				1,2024-01-01,+1234.56,USD,Salary`,
			),
			expectedErr: &time.ParseError{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			records, err := GetCSVRows(tc.data)
			if tc.expectedErr != nil {
				assert.Error(t, err)
			} else {
				assert.Equal(t, tc.expected[0].Description, records[0].Description)
				assert.Equal(t, tc.expected[0].Transaction, records[0].Transaction)
				assert.Equal(t, tc.expected[0].Currency, records[0].Currency)
				assert.Equal(t, tc.expected[0].Date.Time, records[0].Date.Time)
			}
		})
	}
}
