package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/apsdehal/go-logger"
	"github.com/faangbait/go-allyinvest/ally"
)

type ConfFlag *bool

var (
	Settings          = [6]bool{}
	Auto     ConfFlag = &Settings[0]
	Live     ConfFlag = &Settings[1]
	Sell     ConfFlag = &Settings[2]
	Buy      ConfFlag = &Settings[3]
	Report   ConfFlag = &Settings[4]
	Viewer   ConfFlag = &Settings[5]

	ActiveAccounts = strings.Split(os.Getenv("ALLY_ACCOUNTS"), ":")
	IgnoreList     = []string{"YY", "SUMO"}
)

const (
	MinCapitalization = 100
	MaxCapitalization = MinCapitalization * 5
)

// "report only" mode -- ./ally --report
// "view reports" mode -- ./ally --view
// "full auto" mode -- ./ally --auto --live --sell --buy
func ParseFlags() []bool {

	flag.BoolVar(&Settings[0], "auto", false, "make trades automatically")
	flag.BoolVar(&Settings[1], "live", false, "transmit real orders")
	flag.BoolVar(&Settings[2], "sell", false, "selling enabled")
	flag.BoolVar(&Settings[3], "buy", false, "buying enabled")
	flag.BoolVar(&Settings[4], "report", false, "generate reports, then exit")
	flag.BoolVar(&Settings[5], "view", false, "launch in interactive mode")
	flag.Parse()

	return Settings[:]
}

func ExecuteTradingStrategy() error {
	return fmt.Errorf("not implemented")
}

func ManageWarnings() {
	// this map determines whether an order warning should be "overridden" and resubmitted
	// if continuing requires changing order type (MKT -> LMT), this is automatic if "true"

	ally.OverrideMap[ally.DuplicateOrder] = false
	ally.OverrideMap[ally.UnsettledFunds] = true
	ally.OverrideMap[ally.HigherMarginReq] = true
	ally.OverrideMap[ally.NotMarginable] = true
	ally.OverrideMap[ally.MktOrderWhileClosed] = false
	ally.OverrideMap[ally.ExchangeClosed] = false
	ally.OverrideMap[ally.ForeignSettlementFee] = false
	ally.OverrideMap[ally.NoMarketOrder1] = true
	ally.OverrideMap[ally.NoMarketOrder2] = true
}

type StockType int

const (
	Regular StockType = iota
	Unspecified
	ClassA
	ClassB
	ETFM
	NewIssue
	Delinquent
	ForeignIssue
	Convertible1
	Convertible2
	Convertible3
	Voting
	NonVoting
	Misc
	Preferred1
	Preferred2
	Preferred3
	Preferred4
	Bankrupt
	RightsIssue
	BeneficialInterest
	WarrantsRights
	Units
	CorporateAction
	Warrants
	MFQS
	ADR
	OptionContract
)

var StockTypeMap = map[string]StockType{
	"A": ClassA,
	"B": ClassB,
	"C": ETFM,
	"D": NewIssue,
	"E": Delinquent,
	"F": ForeignIssue,
	"G": Convertible1,
	"H": Convertible2,
	"I": Convertible3,
	"J": Voting,
	"K": NonVoting,
	"L": Misc,
	"M": Preferred4,
	"N": Preferred3,
	"O": Preferred2,
	"P": Preferred1,
	"Q": Bankrupt,
	"R": RightsIssue,
	"S": BeneficialInterest,
	"T": WarrantsRights,
	"U": Units,
	"V": CorporateAction,
	"W": Warrants,
	"X": MFQS,
	"Y": ADR,
	"Z": Misc,
}

// Looks up the classification of a ticker based on NASDAQ symbol
func GetStockType(ticker string) StockType {

	switch len(ticker) {
	case 1, 2, 3, 4:
		return Regular
	case 5:
		last_letter := strings.ToUpper(ticker[len(ticker)-1:])
		res, ok := StockTypeMap[last_letter]
		if ok {
			return res
		} else {
			return Unspecified
		}

	default:
		return OptionContract
	}
}

// Helper function for logging order errors
func logOrderError(o *ally.Order, err error) {
	log.Printf("got error posting %s order: [%s] %s %0.2f %s @%.2f %s > %s\n", o.OrderType, o.Account, o.Action, o.TotalQuantity, o.Symbol, o.LmtPrice, o.Tif, err.Error())
}

// Helper function for logging order successes
func logOrderSuccess(o *ally.Order, res ally.IPostOrder) string {
	return fmt.Sprintf("%s [%s] %s %0.2f %s @%0.2f %s > %s\n", o.OrderType, o.Account, o.Action, o.TotalQuantity, o.Symbol, o.LmtPrice, o.Tif, res.SVIOrderID)
}

func maskIfTrue(val byte, mask byte, predicate bool) byte {
	if predicate {
		return val
	}
	return mask
}

func ConfigureLogger() (*logger.Logger, *os.File) {
	t := time.Now()
	flags := [4]byte{
		maskIfTrue('A', 'X', *Auto),
		maskIfTrue('L', 'X', *Live),
		maskIfTrue('S', 'X', *Sell),
		maskIfTrue('B', 'X', *Buy),
	}

	logger.SetDefaultFormat("%{time:15:14}: %{message}")

	LOG_FILE := "webserver/www/noop.log"

	if *Auto {
		LOG_FILE = fmt.Sprintf("webserver/www/dev.%s.%s.log", flags, t.Format("2006-01-02-Mon"))

		if *Live && (*Sell || *Buy) {
			LOG_FILE = fmt.Sprintf("webserver/www/live.%s.%s.log", flags, t.Format("2006-01-02-Mon"))
		}
	}

	logFile, err := os.Create(LOG_FILE)

	if err != nil {
		log.Panic(err)
	}

	mw := io.MultiWriter(os.Stdout, logFile)
	log, err := logger.New("main", 0, mw)

	if err != nil {
		log.Panic(err.Error())
	}

	return log, logFile
}
