package feecalculator

// FeeCalculator calculates the monthly fee for an account.
// Defined on the consumer side — this package decides what it needs.
type FeeCalculator interface {
	MonthlyFee() int64
}

// SavingsAccount charges a flat $5/month fee.
// Use a VALUE receiver — both SavingsAccount and *SavingsAccount will satisfy FeeCalculator.
type SavingsAccount struct {
	Balance int64
}

// MonthlyFee returns the monthly fee for a SavingsAccount.
// TODO: implement — return 500 (representing $5.00)
func (s SavingsAccount) MonthlyFee() int64 {
	return 500
}

// PremiumAccount charges a flat $25/month fee.
// Use a POINTER receiver — only *PremiumAccount will satisfy FeeCalculator.
// If you use a value receiver here, the test line `var calc FeeCalculator = &PremiumAccount{...}`
// will compile but `var calc FeeCalculator = PremiumAccount{...}` won't — and that's the lesson.
type PremiumAccount struct {
	Balance int64
}

// MonthlyFee returns the monthly fee for a PremiumAccount.
// TODO: implement using a POINTER receiver — return 2500 (representing $25.00)
func (p *PremiumAccount) MonthlyFee() int64 {
	return 2500
}

// TotalFees sums the monthly fees for all accounts.
// Do not modify this function.
func TotalFees(accounts []FeeCalculator) int64 {
	var total int64
	for _, a := range accounts {
		total += a.MonthlyFee()
	}
	return total
}
