package interestbug

import "errors"

// ErrNegativeRate is returned when the interest rate is negative.
var ErrNegativeRate = errors.New("interest rate cannot be negative")

// Calculate computes compound interest.
//
//	principal: starting amount (must be > 0)
//	rate:      annual interest rate as a decimal (e.g. 0.05 for 5%)
//	years:     number of years
//
// Returns the final amount after compound interest, or an error.
//
// BUG: negative rate is not validated — Calculate returns a result instead of ErrNegativeRate.
// (Do not fix this — your job is to write tests that CATCH this bug.)
func Calculate(principal, rate float64, years int) (float64, error) {
	if years < 0 {
		return 0, errors.New("years cannot be negative")
	}
	// BUG: missing: if rate < 0 { return 0, ErrNegativeRate }
	result := principal
	for i := 0; i < years; i++ {
		result *= 1 + rate
	}
	return result, nil
}
