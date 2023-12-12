package currency

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestStrategy(t *testing.T) {
	start := USD
	currencies := []Currency{USD, EUR, CHF}
	rates := []ExchangeRate{
		{
			From: USD,
			To:   EUR,
			Rate: 0.9,
		},
		{
			From: EUR,
			To:   CHF,
			Rate: 1,
		},
		{
			From: CHF,
			To:   USD,
			Rate: 1.12,
		},

		{
			From: EUR,
			To:   USD,
			Rate: 1.12,
		},
		{
			From: CHF,
			To:   EUR,
			Rate: 0.99,
		},
		{
			From: USD,
			To:   CHF,
			Rate: 0.89,
		},
	}

	strategy := FindStrategy(start, currencies, rates)
	require.Len(t, strategy, 4)
	require.Equal(t, USD, strategy[0])
	require.Equal(t, EUR, strategy[1])
	require.Equal(t, CHF, strategy[2])
	require.Equal(t, USD, strategy[3])
}

func TestNoStrategy(t *testing.T) {
	start := USD
	currencies := []Currency{USD, EUR, CHF}
	rates := []ExchangeRate{
		{
			From: USD,
			To:   EUR,
			Rate: 0.8,
		},
		{
			From: EUR,
			To:   CHF,
			Rate: 1,
		},
		{
			From: CHF,
			To:   USD,
			Rate: 1.11,
		},

		{
			From: EUR,
			To:   USD,
			Rate: 1.12,
		},
		{
			From: CHF,
			To:   EUR,
			Rate: 0.99,
		},
		{
			From: USD,
			To:   CHF,
			Rate: 0.89,
		},
	}

	strategy := FindStrategy(start, currencies, rates)
	require.Nil(t, strategy)
}
