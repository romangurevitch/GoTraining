package transfer_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/romangurevitch/go-training/internal/bank/api/middleware"
	"github.com/romangurevitch/go-training/internal/bank/api/transfer"
	"github.com/romangurevitch/go-training/internal/bank/domain"
	"github.com/romangurevitch/go-training/internal/bank/service/mocks"
)

const testSecret = "test-secret"

// testToken issues a signed JWT for test Authorization headers.
func testToken(t *testing.T, sub, scope string) string {
	t.Helper()
	claims := middleware.Claims{
		Scope: scope,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   sub,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(testSecret))
	require.NoError(t, err)
	return signed
}

// setupRouter builds a minimal Gin engine with transfer and durable-transfer routes.
func setupRouter(svc *mocks.Service) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := transfer.New(svc)

	v1 := r.Group("/v1/transfers")
	v1.Use(middleware.JWTMiddleware(testSecret))
	{
		v1.POST("", middleware.RequireScope("transfers:write"), h.CreateTransfer)
	}

	durable := r.Group("/v1/durable-transfers")
	durable.Use(middleware.JWTMiddleware(testSecret))
	{
		durable.POST("", middleware.RequireScope("transfers:write"), h.StartDurableTransfer)
		durable.POST("/:id/approve", middleware.RequireScope("transfers:write"), h.ApproveTransfer)
		durable.POST("/:id/reject", middleware.RequireScope("transfers:write"), h.RejectTransfer)
	}
	return r
}

var aliceAccount = &domain.Account{ID: "ACC-001", Owner: "alice", Balance: 50000, Status: domain.StatusOpen}

func TestCreateTransfer(t *testing.T) {
	type fields struct {
		svc func(t *testing.T) *mocks.Service
	}
	tests := []struct {
		name     string
		fields   fields
		body     any
		tokenSub string
		wantCode int
	}{
		{
			name: "success — 200",
			fields: fields{
				svc: func(t *testing.T) *mocks.Service {
					m := mocks.NewService(t)
					m.EXPECT().GetAccount(mock.Anything, "ACC-001").Return(aliceAccount, nil).Once()
					m.EXPECT().Transfer(mock.Anything, "ACC-001", "ACC-002", int64(5000)).Return(nil).Once()
					return m
				},
			},
			body:     map[string]any{"from_account_id": "ACC-001", "to_account_id": "ACC-002", "amount": 5000},
			tokenSub: "alice",
			wantCode: http.StatusOK,
		},
		{
			name: "missing amount — 400",
			fields: fields{
				svc: func(t *testing.T) *mocks.Service {
					return mocks.NewService(t)
				},
			},
			body:     map[string]string{"from_account_id": "ACC-001", "to_account_id": "ACC-002"},
			tokenSub: "alice",
			wantCode: http.StatusBadRequest,
		},

		{
			name: "wrong owner — 403",
			fields: fields{
				svc: func(t *testing.T) *mocks.Service {
					m := mocks.NewService(t)
					bobAccount := &domain.Account{ID: "ACC-001", Owner: "bob", Balance: 50000, Status: domain.StatusOpen}
					m.EXPECT().GetAccount(mock.Anything, "ACC-001").Return(bobAccount, nil).Once()
					return m
				},
			},
			body:     map[string]any{"from_account_id": "ACC-001", "to_account_id": "ACC-002", "amount": 5000},
			tokenSub: "alice",
			wantCode: http.StatusForbidden,
		},
		{
			name: "insufficient funds — 422",
			fields: fields{
				svc: func(t *testing.T) *mocks.Service {
					m := mocks.NewService(t)
					m.EXPECT().GetAccount(mock.Anything, "ACC-001").Return(aliceAccount, nil).Once()
					m.EXPECT().Transfer(mock.Anything, "ACC-001", "ACC-002", int64(5000)).Return(domain.ErrInsufficientFunds).Once()
					return m
				},
			},
			body:     map[string]any{"from_account_id": "ACC-001", "to_account_id": "ACC-002", "amount": 5000},
			tokenSub: "alice",
			wantCode: http.StatusUnprocessableEntity,
		},
		{
			name: "source account not found — 404",
			fields: fields{
				svc: func(t *testing.T) *mocks.Service {
					m := mocks.NewService(t)
					m.EXPECT().GetAccount(mock.Anything, "ACC-001").Return(nil, domain.ErrAccountNotFound).Once()
					return m
				},
			},
			body:     map[string]any{"from_account_id": "ACC-001", "to_account_id": "ACC-002", "amount": 5000},
			tokenSub: "alice",
			wantCode: http.StatusNotFound,
		},
		{
			name: "account locked — 422",
			fields: fields{
				svc: func(t *testing.T) *mocks.Service {
					m := mocks.NewService(t)
					m.EXPECT().GetAccount(mock.Anything, "ACC-001").Return(aliceAccount, nil).Once()
					m.EXPECT().Transfer(mock.Anything, "ACC-001", "ACC-002", int64(5000)).Return(domain.ErrAccountLocked).Once()
					return m
				},
			},
			body:     map[string]any{"from_account_id": "ACC-001", "to_account_id": "ACC-002", "amount": 5000},
			tokenSub: "alice",
			wantCode: http.StatusUnprocessableEntity,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := setupRouter(tt.fields.svc(t))
			bodyBytes, _ := json.Marshal(tt.body)
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/v1/transfers", bytes.NewBuffer(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+testToken(t, tt.tokenSub, "transfers:write"))
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantCode, w.Code)
		})
	}
}

func TestStartDurableTransfer(t *testing.T) {
	tests := []struct {
		name     string
		svc      func(t *testing.T) *mocks.Service
		body     any
		tokenSub string
		wantCode int
	}{
		{
			name: "success — 202",
			svc: func(t *testing.T) *mocks.Service {
				m := mocks.NewService(t)
				m.EXPECT().GetAccount(mock.Anything, "ACC-001").Return(aliceAccount, nil).Once()
				m.EXPECT().StartDurableTransfer(mock.Anything, "ACC-001", "ACC-002", int64(5000), "ref-1").Return("TRX-123", nil).Once()
				return m
			},
			body:     map[string]any{"from_account_id": "ACC-001", "to_account_id": "ACC-002", "amount": 5000, "reference": "ref-1"},
			tokenSub: "alice",
			wantCode: http.StatusAccepted,
		},
		{
			name: "wrong owner — 403",
			svc: func(t *testing.T) *mocks.Service {
				m := mocks.NewService(t)
				bobAccount := &domain.Account{ID: "ACC-001", Owner: "bob", Balance: 50000, Status: domain.StatusOpen}
				m.EXPECT().GetAccount(mock.Anything, "ACC-001").Return(bobAccount, nil).Once()
				return m
			},
			body:     map[string]any{"from_account_id": "ACC-001", "to_account_id": "ACC-002", "amount": 5000},
			tokenSub: "alice",
			wantCode: http.StatusForbidden,
		},
		{
			name: "account not found — 404",
			svc: func(t *testing.T) *mocks.Service {
				m := mocks.NewService(t)
				m.EXPECT().GetAccount(mock.Anything, "ACC-001").Return(nil, domain.ErrAccountNotFound).Once()
				return m
			},
			body:     map[string]any{"from_account_id": "ACC-001", "to_account_id": "ACC-002", "amount": 5000},
			tokenSub: "alice",
			wantCode: http.StatusNotFound,
		},
		{
			name: "missing body — 400",
			svc: func(t *testing.T) *mocks.Service {
				return mocks.NewService(t)
			},
			body:     map[string]string{},
			tokenSub: "alice",
			wantCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := setupRouter(tt.svc(t))
			bodyBytes, _ := json.Marshal(tt.body)
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/v1/durable-transfers", bytes.NewBuffer(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+testToken(t, tt.tokenSub, "transfers:write"))
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantCode, w.Code)
		})
	}
}

func TestApproveTransfer(t *testing.T) {
	t.Run("success — 200", func(t *testing.T) {
		m := mocks.NewService(t)
		m.EXPECT().ApproveTransfer(mock.Anything, "TRX-123").Return(nil).Once()

		router := setupRouter(m)
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/v1/durable-transfers/TRX-123/approve", nil)
		req.Header.Set("Authorization", "Bearer "+testToken(t, "alice", "transfers:write"))
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("service error — 500", func(t *testing.T) {
		m := mocks.NewService(t)
		m.EXPECT().ApproveTransfer(mock.Anything, "TRX-456").Return(assert.AnError).Once()

		router := setupRouter(m)
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/v1/durable-transfers/TRX-456/approve", nil)
		req.Header.Set("Authorization", "Bearer "+testToken(t, "alice", "transfers:write"))
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestRejectTransfer(t *testing.T) {
	t.Run("success — 200", func(t *testing.T) {
		m := mocks.NewService(t)
		m.EXPECT().RejectTransfer(mock.Anything, "TRX-789").Return(nil).Once()

		router := setupRouter(m)
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/v1/durable-transfers/TRX-789/reject", nil)
		req.Header.Set("Authorization", "Bearer "+testToken(t, "alice", "transfers:write"))
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("service error — 500", func(t *testing.T) {
		m := mocks.NewService(t)
		m.EXPECT().RejectTransfer(mock.Anything, "TRX-000").Return(assert.AnError).Once()

		router := setupRouter(m)
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/v1/durable-transfers/TRX-000/reject", nil)
		req.Header.Set("Authorization", "Bearer "+testToken(t, "alice", "transfers:write"))
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
