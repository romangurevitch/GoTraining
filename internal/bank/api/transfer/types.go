package transfer

// CreateTransferRequest is the JSON body for POST /v1/transfers.
type CreateTransferRequest struct {
	FromAccountID string `json:"from_account_id" binding:"required"`
	ToAccountID   string `json:"to_account_id"   binding:"required"`
	Amount        int64  `json:"amount"          binding:"required,gte=1"`
}

// TransferResponse is the JSON body returned on successful transfer.
type TransferResponse struct {
	Status string `json:"status"`
}

// DurableTransferRequest is the JSON body for POST /v1/durable-transfers.
type DurableTransferRequest struct {
	FromAccountID string `json:"from_account_id" binding:"required"`
	ToAccountID   string `json:"to_account_id"   binding:"required"`
	Amount        int64  `json:"amount"          binding:"required,gt=0"`
	Reference     string `json:"reference"`
}

// DurableTransferResponse is the JSON body returned when a workflow is started.
type DurableTransferResponse struct {
	TransferID string `json:"transfer_id"`
	Status     string `json:"status"`
}
