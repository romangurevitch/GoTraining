package silenterrors

import (
	"errors"
	"testing"
)

func TestWithdraw_Success(t *testing.T) {
	newBalance, err := Withdraw(100.0, 40.0)
	if err != nil {
		t.Fatalf("expected no error for valid withdrawal, got: %v", err)
	}
	if newBalance != 60.0 {
		t.Fatalf("expected new balance 60.0, got %.2f", newBalance)
	}
}

func TestWithdraw_InsufficientFunds(t *testing.T) {
	_, err := Withdraw(50.0, 100.0)
	if !errors.Is(err, ErrInsufficientFunds) {
		t.Fatalf(
			"expected ErrInsufficientFunds when balance < amount, got: %v\n"+
				"  Hint: use errors.Is, not == or string comparison.",
			err,
		)
	}
}

func TestWithdraw_NegativeAmount(t *testing.T) {
	_, err := Withdraw(100.0, -10.0)
	if !errors.Is(err, ErrNegativeAmount) {
		t.Fatalf("expected ErrNegativeAmount for amount <= 0, got: %v", err)
	}
}

func TestWithdraw_ZeroAmount(t *testing.T) {
	_, err := Withdraw(100.0, 0)
	if !errors.Is(err, ErrNegativeAmount) {
		t.Fatalf("expected ErrNegativeAmount for amount=0, got: %v", err)
	}
}

func TestWithdraw_ExactBalance(t *testing.T) {
	newBalance, err := Withdraw(100.0, 100.0)
	if err != nil {
		t.Fatalf("expected no error when withdrawing exact balance, got: %v", err)
	}
	if newBalance != 0.0 {
		t.Fatalf("expected new balance 0.0, got %.2f", newBalance)
	}
}
