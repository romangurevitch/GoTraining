package silenterrors

import "errors"

// ErrInsufficientFunds is returned when the withdrawal amount exceeds the balance.
var ErrInsufficientFunds = errors.New("insufficient funds")

// ErrNegativeAmount is returned when the withdrawal amount is zero or negative.
var ErrNegativeAmount = errors.New("negative or zero amount")

// Withdraw deducts amount from balance and returns the new balance.
//
// Returns (0, ErrNegativeAmount) if amount <= 0.
// Returns (0, ErrInsufficientFunds) if balance < amount.
// Returns (balance - amount, nil) on success.
func Withdraw(balance, amount int64) (int64, error) {
	panic("implement me")
}
