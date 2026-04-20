package workflows

import (
	"github.com/google/uuid"
	"go.temporal.io/sdk/workflow"
)

type PaymentDetails struct {
	PayID  uuid.UUID
	Amount int64 // in cents
}

func ProcessPayment(ctx workflow.Context, in PaymentDetails) error {
	workflow.GetLogger(ctx).Info("Processing payment", "pay_id", in.PayID, "amount_cents", in.Amount)
	return nil
}
