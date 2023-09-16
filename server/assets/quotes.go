package assets

import (
	"bytes"
	_ "embed"
	"encoding/csv"
	"fmt"

	"github.com/pkg/errors"
)

//go:embed quotes.csv
var quotesCSV []byte

type Quote struct {
	Author string
	Msg    string
}

func (q Quote) String() string {
	return fmt.Sprintf("%s: %s", q.Author, q.Msg)

}
func GetQuotes() ([]Quote, error) {
	r := csv.NewReader(bytes.NewReader(quotesCSV))
	records, err := r.ReadAll()
	if err != nil {
		return nil, errors.New("could not read quotes CSV")
	}

	var quotes []Quote
	for _, record := range records[1:] { // Skip header
		quotes = append(quotes, Quote{
			Author: record[0],
			Msg:    record[1],
		})
	}
	return quotes, nil
}
