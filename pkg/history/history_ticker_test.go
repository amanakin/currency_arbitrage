package history

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestHistoryTicker(t *testing.T) {
	ticker := NewTicker("a", "b")
	ticker.AddRow(time.Date(2020, 1, 1, 0, 0, 30, 0, time.UTC), 1)
	ticker.AddRow(time.Date(2020, 1, 1, 0, 0, 31, 0, time.UTC), 2)
	ticker.AddRow(time.Date(2020, 1, 1, 0, 0, 31, 500_000_000, time.UTC), 3)
	ticker.PrepareRows()

	require.Equal(t, 1., ticker.GetLast(
		time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)).Rate)
	require.Equal(t, 1., ticker.GetLast(
		time.Date(2020, 1, 1, 0, 0, 30, 0, time.UTC)).Rate)
	require.Equal(t, 1., ticker.GetLast(
		time.Date(2020, 1, 1, 0, 0, 30, 50_000_000, time.UTC)).Rate)
	require.Equal(t, 2., ticker.GetLast(
		time.Date(2020, 1, 1, 0, 0, 31, 0, time.UTC)).Rate)
	require.Equal(t, 2., ticker.GetLast(
		time.Date(2020, 1, 1, 0, 0, 31, 50_000_000, time.UTC)).Rate)
	require.Equal(t, 3., ticker.GetLast(
		time.Date(2020, 1, 1, 0, 0, 31, 500_000_000, time.UTC)).Rate)
	require.Equal(t, 3., ticker.GetLast(
		time.Date(2020, 1, 1, 0, 0, 32, 0, time.UTC)).Rate)
	require.Equal(t, 3., ticker.GetLast(
		time.Date(2021, 1, 1, 0, 0, 32, 0, time.UTC)).Rate)
}
