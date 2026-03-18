package transfer_test

// PARTICIPANT QUEST — complete the test cases marked TODO.
//
// Before starting: read api/account/handler_test.go in full.
// The setup helpers (testToken, setupRouter) follow the identical pattern.
//
// Time estimate: 20-30 minutes for Step 4 of the quest.

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
// Pattern: identical to api/account/handler_test.go — copy it here.
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

// setupRouter builds a minimal Gin engine with the transfer route.
// Pattern: identical to api/account/handler_test.go — copy and adapt for transfers.
func setupRouter(svc *mocks.Service) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := transfer.New(svc)

	v1 := r.Group("/v1/transfers")
	v1.Use(middleware.JWTMiddleware(testSecret))
	{
		v1.POST("", middleware.RequireScope("transfers:write"), h.CreateTransfer)
	}
	return r
}

var aliceAccount = &domain.Account{ID: "ACC-001", Owner: "alice", Balance: 50000, Status: domain.StatusOpen}

// var bobAccount = &domain.Account{ID: "ACC-002", Owner: "bob", Balance: 0, Status: domain.StatusOpen}

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
		// PRE-WRITTEN: Happy path — transfer succeeds
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
		// PRE-WRITTEN: Invalid body — 400
		{
			name: "missing amount — 400",
			fields: fields{
				svc: func(t *testing.T) *mocks.Service {
					return mocks.NewService(t) // no calls expected
				},
			},
			body:     map[string]string{"from_account_id": "ACC-001", "to_account_id": "ACC-002"}, // amount missing
			tokenSub: "alice",
			wantCode: http.StatusBadRequest,
		},

		// TODO: Wrong owner — sub is "alice" but from_account is owned by "bob"
		// Expected: 403
		// Mock setup: GetAccount("ACC-002") returns bobAccount (owner: "bob")
		// Token sub: "alice"
		// {
		//     name: "wrong owner — 403",
		//     ...
		// },

		// TODO: Insufficient funds — transfer amount exceeds from_account balance
		// Expected: 422
		// Mock setup: GetAccount("ACC-001") returns aliceAccount; Transfer returns domain.ErrInsufficientFunds
		// Token sub: "alice"
		// {
		//     name: "insufficient funds — 422",
		//     ...
		// },

		// TODO: Source account not found
		// Expected: 404
		// Mock setup: GetAccount("MISSING") returns nil, domain.ErrAccountNotFound
		// Token sub: "alice"
		// {
		//     name: "source account not found — 404",
		//     ...
		// },
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
