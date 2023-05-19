package webserver

import (
	"errors"
	"fmt"
	"net/http"
	"os"
)

const (
	logTrades   = "webserver/www/trades.log"
	repAcct     = "webserver/www/report-account.html"
	repAlpha    = "webserver/www/report-symbol.html"
	repDay      = "webserver/www/report-dailygainandloss.html"
	repGL       = "webserver/www/report-totalgainandloss.html"
	repLastPx   = "webserver/www/report-lastprice.html"
	repPosition = "webserver/www/report-positionsize.html"
)

func Listen() {

	http.HandleFunc("/", func(resp http.ResponseWriter, req *http.Request) {
		http.ServeFile(resp, req, repDay)
	})
	http.HandleFunc("/trades", func(resp http.ResponseWriter, req *http.Request) {
		http.ServeFile(resp, req, logTrades)
	})
	http.HandleFunc("/reports/account", func(resp http.ResponseWriter, req *http.Request) {
		http.ServeFile(resp, req, repAcct)
	})
	http.HandleFunc("/reports/alphabetical", func(resp http.ResponseWriter, req *http.Request) {
		http.ServeFile(resp, req, repAlpha)
	})
	http.HandleFunc("/reports/daily", func(resp http.ResponseWriter, req *http.Request) {
		http.ServeFile(resp, req, repDay)
	})
	http.HandleFunc("/reports/gainloss", func(resp http.ResponseWriter, req *http.Request) {
		http.ServeFile(resp, req, repGL)
	})
	http.HandleFunc("/reports/lastprice", func(resp http.ResponseWriter, req *http.Request) {
		http.ServeFile(resp, req, repLastPx)
	})
	http.HandleFunc("/reports/position", func(resp http.ResponseWriter, req *http.Request) {
		http.ServeFile(resp, req, repPosition)
	})

	err := http.ListenAndServe(":3333", nil)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
