package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/faangbait/go-allyinvest/ally"
	"github.com/faangbait/go-allyinvest/bot"
	"github.com/faangbait/go-allyinvest/webserver"
)

const (
	ClientVersion = "1.0"
)

func main() {
	ally.TrySetEnv()
	t := time.Now()

	bot.ParseFlags()

	if *bot.Viewer { // "View Mode" preempts all other operations
		webserver.Listen()
	} else {
		if ally.APIVersion().Version != ClientVersion {
			fmt.Print(fmt.Errorf("wrong version: %s!=%s", ally.APIVersion().Version, ClientVersion))
		}

		if os.Getenv("ALLY_ACCOUNTS") == "" || len(strings.Split(os.Getenv("ALLY_ACCOUNTS"), ":")) == 0 {
			fmt.Print(fmt.Errorf("no accounts found! set the ALLY_ACCOUNTS environment variable"))
		}

		if *bot.Report { // "Report Mode" preempts all automated trading operations
			genReports()
		} else {
			err := bot.ExecuteTradingStrategy()

			if err != nil {
				fmt.Print(fmt.Errorf("error running application: %v", err.Error()))
			}
			t2 := time.Now()
			fmt.Printf("Ran for %dms", t2.Sub(t).Milliseconds())
		}
	}
}
