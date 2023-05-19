package ally

import (
	"log"
	"testing"
)

func TestGetWatchlists(t *testing.T) {
	wl := GetWatchlists()
	if len(wl.Watchlists) < 1 {
		t.Errorf("GetWatchlists() returned zero watchlists")
	}
}

func TestGetWatchlist(t *testing.T) {
	wl := GetWatchlist("ANY")
	if len(wl.Items) < 1 {
		t.Errorf("GetWatchlist('ANY') returned zero items")
	}
}

func TestCreateDeleteWatchlist(t *testing.T) {
	pre := len(GetWatchlists().Watchlists)

	cl := CreateWatchlist("watchlist-test", []string{"AAPL", "MSFT"})
	log.Print(cl.Watchlists)
	post := len(GetWatchlists().Watchlists)

	if post <= pre {
		t.Errorf("No watchlist created")
	}

	dl := DeleteWatchlist("watchlist-test")

	if dl.Error != "Success" {
		t.Error("Failed to delete watchlist")
	}
}

func TestAppendRemoveWatchlist(t *testing.T) {
	wl := AppendWatchlist("ANY", []string{"HRL"})
	if wl.Error != "Success" {
		t.Error("Failed to append to watchlist")
	}

	dl := RemoveFromWatchlist("ANY", "HRL")
	if dl.Error != "Success" {
		t.Error("Failed to remove from watchlist")
	}

}
