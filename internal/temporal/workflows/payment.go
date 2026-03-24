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
	// TODO: add payment processing logic.
	workflow.GetLogger(ctx).Info("Did some processing and recieved the payment", "amount", in.Amount)

	return nil
}
