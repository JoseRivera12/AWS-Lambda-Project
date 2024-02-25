package fileutils

import (
	"bytes"
	"fmt"
	"time"

	"github.com/gocarina/gocsv"
	"github.com/shopspring/decimal"
)

type DateTime struct {
	time.Time
}

func (date *DateTime) MarshalCSV() (string, error) {
	return date.Time.Format("2006-01-02T15:04:05"), nil
}

func (date *DateTime) UnmarshalCSV(csv string) (err error) {
	date.Time, err = time.Parse("2006-01-02T15:04:05", csv)
	return err
}

type Record struct {
	Date        DateTime        `csv:"Date"`
	Currency    string          `csv:"Currency"`
	Description string          `csv:"Description"`
	Transaction decimal.Decimal `csv:"Transaction"`
}

func GetCSVRows(data []byte) ([]*Record, error) {
	var records []*Record
	reader := bytes.NewReader(data)
	if err := gocsv.Unmarshal(reader, &records); err != nil {
		return nil, fmt.Errorf("error in csv file: %w", err)
	}

	return records, nil
}
