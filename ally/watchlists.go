package ally

import (
	"fmt"
	"net/url"
	"strings"
)

// GET watchlists
func GetWatchlists() IGetWatchlists {
	return get[IGetWatchlists]("watchlists", url.Values{}, authedRL)
}

// GET watchlists/:id
func GetWatchlist(id string) IGetWatchlist {
	return get[IGetWatchlist](fmt.Sprintf("watchlists/%s", id), url.Values{}, authedRL)
}

// POST watchlists
func CreateWatchlist(id string, symbols []string) IGetWatchlists {
	data := url.Values{}
	data.Add("id", id)
	data.Add("symbols", strings.Join(symbols, ","))
	return postform[IGetWatchlists]("watchlists", data, authedRL)
}

// DELETE watchlists/:id
func DeleteWatchlist(id string) IGetSingleWatchlist {
	return delete[IGetSingleWatchlist](fmt.Sprintf("watchlists/%s", id), authedRL)
}

// POST watchlists/:id/symbols
func AppendWatchlist(id string, symbols []string) IGetWatchlists {
	data := url.Values{}
	data.Add("id", id)
	data.Add("symbols", strings.Join(symbols, ","))
	return postform[IGetWatchlists](fmt.Sprintf("watchlists/%s/symbols", id), data, authedRL)
}

// DELETE watchlists/:id/symbols/:sym
func RemoveFromWatchlist(id string, symbol string) IGetWatchlists {
	return delete[IGetWatchlists](fmt.Sprintf("watchlists/%s/symbols/%s", id, symbol), authedRL)
}
