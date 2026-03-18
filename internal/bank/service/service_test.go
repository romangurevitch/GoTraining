package service_test

import (
	"context"
	"testing"

	"github.com/romangurevitch/go-training/internal/bank/service"
	"github.com/romangurevitch/go-training/internal/bank/store"
	"github.com/stretchr/testify/assert"
)

func TestBankService_Deposit(t *testing.T) {
	// QUEST 6: Participants should implement this table-driven test.

	type args struct {
		accountOwner  string
		depositAmount float64
	}

	tests := []struct {
		name            string
		args            args
		wantErr         bool
		expectedBalance float64
	}{
		{
			name: "Successful deposit",
			args: args{
				accountOwner:  "John Doe",
				depositAmount: 100.0,
			},
			wantErr:         false,
			expectedBalance: 100.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := store.NewMemoryStore()
			bs := service.NewBankService(s)
			ctx := context.Background()

			// Create account
			acc, err := bs.CreateAccount(ctx, tt.args.accountOwner)
			assert.NoError(t, err)

			// Perform deposit
			err = bs.Deposit(ctx, acc.ID, tt.args.depositAmount)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				updatedAcc, _ := bs.GetAccount(ctx, acc.ID)
				assert.Equal(t, tt.expectedBalance, updatedAcc.Balance)
			}
		})
	}
}
