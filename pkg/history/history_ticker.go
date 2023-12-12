package history

import (
	"github.com/amanakin/currency_arbitrage/pkg/currency"
	"sort"
	"time"
)

type historyTickerRow struct {
	ts   time.Time
	rate float64
}

type Ticker struct {
	from currency.Currency
	to   currency.Currency

	rows []historyTickerRow
	curr int
}

func NewTicker(from currency.Currency, to currency.Currency) *Ticker {
	return &Ticker{
		from: from,
		to:   to,
	}
}

func (t *Ticker) AddRow(ts time.Time, rate float64) {
	t.rows = append(t.rows, historyTickerRow{
		ts:   ts,
		rate: rate,
	})
}

func (t *Ticker) PrepareRows() {
	sort.Slice(t.rows, func(i, j int) bool {
		return t.rows[i].ts.Before(t.rows[j].ts)
	})
	t.curr = 0
}

func (t *Ticker) GetLast(ts time.Time) currency.ExchangeRate {
	next := t.curr + 1
	if next >= len(t.rows) {
		return currency.ExchangeRate{
			From: t.from,
			To:   t.to,
			Rate: t.rows[len(t.rows)-1].rate,
		}
	}

	if !t.rows[next].ts.After(ts) {
		t.curr++
	}

	return currency.ExchangeRate{
		From: t.from,
		To:   t.to,
		Rate: t.rows[t.curr].rate,
	}
}

func (t *Ticker) Configure(time time.Time) {
	i := sort.Search(len(t.rows), func(i int) bool { return time.Before(t.rows[i].ts) })
	if i == len(t.rows) {
		i = len(t.rows) - 1
	}
	t.curr = i
}

func (t *Ticker) NextDiff() time.Duration {
	next := t.curr + 1
	if next >= len(t.rows) {
		return time.Hour
	}

	return t.rows[next].ts.Sub(t.rows[t.curr].ts)
}