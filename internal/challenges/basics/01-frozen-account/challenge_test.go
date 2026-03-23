package frozenaccount

import "testing"

func TestDeposit_UpdatesBalance(t *testing.T) {
	account := &Account{Balance: 0}
	account.Deposit(100)

	if account.Balance != 100 {
		t.Fatalf(
			"Balance should be 100 after Deposit(100) — got %.2f.\n"+
				"  Hint: check whether Deposit uses a value or pointer receiver.",
			account.Balance,
		)
	}
}

func TestDeposit_MultipleDeposits(t *testing.T) {
	account := &Account{Balance: 50}
	account.Deposit(25)
	account.Deposit(25)

	if account.Balance != 100 {
		t.Fatalf(
			"Balance should be 100 after two Deposit(25) calls — got %.2f.",
			account.Balance,
		)
	}
}
