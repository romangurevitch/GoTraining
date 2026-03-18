package transfer

// CreateTransferRequest is the JSON body for POST /v1/transfers.
// TODO (Step 1 — OpenAPI): You defined these fields in transfers.yaml — use the same names here.
type CreateTransferRequest struct {
	FromAccountID string `json:"from_account_id" binding:"required"`
	ToAccountID   string `json:"to_account_id"   binding:"required"`
	Amount        int64  `json:"amount"           binding:"required,gte=1"`
}

// TransferResponse is the JSON body returned on successful transfer.
// TODO (Step 1 — OpenAPI): Matches the 200 response you defined in transfers.yaml.
type TransferResponse struct {
	Status string `json:"status"`
}
