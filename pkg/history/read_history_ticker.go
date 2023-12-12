package history

import (
	"encoding/csv"
	"fmt"
	"github.com/amanakin/currency_arbitrage/pkg/currency"
	"io"
	"os"
	"strconv"
	"time"
)

type RawHistoryDataRow struct {
	DateTime string
	High     float64 // Buy for High
	Low      float64 // Sell for Low
}

func ReadHistoryTickers(filename string, from currency.Currency, to currency.Currency) (*Ticker, *Ticker, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, nil, fmt.Errorf("read (%s): %w", filename, err)
	}

	reader := csv.NewReader(file)

	forwardTicker := NewTicker(from, to)
	backwardTicker := NewTicker(to, from)
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}

		t := parseDateTime(row[0])
		bid, err := strconv.ParseFloat(row[1], 64)
		if err != nil {
			panic(err)
		}

		ask, err := strconv.ParseFloat(row[2], 64)
		if err != nil {
			panic(err)
		}

		forwardTicker.AddRow(t, bid)
		backwardTicker.AddRow(t, 1/ask)
	}

	forwardTicker.PrepareRows()
	backwardTicker.PrepareRows()

	return forwardTicker, backwardTicker, nil
}

// parseDateTime parses datetime from format 20021001 000054000
func parseDateTime(str string) time.Time {
	dateTime := str[:len("20021001 000054")]
	ms := str[len("20021001 000054"):len("20021001 000054000")]

	loc, err := time.LoadLocation("EST")
	if err != nil {
		panic(err)
	}
	t, err := time.ParseInLocation("20060102 150405", dateTime, loc)
	if err != nil {
		panic(err)
	}
	msParsed, err := strconv.Atoi(ms)
	if err != nil {
		panic(err)
	}

	t = t.Add(time.Millisecond * time.Duration(msParsed))
	return t
}
