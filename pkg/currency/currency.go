package currency

import (
	"github.com/amanakin/currency_arbitrage/pkg/graph"
	"math"
)

type Currency string

const (
	USD Currency = "usd"
	CHF Currency = "chf"
	EUR Currency = "eur"
	GBP Currency = "gbp"
)

// ExchangeRate is From/To forex rate
type ExchangeRate struct {
	From Currency
	To   Currency
	Rate float64 // How much To you can get for 1 From
}

func FindStrategy(start Currency, currencies []Currency, rates []ExchangeRate) []Currency {
	g := graph.NewGraph()
	for _, currency := range currencies {
		g.InsertVertex(string(currency))
	}
	for _, rate := range rates {
		g.InsertEdge(string(rate.From), string(rate.To), -math.Log10(rate.Rate))
	}

	cycle := g.FindNegativeCycleFromStart(string(start))
	if len(cycle) == 0 {
		return nil
	}

	var result []Currency
	for _, v := range cycle {
		result = append(result, Currency(v.Name))
	}

	return result
}
