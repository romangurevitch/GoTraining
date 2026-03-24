package silenterrors

import (
	"errors"
	"testing"
)

func TestWithdraw_Success(t *testing.T) {
	newBalance, err := Withdraw(100, 40)
	if err != nil {
		t.Fatalf("expected no error for valid withdrawal, got: %v", err)
	}
	if newBalance != 60 {
		t.Fatalf("expected new balance 60, got %d", newBalance)
	}
}

func TestWithdraw_InsufficientFunds(t *testing.T) {
	_, err := Withdraw(50, 100)
	if !errors.Is(err, ErrInsufficientFunds) {
		t.Fatalf(
			"expected ErrInsufficientFunds when balance < amount, got: %v\n"+
				"  Hint: use errors.Is, not == or string comparison.",
			err,
		)
	}
}

func TestWithdraw_NegativeAmount(t *testing.T) {
	_, err := Withdraw(100, -10)
	if !errors.Is(err, ErrNegativeAmount) {
		t.Fatalf("expected ErrNegativeAmount for amount <= 0, got: %v", err)
	}
}

func TestWithdraw_ZeroAmount(t *testing.T) {
	_, err := Withdraw(100, 0)
	if !errors.Is(err, ErrNegativeAmount) {
		t.Fatalf("expected ErrNegativeAmount for amount=0, got: %v", err)
	}
}

func TestWithdraw_ExactBalance(t *testing.T) {
	newBalance, err := Withdraw(100, 100)
	if err != nil {
		t.Fatalf("expected no error when withdrawing exact balance, got: %v", err)
	}
	if newBalance != 0 {
		t.Fatalf("expected new balance 0, got %d", newBalance)
	}
}
