package main

import (
	"context"
	"fmt"
	"github.com/amanakin/currency_arbitrage/pkg/currency"
	"github.com/amanakin/currency_arbitrage/pkg/history"
	"github.com/amanakin/currency_arbitrage/pkg/service"
	"github.com/amanakin/currency_arbitrage/pkg/ticker"
	"github.com/spf13/cobra"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
	"path/filepath"
	"time"
)

var flags struct {
	dataDir string
	year    int
	month   int
}

var simulateCmd = cobra.Command{
	Use:   "simulate",
	Short: "start simulation",
	Args:  cobra.NoArgs,
	Run:   simulate,
}

func init() {
	f := simulateCmd.Flags()
	f.StringVar(&flags.dataDir, "data", "", "directory for data")
	f.IntVar(&flags.year, "year", 0, "year of data")
	f.IntVar(&flags.month, "month", 0, "month of data")
	/*f.StringVar(&f.listenAddr, "addr", "127.0.0.1:0", "listen addr")
	f.StringArrayVar(&serveFlags.backends, "backends", []string{}, "addresses of backends")
	f.IntVar(&serveFlags.concurrency, "concurrency", 4, "number of concurrent requests to backends")*/

	rootCmd.AddCommand(&simulateCmd)
}

func simulate(cmd *cobra.Command, args []string) {
	dir := flags.dataDir
	files := []struct {
		filename string
		from     currency.Currency
		to       currency.Currency
	}{
		{
			filename: "EURCHF.csv",
			from:     currency.EUR,
			to:       currency.CHF,
		},
		{
			filename: "EURGBP.csv",
			from:     currency.EUR,
			to:       currency.GBP,
		},
		{
			filename: "EURUSD.csv",
			from:     currency.EUR,
			to:       currency.USD,
		},
		{
			filename: "GBPCHF.csv",
			from:     currency.GBP,
			to:       currency.CHF,
		},
		{
			filename: "GBPUSD.csv",
			from:     currency.GBP,
			to:       currency.USD,
		},
		{
			filename: "USDCHF.csv",
			from:     currency.USD,
			to:       currency.CHF,
		},
	}

	var tickers []ticker.Ticker
	for _, file := range files {
		currFile := filepath.Join(dir, file.filename)
		fmt.Printf("read %s\n", currFile)
		ticker1, ticker2, err := history.ReadHistoryTickers(currFile, file.from, file.to)
		if err != nil {
			panic(err)
		}
		fmt.Printf("read successfully %s\n", currFile)

		tickers = append(tickers, ticker1, ticker2)
	}

	loc, err := time.LoadLocation("EST")
	if err != nil {
		panic(err)
	}

	//timeStart := time.Date(2010, 8, 1, 19, 55, 0, 0, loc)
	//timeEnd := time.Date(2010, 8, 1, 19, 55, 1, 0, loc)

	timeStart := time.Date(flags.year, time.Month(flags.month), 1, 8, 00, 0, 0, loc)
	timeEnd := time.Date(flags.year, time.Month(flags.month), 29, 16, 00, 0, 0, loc)
	for _, t := range tickers {
		t.Configure(timeStart)
	}
	currencies := []currency.Currency{currency.USD, currency.EUR, currency.CHF, currency.GBP}
	s := service.NewService(
		10000,
		tickers,
		timeStart,
		500*time.Millisecond,
		currencies,
		currency.USD,
		10*time.Second)

	levels := s.Simulate(context.Background(), timeEnd)

	xticks := plot.TimeTicks{Format: "2006-01-02\n15:04"}
	var line plotter.XYs
	for i, level := range levels {
		line = append(line, plotter.XY{
			X: float64(timeStart.Add(time.Duration(i) * 10 * time.Second).Unix()),
			Y: level,
		})
	}

	plt := plot.New()
	plt.X.Tick.Marker = xticks
	plt.Title.Text = "Budget with time"
	plt.X.Label.Text = "time"
	plt.Y.Label.Text = "money"

	err = plotutil.AddLinePoints(plt, "budget", line)
	if err != nil {
		panic(err)
	}
	if err = plt.Save(4*vg.Inch, 4*vg.Inch, "budget.png"); err != nil {
		panic(err)
	}
}
