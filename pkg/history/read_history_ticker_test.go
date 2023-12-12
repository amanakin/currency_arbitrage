package history

import (
	"fmt"
	"github.com/amanakin/currency_arbitrage/pkg/graph"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestParseDateTime(t *testing.T) {
	loc, err := time.LoadLocation("EST")
	require.NoError(t, err)

	cases := []struct {
		str      string
		expected time.Time
	}{
		{
			str:      "20100401 000003000",
			expected: time.Date(2010, 4, 1, 0, 0, 3, 0, loc),
		},
		{
			str:      "20100430 165104050",
			expected: time.Date(2010, 4, 30, 16, 51, 4, 50_000_000, loc),
		},
	}

	for i := range cases {
		t.Run(fmt.Sprintf("%d/%d", i+1, len(cases)), func(t *testing.T) {
			require.Equal(t, cases[i].expected, parseDateTime(cases[i].str))
		})
	}
}

const fileexample = `20100401 000003000,1.423300,1.423600,0
20100401 000014000,1.423400,1.423700,0
20100401 000059000,1.423500,1.423800,0
20100401 000103000,1.423600,1.423900,0
20100401 000103000,1.423500,1.423800,0
20100401 000104000,1.423400,1.423800,0
20100401 000104000,1.423500,1.423800,0
20100401 000104000,1.423400,1.423800,0
20100401 000104000,1.423400,1.423700,0
20100401 000105000,1.423400,1.423800,0
20100401 000113000,1.423400,1.423700,0
20100401 000120000,1.423400,1.423800,0`

func TestReadHistoryTickers(t *testing.T) {
	dir := t.TempDir()
	filename := filepath.Join(dir, "eur-chf.csv")

	err := os.WriteFile(filename, []byte(fileexample), 0666)
	require.NoError(t, err)

	EURToCHFTicker, CHFToEURTicker, err := ReadHistoryTickers(filename, "eur", "chf")
	require.NoError(t, err)

	require.Len(t, EURToCHFTicker.rows, 12)
	require.Len(t, CHFToEURTicker.rows, 12)

	for i := 0; i < 12; i++ {
		require.Less(t, EURToCHFTicker.rows[i].rate, 1/CHFToEURTicker.rows[i].rate)
	}

	for i := 1; i < 12; i++ {
		require.False(t, EURToCHFTicker.rows[i-1].ts.After(EURToCHFTicker.rows[i].ts))
	}
}

func TestTickersWithGraph(t *testing.T) {
	dir := t.TempDir()
	filename := filepath.Join(dir, "eur-chf.csv")

	err := os.WriteFile(filename, []byte(fileexample), 0666)
	require.NoError(t, err)

	EURToCHFTicker, CHFToEURTicker, err := ReadHistoryTickers(filename, "eur", "chf")
	require.NoError(t, err)

	g := graph.NewGraph()
	g.InsertVertex("eur")
	g.InsertVertex("chf")

	g.InsertEdge("eur", "chf", EURToCHFTicker.rows[0].rate)
	g.InsertEdge("chf", "eur", CHFToEURTicker.rows[0].rate)

	cycles := g.FindNegativeCycleFromStart("eur")
	require.Empty(t, cycles)
	cycles = g.FindNegativeCycleFromStart("chf")
	require.Empty(t, cycles)
}
