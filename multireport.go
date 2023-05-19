package main

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/faangbait/go-allyinvest/ally"
)

type Color string

const (
	Light         Color = "light"
	Danger        Color = "danger"
	Success       Color = "success"
	Info          Color = "info"
	Warning       Color = "warning"
	StrongDanger  Color = "btn-danger"
	StrongSuccess Color = "btn-success"
	StrongInfo    Color = "btn-info"
	StrongWarning Color = "btn-warning"
)

type Sorter string

const (
	LastPrice     Sorter = "Last Price"
	DailyGainLoss Sorter = "Daily Gain and Loss"
	GainLoss      Sorter = "Total Gain and Loss"
	PositionValue Sorter = "Position Size"
	ByAccount     Sorter = "Account"
	Alphabetical  Sorter = "Symbol"
)

func genReports() {
	var securities = map[string]ally.Security{}

	// Load securities map
	ally.RefreshAccounts(securities)

	// Get account data
	accounts := ally.GetAccounts()
	balances := ally.GetAccountBalances()
	summaries := accounts.AccountSummaries
	sort.SliceStable(summaries, func(i, j int) bool {
		return summaries[i].Account < summaries[j].Account
	})

	// Generate accounts table; this is the same for every sorter
	var accountsInner []string
	for _, account := range summaries {
		accountsInner = append(accountsInner, htmlAccountsLine(account))
	}
	tableAccounts := htmlAccountsTable(strings.Join(accountsInner, ""), balances.TotalBalance)

	// Generate list of securities, we don't need symbol lookup
	secList := make([]ally.Security, 0, len(securities))
	for _, sec := range securities {
		secList = append(secList, sec)
	}

	// Write an output file for each sorter
	for _, sorter := range []Sorter{LastPrice, DailyGainLoss, GainLoss, PositionValue, ByAccount, Alphabetical} {
		var positionsInner []string
		var lastAcctId string

		outFile := fmt.Sprintf("logs/report-%s.html", strings.ToLower(strings.ReplaceAll(string(sorter), " ", "")))

		// sort securities by current sorter
		sort.SliceStable(secList, func(i, j int) bool {
			switch sorter {
			case LastPrice:
				return secList[i].Holding.CurrentPrice > secList[j].Holding.CurrentPrice
			case DailyGainLoss:
				return secList[i].Holding.DailyGainLoss > secList[j].Holding.DailyGainLoss
			case GainLoss:
				return (secList[i].Holding.GainLoss / secList[i].Holding.CostBasis) > (secList[j].Holding.GainLoss / secList[j].Holding.CostBasis)
			case PositionValue:
				return secList[i].Holding.NetAssetValue > secList[j].Holding.NetAssetValue
			case ByAccount:
				return strings.Join([]string{secList[i].AccountID, secList[i].Symbol}, "") < strings.Join([]string{secList[j].AccountID, secList[j].Symbol}, "")
			default:
				return secList[i].Symbol < secList[j].Symbol
			}
		})

		for idx, sec := range secList {
			if sorter == ByAccount && sec.AccountID != lastAcctId {
				positionsInner = append(positionsInner, htmlAccountBanner(sec.AccountID))
			}
			positionsInner = append(positionsInner, htmlPositionLine(idx, sec, sorter))
			lastAcctId = sec.AccountID
		}

		// Generate the entire positions table
		tablePositions := htmlPositionTable(strings.Join(positionsInner, ""), sorter)

		// Render the final html
		htmlResult := htmlWrapper(tablePositions, tableAccounts)

		// Write the html to a file
		err := os.WriteFile(outFile, []byte(htmlResult), 0644)

		if err != nil {
			panic(err)
		}
	}

}

func htmlWrapper(tablePositions string, tableAccounts string) string {
	return fmt.Sprintf(`
<html lang="en">
	<head>
		<meta charset="utf-8">
		<meta http-equiv="X-UA-Compatible" content="IE=edge">
		<meta name="viewport" content="width=device-width, initial-scale=1">
		<title>Autotrader Report</title>
		<link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/bootstrap@3.4.1/dist/css/bootstrap.min.css" integrity="sha384-HSMxcRTRxnN+Bdg0JdbxYKrThecOKuH5zCYotlSAcp1+c8xmyTe9GYg1l9a69psu" crossorigin="anonymous">
		<link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/bootstrap@3.4.1/dist/css/bootstrap-theme.min.css" integrity="sha384-6pzBo3FDv/PJ8r2KRkGHifhEocL+1X2rVCTTkUfGk7/0pbek5mMa1upzvWbrUbOZ" crossorigin="anonymous">
	</head>
	<body>
		<div class="container">
		%s%s%s
		</div>
		<script src="https://code.jquery.com/jquery-1.12.4.min.js" integrity="sha384-nvAa0+6Qg9clwYCGGPpDQLVpLNn0fRaROjHqs13t4Ggj3Ez50XnGQqc/r8MhnRDZ" crossorigin="anonymous"></script>
		<script src="https://cdn.jsdelivr.net/npm/bootstrap@3.4.1/dist/js/bootstrap.min.js" integrity="sha384-aJ21OjlMXNL5UyIl/XNwTMqvzeRMZH2w8c5cRVpzpU8Y5bApTppSuUkhZXN0VxHd" crossorigin="anonymous"></script>
	</body>
</html>`,
		htmlNavBar(),
		tableAccounts,
		tablePositions)
}

func htmlNavBar() string {
	return `	<nav class="navbar navbar-default">
				<div class="container-fluid">
					<div class="navbar-header">
						<button type="button" class="navbar-toggle collapsed" data-toggle="collapse" data-target="#bs-example-navbar-collapse-1" aria-expanded="false">
							<span class="sr-only">Toggle navigation</span>
							<span class="icon-bar"></span>
							<span class="icon-bar"></span>
							<span class="icon-bar"></span>
						</button>
						<a class="navbar-brand" href="/">AutoTrader</a>
					</div>
					<div class="collapse navbar-collapse" id="bs-example-navbar-collapse-1">
						<ul class="nav navbar-nav">
							<li class="dropdown">
								<a href="#" class="dropdown-toggle" data-toggle="dropdown" role="button" aria-haspopup="true" aria-expanded="false">Reports <span class="caret"></span></a>
								<ul class="dropdown-menu">
									<li><a href="/reports/account">By Account</a></li>
									<li><a href="/reports/lastprice">By Last Price</a></li>
									<li><a href="/reports/alphabetical">By Symbol</a></li>
									<li role="separator" class="divider"></li>
									<li><a href="/reports/position">By Position Size</a></li>
									<li><a href="/reports/gainloss">By Total G/L</a></li>
									<li><a href="/reports/daily">By Daily G/L</a></li>
								</ul>
						  	</li>
							<li><a href="/trades">Trade Log</a></li>
						</ul>
					</div>
				</div>
			</nav>`
}

func htmlPositionTable(html string, sorter Sorter) string {
	return fmt.Sprintf(`
			<h1>Results by %s</h1>
			<table class='table table-condensed'>
				<thead>
					<tr>
						<th></th>
						<th class='text-center'>SYMBOL</th>
						<th class='text-right'>LAST PRICE</th>
						<th class='text-right'>DAILY G/L</th>
						<th class='text-right'>QTY</th>
						<th class='text-right'>MKT VALUE</th>
						<th class='text-right'>TOTAL G/L %%</th>
					</tr>
				</thead>
				<tbody>%s
				</tbody>
			</table>`, string(sorter), html)
}

func htmlAccountBanner(account string) string {
	return fmt.Sprintf(`
					<tr class='text-center bg-primary'>
						<td colspan=7><strong>%s</strong></td>
					</tr>`, account)
}

func htmlPositionLine(idx int, sec ally.Security, sorter Sorter) string {
	lineColor := Light
	changePct := (sec.Holding.GainLoss / sec.Holding.CostBasis) * 100
	lastPx := sec.Holding.Quote.LastPx

	var use float64

	switch sorter { // In special cases, use total value
	case GainLoss:
		use = changePct
	case PositionValue:
		use = changePct
	case ByAccount:
		use = changePct
	default: // Default to using the Daily statistics
		use = sec.Holding.DailyGainLoss
	}

	if use < 0 {
		lineColor = Danger
	} else if use > 0 {
		lineColor = Success
	}

	if lastPx > 0 { // Override coloring when stock is near minimum limit
		if lastPx < ally.MINIMUM_STOCK_PRICE*1.25 {
			lineColor = StrongWarning
		}

		if lastPx < ally.MINIMUM_STOCK_PRICE {
			lineColor = StrongDanger
		}
	}

	return fmt.Sprintf(`
					<tr class='%v text-right'>
						<th class='text-right'>%d</th>
						<th class='text-center'>%s</th>
						<td>$%.2f</td>
						<td>$%.2f</td>
						<td>%.f</td>
						<td>$%.2f</td>
						<td>%.2f%%</td>
					</tr>`,
		lineColor,
		idx,
		sec.Symbol,
		sec.Holding.CurrentPrice,
		sec.Holding.DailyGainLoss,
		sec.Holding.Qty,
		sec.Holding.NetAssetValue,
		changePct)
}

func htmlAccountsTable(html string, accountvalue float64) string {
	return fmt.Sprintf(`
			<h1>Accounts <span class='text-muted'>(NAV: $%.0f)</span></h1>
			<table class='table table-condensed'>
				<thead>
					<tr>
						<th class='text-center'>Account</th>
						<th class='text-right'>Daily Chg</th>
						<th class='text-right'>Cash</th>
						<th class='text-right'>Securities</th>
						<th class='text-right'>Margin</th>
						<th class='text-right'>Mgn Calls</th>
					</tr>
				</thead>
				<tbody>%s
				</tbody>
			</table>`, accountvalue, html)
}

func htmlAccountsLine(a ally.IAccountSummary) string {
	daily_gl := 0.0
	for _, h := range a.Holdings {
		daily_gl += h.DailyGainLoss
	}

	return fmt.Sprintf(`
					<tr class='text-right'>
						<th class='text-center'>%s</th>
						<td>$%.2f</td>
						<td>$%.2f</td>
						<td>$%.2f</td>
						<td>$%.2f</td>
						<td>$%.2f</td>
					</tr>`,
		a.Account,
		daily_gl,
		a.Balance.MoneyValue.TotalCash,
		a.TotalSecurities,
		a.Balance.MoneyValue.MarginBalance,
		a.Balance.FedCall+a.Balance.HouseCall)
}
