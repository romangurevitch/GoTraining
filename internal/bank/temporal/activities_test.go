package temporal

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/romangurevitch/go-training/internal/bank/domain"
	"github.com/romangurevitch/go-training/internal/bank/repository/mocks"
)

func openAccount(id, owner string, balance int64) *domain.Account {
	return &domain.Account{
		ID:        id,
		Owner:     owner,
		Balance:   balance,
		Status:    domain.StatusOpen,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func lockedAccount(id, owner string) *domain.Account {
	return &domain.Account{
		ID:     id,
		Owner:  owner,
		Status: domain.StatusLocked,
	}
}

// --- ValidateAccounts Tests ---

func TestValidateAccounts_Success(t *testing.T) {
	repo := mocks.NewRepository(t)
	acts := NewActivities(repo)

	req := TransferRequest{TransferID: "tx-1", FromAccountID: "acc-1", ToAccountID: "acc-2", Amount: 500}
	repo.EXPECT().GetAccount(mock.Anything, "acc-1").Return(openAccount("acc-1", "alice", 10000), nil)
	repo.EXPECT().GetAccount(mock.Anything, "acc-2").Return(openAccount("acc-2", "bob", 5000), nil)

	err := acts.ValidateAccounts(context.Background(), req)
	assert.NoError(t, err)
}

func TestValidateAccounts_SourceNotFound(t *testing.T) {
	repo := mocks.NewRepository(t)
	acts := NewActivities(repo)

	req := TransferRequest{TransferID: "tx-1", FromAccountID: "acc-1", ToAccountID: "acc-2", Amount: 500}
	repo.EXPECT().GetAccount(mock.Anything, "acc-1").Return(nil, domain.ErrAccountNotFound)

	err := acts.ValidateAccounts(context.Background(), req)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "source account not found")
}

func TestValidateAccounts_DestinationNotFound(t *testing.T) {
	repo := mocks.NewRepository(t)
	acts := NewActivities(repo)

	req := TransferRequest{TransferID: "tx-1", FromAccountID: "acc-1", ToAccountID: "acc-2", Amount: 500}
	repo.EXPECT().GetAccount(mock.Anything, "acc-1").Return(openAccount("acc-1", "alice", 10000), nil)
	repo.EXPECT().GetAccount(mock.Anything, "acc-2").Return(nil, domain.ErrAccountNotFound)

	err := acts.ValidateAccounts(context.Background(), req)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "destination account not found")
}

func TestValidateAccounts_SourceLocked(t *testing.T) {
	repo := mocks.NewRepository(t)
	acts := NewActivities(repo)

	req := TransferRequest{TransferID: "tx-1", FromAccountID: "acc-1", ToAccountID: "acc-2", Amount: 500}
	repo.EXPECT().GetAccount(mock.Anything, "acc-1").Return(lockedAccount("acc-1", "alice"), nil)
	repo.EXPECT().GetAccount(mock.Anything, "acc-2").Return(openAccount("acc-2", "bob", 5000), nil)

	err := acts.ValidateAccounts(context.Background(), req)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "source account cannot perform transaction")
}

func TestValidateAccounts_DestinationLocked(t *testing.T) {
	repo := mocks.NewRepository(t)
	acts := NewActivities(repo)

	req := TransferRequest{TransferID: "tx-1", FromAccountID: "acc-1", ToAccountID: "acc-2", Amount: 500}
	repo.EXPECT().GetAccount(mock.Anything, "acc-1").Return(openAccount("acc-1", "alice", 10000), nil)
	repo.EXPECT().GetAccount(mock.Anything, "acc-2").Return(lockedAccount("acc-2", "bob"), nil)

	err := acts.ValidateAccounts(context.Background(), req)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "destination account cannot perform transaction")
}

func TestValidateAccounts_SourceClosed(t *testing.T) {
	repo := mocks.NewRepository(t)
	acts := NewActivities(repo)

	closedAcc := &domain.Account{ID: "acc-1", Owner: "alice", Status: domain.StatusClosed}
	req := TransferRequest{TransferID: "tx-1", FromAccountID: "acc-1", ToAccountID: "acc-2", Amount: 500}
	repo.EXPECT().GetAccount(mock.Anything, "acc-1").Return(closedAcc, nil)
	repo.EXPECT().GetAccount(mock.Anything, "acc-2").Return(openAccount("acc-2", "bob", 5000), nil)

	err := acts.ValidateAccounts(context.Background(), req)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "source account cannot perform transaction")
}

// --- DebitAccount Tests ---

func TestDebitAccount_Success(t *testing.T) {
	repo := mocks.NewRepository(t)
	acts := NewActivities(repo)

	repo.EXPECT().ListTransactions(mock.Anything, "acc-1").Return(nil, nil)
	repo.EXPECT().GetAccount(mock.Anything, "acc-1").Return(openAccount("acc-1", "alice", 10000), nil)
	repo.EXPECT().SaveAccount(mock.Anything, mock.MatchedBy(func(a *domain.Account) bool {
		return a.ID == "acc-1" && a.Balance == 9500
	})).Return(nil)
	repo.EXPECT().SaveTransaction(mock.Anything, mock.MatchedBy(func(tx *domain.Transaction) bool {
		return tx.ID == "tx-1-debit" && tx.Amount == 500 && tx.Type == domain.TypeWithdrawal
	})).Return(nil)

	err := acts.DebitAccount(context.Background(), "acc-1", 500, "tx-1")
	assert.NoError(t, err)
}

func TestDebitAccount_Idempotent(t *testing.T) {
	repo := mocks.NewRepository(t)
	acts := NewActivities(repo)

	existing := []domain.Transaction{
		{ID: "tx-1-debit", AccountID: "acc-1", Amount: 500, Type: domain.TypeWithdrawal},
	}
	repo.EXPECT().ListTransactions(mock.Anything, "acc-1").Return(existing, nil)

	err := acts.DebitAccount(context.Background(), "acc-1", 500, "tx-1")
	assert.NoError(t, err)
}

func TestDebitAccount_InsufficientFunds(t *testing.T) {
	repo := mocks.NewRepository(t)
	acts := NewActivities(repo)

	repo.EXPECT().ListTransactions(mock.Anything, "acc-1").Return(nil, nil)
	repo.EXPECT().GetAccount(mock.Anything, "acc-1").Return(openAccount("acc-1", "alice", 100), nil)

	err := acts.DebitAccount(context.Background(), "acc-1", 500, "tx-1")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "insufficient funds")
}

func TestDebitAccount_AccountLocked(t *testing.T) {
	repo := mocks.NewRepository(t)
	acts := NewActivities(repo)

	repo.EXPECT().ListTransactions(mock.Anything, "acc-1").Return(nil, nil)
	repo.EXPECT().GetAccount(mock.Anything, "acc-1").Return(lockedAccount("acc-1", "alice"), nil)

	err := acts.DebitAccount(context.Background(), "acc-1", 500, "tx-1")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "source account cannot perform transaction")
}

// --- CreditAccount Tests ---

func TestCreditAccount_Success(t *testing.T) {
	repo := mocks.NewRepository(t)
	acts := NewActivities(repo)

	repo.EXPECT().ListTransactions(mock.Anything, "acc-2").Return(nil, nil)
	repo.EXPECT().GetAccount(mock.Anything, "acc-2").Return(openAccount("acc-2", "bob", 5000), nil)
	repo.EXPECT().SaveAccount(mock.Anything, mock.MatchedBy(func(a *domain.Account) bool {
		return a.ID == "acc-2" && a.Balance == 5500
	})).Return(nil)
	repo.EXPECT().SaveTransaction(mock.Anything, mock.MatchedBy(func(tx *domain.Transaction) bool {
		return tx.ID == "tx-1-credit" && tx.Amount == 500 && tx.Type == domain.TypeDeposit
	})).Return(nil)

	err := acts.CreditAccount(context.Background(), "acc-2", 500, "tx-1")
	assert.NoError(t, err)
}

func TestCreditAccount_Idempotent(t *testing.T) {
	repo := mocks.NewRepository(t)
	acts := NewActivities(repo)

	existing := []domain.Transaction{
		{ID: "tx-1-credit", AccountID: "acc-2", Amount: 500, Type: domain.TypeDeposit},
	}
	repo.EXPECT().ListTransactions(mock.Anything, "acc-2").Return(existing, nil)

	err := acts.CreditAccount(context.Background(), "acc-2", 500, "tx-1")
	assert.NoError(t, err)
}

func TestCreditAccount_AccountLocked(t *testing.T) {
	repo := mocks.NewRepository(t)
	acts := NewActivities(repo)

	repo.EXPECT().ListTransactions(mock.Anything, "acc-2").Return(nil, nil)
	repo.EXPECT().GetAccount(mock.Anything, "acc-2").Return(lockedAccount("acc-2", "bob"), nil)

	err := acts.CreditAccount(context.Background(), "acc-2", 500, "tx-1")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "account state invalid for credit")
}

// --- RefundDebitActivity Tests ---

func TestRefundDebitActivity_Success(t *testing.T) {
	repo := mocks.NewRepository(t)
	acts := NewActivities(repo)

	repo.EXPECT().ListTransactions(mock.Anything, "acc-1").Return(nil, nil)
	repo.EXPECT().GetAccount(mock.Anything, "acc-1").Return(openAccount("acc-1", "alice", 9500), nil)
	repo.EXPECT().SaveAccount(mock.Anything, mock.MatchedBy(func(a *domain.Account) bool {
		return a.ID == "acc-1" && a.Balance == 10000
	})).Return(nil)
	repo.EXPECT().SaveTransaction(mock.Anything, mock.MatchedBy(func(tx *domain.Transaction) bool {
		return tx.ID == "tx-1-refund" && tx.Amount == 500 && tx.Type == domain.TypeDeposit
	})).Return(nil)

	err := acts.RefundDebitActivity(context.Background(), "acc-1", 500, "tx-1")
	assert.NoError(t, err)
}

func TestRefundDebitActivity_Idempotent(t *testing.T) {
	repo := mocks.NewRepository(t)
	acts := NewActivities(repo)

	existing := []domain.Transaction{
		{ID: "tx-1-refund", AccountID: "acc-1", Amount: 500, Type: domain.TypeDeposit},
	}
	repo.EXPECT().ListTransactions(mock.Anything, "acc-1").Return(existing, nil)

	err := acts.RefundDebitActivity(context.Background(), "acc-1", 500, "tx-1")
	assert.NoError(t, err)
}
