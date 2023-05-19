package ally

import (
	"encoding/xml"
	"fmt"
	"net/url"
	"strings"
	"time"
)

// GET market/clock
func Clock() IClock {
	return get[IClock]("market/clock", url.Values{}, infiniteRL)
}

// GET market/ext/quotes
func GetQuotes(symbols []string, fids []string) IGetQuotes {

	query := url.Values{}
	query.Add("symbols", strings.Join(symbols, ","))
	query.Add("fids", strings.Join(fids, ","))

	return get[IGetQuotes]("market/ext/quotes", query, quotesRL)
}

// GET market/news/search
func NewsSearch(symbol string, max uint8, start time.Time, end time.Time) INewsSearch {
	query := url.Values{}
	query.Add("symbols", symbol)
	query.Add("maxhits", fmt.Sprint(max))
	query.Add("startdate", fmt.Sprint(start))
	query.Add("enddate", fmt.Sprint(end))

	return get[INewsSearch]("market/news/search", query, infiniteRL)
}

// GET market/news/:id
func NewsById(id string) INewsById {
	return get[INewsById](fmt.Sprintf("market/news/%s", id), url.Values{}, infiniteRL)
}

// GET market/options/search
func OptionsSearch(symbol string, querystring string, fids []string) IOptionsSearch {
	query := url.Values{}
	query.Add("symbol", symbol)
	query.Add("query", querystring)
	query.Add("fids", strings.Join(fids, ","))

	return get[IOptionsSearch]("market/options/search", query, quotesRL)
}

// GET market/options/strikes
func OptionsStrike(symbol string) IOptionsStrike {
	query := url.Values{}
	query.Add("symbol", symbol)

	return get[IOptionsStrike]("market/options/strikes", query, quotesRL)
}

// GET market/options/expirations
func OptionsExpire(symbol string) IOptionsExpire {
	query := url.Values{}
	query.Add("symbol", symbol)

	return get[IOptionsExpire]("market/options/expirations", query, quotesRL)
}

// GET market/timesales
func TimeSales(symbol string, interval uint8, rpp uint8, index uint8, start time.Time, end time.Time) ITimeSales {
	query := url.Values{}
	query.Add("symbol", symbol)

	if interval > 0 {
		query.Add("interval", fmt.Sprintf("%dmin", interval))
	} else {
		query.Add("interval", "tick")
		query.Add("rpp", fmt.Sprint(rpp))
		query.Add("index", fmt.Sprint(index))
	}

	query.Add("startdate", start.Format("2006-01-02"))
	query.Add("enddate", end.Format("2006-01-02"))

	marketOpen := TruncateToMarketStart(time.Now())
	if start.After(marketOpen) {
		query.Add("starttime", start.Format("15:04"))
	}

	return get[ITimeSales]("market/timesales", query, quotesRL)

}

// GET market/toplists/:name
func TopLists(listname string, exchange string) ITopLists {
	query := url.Values{}
	query.Add("exchange", exchange)

	return get[ITopLists](fmt.Sprintf("market/toplists/%s", listname), query, authedRL)

}

type clockTime struct {
	time.Time
}

func (t *clockTime) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	const format = "2006-01-02 15:04:05-07:00"
	var v string
	d.DecodeElement(&v, &start)
	parse, err := time.Parse(format, v)
	if err != nil {
		return err
	}
	t.Time = parse
	return nil
}
