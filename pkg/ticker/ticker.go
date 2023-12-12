package ticker

import (
	"github.com/amanakin/currency_arbitrage/pkg/currency"
	"time"
)

type Ticker interface {
	Configure(time time.Time)
	GetLast(time time.Time) currency.ExchangeRate
	NextDiff() time.Duration
}
