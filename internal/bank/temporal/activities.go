package temporal

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/romangurevitch/go-training/internal/bank/domain"
	"github.com/romangurevitch/go-training/internal/bank/repository"
	"go.temporal.io/sdk/temporal"
)

type Activities struct {
	repo repository.Repository
}

func NewActivities(repo repository.Repository) *Activities {
	return &Activities{repo: repo}
}

func (a *Activities) ValidateAccounts(ctx context.Context, req TransferRequest) error {
	from, err := a.repo.GetAccount(ctx, req.FromAccountID)
	if err != nil {
		if errors.Is(err, domain.ErrAccountNotFound) {
			return temporal.NewNonRetryableApplicationError("source account not found", "ACCOUNT_NOT_FOUND", err)
		}
		return err
	}

	to, err := a.repo.GetAccount(ctx, req.ToAccountID)
	if err != nil {
		if errors.Is(err, domain.ErrAccountNotFound) {
			return temporal.NewNonRetryableApplicationError("destination account not found", "ACCOUNT_NOT_FOUND", err)
		}
		return err
	}

	if err := from.CanPerformTransaction(); err != nil {
		return temporal.NewNonRetryableApplicationError("source account cannot perform transaction", "ACCOUNT_LOCKED", err)
	}

	if err := to.CanPerformTransaction(); err != nil {
		return temporal.NewNonRetryableApplicationError("destination account cannot perform transaction", "ACCOUNT_LOCKED", err)
	}

	return nil
}

func (a *Activities) DebitAccount(ctx context.Context, accountID string, amount int64, transferID string) error {
	transactionID := fmt.Sprintf("%s-debit", transferID)

	// Check if already processed (Idempotency)
	txs, err := a.repo.ListTransactions(ctx, accountID)
	if err == nil {
		for _, tx := range txs {
			if tx.ID == transactionID {
				return nil // Already processed
			}
		}
	}

	acc, err := a.repo.GetAccount(ctx, accountID)
	if err != nil {
		return err
	}

	if err := acc.CanPerformTransaction(); err != nil {
		return temporal.NewNonRetryableApplicationError("source account cannot perform transaction", "ACCOUNT_STATE", err)
	}

	if acc.Balance < amount {
		return temporal.NewNonRetryableApplicationError("insufficient funds", "INSUFFICIENT_FUNDS", domain.ErrInsufficientFunds)
	}

	acc.Balance -= amount
	acc.UpdatedAt = time.Now()

	if err := a.repo.SaveAccount(ctx, acc); err != nil {
		return err
	}

	tx := &domain.Transaction{
		ID:        transactionID,
		AccountID: accountID,
		Amount:    amount,
		Type:      domain.TypeWithdrawal,
		CreatedAt: time.Now(),
	}

	return a.repo.SaveTransaction(ctx, tx)
}

func (a *Activities) CreditAccount(ctx context.Context, accountID string, amount int64, transferID string) error {
	transactionID := fmt.Sprintf("%s-credit", transferID)

	// Check if already processed (Idempotency)
	txs, err := a.repo.ListTransactions(ctx, accountID)
	if err == nil {
		for _, tx := range txs {
			if tx.ID == transactionID {
				return nil // Already processed
			}
		}
	}

	acc, err := a.repo.GetAccount(ctx, accountID)
	if err != nil {
		return err
	}

	if err := acc.CanPerformTransaction(); err != nil {
		return temporal.NewNonRetryableApplicationError("account state invalid for credit", "INVALID_STATE", err)
	}

	acc.Balance += amount
	acc.UpdatedAt = time.Now()

	if err := a.repo.SaveAccount(ctx, acc); err != nil {
		return err
	}

	tx := &domain.Transaction{
		ID:        transactionID,
		AccountID: accountID,
		Amount:    amount,
		Type:      domain.TypeDeposit,
		CreatedAt: time.Now(),
	}

	return a.repo.SaveTransaction(ctx, tx)
}

func (a *Activities) RefundDebitActivity(ctx context.Context, accountID string, amount int64, transferID string) error {
	transactionID := fmt.Sprintf("%s-refund", transferID)

	// Check if already processed (Idempotency)
	txs, err := a.repo.ListTransactions(ctx, accountID)
	if err == nil {
		for _, tx := range txs {
			if tx.ID == transactionID {
				return nil // Already processed
			}
		}
	}

	acc, err := a.repo.GetAccount(ctx, accountID)
	if err != nil {
		return err
	}

	acc.Balance += amount
	acc.UpdatedAt = time.Now()

	if err := a.repo.SaveAccount(ctx, acc); err != nil {
		return err
	}

	tx := &domain.Transaction{
		ID:        transactionID,
		AccountID: accountID,
		Amount:    amount,
		Type:      domain.TypeDeposit,
		CreatedAt: time.Now(),
	}

	return a.repo.SaveTransaction(ctx, tx)
}
