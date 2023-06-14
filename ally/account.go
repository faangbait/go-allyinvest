package ally

import (
	"fmt"
	"net/url"
)

// GET accounts
func GetAccounts() IGetAccounts {
	return get[IGetAccounts]("accounts", url.Values{}, authedRL)
}

// GET accounts/balances
func GetAccountBalances() IGetAccountBalances {
	return get[IGetAccountBalances]("accounts/balances", url.Values{}, authedRL)
}

// GET accounts/:id
func GetAccountByID(id string) IGetAccount {
	return get[IGetAccount](fmt.Sprintf("accounts/%s", id), url.Values{}, authedRL)
}

// GET accounts/:id/balances
func GetAccountBalancesByID(id string) IGetAccountBalance {
	return get[IGetAccountBalance](fmt.Sprintf("accounts/%s/balances", id), url.Values{}, authedRL)
}

// GET accounts/:id/history
func GetAccountHistory(id string, daterange string, transactions string) IGetAccountHistory {
	params := url.Values{}
	params.Add("range", daterange)
	params.Add("transactions", transactions)
	return get[IGetAccountHistory](fmt.Sprintf("accounts/%s/history", id), params, authedRL)
}

// GET accounts/:id/holdings
func GetAccountHoldings(id string) IGetAccountHoldings {
	return get[IGetAccountHoldings](fmt.Sprintf("accounts/%s/holdings", id), url.Values{}, authedRL)
}
