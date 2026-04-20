package temporal

import (
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type TransferRequest struct {
	TransferID    string `json:"transfer_id"`
	FromAccountID string `json:"from_account_id"`
	ToAccountID   string `json:"to_account_id"`
	Amount        int64  `json:"amount"` // in cents
	Reference     string `json:"reference"`
}

type TransferResponse struct {
	TransferID string `json:"transfer_id"`
	Status     string `json:"status"`
}

const (
	ApprovalSignal = "approval-signal"
	RejectSignal   = "reject-signal"
)

func DurableTransferWorkflow(ctx workflow.Context, req TransferRequest) (TransferResponse, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("Durable transfer workflow started", "transfer_id", req.TransferID)

	options := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second,
			BackoffCoefficient: 2.0,
			MaximumInterval:    time.Minute,
		},
	}
	ctx = workflow.WithActivityOptions(ctx, options)

	var a *Activities

	// 1. Validation
	err := workflow.ExecuteActivity(ctx, a.ValidateAccounts, req).Get(ctx, nil)
	if err != nil {
		return TransferResponse{}, err
	}

	// 2. Approval Gate
	if req.Amount > 1000 {
		var approved bool
		var rejected bool

		workflow.Go(ctx, func(ctx workflow.Context) {
			signalChan := workflow.GetSignalChannel(ctx, ApprovalSignal)
			for {
				signalChan.Receive(ctx, nil)
				approved = true
				workflow.GetLogger(ctx).Info("Received approval signal")
			}
		})

		workflow.Go(ctx, func(ctx workflow.Context) {
			rejectChan := workflow.GetSignalChannel(ctx, RejectSignal)
			for {
				rejectChan.Receive(ctx, nil)
				rejected = true
				workflow.GetLogger(ctx).Info("Received reject signal")
			}
		})

		// Wait for signal or timeout (24h)
		timeout := 24 * time.Hour
		ok, _ := workflow.AwaitWithTimeout(ctx, timeout, func() bool {
			return approved || rejected
		})

		if !ok {
			logger.Warn("Approval timed out")
			return TransferResponse{
				TransferID: req.TransferID,
				Status:     "FAILED",
			}, temporal.NewNonRetryableApplicationError("approval timed out", "TIMEOUT", nil)
		}

		if rejected {
			logger.Warn("Transfer rejected")
			return TransferResponse{
				TransferID: req.TransferID,
				Status:     "REJECTED",
			}, temporal.NewNonRetryableApplicationError("transfer rejected", "REJECTED", nil)
		}
		logger.Info("Transfer approved")
	}

	// 3. Debit
	err = workflow.ExecuteActivity(ctx, a.DebitAccount, req.FromAccountID, req.Amount, req.TransferID).Get(ctx, nil)
	if err != nil {
		return TransferResponse{}, err
	}

	// 4. Credit
	err = workflow.ExecuteActivity(ctx, a.CreditAccount, req.ToAccountID, req.Amount, req.TransferID).Get(ctx, nil)
	if err != nil {
		logger.Error("Credit failed, initiating compensation", "error", err)
		newCtx, _ := workflow.NewDisconnectedContext(ctx)
		compensationErr := workflow.ExecuteActivity(newCtx, a.RefundDebitActivity, req.FromAccountID, req.Amount, req.TransferID).Get(newCtx, nil)
		if compensationErr != nil {
			logger.Error("Compensation (refund) failed — manual intervention required", "error", compensationErr)
		}
		return TransferResponse{
			TransferID: req.TransferID,
			Status:     "FAILED",
		}, err
	}

	return TransferResponse{
		TransferID: req.TransferID,
		Status:     "COMPLETED",
	}, nil
}
