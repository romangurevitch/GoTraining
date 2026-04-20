package transfererror

import "fmt"

// TransferError carries structured context about a failed transfer.
// Use errors.As to extract this from a returned error.
type TransferError struct {
	FromID string
	ToID   string
	Amount int64
	Reason string
}

// Error implements the error interface.
// TODO: return a descriptive string, e.g.:
//
//	"transfer from acc-001 to acc-002 of 10000 cents failed: <reason>"
func (e *TransferError) Error() string {
	return fmt.Sprintf("transfer from %s to %s of %d cents failed: %s", e.FromID, e.ToID, e.Amount, e.Reason)
}

// Transfer validates and records a fund transfer.
//
// Rules:
//   - amount must be > 0, otherwise return *TransferError with Reason "amount must be positive"
//   - fromID and toID must differ, otherwise return *TransferError with Reason "cannot transfer to same account"
//   - on success, return nil
//
// TODO: implement the validation logic.
func Transfer(fromID, toID string, amount int64) error {
	if amount <= 0 {
		return &TransferError{FromID: fromID, ToID: toID, Amount: amount, Reason: "amount must be positive"}
	}
	if fromID == toID {
		return &TransferError{FromID: fromID, ToID: toID, Amount: amount, Reason: "cannot transfer to same account"}
	}
	return nil
}
