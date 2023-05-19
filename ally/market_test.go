package ally

import (
	"testing"
	"time"
)

func TestClock(t *testing.T) {
	est, _ := time.LoadLocation("America/New_York")
	now := time.Now().In(est)
	got := Clock()

	if now.After(TruncateToMarketEnd(now)) {
		if got.Status.Current != "close" && got.Status.Current != "after" {
			t.Errorf("Clock() = %v, want %v", got.Status.Current, "close|after")
		}
	}

	if now.Before(TruncateToMarketStart(now)) {
		if got.Status.Current != "close" && got.Status.Current != "pre" {
			t.Errorf("Clock() = %v, want %v", got.Status.Current, "close|pre")
		}
	}

	if got.Date.Hour() != now.Hour() {
		t.Errorf("Clock() = %v, want %v", got.Date.Hour(), now.Hour())
	}
}

func TestGetQuotes(t *testing.T) {
	quotes := GetQuotes([]string{"AAPL", "MSFT"}, []string{"bid", "ask", "symbol"})

	if len(quotes.Quotes.Quote) < 2 {
		t.Errorf("Didn't get all my quotes")
	}
}

func TestNewsSearch(t *testing.T)    {}
func TestNewsById(t *testing.T)      {}
func TestOptionsSearch(t *testing.T) {}
func TestOptionsStrike(t *testing.T) {}
func TestOptionsExpire(t *testing.T) {}
func TestTimeSales(t *testing.T)     {}
func TestTopLists(t *testing.T)      {}
