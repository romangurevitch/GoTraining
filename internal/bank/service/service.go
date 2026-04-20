package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/romangurevitch/go-training/internal/bank/config"
	"github.com/romangurevitch/go-training/internal/bank/domain"
	"github.com/romangurevitch/go-training/internal/bank/repository"
	"github.com/romangurevitch/go-training/internal/bank/temporal"
	"go.temporal.io/sdk/client"
)

// Service is the business logic interface. Enables mock injection in handler tests.
type Service interface {
	CreateAccount(ctx context.Context, owner string) (*domain.Account, error)
	GetAccount(ctx context.Context, id string) (*domain.Account, error)
	Deposit(ctx context.Context, accountID string, amount int64) error
	Withdraw(ctx context.Context, accountID string, amount int64) error
	Transfer(ctx context.Context, fromID, toID string, amount int64) error
	StartDurableTransfer(ctx context.Context, fromID, toID string, amount int64, reference string) (string, error)
	ApproveTransfer(ctx context.Context, transferID string) error
	RejectTransfer(ctx context.Context, transferID string) error
}

// BankService implements Service backed by a Repository.
type BankService struct {
	repo           repository.Repository
	temporalClient client.Client
}

// Ensure BankService implements Service at compile time.
var _ Service = (*BankService)(nil)

func NewBankService(repo repository.Repository, temporalClient client.Client) *BankService {
	return &BankService{
		repo:           repo,
		temporalClient: temporalClient,
	}
}

func (s *BankService) CreateAccount(ctx context.Context, owner string) (*domain.Account, error) {
	acc := &domain.Account{
		ID:        fmt.Sprintf("ACC-%d", time.Now().UnixNano()),
		Owner:     owner,
		Balance:   0,
		Status:    domain.StatusOpen,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.repo.SaveAccount(ctx, acc); err != nil {
		return nil, fmt.Errorf("failed to save account: %w", err)
	}

	slog.InfoContext(ctx, "account created", slog.String("id", acc.ID), slog.String("owner", acc.Owner))
	return acc, nil
}

func (s *BankService) GetAccount(ctx context.Context, id string) (*domain.Account, error) {
	return s.repo.GetAccount(ctx, id)
}

func (s *BankService) Deposit(ctx context.Context, accountID string, amount int64) error {
	if amount <= 0 {
		return domain.ErrInvalidAmount
	}

	acc, err := s.repo.GetAccount(ctx, accountID)
	if err != nil {
		return err
	}

	if err := acc.CanPerformTransaction(); err != nil {
		return err
	}

	acc.Balance += amount
	acc.UpdatedAt = time.Now()

	if err := s.repo.SaveAccount(ctx, acc); err != nil {
		return err
	}

	t := &domain.Transaction{
		ID:        fmt.Sprintf("TRX-%d", time.Now().UnixNano()),
		AccountID: accountID,
		Amount:    amount,
		Type:      domain.TypeDeposit,
		CreatedAt: time.Now(),
	}

	return s.repo.SaveTransaction(ctx, t)
}

func (s *BankService) Withdraw(ctx context.Context, accountID string, amount int64) error {
	if amount <= 0 {
		return domain.ErrInvalidAmount
	}

	acc, err := s.repo.GetAccount(ctx, accountID)
	if err != nil {
		return err
	}

	if err := acc.CanPerformTransaction(); err != nil {
		return err
	}

	if acc.Balance < amount {
		return domain.ErrInsufficientFunds
	}

	acc.Balance -= amount
	acc.UpdatedAt = time.Now()

	if err := s.repo.SaveAccount(ctx, acc); err != nil {
		return err
	}

	t := &domain.Transaction{
		ID:        fmt.Sprintf("TRX-%d", time.Now().UnixNano()),
		AccountID: accountID,
		Amount:    amount,
		Type:      domain.TypeWithdrawal,
		CreatedAt: time.Now(),
	}

	return s.repo.SaveTransaction(ctx, t)
}

// Transfer moves funds from one account to another.
// Pre-built for participants — they call this from the transfer handler.
func (s *BankService) Transfer(ctx context.Context, fromID, toID string, amount int64) error {
	if amount <= 0 {
		return domain.ErrInvalidAmount
	}

	from, err := s.repo.GetAccount(ctx, fromID)
	if err != nil {
		return err
	}

	to, err := s.repo.GetAccount(ctx, toID)
	if err != nil {
		return err
	}

	if err := from.CanPerformTransaction(); err != nil {
		return err
	}

	if err := to.CanPerformTransaction(); err != nil {
		return err
	}

	if from.Balance < amount {
		return domain.ErrInsufficientFunds
	}

	from.Balance -= amount
	from.UpdatedAt = time.Now()

	to.Balance += amount
	to.UpdatedAt = time.Now()

	if err := s.repo.SaveAccount(ctx, from); err != nil {
		return fmt.Errorf("failed to debit source account: %w", err)
	}

	if err := s.repo.SaveAccount(ctx, to); err != nil {
		return fmt.Errorf("failed to credit destination account: %w", err)
	}

	debit := &domain.Transaction{
		ID:        fmt.Sprintf("TRX-%d-D", time.Now().UnixNano()),
		AccountID: fromID,
		Amount:    amount,
		Type:      domain.TypeWithdrawal,
		CreatedAt: time.Now(),
	}
	if err := s.repo.SaveTransaction(ctx, debit); err != nil {
		return err
	}

	credit := &domain.Transaction{
		ID:        fmt.Sprintf("TRX-%d-C", time.Now().UnixNano()),
		AccountID: toID,
		Amount:    amount,
		Type:      domain.TypeDeposit,
		CreatedAt: time.Now(),
	}

	slog.InfoContext(ctx, "transfer completed",
		slog.String("from_account_id", fromID),
		slog.String("to_account_id", toID),
		slog.Int64("amount", amount),
	)

	return s.repo.SaveTransaction(ctx, credit)
}

func (s *BankService) StartDurableTransfer(ctx context.Context, fromID, toID string, amount int64, reference string) (string, error) {
	// 1. Validation before starting workflow
	from, err := s.repo.GetAccount(ctx, fromID)
	if err != nil {
		return "", err
	}
	to, err := s.repo.GetAccount(ctx, toID)
	if err != nil {
		return "", err
	}
	if err := from.CanPerformTransaction(); err != nil {
		return "", err
	}
	if err := to.CanPerformTransaction(); err != nil {
		return "", err
	}
	if from.Balance < amount {
		return "", domain.ErrInsufficientFunds
	}

	// 2. Start Temporal Workflow
	transferID := fmt.Sprintf("TRX-%d", time.Now().UnixNano())
	workflowID := "transfer-" + transferID

	options := client.StartWorkflowOptions{
		ID:        workflowID,
		TaskQueue: config.Values.TemporalTaskQueue,
	}

	req := temporal.TransferRequest{
		TransferID:    transferID,
		FromAccountID: fromID,
		ToAccountID:   toID,
		Amount:        amount,
		Reference:     reference,
	}

	we, err := s.temporalClient.ExecuteWorkflow(ctx, options, temporal.DurableTransferWorkflow, req)
	if err != nil {
		return "", fmt.Errorf("failed to start durable transfer workflow: %w", err)
	}

	slog.InfoContext(ctx, "durable transfer workflow started",
		slog.String("workflow_id", we.GetID()),
		slog.String("run_id", we.GetRunID()),
		slog.String("transfer_id", transferID),
	)

	return transferID, nil
}

func (s *BankService) ApproveTransfer(ctx context.Context, transferID string) error {
	workflowID := "transfer-" + transferID
	return s.temporalClient.SignalWorkflow(ctx, workflowID, "", temporal.ApprovalSignal, nil)
}

func (s *BankService) RejectTransfer(ctx context.Context, transferID string) error {
	workflowID := "transfer-" + transferID
	return s.temporalClient.SignalWorkflow(ctx, workflowID, "", temporal.RejectSignal, nil)
}
