package api

import (
	"log/slog"

	"github.com/gin-gonic/gin"

	"github.com/romangurevitch/go-training/internal/bank/api/account"
	apiauth "github.com/romangurevitch/go-training/internal/bank/api/auth"
	"github.com/romangurevitch/go-training/internal/bank/api/middleware"
	"github.com/romangurevitch/go-training/internal/bank/api/transfer"
	"github.com/romangurevitch/go-training/internal/bank/service"
)

// Config holds the API server configuration.
type Config struct {
	JWTSecret   string
	ServiceName string // used by OTel span names
}

// NewServer builds the Gin engine with all middleware and routes pre-wired.
//
// Source: go-training-cba-solution/internal/server/rest/account/server.go
// Changes vs reference:
//   - gin.BasicAuth → JWTMiddleware + RequireScope per route
//   - JSONLogMiddleware (logrus) → LoggingMiddleware (slog-gin)
//   - Added TracingMiddleware (otelgin) and RequestIDMiddleware (uuid)
//   - Added transfer group (pre-wired reference; TODO for participants in server.go)
func NewServer(svc service.Service, logger *slog.Logger, cfg Config) *gin.Engine {
	r := gin.New()

	// Middleware stack — applied to all routes in order:
	r.Use(middleware.RequestIDMiddleware())              // 1. Generate + inject X-Request-Id
	r.Use(middleware.TracingMiddleware(cfg.ServiceName)) // 2. Start OTel span (otelgin)
	r.Use(middleware.LoggingMiddleware(logger))          // 3. Structured request log (slog-gin)
	r.Use(gin.Recovery())                                // 4. Recover panics → 500

	// Auth handler — no JWT required (issues tokens)
	authHandler := apiauth.New(cfg.JWTSecret)
	r.POST("/v1/token", authHandler.IssueToken)

	// Accounts — REFERENCE: fully wired, participants read this
	accountHandler := account.New(svc)
	accounts := r.Group("/v1/accounts")
	accounts.Use(middleware.JWTMiddleware(cfg.JWTSecret))
	{
		accounts.GET("/:id", middleware.RequireScope("accounts:read"), accountHandler.GetAccount)
		accounts.POST("", middleware.RequireScope("accounts:write"), accountHandler.CreateAccount)
	}

	// Transfers — TODO for participants (Step 2 of quest)
	// Pattern: identical to the accounts group above.
	transferHandler := transfer.New(svc)
	transfers := r.Group("/v1/transfers")
	transfers.Use(middleware.JWTMiddleware(cfg.JWTSecret))
	{
	    transfers.POST("", middleware.RequireScope("transfers:write"), transferHandler.CreateTransfer)
	}

	return r
}
