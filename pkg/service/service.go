package service

import (
	"context"
	"fmt"
	"github.com/amanakin/currency_arbitrage/pkg/currency"
	"github.com/amanakin/currency_arbitrage/pkg/ticker"
	"strings"
	"time"
)

type Service struct {
	Money   float64
	Tickers []ticker.Ticker

	StartTime time.Time
	Step      time.Duration

	Currencies    []currency.Currency
	StartCurrency currency.Currency

	PlotStep time.Duration
}

func NewService(
	money float64,
	tickers []ticker.Ticker,
	start time.Time,
	step time.Duration,
	currencies []currency.Currency,
	startCurrency currency.Currency,
	plotStep time.Duration) *Service {
	return &Service{
		Money:   money,
		Tickers: tickers,

		StartTime: start,
		Step:      step,

		Currencies:    currencies,
		StartCurrency: startCurrency,

		PlotStep: plotStep,
	}
}

func (s *Service) Simulate(ctx context.Context, last time.Time) []float64 {
	var levels []float64

	t := s.StartTime
	for i := 0; true; i++ {
		res := func() bool {
			defer func() {
				// Add metrics
				batchSize := int(s.PlotStep / s.Step)
				if i%batchSize == 0 {
					for j := range s.Tickers {
						s.Tickers[j].Configure(t)
					}
					levels = append(levels, s.Money)
				}
			}()

			// Setup rates
			var rates []currency.ExchangeRate
			for j := range s.Tickers {
				rates = append(rates, s.Tickers[j].GetLast(t))
			}

			// Find strategy
			strategy := currency.FindStrategy(s.StartCurrency, s.Currencies, rates)

			// Add ts
			t = t.Add(s.Step)
			if t.After(last) {
				return true
			}

			// Try restart if we have big gap
			var restart bool
			for j := range s.Tickers {
				if s.Tickers[j].NextDiff() > 1*time.Minute {
					restart = true
				}
			}
			if restart {
				return false
			}

			// Update rates with new time
			for j := range s.Tickers {
				rates[j] = s.Tickers[j].GetLast(t)
			}

			// Follow strategy
			if len(strategy) > 0 {
				s.FollowStrategy(strategy, rates, t)
			}

			if ctx.Err() != nil {
				return true
			}
			return false
		}()

		if res {
			break
		}
	}

	return levels
}

type key struct {
	from currency.Currency
	to   currency.Currency
}

func (s *Service) FollowStrategy(
	strategy []currency.Currency,
	rates []currency.ExchangeRate,
	t time.Time) {
	m := make(map[key]float64)
	for _, rate := range rates {
		m[key{
			from: rate.From,
			to:   rate.To,
		}] = rate.Rate
	}

	// TODO do something smarted
	spend := float64(int(s.Money / 10))
	s.Money -= spend

	was := spend
	for i := 1; i < len(strategy); i++ {
		from := strategy[i-1]
		to := strategy[i]
		k := key{
			from: from,
			to:   to,
		}

		rate := m[k]
		if rate == 0 {
			panic("Loose all money")
		}

		spend = spend * rate
	}

	var arr []string
	for i := 0; i < len(strategy); i++ {
		arr = append(arr, string(strategy[i]))
	}
	fmt.Printf("%f->%f: %s(%s)\n", was, spend, t.String(), strings.Join(arr, "->"))

	s.Money += spend
	return
}
