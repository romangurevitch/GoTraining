package temporal

import (
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/testsuite"
)

func TestWorkflowTestSuite(t *testing.T) {
	suite.Run(t, new(WorkflowTestSuite))
}

type WorkflowTestSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite

	env *testsuite.TestWorkflowEnvironment

	activities *Activities
}

func (s *WorkflowTestSuite) SetupTest() {
	s.env = s.NewTestWorkflowEnvironment()
	s.activities = &Activities{}
}

func (s *WorkflowTestSuite) AfterTest(suiteName, testName string) {
	s.env.AssertExpectations(s.T())
}

func (s *WorkflowTestSuite) Test_DurableTransferWorkflow_Success_SmallAmount() {
	req := TransferRequest{
		TransferID:    "tx-123",
		FromAccountID: "acc-1",
		ToAccountID:   "acc-2",
		Amount:        500,
		Reference:     "Small transfer",
	}

	s.env.OnActivity(s.activities.ValidateAccounts, mock.Anything, req).Return(nil)
	s.env.OnActivity(s.activities.DebitAccount, mock.Anything, "acc-1", int64(500), "tx-123").Return(nil)
	s.env.OnActivity(s.activities.CreditAccount, mock.Anything, "acc-2", int64(500), "tx-123").Return(nil)

	s.env.ExecuteWorkflow(DurableTransferWorkflow, req)

	s.True(s.env.IsWorkflowCompleted())
	s.NoError(s.env.GetWorkflowError())

	var res TransferResponse
	err := s.env.GetWorkflowResult(&res)
	s.NoError(err)
	s.Equal("tx-123", res.TransferID)
	s.Equal("COMPLETED", res.Status)
}

func (s *WorkflowTestSuite) Test_DurableTransferWorkflow_Success_LargeAmount_Approved() {
	req := TransferRequest{
		TransferID:    "tx-large",
		FromAccountID: "acc-1",
		ToAccountID:   "acc-2",
		Amount:        2000,
		Reference:     "Large transfer",
	}

	s.env.OnActivity(s.activities.ValidateAccounts, mock.Anything, req).Return(nil)

	s.env.RegisterDelayedCallback(func() {
		s.env.SignalWorkflow(ApprovalSignal, nil)
	}, time.Hour)

	s.env.OnActivity(s.activities.DebitAccount, mock.Anything, "acc-1", int64(2000), "tx-large").Return(nil)
	s.env.OnActivity(s.activities.CreditAccount, mock.Anything, "acc-2", int64(2000), "tx-large").Return(nil)

	s.env.ExecuteWorkflow(DurableTransferWorkflow, req)

	s.True(s.env.IsWorkflowCompleted())
	s.NoError(s.env.GetWorkflowError())

	var res TransferResponse
	err := s.env.GetWorkflowResult(&res)
	s.NoError(err)
	s.Equal("COMPLETED", res.Status)
}

func (s *WorkflowTestSuite) Test_DurableTransferWorkflow_Rejected() {
	req := TransferRequest{
		TransferID:    "tx-reject",
		FromAccountID: "acc-1",
		ToAccountID:   "acc-2",
		Amount:        2000,
		Reference:     "Reject me",
	}

	s.env.OnActivity(s.activities.ValidateAccounts, mock.Anything, req).Return(nil)

	s.env.RegisterDelayedCallback(func() {
		s.env.SignalWorkflow(RejectSignal, nil)
	}, time.Hour)

	s.env.ExecuteWorkflow(DurableTransferWorkflow, req)

	s.True(s.env.IsWorkflowCompleted())
	s.Error(s.env.GetWorkflowError())
	s.Contains(s.env.GetWorkflowError().Error(), "transfer rejected")
}

func (s *WorkflowTestSuite) Test_DurableTransferWorkflow_Compensation() {
	req := TransferRequest{
		TransferID:    "tx-compensate",
		FromAccountID: "acc-1",
		ToAccountID:   "acc-2",
		Amount:        500,
		Reference:     "Failure test",
	}

	s.env.OnActivity(s.activities.ValidateAccounts, mock.Anything, req).Return(nil)
	s.env.OnActivity(s.activities.DebitAccount, mock.Anything, "acc-1", int64(500), "tx-compensate").Return(nil)
	s.env.OnActivity(s.activities.CreditAccount, mock.Anything, "acc-2", int64(500), "tx-compensate").Return(
		temporal.NewNonRetryableApplicationError("destination account closed", "ACCOUNT_CLOSED", nil),
	)
	s.env.OnActivity(s.activities.RefundDebitActivity, mock.Anything, "acc-1", int64(500), "tx-compensate").Return(nil)

	s.env.ExecuteWorkflow(DurableTransferWorkflow, req)

	s.True(s.env.IsWorkflowCompleted())
	s.Error(s.env.GetWorkflowError())
	s.Contains(s.env.GetWorkflowError().Error(), "destination account closed")
}

func (s *WorkflowTestSuite) Test_DurableTransferWorkflow_ApprovalTimeout() {
	req := TransferRequest{
		TransferID:    "tx-timeout",
		FromAccountID: "acc-1",
		ToAccountID:   "acc-2",
		Amount:        2000,
	}

	s.env.OnActivity(s.activities.ValidateAccounts, mock.Anything, req).Return(nil)

	// No signal sent.

	s.env.ExecuteWorkflow(DurableTransferWorkflow, req)

	s.True(s.env.IsWorkflowCompleted())
	s.Error(s.env.GetWorkflowError())
	s.Contains(s.env.GetWorkflowError().Error(), "approval timed out")
}
