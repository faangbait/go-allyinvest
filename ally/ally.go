package ally

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/meirf/gopart"
)

const (
	MINIMUM_STOCK_PRICE = 2.0
)

func RequestHeaders() map[string][]string {
	headers := map[string][]string{}

	trade_password := os.Getenv("ALLY_TRADE_PASS")

	if trade_password != "" {
		headers["TKI_TRADEPASS"] = []string{trade_password}
	}

	return headers
}

func get[O any](uri string, query url.Values, limit rateLimitType) O {
	s := NewClient(RateLimitMap[limit])
	var iface O

	req, err := s.Do(&http.Request{
		Method: http.MethodGet,
		URL: &url.URL{
			Scheme:   "https",
			Host:     "devapi.invest.ally.com",
			Path:     fmt.Sprintf("/v1/%s%s", uri, ".xml"),
			RawQuery: query.Encode(),
		},
	})

	if err != nil {
		log.Panicf("HTTP Error: %s", err.Error())
	}

	defer req.Body.Close()
	body, err := io.ReadAll(req.Body)

	if err != nil {
		log.Panicf("HTTP Read Error: %s", err.Error())
	}

	err = xml.Unmarshal(body, &iface)

	if err != nil {
		log.Panicf("Unmarshal Error: %s", err.Error())
	}

	return iface
}

func post[O any](uri string, postData []byte, headers map[string][]string, query url.Values, limit rateLimitType) (O, []byte, error) {
	s := NewClient(RateLimitMap[limit])
	var iface O

	bodyReader := io.NopCloser(bytes.NewReader(postData))

	req, err := s.Do(&http.Request{
		Method: http.MethodPost,
		URL: &url.URL{
			Scheme:   "https",
			Host:     "devapi.invest.ally.com",
			Path:     fmt.Sprintf("/v1/%s%s", uri, ".xml"),
			RawQuery: query.Encode(),
		},
		Body:   bodyReader,
		Header: headers,
	})

	if err != nil {
		log.Panicf("HTTP Error: %s", err.Error())
	}

	defer req.Body.Close()
	body, err := io.ReadAll(req.Body)

	if err != nil {
		log.Panicf("HTTP Read Error: %s", err.Error())
	}

	err = xml.Unmarshal(body, &iface)

	return iface, body, err
}

func postform[O any](uri string, data url.Values, limit rateLimitType) O {
	s := NewClient(RateLimitMap[limit])
	var iface O

	bodyReader := io.NopCloser(strings.NewReader(data.Encode()))

	r, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("https://devapi.invest.ally.com/v1/%s%s", uri, ".xml"), bodyReader)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req, err := s.Do(r)

	if err != nil {
		log.Panicf("HTTP Error: %s", err.Error())
	}

	defer req.Body.Close()
	body, err := io.ReadAll(req.Body)

	if err != nil {
		log.Panicf("HTTP Read Error: %s", err.Error())
	}

	err = xml.Unmarshal(body, &iface)

	if err != nil {
		log.Panicf("Unmarshal Error: %s", err.Error())
	}

	return iface
}

/*
	func postjson[O any](uri string, data url.Values, limit rateLimitType) O {
		s := setup()
		var iface O

		bodyReader := io.NopCloser(bytes.NewReader([]byte(data.Encode())))

		req, err := s.Do(&http.Request{
			Method: http.MethodPost,
			URL: &url.URL{
				Scheme: "https",
				Host:   "devapi.invest.ally.com",
				Path:   fmt.Sprintf("/v1/%s%s", uri, ".xml"),
			},
			Body: bodyReader,
		})

		if err != nil {
			log.Panicf("HTTP Error: %s", err.Error())
		}

		defer req.Body.Close()
		body, err := io.ReadAll(req.Body)

		if err != nil {
			log.Panicf("HTTP Read Error: %s", err.Error())
		}

		log.Print(string(body))

		err = json.Unmarshal(body, &iface)

		if err != nil {
			log.Panicf("Unmarshal Error: %s", err.Error())
		}

		return iface
	}
*/

func delete[O any](uri string, limit rateLimitType) O {
	s := NewClient(RateLimitMap[limit])
	var iface O

	req, err := s.Do(&http.Request{
		Method: http.MethodDelete,
		URL: &url.URL{
			Scheme: "https",
			Host:   "devapi.invest.ally.com",
			Path:   fmt.Sprintf("/v1/%s%s", uri, ".xml"),
		},
	})

	if err != nil {
		log.Panicf("HTTP Error: %s", err.Error())
	}

	defer req.Body.Close()
	body, err := io.ReadAll(req.Body)

	if err != nil {
		log.Panicf("HTTP Read Error: %s", err.Error())
	}

	err = xml.Unmarshal(body, &iface)

	if err != nil {
		log.Panicf("Unmarshal Error: %s", err.Error())
	}

	return iface
}

func RefreshPrices(securities map[string]Security) {
	keys := make([]string, 0, len(securities))
	for k := range securities {
		keys = append(keys, k)
	}

	for idxRange := range gopart.Partition(len(keys), 85) {
		symbols := keys[idxRange.Low:idxRange.High]
		quotes := GetQuotes(symbols, []string{"bid", "ask", "symbol", "last"})
		for _, quote := range quotes.Quotes.Quote {
			s := securities[quote.Symbol]
			s.Quote = quote
			if s.Quote.Bid == 0 {
				s.Quote.Bid = s.Quote.Last
			}
			if s.Quote.Ask == 0 {
				s.Quote.Ask = s.Quote.Last
			}

			securities[quote.Symbol] = s
		}
	}
}

func ClearHoldings(securities map[string]Security) {
	for k, s := range securities {
		s.Holding = IHolding{}
		s.AccountID = ""
		securities[k] = s
	}
}

func RefreshAccounts(securities map[string]Security) {
	accounts := GetAccounts()

	ClearHoldings(securities)

	for _, acct := range accounts.AccountSummaries {
		for _, holding := range acct.Holdings {
			var s Security

			sym := strings.ToUpper(holding.Instrument.FIXMLSymbol)

			s, ok := securities[sym]

			if !ok {
				s = Security{}
			}

			s.Holding = holding
			s.AccountID = acct.Account
			s.Symbol = sym
			securities[sym] = s
		}
	}
}
