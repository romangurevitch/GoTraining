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
	apierror "github.com/romangurevitch/go-training/pkg/api/error"
)

type Handler struct {
	svc service.Service
}

func New(svc service.Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) CreateTransfer(c *gin.Context) {
	ctx := c.Request.Context()

	var req CreateTransferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(apierror.NewAPIError(ctx, http.StatusBadRequest, "bad request", err))
		return
	}

	ctx, span := otel.Tracer("bank").Start(ctx, "transfer.create")
	defer span.End()
	span.SetAttributes(
		attribute.String("from_account_id", req.FromAccountID),
		attribute.String("to_account_id", req.ToAccountID),
		attribute.Int64("amount", req.Amount),
	)

	claims := middleware.ClaimsFromCtx(ctx)
	fromAccount, err := h.svc.GetAccount(ctx, req.FromAccountID)
	switch {
	case errors.Is(err, domain.ErrAccountNotFound):
		c.JSON(apierror.NewAPIError(ctx, http.StatusNotFound, "account not found", err))
		return
	case err != nil:
		c.JSON(apierror.NewAPIError(ctx, http.StatusInternalServerError, "could not get account", err))
		return
	}
	if fromAccount.Owner != claims.Subject {
		c.JSON(apierror.NewAPIError(ctx, http.StatusForbidden, "forbidden: not account owner", nil))
		return
	}

	err = h.svc.Transfer(ctx, req.FromAccountID, req.ToAccountID, req.Amount)
	switch {
	case errors.Is(err, domain.ErrAccountNotFound):
		c.JSON(apierror.NewAPIError(ctx, http.StatusNotFound, "account not found", err))
	case errors.Is(err, domain.ErrInsufficientFunds):
		c.JSON(apierror.NewAPIError(ctx, http.StatusUnprocessableEntity, "insufficient funds", err))
	case errors.Is(err, domain.ErrAccountLocked):
		c.JSON(apierror.NewAPIError(ctx, http.StatusUnprocessableEntity, "account locked", err))
	case err != nil:
		c.JSON(apierror.NewAPIError(ctx, http.StatusInternalServerError, "transfer failed", err))
	default:
		slog.InfoContext(ctx, "transfer completed",
			slog.String("from_account_id", req.FromAccountID),
			slog.String("to_account_id", req.ToAccountID),
			slog.Int64("amount", req.Amount),
		)
		c.JSON(http.StatusOK, TransferResponse{Status: "completed"})
	}
}
