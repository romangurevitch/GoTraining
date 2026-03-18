package account

// REFERENCE IMPLEMENTATION — participants read this before implementing api/transfer/handler.go.
//
// Source: derived from reference solution
// Changes vs original:
//   - logrus logger.WithContext(ctx).WithField(...).Info() → slog.InfoContext(ctx, ...)
//   - OTel span added per handler (original had no per-handler tracing)
//   - Direct store access → service.Service interface

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"

	"github.com/romangurevitch/go-training/internal/bank/domain"
	"github.com/romangurevitch/go-training/internal/bank/service"
	apierror "github.com/romangurevitch/go-training/pkg/api/error"
)

// Handler handles account-related HTTP requests.
// Identical pattern to reference accountServer — just wired to service.Service.
type Handler struct {
	svc service.Service
}

func New(svc service.Service) *Handler {
	return &Handler{svc: svc}
}

// GetAccount handles GET /v1/accounts/:id
//
// Pattern demonstrated (in order):
//  1. Extract ctx from request
//  2. Start OTel span — defer span.End()
//  3. Read URL param
//  4. slog.InfoContext — structured log with context (trace_id auto-injected)
//  5. Call service method
//  6. Map errors with errors.Is — return apierror.NewAPIError tuple to c.JSON
//  7. On success: set span attribute + return 200
func (h *Handler) GetAccount(c *gin.Context) {
	ctx := c.Request.Context()

	ctx, span := otel.Tracer("bank").Start(ctx, "account.get")
	defer span.End()

	id := c.Param("id")
	slog.InfoContext(ctx, "get account", slog.String("account_id", id))

	result, err := h.svc.GetAccount(ctx, id)
	switch {
	case errors.Is(err, domain.ErrAccountNotFound):
		c.JSON(apierror.NewAPIError(ctx, http.StatusNotFound, "account not found", err))
	case err != nil:
		c.JSON(apierror.NewAPIError(ctx, http.StatusInternalServerError, "could not get account", err))
	default:
		span.SetAttributes(attribute.String("account.owner", result.Owner))
		c.JSON(http.StatusOK, toAccountResponse(result))
	}
}

// CreateAccount handles POST /v1/accounts
//
// Pattern demonstrated (in addition to GetAccount):
//  1. ShouldBindJSON — bind + validate request body
//  2. On bind error: 400 with apierror
//  3. Call service, map domain errors
//  4. On success: log created entity + return 201
func (h *Handler) CreateAccount(c *gin.Context) {
	ctx := c.Request.Context()

	ctx, span := otel.Tracer("bank").Start(ctx, "account.create")
	defer span.End()

	var req CreateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(apierror.NewAPIError(ctx, http.StatusBadRequest, "bad request", err))
		return
	}

	result, err := h.svc.CreateAccount(ctx, req.Owner)
	switch {
	case errors.Is(err, domain.ErrAccountAlreadyExists):
		c.JSON(apierror.NewAPIError(ctx, http.StatusConflict, "account already exists", err))
	case err != nil:
		c.JSON(apierror.NewAPIError(ctx, http.StatusInternalServerError, "could not create account", err))
	default:
		span.SetAttributes(attribute.String("account.id", result.ID))
		slog.InfoContext(ctx, "account created", slog.String("account_id", result.ID))
		c.JSON(http.StatusCreated, toAccountResponse(result))
	}
}
