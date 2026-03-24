package feecalculator

import "testing"

func TestSavingsAccount_MonthlyFee(t *testing.T) {
	// SavingsAccount uses a value receiver — both T and *T satisfy the interface.
	var calc FeeCalculator = SavingsAccount{Balance: 100000}
	if got := calc.MonthlyFee(); got != 500 {
		t.Fatalf("SavingsAccount.MonthlyFee() = %d, want 500", got)
	}
}

func TestPremiumAccount_MonthlyFee(t *testing.T) {
	// PremiumAccount must use a pointer receiver — only *PremiumAccount satisfies the interface.
	// If you use a value receiver, this line will not compile.
	var calc FeeCalculator = &PremiumAccount{Balance: 500000}
	if got := calc.MonthlyFee(); got != 2500 {
		t.Fatalf("PremiumAccount.MonthlyFee() = %d, want 2500", got)
	}
}

func TestTotalFees(t *testing.T) {
	accounts := []FeeCalculator{
		SavingsAccount{Balance: 100000},
		&PremiumAccount{Balance: 500000},
	}
	total := TotalFees(accounts)
	if total != 3000 {
		t.Fatalf("TotalFees() = %d, want 3000", total)
	}
}
