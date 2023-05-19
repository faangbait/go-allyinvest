package ally

import (
	"os"
	"strings"
	"testing"
)

func TestGetAccounts(t *testing.T) {
	got := GetAccounts().AccountSummaries

	for _, acct := range got {
		if acct.Balance.AccountValue <= 0 {
			t.Errorf("GetAccounts().value = %v", acct.Balance.AccountValue)
		}
	}

}
func TestGetAccountBalances(t *testing.T) {
	accts := GetAccountBalances()
	acct_last := accts.AccountBalances[len(accts.AccountBalances)-1]

	if acct_last.Accountvalue <= 0 {
		t.Errorf("GetAccountBalances got invalid response: %v", acct_last.Accountvalue)
	}
}
func TestGetAccountByID(t *testing.T) {
	accts := strings.Split(os.Getenv("ALLY_ACCOUNTS"), ":")
	acct1 := accts[len(accts)-1]

	bal := GetAccountByID(acct1)

	if bal.Balance.AccountID[:3] != "3LB" {
		t.Errorf("GetAccountsByID got invalid response: %v", bal.Balance.AccountID)
	}

}
func TestGetAccountBalancesByID(t *testing.T) {
	accts := strings.Split(os.Getenv("ALLY_ACCOUNTS"), ":")
	acct1 := accts[len(accts)-1]

	bal := GetAccountBalancesByID(acct1)

	if bal.Balance.AccountID[:3] != "3LB" {
		t.Errorf("GetAccountsByID got invalid response: %v", bal.Balance.AccountID)
	}
}
func TestGetAccountHistory(t *testing.T) {}
func TestGetAccountHoldings(t *testing.T) {
	accts := strings.Split(os.Getenv("ALLY_ACCOUNTS"), ":")
	acct1 := accts[len(accts)-1]

	bal := GetAccountHoldings(acct1)

	if len(bal.Holdings) == 0 {
		t.Errorf("GetAccountHoldings got invalid response: %v", bal.Holdings)
	}
}
