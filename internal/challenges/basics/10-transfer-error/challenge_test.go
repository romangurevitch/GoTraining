package transfererror

import (
	"errors"
	"testing"
)

func TestTransfer_Success(t *testing.T) {
	err := Transfer("acc-001", "acc-002", 100.0)
	if err != nil {
		t.Fatalf("expected no error for valid transfer, got: %v", err)
	}
}

func TestTransfer_ZeroAmount_ReturnsTransferError(t *testing.T) {
	err := Transfer("acc-001", "acc-002", 0)

	var te *TransferError
	if !errors.As(err, &te) {
		t.Fatalf(
			"expected *TransferError, got: %T %v\n"+
				"  Hint: Transfer must return *TransferError for validation failures.",
			err, err,
		)
	}
	if te.FromID != "acc-001" {
		t.Errorf("TransferError.FromID = %q, want %q", te.FromID, "acc-001")
	}
	if te.Amount != 0 {
		t.Errorf("TransferError.Amount = %.2f, want 0", te.Amount)
	}
}

func TestTransfer_NegativeAmount_ReturnsTransferError(t *testing.T) {
	err := Transfer("acc-X", "acc-Y", -50.0)

	var te *TransferError
	if !errors.As(err, &te) {
		t.Fatalf("expected *TransferError for negative amount, got: %T %v", err, err)
	}
	if te.Amount != -50.0 {
		t.Errorf("TransferError.Amount = %.2f, want -50.00", te.Amount)
	}
}

func TestTransfer_SameAccount_ReturnsTransferError(t *testing.T) {
	err := Transfer("acc-001", "acc-001", 100.0)

	var te *TransferError
	if !errors.As(err, &te) {
		t.Fatalf("expected *TransferError for same-account transfer, got: %T %v", err, err)
	}
	if te.Reason == "" {
		t.Error("TransferError.Reason should describe why the transfer failed")
	}
}

func TestTransferError_ErrorString(t *testing.T) {
	te := &TransferError{
		FromID: "acc-A",
		ToID:   "acc-B",
		Amount: 500.0,
		Reason: "amount must be positive",
	}
	msg := te.Error()
	if msg == "" {
		t.Fatal("TransferError.Error() must return a non-empty string")
	}
}
