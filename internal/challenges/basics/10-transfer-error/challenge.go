package transfererror

import "fmt"

// TransferError carries structured context about a failed transfer.
// Use errors.As to extract this from a returned error.
type TransferError struct {
	FromID string
	ToID   string
	Amount float64
	Reason string
}

// Error implements the error interface.
// TODO: return a descriptive string, e.g.:
//
//	"transfer from acc-001 to acc-002 of $100.00 failed: <reason>"
func (e *TransferError) Error() string {
	_ = fmt.Sprintf // hint: use this
	panic("implement me")
}

// Transfer validates and records a fund transfer.
//
// Rules:
//   - amount must be > 0, otherwise return *TransferError with Reason "amount must be positive"
//   - fromID and toID must differ, otherwise return *TransferError with Reason "cannot transfer to same account"
//   - on success, return nil
//
// TODO: implement the validation logic.
func Transfer(fromID, toID string, amount float64) error {
	panic("implement me")
}
