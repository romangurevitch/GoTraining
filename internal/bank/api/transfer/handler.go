package transfer

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"

	"github.com/romangurevitch/go-training/internal/bank/api/middleware"
	"github.com/romangurevitch/go-training/internal/bank/domain"
	"github.com/romangurevitch/go-training/internal/bank/service"
	"github.com/romangurevitch/go-training/pkg/api/apierror"
)

// Handler handles transfer-related HTTP requests.
type Handler struct {
	svc service.Service
}

// New creates a new transfer handler.
func New(svc service.Service) *Handler {
	return &Handler{svc: svc}
}

// CreateTransfer handles POST /v1/transfers.
func (h *Handler) CreateTransfer(c *gin.Context) {
	ctx := c.Request.Context()

	ctx, span := otel.Tracer("bank").Start(ctx, "transfer.create")
	defer span.End()

	var req CreateTransferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(apierror.NewAPIError(ctx, http.StatusBadRequest, "invalid request body", err))
		return
	}

	span.SetAttributes(
		attribute.String("transfer.from_account_id", req.FromAccountID),
		attribute.String("transfer.to_account_id", req.ToAccountID),
		attribute.Int64("transfer.amount", req.Amount),
	)

	claims := middleware.ClaimsFromCtx(ctx)
	if claims == nil {
		c.JSON(apierror.NewUnauthorizedError())
		return
	}

	fromAccount, err := h.svc.GetAccount(ctx, req.FromAccountID)
	if err != nil {
		if errors.Is(err, domain.ErrAccountNotFound) {
			c.JSON(apierror.NewAPIError(ctx, http.StatusNotFound, "account not found", err))
			return
		}
		c.JSON(apierror.NewAPIError(ctx, http.StatusInternalServerError, "failed to get account", err))
		return
	}

	if fromAccount.Owner != claims.Subject {
		c.JSON(http.StatusForbidden, &apierror.APIError{Message: "forbidden: you do not own the source account"})
		return
	}

	if err := h.svc.Transfer(ctx, req.FromAccountID, req.ToAccountID, req.Amount); err != nil {
		if errors.Is(err, domain.ErrAccountNotFound) {
			c.JSON(apierror.NewAPIError(ctx, http.StatusNotFound, "account not found", err))
			return
		}
		if errors.Is(err, domain.ErrInsufficientFunds) {
			c.JSON(apierror.NewAPIError(ctx, http.StatusUnprocessableEntity, "insufficient funds", err))
			return
		}
		if errors.Is(err, domain.ErrAccountLocked) {
			c.JSON(apierror.NewAPIError(ctx, http.StatusUnprocessableEntity, "account is locked", err))
			return
		}
		c.JSON(apierror.NewAPIError(ctx, http.StatusInternalServerError, "failed to transfer", err))
		return
	}

	slog.InfoContext(ctx, "transfer completed",
		slog.String("from_account_id", req.FromAccountID),
		slog.String("to_account_id", req.ToAccountID),
		slog.Int64("amount", req.Amount),
	)
	c.JSON(http.StatusOK, TransferResponse{Status: "SUCCESS"})
}

// StartDurableTransfer handles POST /v1/durable-transfers.
func (h *Handler) StartDurableTransfer(c *gin.Context) {
	ctx := c.Request.Context()

	ctx, span := otel.Tracer("bank").Start(ctx, "transfer.start_durable")
	defer span.End()

	var req DurableTransferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(apierror.NewAPIError(ctx, http.StatusBadRequest, "invalid request body", err))
		return
	}

	span.SetAttributes(
		attribute.String("transfer.from_account_id", req.FromAccountID),
		attribute.String("transfer.to_account_id", req.ToAccountID),
		attribute.Int64("transfer.amount", req.Amount),
	)

	claims := middleware.ClaimsFromCtx(ctx)
	if claims == nil {
		c.JSON(apierror.NewUnauthorizedError())
		return
	}

	fromAccount, err := h.svc.GetAccount(ctx, req.FromAccountID)
	if err != nil {
		if errors.Is(err, domain.ErrAccountNotFound) {
			c.JSON(apierror.NewAPIError(ctx, http.StatusNotFound, "account not found", err))
			return
		}
		c.JSON(apierror.NewAPIError(ctx, http.StatusInternalServerError, "failed to get account", err))
		return
	}

	if fromAccount.Owner != claims.Subject {
		c.JSON(http.StatusForbidden, &apierror.APIError{Message: "forbidden: you do not own the source account"})
		return
	}

	transferID, err := h.svc.StartDurableTransfer(ctx, req.FromAccountID, req.ToAccountID, req.Amount, req.Reference)
	if err != nil {
		if errors.Is(err, domain.ErrAccountNotFound) {
			c.JSON(apierror.NewAPIError(ctx, http.StatusNotFound, "account not found", err))
			return
		}
		if errors.Is(err, domain.ErrInsufficientFunds) {
			c.JSON(apierror.NewAPIError(ctx, http.StatusUnprocessableEntity, "insufficient funds", err))
			return
		}
		if errors.Is(err, domain.ErrAccountLocked) {
			c.JSON(apierror.NewAPIError(ctx, http.StatusUnprocessableEntity, "account is locked", err))
			return
		}
		c.JSON(apierror.NewAPIError(ctx, http.StatusInternalServerError, "failed to start durable transfer", err))
		return
	}

	slog.InfoContext(ctx, "durable transfer started",
		slog.String("transfer_id", transferID),
		slog.String("from_account_id", req.FromAccountID),
		slog.String("to_account_id", req.ToAccountID),
		slog.Int64("amount", req.Amount),
	)
	c.JSON(http.StatusAccepted, DurableTransferResponse{
		TransferID: transferID,
		Status:     "PENDING",
	})
}

// ApproveTransfer handles POST /v1/durable-transfers/:id/approve.
func (h *Handler) ApproveTransfer(c *gin.Context) {
	ctx := c.Request.Context()
	transferID := c.Param("id")

	if err := h.svc.ApproveTransfer(ctx, transferID); err != nil {
		c.JSON(apierror.NewAPIError(ctx, http.StatusInternalServerError, "failed to approve transfer", err))
		return
	}

	slog.InfoContext(ctx, "durable transfer approved", slog.String("transfer_id", transferID))
	c.JSON(http.StatusOK, TransferResponse{Status: "APPROVED"})
}

// RejectTransfer handles POST /v1/durable-transfers/:id/reject.
func (h *Handler) RejectTransfer(c *gin.Context) {
	ctx := c.Request.Context()
	transferID := c.Param("id")

	if err := h.svc.RejectTransfer(ctx, transferID); err != nil {
		c.JSON(apierror.NewAPIError(ctx, http.StatusInternalServerError, "failed to reject transfer", err))
		return
	}

	slog.InfoContext(ctx, "durable transfer rejected", slog.String("transfer_id", transferID))
	c.JSON(http.StatusOK, TransferResponse{Status: "REJECTED"})
}
