package frozenaccount

// Account holds a bank account balance.
type Account struct {
	Balance float64
}

// Deposit adds amount to the account balance.
//
// BUG: this method uses a value receiver — it operates on a copy of Account,
// so the original balance is never updated.
//
// TODO: fix the receiver so deposits actually change the balance.
func (a Account) Deposit(amount float64) {
	a.Balance += amount
}
