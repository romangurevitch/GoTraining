package transfer

// PARTICIPANT QUEST — implement this handler.
//
// Before starting: read api/account/handler.go in full.
// Every TODO below points to the exact line in account/handler.go that demonstrates the pattern.
//
// Time estimate: 60-80 minutes for Step 3 of the quest.

import (
	"errors"

	"github.com/gin-gonic/gin"

	"github.com/romangurevitch/go-training/internal/bank/api/middleware"
	"github.com/romangurevitch/go-training/internal/bank/domain"
	"github.com/romangurevitch/go-training/internal/bank/service"
	apierror "github.com/romangurevitch/go-training/pkg/api/error"
	// TODO 2: uncomment when implementing OTel span
	// "go.opentelemetry.io/otel"
	// "go.opentelemetry.io/otel/attribute"
	// "log/slog"
	// "net/http"
)

// Handler handles transfer-related HTTP requests.
// Same shape as api/account/handler.go — see Handler and New there.
type Handler struct {
	svc service.Service
}

func New(svc service.Service) *Handler {
	return &Handler{svc: svc}
}

// CreateTransfer handles POST /v1/transfers.
//
// Quest steps:
//
//	TODO 1: Parse request body — api/account/handler.go CreateAccount, line ~40
//	TODO 2: Start OTel span — api/account/handler.go GetAccount, line ~20 (NEW: not in reference solution)
//	TODO 3: Verify ownership (JWT sub must be account owner) — NEW: not in reference solution
//	TODO 4: Call service.Transfer + map errors — api/account/handler.go CreateAccount, lines ~50-60
//	TODO 5: Log success + return 200 — api/account/handler.go CreateAccount, line ~65
func (h *Handler) CreateTransfer(c *gin.Context) {
	ctx := c.Request.Context()

	// TODO 1: Parse and validate request body.
	//   Pattern: api/account/handler.go CreateAccount — ShouldBindJSON + 400 on error
	//
	//   var req CreateTransferRequest
	//   if err := c.ShouldBindJSON(&req); err != nil {
	//       c.JSON(apierror.NewAPIError(ctx, http.StatusBadRequest, "bad request", err))
	//       return
	//   }

	// TODO 2: Start an OTel span.
	//   Pattern: api/account/handler.go GetAccount — otel.Tracer("bank").Start + defer span.End()
	//   NEW: the reference solution had no per-handler tracing — this is the upgrade.
	//
	//   ctx, span := otel.Tracer("bank").Start(ctx, "transfer.create")
	//   defer span.End()
	//   span.SetAttributes(
	//       attribute.String("from_account_id", req.FromAccountID),
	//       attribute.String("to_account_id",   req.ToAccountID),
	//       attribute.Float64("amount",         req.Amount),
	//   )

	// TODO 3: Verify ownership — JWT sub claim must match the from_account owner.
	//   NEW: not in reference solution.
	//
	//   claims := middleware.ClaimsFromCtx(ctx)
	//   fromAccount, err := h.svc.GetAccount(ctx, req.FromAccountID)
	//   switch {
	//   case errors.Is(err, domain.ErrAccountNotFound):
	//       c.JSON(apierror.NewAPIError(ctx, http.StatusNotFound, "account not found", err))
	//       return
	//   case err != nil:
	//       c.JSON(apierror.NewAPIError(ctx, http.StatusInternalServerError, "could not get account", err))
	//       return
	//   }
	//   if fromAccount.Owner != claims.Subject {
	//       c.JSON(apierror.NewAPIError(ctx, http.StatusForbidden, "forbidden: not account owner", nil))
	//       return
	//   }

	// TODO 4: Call service and map errors with errors.Is.
	//   Pattern: api/account/handler.go CreateAccount — switch errors.Is pattern
	//
	//   err = h.svc.Transfer(ctx, req.FromAccountID, req.ToAccountID, req.Amount)
	//   switch {
	//   case errors.Is(err, domain.ErrAccountNotFound):
	//       c.JSON(apierror.NewAPIError(ctx, http.StatusNotFound, "account not found", err))
	//   case errors.Is(err, domain.ErrInsufficientFunds):
	//       c.JSON(apierror.NewAPIError(ctx, http.StatusUnprocessableEntity, "insufficient funds", err))
	//   case errors.Is(err, domain.ErrAccountLocked):
	//       c.JSON(apierror.NewAPIError(ctx, http.StatusUnprocessableEntity, "account locked", err))
	//   case err != nil:
	//       c.JSON(apierror.NewAPIError(ctx, http.StatusInternalServerError, "transfer failed", err))
	//   default:
	//       // TODO 5 goes here
	//   }

	// TODO 5: Log success and return 200.
	//   Pattern: api/account/handler.go CreateAccount — slog.InfoContext + c.JSON
	//   Note: trace_id and span_id are injected automatically by OtelHandler — no extra code.
	//
	//   slog.InfoContext(ctx, "transfer completed",
	//       slog.String("from_account_id", req.FromAccountID),
	//       slog.String("to_account_id",   req.ToAccountID),
	//       slog.Float64("amount",         req.Amount),
	//   )
	//   c.JSON(http.StatusOK, TransferResponse{Status: "completed"})

	// REMOVE THIS LINE when you implement TODO 1:
	_, _ = errors.New(""), middleware.ClaimsFromCtx(ctx) // silence unused imports
	_ = domain.ErrAccountNotFound                        // silence unused import
	_ = apierror.ErrInternalServerError                  // silence unused import
}
