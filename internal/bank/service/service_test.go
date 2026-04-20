package service_test

import (
	"context"
	"testing"

	"github.com/romangurevitch/go-training/internal/bank/domain"
	"github.com/romangurevitch/go-training/internal/bank/repository/mocks"
	"github.com/romangurevitch/go-training/internal/bank/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func openAccount(id, owner string, balance int64) *domain.Account {
	return &domain.Account{
		ID:      id,
		Owner:   owner,
		Balance: balance,
		Status:  domain.StatusOpen,
	}
}

func TestBankService_Deposit(t *testing.T) {
	tests := []struct {
		name            string
		accountOwner    string
		depositAmount   int64
		accountStatus   domain.AccountStatus
		wantErr         error
		expectedBalance int64
	}{
		{
			name:            "Successful deposit",
			accountOwner:    "John Doe",
			depositAmount:   10000,
			accountStatus:   domain.StatusOpen,
			expectedBalance: 10000,
		},
		{
			name:          "Zero amount returns ErrInvalidAmount",
			accountOwner:  "John Doe",
			depositAmount: 0,
			accountStatus: domain.StatusOpen,
			wantErr:       domain.ErrInvalidAmount,
		},
		{
			name:          "Negative amount returns ErrInvalidAmount",
			accountOwner:  "John Doe",
			depositAmount: -100,
			accountStatus: domain.StatusOpen,
			wantErr:       domain.ErrInvalidAmount,
		},
		{
			name:          "Locked account returns ErrAccountLocked",
			accountOwner:  "John Doe",
			depositAmount: 1000,
			accountStatus: domain.StatusLocked,
			wantErr:       domain.ErrAccountLocked,
		},
		{
			name:          "Closed account returns ErrAccountClosed",
			accountOwner:  "John Doe",
			depositAmount: 1000,
			accountStatus: domain.StatusClosed,
			wantErr:       domain.ErrAccountClosed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mocks.NewRepository(t)
			bs := service.NewBankService(repo, nil)
			ctx := context.Background()

			acc := &domain.Account{
				ID:      "ACC-1",
				Owner:   tt.accountOwner,
				Balance: 0,
				Status:  tt.accountStatus,
			}

			if tt.wantErr == nil || (tt.wantErr != domain.ErrInvalidAmount) {
				repo.EXPECT().GetAccount(ctx, "ACC-1").Return(acc, nil)
			}
			if tt.wantErr == nil {
				repo.EXPECT().SaveAccount(ctx, mock.Anything).Return(nil)
				repo.EXPECT().SaveTransaction(ctx, mock.Anything).Return(nil)
			}

			err := bs.Deposit(ctx, acc.ID, tt.depositAmount)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedBalance, acc.Balance)
			}
		})
	}
}

func TestBankService_Withdraw(t *testing.T) {
	tests := []struct {
		name    string
		balance int64
		amount  int64
		wantErr error
	}{
		{
			name:    "Successful withdrawal",
			balance: 10000,
			amount:  5000,
		},
		{
			name:    "Insufficient funds",
			balance: 1000,
			amount:  5000,
			wantErr: domain.ErrInsufficientFunds,
		},
		{
			name:    "Zero amount",
			balance: 10000,
			amount:  0,
			wantErr: domain.ErrInvalidAmount,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mocks.NewRepository(t)
			bs := service.NewBankService(repo, nil)
			ctx := context.Background()

			acc := openAccount("ACC-1", "alice", tt.balance)

			if tt.wantErr != domain.ErrInvalidAmount {
				repo.EXPECT().GetAccount(ctx, "ACC-1").Return(acc, nil)
			}
			if tt.wantErr == nil {
				repo.EXPECT().SaveAccount(ctx, mock.Anything).Return(nil)
				repo.EXPECT().SaveTransaction(ctx, mock.Anything).Return(nil)
			}

			err := bs.Withdraw(ctx, "ACC-1", tt.amount)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.balance-tt.amount, acc.Balance)
			}
		})
	}
}

func TestBankService_Transfer(t *testing.T) {
	tests := []struct {
		name     string
		fromBal  int64
		toBal    int64
		amount   int64
		wantErr  error
		fromErr  error
		toErr    error
	}{
		{
			name:    "Successful transfer",
			fromBal: 10000,
			toBal:   5000,
			amount:  3000,
		},
		{
			name:    "Insufficient funds",
			fromBal: 1000,
			toBal:   5000,
			amount:  5000,
			wantErr: domain.ErrInsufficientFunds,
		},
		{
			name:    "Zero amount",
			fromBal: 10000,
			toBal:   5000,
			amount:  0,
			wantErr: domain.ErrInvalidAmount,
		},
		{
			name:    "Source account not found",
			fromBal: 10000,
			toBal:   5000,
			amount:  1000,
			fromErr: domain.ErrAccountNotFound,
			wantErr: domain.ErrAccountNotFound,
		},
		{
			name:    "Destination account not found",
			fromBal: 10000,
			toBal:   5000,
			amount:  1000,
			toErr:   domain.ErrAccountNotFound,
			wantErr: domain.ErrAccountNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mocks.NewRepository(t)
			bs := service.NewBankService(repo, nil)
			ctx := context.Background()

			from := openAccount("ACC-1", "alice", tt.fromBal)
			to := openAccount("ACC-2", "bob", tt.toBal)

			if tt.wantErr != domain.ErrInvalidAmount {
				if tt.fromErr != nil {
					repo.EXPECT().GetAccount(ctx, "ACC-1").Return(nil, tt.fromErr)
				} else {
					repo.EXPECT().GetAccount(ctx, "ACC-1").Return(from, nil)
					if tt.toErr != nil {
						repo.EXPECT().GetAccount(ctx, "ACC-2").Return(nil, tt.toErr)
					} else {
						repo.EXPECT().GetAccount(ctx, "ACC-2").Return(to, nil)
					}
				}
			}

			if tt.wantErr == nil {
				repo.EXPECT().SaveAccount(ctx, mock.Anything).Return(nil).Times(2)
				repo.EXPECT().SaveTransaction(ctx, mock.Anything).Return(nil).Times(2)
			}

			err := bs.Transfer(ctx, "ACC-1", "ACC-2", tt.amount)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.fromBal-tt.amount, from.Balance)
				assert.Equal(t, tt.toBal+tt.amount, to.Balance)
			}
		})
	}
}

func TestBankService_CreateAccount(t *testing.T) {
	repo := mocks.NewRepository(t)
	bs := service.NewBankService(repo, nil)
	ctx := context.Background()

	repo.EXPECT().SaveAccount(ctx, mock.MatchedBy(func(acc *domain.Account) bool {
		return acc.Owner == "alice" && acc.Balance == 0 && acc.Status == domain.StatusOpen
	})).Return(nil)

	acc, err := bs.CreateAccount(ctx, "alice")
	require.NoError(t, err)
	assert.Equal(t, "alice", acc.Owner)
	assert.Equal(t, int64(0), acc.Balance)
	assert.Equal(t, domain.StatusOpen, acc.Status)
}

func TestBankService_GetAccount(t *testing.T) {
	repo := mocks.NewRepository(t)
	bs := service.NewBankService(repo, nil)
	ctx := context.Background()

	expected := openAccount("ACC-1", "alice", 5000)
	repo.EXPECT().GetAccount(ctx, "ACC-1").Return(expected, nil)

	acc, err := bs.GetAccount(ctx, "ACC-1")
	require.NoError(t, err)
	assert.Equal(t, expected, acc)
}
