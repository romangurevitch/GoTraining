package feecalculator

import "testing"

func TestSavingsAccount_MonthlyFee(t *testing.T) {
	// SavingsAccount uses a value receiver — both T and *T satisfy the interface.
	var calc FeeCalculator = SavingsAccount{Balance: 1000}
	if got := calc.MonthlyFee(); got != 5.0 {
		t.Fatalf("SavingsAccount.MonthlyFee() = %.2f, want 5.00", got)
	}
}

func TestPremiumAccount_MonthlyFee(t *testing.T) {
	// PremiumAccount must use a pointer receiver — only *PremiumAccount satisfies the interface.
	// If you use a value receiver, this line will not compile.
	var calc FeeCalculator = &PremiumAccount{Balance: 5000}
	if got := calc.MonthlyFee(); got != 25.0 {
		t.Fatalf("PremiumAccount.MonthlyFee() = %.2f, want 25.00", got)
	}
}

func TestTotalFees(t *testing.T) {
	accounts := []FeeCalculator{
		SavingsAccount{Balance: 1000},
		&PremiumAccount{Balance: 5000},
	}
	total := TotalFees(accounts)
	if total != 30.0 {
		t.Fatalf("TotalFees() = %.2f, want 30.00", total)
	}
}
