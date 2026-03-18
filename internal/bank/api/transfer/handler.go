package transfer

// PARTICIPANT QUEST — implement this handler.
//
// Before starting: read api/account/handler.go in full.
// Every TODO below points to the exact line in account/handler.go that demonstrates the pattern.

import (
	"errors"

	"github.com/gin-gonic/gin"

	"github.com/romangurevitch/go-training/internal/bank/api/middleware"
	"github.com/romangurevitch/go-training/internal/bank/domain"
	"github.com/romangurevitch/go-training/internal/bank/service"
	apierror "github.com/romangurevitch/go-training/pkg/api/error"
)

// Handler handles transfer-related HTTP requests.
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
//	TODO 1: Parse and validate request body (refer to account/handler.go CreateAccount)
//	TODO 2: Start OTel span (refer to account/handler.go GetAccount)
//	TODO 3: Verify ownership — JWT sub claim must match the from_account owner
//	TODO 4: Call service.Transfer and map domain errors (refer to account/handler.go CreateAccount)
//	TODO 5: Log success and return 200 (refer to account/handler.go CreateAccount)
func (h *Handler) CreateTransfer(c *gin.Context) {
	ctx := c.Request.Context()

	// TODO 1: Parse and validate request body.

	// TODO 2: Start an OTel span and set attributes.

	// TODO 3: Verify ownership — JWT sub claim must match the from_account owner.

	// TODO 4: Call service and map errors with errors.Is.

	// TODO 5: Log success and return 200.

	// REMOVE THIS LINE when you implement TODO 1:
	_, _ = errors.New(""), middleware.ClaimsFromCtx(ctx) // silence unused imports
	_ = domain.ErrAccountNotFound                        // silence unused import
	_ = apierror.ErrInternalServerError                  // silence unused import
}
