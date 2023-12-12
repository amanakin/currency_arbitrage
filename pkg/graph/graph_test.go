package graph

import (
	"github.com/stretchr/testify/require"
	"testing"
)

const (
	usd = "USD"
	eur = "EUR"
	chf = "CHF"
)

func TestFindNegativeCycleFromStart(t *testing.T) {
	g := NewGraph()
	g.InsertVertex(usd)
	g.InsertVertex(chf)
	g.InsertVertex(eur)

	g.InsertEdge(usd, eur, 1)
	g.InsertEdge(eur, usd, 1)

	g.InsertEdge(eur, chf, -1)
	g.InsertEdge(chf, eur, 1)

	g.InsertEdge(chf, usd, -0.01)
	g.InsertEdge(usd, chf, 1)

	cycle := g.FindNegativeCycleFromStart(usd)
	require.Len(t, cycle, 4)
	require.Equal(t, usd, cycle[0].Name)
	require.Equal(t, eur, cycle[1].Name)
	require.Equal(t, chf, cycle[2].Name)
	require.Equal(t, usd, cycle[3].Name)
}

func TestFindNegativeCycleFromStartNoCycle(t *testing.T) {
	g := NewGraph()
	g.InsertVertex(usd)
	g.InsertVertex(chf)
	g.InsertVertex(eur)

	g.InsertEdge(usd, eur, 1)
	g.InsertEdge(eur, usd, 1)

	g.InsertEdge(eur, chf, -1)
	g.InsertEdge(chf, eur, 1)

	g.InsertEdge(chf, usd, 0)
	g.InsertEdge(usd, chf, 1)

	cycle := g.FindNegativeCycleFromStart(usd)
	require.Len(t, cycle, 0)
}
