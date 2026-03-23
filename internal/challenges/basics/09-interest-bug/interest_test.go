package interestbug

import (
	"errors"
	"math"
	"testing"
)

// checkResult is a pre-built test helper. Use it in your table cases.
// t.Helper() ensures that when a test fails, the error points to the
// table row that failed — not to this function.
func checkResult(t *testing.T, got, want float64) {
	t.Helper()
	if math.Abs(got-want) > 0.001 {
		t.Errorf("got %.4f, want %.4f", got, want)
	}
}

// TestCalculate exercises Calculate with multiple scenarios.
// Your job: add at least 5 test cases. Include the case that exposes the hidden bug.
//
// The bug: Calculate does not validate negative rates — it should return ErrNegativeRate.
// One of your test cases must catch this bug and fail until interest.go is fixed.
func TestCalculate(t *testing.T) {
	tests := []struct {
		name      string
		principal float64
		rate      float64
		years     int
		want      float64
		wantErr   error
	}{
		// TODO: add at least 5 test cases.
		// Include:
		//   - a normal case (e.g. 1000 at 5% for 1 year = 1050)
		//   - zero years (result == principal)
		//   - zero rate (result == principal)
		//   - negative years (want error)
		//   - negative rate (want ErrNegativeRate — this one will fail until the bug is fixed!)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Calculate(tt.principal, tt.rate, tt.years)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("expected error %v, got %v", tt.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			checkResult(t, got, tt.want)
		})
	}
}
