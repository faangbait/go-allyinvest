package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/faangbait/go-allyinvest/ally"
	"github.com/faangbait/go-allyinvest/webserver"
)

const (
	ClientVersion = "1.0"
)

func main() {
	ally.TrySetEnv()
	t := time.Now()

	ParseFlags()

	if *Viewer { // "View Mode" preempts all other operations
		webserver.Listen()
	} else {
		if ally.APIVersion().Version != ClientVersion {
			fmt.Print(fmt.Errorf("wrong version: %s!=%s", ally.APIVersion().Version, ClientVersion))
		}

		if os.Getenv("ALLY_ACCOUNTS") == "" || len(strings.Split(os.Getenv("ALLY_ACCOUNTS"), ":")) == 0 {
			fmt.Print(fmt.Errorf("no accounts found! set the ALLY_ACCOUNTS environment variable"))
		}

		if *Report { // "Report Mode" preempts all automated trading operations
			genReports()
		} else {
			f, _ := os.OpenFile("webserver/www/trades.log", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
			log.SetOutput(f)
			log.Printf("Launched with flags: Auto[%v] Live[%v] Sell[%v] Buy[%v] Report[%v] Viewer[%v] ", *Auto, *Live, *Sell, *Buy, *Report, *Viewer)
			defer f.Close()

			err := ExecuteTradingStrategy()

			if err != nil {
				fmt.Print(fmt.Errorf("error running application: %v", err.Error()))
			}
			t2 := time.Now()
			fmt.Printf("Ran for %dms", t2.Sub(t).Milliseconds())
		}
	}
}
