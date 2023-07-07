package ally

import (
	"encoding/xml"
	"time"
)

type AllyResponse struct {
	XMLName     xml.Name `xml:"response"`
	ResponseId  string   `xml:"id,attr"`
	ElapsedTime float64  `xml:"elapsedtime"`
	Error       string   `xml:"error"`
}

type FIXMLOrder struct {
	XMLName xml.Name   `xml:"FIXML"`
	Xmlns   string     `xml:"xmlns,attr"`
	Order   fixmlOrder `xml:"Order,omitempty"`
}

type FIXMLChange struct {
	XMLName xml.Name   `xml:"FIXML"`
	Xmlns   string     `xml:"xmlns,attr"`
	Change  fixmlOrder `xml:"OrderCxlRplcReq,omitempty"`
}

type FIXMLCancel struct {
	XMLName xml.Name   `xml:"FIXML"`
	Xmlns   string     `xml:"xmlns,attr"`
	Cancel  fixmlOrder `xml:"OrderCxlReq,omitempty"`
}

type FIXMLFetch struct {
	XMLName xml.Name  `xml:"FIXML"`
	Xmlns   string    `xml:"xmlns,attr"`
	Rpt     fixmlResp `xml:"ExecRpt"`
}

type fixmlResp struct {
	Txt       string        `xml:"Txt,attr,omitempty"`
	TxnTm     time.Time     `xml:"TxnTm,attr,omitempty"`
	TrdDt     time.Time     `xml:"TrdDt,attr,omitempty"`
	LeavesQty float64       `xml:"LeavesQty,attr"`
	CumQty    float64       `xml:"CumQty,attr,omitempty"`
	LastQty   float64       `xml:"LastQty,attr,omitempty"`
	LastPx    float64       `xml:"LastPx,attr,omitempty"`
	Px        float64       `xml:"Px,attr"`
	StopPx    float64       `xml:"StopPx,attr"`
	AvgPx     float64       `xml:"AvgPx,attr,omitempty"`
	TmInForce int           `xml:"TmInForce,attr"`
	Typ       int           `xml:"Typ,attr"`
	Side      int           `xml:"Side,attr"`
	AcctTyp   int           `xml:"AcctTyp,attr"`
	Acct      string        `xml:"Acct,attr"`
	Stat      int           `xml:"Stat,attr"`
	ID        string        `xml:"ID,attr"`
	OrdID     string        `xml:"OrdID,attr"`
	Instrmt   fixmlInstrmt  `xml:"Instrmt,omitempty"`
	OrdQty    fixmlOrdQty   `xml:"OrdQty"`
	FillsGrp  fixmlFillsGrp `xml:"FillsGrp,omitempty"`
}

type fixmlOrder struct {
	Acct        string       `xml:"Acct,attr"`
	AcctTyp     string       `xml:"AcctTyp,attr,omitempty"`
	OrigID      string       `xml:"OrigID,attr,omitempty"`
	OrigClOrdID string       `xml:"OrigClOrdID,attr,omitempty"`
	PosEfct     string       `xml:"PosEfct,attr,omitempty"`
	Px          string       `xml:"Px,attr,omitempty"`
	Side        string       `xml:"Side,attr"`
	StopPx      string       `xml:"StopPx,attr,omitempty"`
	TmInForce   string       `xml:"TmInForce,attr"`
	Typ         string       `xml:"Typ,attr"`
	Instrmt     fixmlInstrmt `xml:"Instrmt,omitempty"`
	OrdQty      fixmlOrdQty  `xml:"OrdQty"`
}

// type fixmlLeg struct {
// 	Text   string `xml:",chardata"`
// 	Side   string `xml:"Side,attr"`
// 	Strk   string `xml:"Strk,attr"`
// 	Mat    string `xml:"Mat,attr"`
// 	MMY    string `xml:"MMY,attr"`
// 	SecTyp string `xml:"SecTyp,attr"`
// 	CFI    string `xml:"CFI,attr"`
// 	Sym    string `xml:"Sym,attr"`
// }

type fixmlInstrmt struct {
	CFI    string `xml:"CFI,attr,omitempty"`
	MatDt  string `xml:"MatDt,attr,omitempty"`
	StrkPx string `xml:"StrkPx,attr,omitempty"`
	SecTyp string `xml:"SecTyp,attr"`
	Sym    string `xml:"Sym,attr"`
}

type fixmlOrdQty struct {
	Qty string `xml:"Qty,attr"`
}

type fixmlFillsGrp struct {
	FillQty string `xml:"FillQty,attr,omitempty"`
	FillPx  string `xml:"FillPx,attr,omitempty"`
}

type WarningCode uint16

const (
	DuplicateOrder       WarningCode = 004  // Instrmt has a pending/open order
	UnsettledFunds       WarningCode = 233  // No closing until funds settle (free riding)
	HigherMarginReq      WarningCode = 255  // Margin calls at 50% instead of 30%
	NotMarginable        WarningCode = 367  // Cash trades only
	MktOrderWhileClosed  WarningCode = 466  // Fills at open; slippage risk
	ExchangeClosed       WarningCode = 548  // Extended hours trading (i think)
	ForeignSettlementFee WarningCode = 563  // Indicates a $50 fee surcharge
	NoMarketOrder1       WarningCode = 2000 // We are not currently accepting Market orders for this security. Please change your order to a Limit order.
	NoMarketOrder2       WarningCode = 2001 // We are not currently accepting Market or Stop orders. Please place a Limit order.
	MaintenanceWindow    WarningCode = 2002 // Due to nightly processing we are unable to accept orders between 11:30 PM and 12:00 AM EST.
)

type IProfile struct {
	AllyResponse
	Account struct {
		AccountID     string `xml:"account"`
		FundEnabled   bool   `xml:"fundtrading"`
		IsIRA         bool   `xml:"ira"`
		MarginEnabled bool   `xml:"margintrading"`
		Nickname      string `xml:"nickname"`
		OptionLevel   uint8  `xml:"optionlevel"`
		IsShared      bool   `xml:"shared"`
		StocksEnabled bool   `xml:"stocktrading"`
	} `xml:"userdata>account"`
	Disabled                    bool `xml:"userdata>disabled"`
	PasswordResetPending        bool `xml:"userdata>resetpassword"`
	PasswordTradingResetPending bool `xml:"userdata>resettradingpassword"`
	User                        struct {
		Entry []struct {
			Name  string `xml:"name"`
			Value string `xml:"value"`
		} `xml:"entry"`
	} `xml:"userdata>userprofile"`
}

type IGetOrders struct {
	AllyResponse
	Orders []IOrder `xml:"orderstatus>order"`
}

type IOrder struct {
	XMLName xml.Name     `xml:"order"`
	FIXML   FixmlMessage `xml:"fixmlmessage"`
}

type FixmlMessage struct {
	Text string `xml:",cdata"`
}

type IPostOrder struct {
	AllyResponse
	SVIOrderID string `xml:"clientorderid"`
	Status     string `xml:"orderstatus"`
}

type IOrderWarn struct {
	Code WarningCode `xml:"warningcode"`
	Text string      `xml:"warningtext"`
}

type IPostOrderWarn struct {
	AllyResponse
	EstCommission float64      `xml:"estcommission"`
	SecFee        float64      `xml:"secfee"`
	MarginReq     float64      `xml:"marginrequirement"`
	Principal     float64      `xml:"principal"`
	FinData       ITradeQuotes `xml:"quotes>instrumentquote"`
	Warning       IOrderWarn   `xml:"warning,omitempty"`
}

type IPostPreview struct {
	AllyResponse
	EstCommission float64      `xml:"estcommission,omitempty"`
	SecFee        float64      `xml:"secfee,omitempty"`
	MarginReq     float64      `xml:"marginrequirement,omitempty"`
	Principal     float64      `xml:"principal,omitempty"`
	WarningText   string       `xml:"warningtext,omitempty"`
	FinData       ITradeQuotes `xml:"quotes>instrumentquote"`
}

type ITradeQuotes struct {
	Greeks     string `xml:"greeks"`
	Instrument struct {
		Description string `xml:"desc"`
		Exchange    string `xml:"exch"`
		SecType     string `xml:"sectyp"`
		Symbol      string `xml:"sym"`
	} `xml:"instrument"`
	Quote struct {
		Ask       float64 `xml:"askprice"`
		Bid       float64 `xml:"bidprice"`
		Change    float64 `xml:"change"`
		ExtraData struct {
			AskSize   uint64 `xml:"asksize"`
			BidSize   uint64 `xml:"bidsize"`
			Dividends struct {
				Amt     float64 `xml:"amt"`
				ExDate  string  `xml:"exdate"`
				PayDate string  `xml:"paydate"`
				Yield   float64 `xml:"yield"`
			} `xml:"dividenddata"`
			EPS           float64 `xml:"earningspershare"`
			HighPxYearly  float64 `xml:"high52price"`
			HighPxDaily   float64 `xml:"highprice"`
			LowPxYearly   float64 `xml:"low52price"`
			LowPxDaily    float64 `xml:"lowprice"`
			OpenPx        float64 `xml:"openprice"`
			PrevClose     float64 `xml:"prevclose"`
			PERatio       float64 `xml:"priceearningsratio"`
			SessionVolume uint64  `xml:"sessionvolume"`
			Volume        uint64  `xml:"volume"`
		} `xml:"extendedquote"`
		LastPx    float64 `xml:"lastprice"`
		ChangePct float64 `xml:"pctchange"`
	} `xml:"quote"`
	UnknownSymbol bool `xml:"unknownsymbol"`
}

type IAPIStatus struct {
	AllyResponse
	Time serverTime `xml:"time"`
}

type IAPIVersion struct {
	AllyResponse
	Version string `xml:"version"`
}

type IGetWatchlists struct {
	AllyResponse
	Watchlists []struct {
		ID string `xml:"id"`
	} `xml:"watchlists>watchlist"`
}

type IGetSingleWatchlist struct {
	AllyResponse
	Watchlists struct {
		ID string `xml:"id"`
	} `xml:"watchlists"`
}
type IGetWatchlist struct {
	AllyResponse
	Items []struct {
		CostBasis  string `xml:"costbasis"`
		Instrument struct {
			Symbol string `xml:"sym"`
		} `xml:"instrument"`
		Qty string `xml:"qty"`
	} `xml:"watchlists>watchlist>watchlistitem"`
}

type IAccountSummary struct {
	Account         string          `xml:"account"`
	Balance         IAccountBalance `xml:"accountbalance"`
	Holdings        []IHolding      `xml:"accountholdings>holding"`
	TotalSecurities float64         `xml:"accountholdings>totalsecurities"`
}

type IGetAccounts struct {
	AllyResponse
	AccountSummaries []IAccountSummary `xml:"accounts>accountsummary"`
}

type IGetAccount struct {
	AllyResponse
	Balance         IAccountBalance `xml:"accountbalance"`
	Holdings        []IHolding      `xml:"accountholdings>holding"`
	TotalSecurities float64         `xml:"accountholdings>totalsecurities"`
}

type IGetAccountBalances struct {
	AllyResponse
	AccountBalances []struct {
		Account      string  `xml:"account"`
		Accountname  string  `xml:"accountname"`
		Accountvalue float64 `xml:"accountvalue"`
	} `xml:"accountbalance"`
	TotalBalance float64 `xml:"totalbalance>accountvalue"`
}

type IGetAccountBalance struct {
	AllyResponse
	Balance IAccountBalance `xml:"accountbalance"`
}

type IGetAccountHistory struct {
	AllyResponse
	Transactions []struct {
		Activity    string                    `xml:"activity"`
		Amount      float64                   `xml:"amount"`
		Date        time.Time                 `xml:"date"`
		Desc        string                    `xml:"desc"`
		Symbol      string                    `xml:"symbol"`
		Transaction IAccountTransactionDetail `xml:"transaction"`
	} `xml:"transactions>transaction"`
}

type IAccountTransactionDetail struct {
	Accounttype uint8   `xml:"accounttype"`
	Commission  float64 `xml:"commission"`
	Description string  `xml:"description"`
	Fee         float64 `xml:"fee"`
	Price       float64 `xml:"price"`
	Quantity    float64 `xml:"quantity"`
	Secfee      float64 `xml:"secfee"`
	Security    struct {
		Cusip  string `xml:"cusip"`
		ID     string `xml:"id"`
		Sectyp string `xml:"sectyp"`
		Sym    string `xml:"sym"`
	} `xml:"security"`
	Settlementdate  time.Time `xml:"settlementdate"`
	Side            uint8     `xml:"side"`
	Source          string    `xml:"source"`
	Tradedate       time.Time `xml:"tradedate"`
	Transactionid   string    `xml:"transactionid"`
	Transactiontype string    `xml:"transactiontype"`
}

type IGetAccountHoldings struct {
	AllyResponse
	TotalSecurities float64    `xml:"accountholdings>totalsecurities"`
	Holdings        []IHolding `xml:"accountholdings>holding"`
}

type IAccountBalance struct {
	AccountID    string  `xml:"account"`
	AccountValue float64 `xml:"accountvalue"`
	BuyingPower  struct {
		Withdrawable  float64 `xml:"cashavailableforwithdrawal,omitempty"`
		DayTrading    float64 `xml:"daytrading,omitempty"`
		EquityPct     float64 `xml:"equitypercentage,omitempty"`
		Options       float64 `xml:"options,omitempty"`
		SodDayTrading float64 `xml:"soddaytrading,omitempty"`
		SodOptions    float64 `xml:"sodoptions,omitempty"`
		SodStock      float64 `xml:"sodstock,omitempty"`
		Stock         float64 `xml:"stock,omitempty"`
	} `xml:"buyingpower"`
	MoneyValue struct {
		AccruedInterest   float64 `xml:"accruedinterest"`
		Cash              float64 `xml:"cash"`
		CashAvailable     float64 `xml:"cashavailable"`
		MarginBalance     float64 `xml:"marginbalance"`
		MoneyMarket       float64 `xml:"mmf"`
		TotalCash         float64 `xml:"total"`
		UnclearedDeposits float64 `xml:"uncleareddeposits"`
		Unsettled         float64 `xml:"unsettledfunds"`
		Yield             float64 `xml:"yield"`
	} `xml:"money"`
	AssetValue struct {
		LongOptions  float64 `xml:"longoptions"`
		LongStocks   float64 `xml:"longstocks"`
		Options      float64 `xml:"options"`
		ShortOptions float64 `xml:"shortoptions"`
		ShortStocks  float64 `xml:"shortstocks"`
		Stocks       float64 `xml:"stocks"`
		Total        float64 `xml:"total"`
	} `xml:"securities"`
	FedCall   float64 `xml:"fedcall"`
	HouseCall float64 `xml:"housecall"`
}

type IHolding struct {
	AccountType uint8   `xml:"accounttype"`
	CostBasis   float64 `xml:"costbasis"`
	GainLoss    float64 `xml:"gainloss"`
	Instrument  struct {
		Cusip        string  `xml:"cusip"`
		Description  string  `xml:"desc"`
		Factor       float64 `xml:"factor"`
		FIXMLSecType string  `xml:"sectyp"`
		FIXMLSymbol  string  `xml:"sym"`
	} `xml:"instrument"`
	NetAssetValue float64 `xml:"marketvalue"`
	DailyGainLoss float64 `xml:"marketvaluechange"`
	CurrentPrice  float64 `xml:"price"`
	PurchasePrice float64 `xml:"purchaseprice"`
	Qty           float64 `xml:"qty"`
	Quote         struct {
		PxGainLoss float64 `xml:"change"`
		LastPx     float64 `xml:"lastprice"`
		Bid        float64
		Ask        float64
	} `xml:"quote"`
	Underlying string `xml:"underlying"`
}

type IClock struct {
	AllyResponse
	Date   clockTime `xml:"date"`
	Status struct {
		Current string `xml:"current"`
	} `xml:"status"`
	Message  string  `xml:"message"`
	Unixtime float64 `xml:"unixtime"`
}

type IGetQuotes struct {
	AllyResponse
	Quotes struct {
		QuoteType string        `xml:"quotetype"`
		Quote     []IStockQuote `xml:"quote"`
	} `xml:"quotes"`
}

type IStockQuote struct {
	Adp100       float64 `xml:"adp_100,omitempty"`
	Adp200       float64 `xml:"adp_200,omitempty"`
	Adp50        float64 `xml:"adp_50,omitempty"`
	Adv21        uint64  `xml:"adv_21,omitempty"`
	Adv30        uint64  `xml:"adv_30,omitempty"`
	Adv90        uint64  `xml:"adv_90,omitempty"`
	Ask          float64 `xml:"ask"`
	AskTime      float64 `xml:"ask_time,omitempty"`
	Asksz        uint32  `xml:"asksz,omitempty"`
	Basis        int     `xml:"basis,omitempty"`
	Beta         float64 `xml:"beta,omitempty"`
	Bid          float64 `xml:"bid"`
	BidTime      string  `xml:"bid_time,omitempty"`
	Bidsz        uint32  `xml:"bidsz,omitempty"`
	Bidtick      uint32  `xml:"bidtick,omitempty"`
	Chg          float64 `xml:"chg,omitempty"`
	ChgSign      string  `xml:"chg_sign,omitempty"`
	ChgT         float64 `xml:"chg_t,omitempty"`
	Cl           float64 `xml:"cl,omitempty"`
	Cusip        string  `xml:"cusip,omitempty"`
	Date         string  `xml:"date,omitempty"`
	Datetime     string  `xml:"datetime,omitempty"`
	Div          float64 `xml:"div,omitempty"`
	Divexdate    string  `xml:"divexdate,omitempty"`
	Divpaydt     string  `xml:"divpaydt,omitempty"`
	DollarValue  float64 `xml:"dollar_value,omitempty"`
	Eps          float64 `xml:"eps,omitempty"`
	Exch         string  `xml:"exch,omitempty"`
	ExchDesc     string  `xml:"exch_desc,omitempty"`
	Hi           float64 `xml:"hi,omitempty"`
	Iad          uint32  `xml:"iad,omitempty"`
	IncrVl       uint32  `xml:"incr_vl,omitempty"`
	Last         float64 `xml:"last,omitempty"`
	Lo           float64 `xml:"lo,omitempty"`
	Name         string  `xml:"name,omitempty"`
	OpFlag       uint32  `xml:"op_flag,omitempty"`
	Opn          float64 `xml:"opn,omitempty"`
	Pchg         string  `xml:"pchg,omitempty"`
	PchgSign     string  `xml:"pchg_sign,omitempty"`
	Pcls         float64 `xml:"pcls,omitempty"`
	Pe           float64 `xml:"pe,omitempty"`
	Phi          float64 `xml:"phi,omitempty"`
	Plo          float64 `xml:"plo,omitempty"`
	Popn         float64 `xml:"popn,omitempty"`
	PrAdp100     float64 `xml:"pr_adp_100,omitempty"`
	PrAdp200     float64 `xml:"pr_adp_200,omitempty"`
	PrAdp50      float64 `xml:"pr_adp_50,omitempty"`
	PrDate       string  `xml:"pr_date,omitempty"`
	Prbook       float64 `xml:"prbook,omitempty"`
	Prchg        float64 `xml:"prchg,omitempty"`
	Pvol         uint64  `xml:"pvol,omitempty"`
	Secclass     uint32  `xml:"secclass,omitempty"`
	Sesn         string  `xml:"sesn,omitempty"`
	Sho          string  `xml:"sho,omitempty"`
	Symbol       string  `xml:"symbol,omitempty"`
	Tcond        string  `xml:"tcond,omitempty"`
	Timestamp    uint64  `xml:"timestamp,omitempty"`
	TrNum        uint64  `xml:"tr_num,omitempty"`
	Tradetick    string  `xml:"tradetick,omitempty"`
	Trend        string  `xml:"trend,omitempty"`
	Vl           uint64  `xml:"vl,omitempty"`
	Volatility12 float64 `xml:"volatility12,omitempty"`
	Vwap         float64 `xml:"vwap,omitempty"`
	Wk52hi       float64 `xml:"wk52hi,omitempty"`
	Wk52hidate   string  `xml:"wk52hidate,omitempty"`
	Wk52lo       float64 `xml:"wk52lo,omitempty"`
	Wk52lodate   string  `xml:"wk52lodate,omitempty"`
}

type INewsSearch struct {
	AllyResponse
	Articles struct {
		Article []INewsArticle `xml:"article"`
	} `xml:"articles"`
}

type INewsById struct {
	AllyResponse
	Article INewsArticle `xml:"article"`
}

type INewsArticle struct {
	Date     time.Time `xml:"date"`
	Headline string    `xml:"headline"`
	ID       string    `xml:"id"`
	Story    string    `xml:"story"`
}

type IOptionsSearch struct {
	AllyResponse
	Quotes struct {
		Quote []struct {
			Ask              string `xml:"ask"`
			AskTime          string `xml:"ask_time"`
			Asksz            string `xml:"asksz"`
			Basis            string `xml:"basis"`
			Bid              string `xml:"bid"`
			BidTime          string `xml:"bid_time"`
			Bidsz            string `xml:"bidsz"`
			Chg              string `xml:"chg"`
			ChgSign          string `xml:"chg_sign"`
			ChgT             string `xml:"chg_t"`
			Cl               string `xml:"cl"`
			ContractSize     string `xml:"contract_size"`
			Date             string `xml:"date"`
			Datetime         string `xml:"datetime"`
			DaysToExpiration string `xml:"days_to_expiration"`
			Exch             string `xml:"exch"`
			ExchDesc         string `xml:"exch_desc"`
			Hi               string `xml:"hi"`
			IncrVl           string `xml:"incr_vl"`
			IssueDesc        string `xml:"issue_desc"`
			Last             string `xml:"last"`
			Lo               string `xml:"lo"`
			OpDelivery       string `xml:"op_delivery"`
			OpFlag           string `xml:"op_flag"`
			OpStyle          string `xml:"op_style"`
			OpSubclass       string `xml:"op_subclass"`
			Opn              string `xml:"opn"`
			Pchg             string `xml:"pchg"`
			PchgSign         string `xml:"pchg_sign"`
			Pcls             string `xml:"pcls"`
			Phi              string `xml:"phi"`
			Plo              string `xml:"plo"`
			Popn             string `xml:"popn"`
			PrOpeninterest   string `xml:"pr_openinterest"`
			Prchg            string `xml:"prchg"`
			PremMult         string `xml:"prem_mult"`
			PutCall          string `xml:"put_call"`
			Pvol             string `xml:"pvol"`
			Rootsymbol       string `xml:"rootsymbol"`
			Secclass         string `xml:"secclass"`
			Sesn             string `xml:"sesn"`
			Strikeprice      string `xml:"strikeprice"`
			Symbol           string `xml:"symbol"`
			Tcond            string `xml:"tcond"`
			Timestamp        string `xml:"timestamp"`
			TrNum            string `xml:"tr_num"`
			Tradetick        string `xml:"tradetick"`
			UnderCusip       string `xml:"under_cusip"`
			Undersymbol      string `xml:"undersymbol"`
			Vl               string `xml:"vl"`
			Vwap             string `xml:"vwap"`
			Wk52hi           string `xml:"wk52hi"`
			Wk52lo           string `xml:"wk52lo"`
			Xdate            string `xml:"xdate"`
			Xday             string `xml:"xday"`
			Xmonth           string `xml:"xmonth"`
			Xyear            string `xml:"xyear"`
			PrDate           string `xml:"pr_date"`
			Wk52hidate       string `xml:"wk52hidate"`
			Wk52lodate       string `xml:"wk52lodate"`
		} `xml:"quote"`
	} `xml:"quotes"`
}

type IOptionsStrike struct {
	AllyResponse
	Prices struct {
		Price []string `xml:"price"`
	} `xml:"prices"`
}

type IOptionsExpire struct {
	AllyResponse
	Expirationdates struct {
		Date []string `xml:"date"`
	} `xml:"expirationdates"`
}

type ITimeSales struct {
	AllyResponse
	Quotes struct {
		Quote []struct {
			Date      string `xml:"date"`
			Datetime  string `xml:"datetime"`
			Hi        string `xml:"hi"`
			IncrVl    string `xml:"incr_vl"`
			Last      string `xml:"last"`
			Lo        string `xml:"lo"`
			Opn       string `xml:"opn"`
			Timestamp string `xml:"timestamp"`
			Vl        string `xml:"vl"`
		} `xml:"quote"`
	} `xml:"quotes"`
}

type ITopLists struct {
	AllyResponse
	Quotes struct {
		Quote []struct {
			Chg     string `xml:"chg"`
			ChgSign string `xml:"chg_sign"`
			Last    string `xml:"last"`
			Name    string `xml:"name"`
			Pchg    string `xml:"pchg"`
			Pcls    string `xml:"pcls"`
			Rank    string `xml:"rank"`
			Symbol  string `xml:"symbol"`
			Vl      string `xml:"vl"`
		} `xml:"quote"`
	} `xml:"quotes"`
}

type Security struct {
	AccountID string
	Holding   IHolding
	Rating    SecurityRating
	Grade     SecurityGrade
	Quote     IStockQuote
	Exchange  string
	Symbol    string
}

type SecurityGrade struct {
	Valuation     float64
	Growth        float64
	Profitability float64
	Momentum      float64
	EPSRevision   float64
}

type SecurityRating struct {
	Quant   float64
	Author  float64
	Analyst float64
}

type Order struct {
	Symbol        string
	SecType       string
	Account       string
	Action        string
	OrderType     string
	Tif           string
	TotalQuantity float64
	LmtPrice      float64
	StopPrice     float64
	Transmit      bool
	Orig_ID       string
	DoNotCoax     bool
}
