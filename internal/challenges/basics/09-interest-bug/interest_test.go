package interestbug

import (
	"errors"
	"testing"
)

// checkResult is a pre-built test helper. Use it in your table cases.
// t.Helper() ensures that when a test fails, the error points to the
// table row that failed — not to this function.
func checkResult(t *testing.T, got, want int64) {
	t.Helper()
	if got != want {
		t.Errorf("got %d, want %d", got, want)
	}
}

// TestCalculate exercises Calculate with multiple scenarios.
// Your job: add at least 5 test cases. Include the case that exposes the hidden bug.
//
// The bug: Calculate does not validate negative rates — it should return ErrNegativeRate.
// One of your test cases must catch this bug and fail until interest.go is fixed.
func TestCalculate(t *testing.T) {
	tests := []struct {
		name       string
		principal  int64
		rate       float64
		years      int
		want       int64
		wantErr    error
		wantAnyErr bool
	}{
		{
			name:      "5% for 1 year",
			principal: 100000,
			rate:      0.05,
			years:     1,
			want:      105000,
		},
		{
			name:      "zero years returns principal",
			principal: 100000,
			rate:      0.05,
			years:     0,
			want:      100000,
		},
		{
			name:      "zero rate returns principal",
			principal: 100000,
			rate:      0.0,
			years:     5,
			want:      100000,
		},
		{
			name:      "10% for 3 years compound",
			principal: 100000,
			rate:      0.10,
			years:     3,
			want:      133100,
		},
		{
			name:       "negative years returns error",
			principal:  100000,
			rate:       0.05,
			years:      -1,
			wantAnyErr: true,
		},
		{
			name:      "negative rate returns ErrNegativeRate",
			principal: 100000,
			rate:      -0.05,
			years:     1,
			wantErr:   ErrNegativeRate,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Calculate(tt.principal, tt.rate, tt.years)
			if tt.wantAnyErr {
				if err == nil {
					t.Errorf("expected an error, got nil")
				}
				return
			}
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
