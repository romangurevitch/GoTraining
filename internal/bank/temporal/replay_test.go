package temporal

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	commonpb "go.temporal.io/api/common/v1"
	enumspb "go.temporal.io/api/enums/v1"
	historypb "go.temporal.io/api/history/v1"
	taskqueuepb "go.temporal.io/api/taskqueue/v1"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/converter"
	"go.temporal.io/sdk/worker"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func encodePayloads(val ...any) *commonpb.Payloads {
	dc := converter.GetDefaultDataConverter()
	p, _ := dc.ToPayloads(val...)
	return p
}

func ts() *timestamppb.Timestamp {
	return timestamppb.New(time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC))
}

// Test_ReplayHistory verifies that the current workflow definition is
// backward-compatible with a previously recorded execution history.
func Test_ReplayHistory(t *testing.T) {
	req := TransferRequest{
		TransferID:    "tx-replay",
		FromAccountID: "acc-1",
		ToAccountID:   "acc-2",
		Amount:        500,
		Reference:     "replay test",
	}

	taskQueue := "bank-transfer-queue"

	history := &historypb.History{
		Events: []*historypb.HistoryEvent{
			{
				EventId:   1,
				EventTime: ts(),
				EventType: enumspb.EVENT_TYPE_WORKFLOW_EXECUTION_STARTED,
				Attributes: &historypb.HistoryEvent_WorkflowExecutionStartedEventAttributes{
					WorkflowExecutionStartedEventAttributes: &historypb.WorkflowExecutionStartedEventAttributes{
						WorkflowType:        &commonpb.WorkflowType{Name: "DurableTransferWorkflow"},
						TaskQueue:           &taskqueuepb.TaskQueue{Name: taskQueue},
						Input:               encodePayloads(req),
						WorkflowRunTimeout:  durationpb.New(time.Hour),
						WorkflowTaskTimeout: durationpb.New(10 * time.Second),
					},
				},
			},
			{
				EventId:   2,
				EventTime: ts(),
				EventType: enumspb.EVENT_TYPE_WORKFLOW_TASK_SCHEDULED,
				Attributes: &historypb.HistoryEvent_WorkflowTaskScheduledEventAttributes{
					WorkflowTaskScheduledEventAttributes: &historypb.WorkflowTaskScheduledEventAttributes{
						TaskQueue: &taskqueuepb.TaskQueue{Name: taskQueue},
					},
				},
			},
			{
				EventId:   3,
				EventTime: ts(),
				EventType: enumspb.EVENT_TYPE_WORKFLOW_TASK_STARTED,
				Attributes: &historypb.HistoryEvent_WorkflowTaskStartedEventAttributes{
					WorkflowTaskStartedEventAttributes: &historypb.WorkflowTaskStartedEventAttributes{
						ScheduledEventId: 2,
					},
				},
			},
			{
				EventId:   4,
				EventTime: ts(),
				EventType: enumspb.EVENT_TYPE_WORKFLOW_TASK_COMPLETED,
				Attributes: &historypb.HistoryEvent_WorkflowTaskCompletedEventAttributes{
					WorkflowTaskCompletedEventAttributes: &historypb.WorkflowTaskCompletedEventAttributes{
						ScheduledEventId: 2,
						StartedEventId:   3,
					},
				},
			},
			// Activity 1: ValidateAccounts (scheduled)
			{
				EventId:   5,
				EventTime: ts(),
				EventType: enumspb.EVENT_TYPE_ACTIVITY_TASK_SCHEDULED,
				Attributes: &historypb.HistoryEvent_ActivityTaskScheduledEventAttributes{
					ActivityTaskScheduledEventAttributes: &historypb.ActivityTaskScheduledEventAttributes{
						ActivityId:   "5",
						ActivityType: &commonpb.ActivityType{Name: "ValidateAccounts"},
						TaskQueue:    &taskqueuepb.TaskQueue{Name: taskQueue},
						Input:        encodePayloads(req),
					},
				},
			},
			{
				EventId:   6,
				EventTime: ts(),
				EventType: enumspb.EVENT_TYPE_ACTIVITY_TASK_STARTED,
				Attributes: &historypb.HistoryEvent_ActivityTaskStartedEventAttributes{
					ActivityTaskStartedEventAttributes: &historypb.ActivityTaskStartedEventAttributes{
						ScheduledEventId: 5,
					},
				},
			},
			{
				EventId:   7,
				EventTime: ts(),
				EventType: enumspb.EVENT_TYPE_ACTIVITY_TASK_COMPLETED,
				Attributes: &historypb.HistoryEvent_ActivityTaskCompletedEventAttributes{
					ActivityTaskCompletedEventAttributes: &historypb.ActivityTaskCompletedEventAttributes{
						ScheduledEventId: 5,
						StartedEventId:   6,
					},
				},
			},
			// Workflow task to process activity result
			{
				EventId:   8,
				EventTime: ts(),
				EventType: enumspb.EVENT_TYPE_WORKFLOW_TASK_SCHEDULED,
				Attributes: &historypb.HistoryEvent_WorkflowTaskScheduledEventAttributes{
					WorkflowTaskScheduledEventAttributes: &historypb.WorkflowTaskScheduledEventAttributes{
						TaskQueue: &taskqueuepb.TaskQueue{Name: taskQueue},
					},
				},
			},
			{
				EventId:   9,
				EventTime: ts(),
				EventType: enumspb.EVENT_TYPE_WORKFLOW_TASK_STARTED,
				Attributes: &historypb.HistoryEvent_WorkflowTaskStartedEventAttributes{
					WorkflowTaskStartedEventAttributes: &historypb.WorkflowTaskStartedEventAttributes{
						ScheduledEventId: 8,
					},
				},
			},
			{
				EventId:   10,
				EventTime: ts(),
				EventType: enumspb.EVENT_TYPE_WORKFLOW_TASK_COMPLETED,
				Attributes: &historypb.HistoryEvent_WorkflowTaskCompletedEventAttributes{
					WorkflowTaskCompletedEventAttributes: &historypb.WorkflowTaskCompletedEventAttributes{
						ScheduledEventId: 8,
						StartedEventId:   9,
					},
				},
			},
			// Activity 2: DebitAccount (scheduled)
			{
				EventId:   11,
				EventTime: ts(),
				EventType: enumspb.EVENT_TYPE_ACTIVITY_TASK_SCHEDULED,
				Attributes: &historypb.HistoryEvent_ActivityTaskScheduledEventAttributes{
					ActivityTaskScheduledEventAttributes: &historypb.ActivityTaskScheduledEventAttributes{
						ActivityId:   "11",
						ActivityType: &commonpb.ActivityType{Name: "DebitAccount"},
						TaskQueue:    &taskqueuepb.TaskQueue{Name: taskQueue},
						Input:        encodePayloads("acc-1", int64(500), "tx-replay"),
					},
				},
			},
			{
				EventId:   12,
				EventTime: ts(),
				EventType: enumspb.EVENT_TYPE_ACTIVITY_TASK_STARTED,
				Attributes: &historypb.HistoryEvent_ActivityTaskStartedEventAttributes{
					ActivityTaskStartedEventAttributes: &historypb.ActivityTaskStartedEventAttributes{
						ScheduledEventId: 11,
					},
				},
			},
			{
				EventId:   13,
				EventTime: ts(),
				EventType: enumspb.EVENT_TYPE_ACTIVITY_TASK_COMPLETED,
				Attributes: &historypb.HistoryEvent_ActivityTaskCompletedEventAttributes{
					ActivityTaskCompletedEventAttributes: &historypb.ActivityTaskCompletedEventAttributes{
						ScheduledEventId: 11,
						StartedEventId:   12,
					},
				},
			},
			// Workflow task to process DebitAccount result
			{
				EventId:   14,
				EventTime: ts(),
				EventType: enumspb.EVENT_TYPE_WORKFLOW_TASK_SCHEDULED,
				Attributes: &historypb.HistoryEvent_WorkflowTaskScheduledEventAttributes{
					WorkflowTaskScheduledEventAttributes: &historypb.WorkflowTaskScheduledEventAttributes{
						TaskQueue: &taskqueuepb.TaskQueue{Name: taskQueue},
					},
				},
			},
			{
				EventId:   15,
				EventTime: ts(),
				EventType: enumspb.EVENT_TYPE_WORKFLOW_TASK_STARTED,
				Attributes: &historypb.HistoryEvent_WorkflowTaskStartedEventAttributes{
					WorkflowTaskStartedEventAttributes: &historypb.WorkflowTaskStartedEventAttributes{
						ScheduledEventId: 14,
					},
				},
			},
			{
				EventId:   16,
				EventTime: ts(),
				EventType: enumspb.EVENT_TYPE_WORKFLOW_TASK_COMPLETED,
				Attributes: &historypb.HistoryEvent_WorkflowTaskCompletedEventAttributes{
					WorkflowTaskCompletedEventAttributes: &historypb.WorkflowTaskCompletedEventAttributes{
						ScheduledEventId: 14,
						StartedEventId:   15,
					},
				},
			},
			// Activity 3: CreditAccount (scheduled)
			{
				EventId:   17,
				EventTime: ts(),
				EventType: enumspb.EVENT_TYPE_ACTIVITY_TASK_SCHEDULED,
				Attributes: &historypb.HistoryEvent_ActivityTaskScheduledEventAttributes{
					ActivityTaskScheduledEventAttributes: &historypb.ActivityTaskScheduledEventAttributes{
						ActivityId:   "17",
						ActivityType: &commonpb.ActivityType{Name: "CreditAccount"},
						TaskQueue:    &taskqueuepb.TaskQueue{Name: taskQueue},
						Input:        encodePayloads("acc-2", int64(500), "tx-replay"),
					},
				},
			},
			{
				EventId:   18,
				EventTime: ts(),
				EventType: enumspb.EVENT_TYPE_ACTIVITY_TASK_STARTED,
				Attributes: &historypb.HistoryEvent_ActivityTaskStartedEventAttributes{
					ActivityTaskStartedEventAttributes: &historypb.ActivityTaskStartedEventAttributes{
						ScheduledEventId: 17,
					},
				},
			},
			{
				EventId:   19,
				EventTime: ts(),
				EventType: enumspb.EVENT_TYPE_ACTIVITY_TASK_COMPLETED,
				Attributes: &historypb.HistoryEvent_ActivityTaskCompletedEventAttributes{
					ActivityTaskCompletedEventAttributes: &historypb.ActivityTaskCompletedEventAttributes{
						ScheduledEventId: 17,
						StartedEventId:   18,
					},
				},
			},
			// Workflow task to complete
			{
				EventId:   20,
				EventTime: ts(),
				EventType: enumspb.EVENT_TYPE_WORKFLOW_TASK_SCHEDULED,
				Attributes: &historypb.HistoryEvent_WorkflowTaskScheduledEventAttributes{
					WorkflowTaskScheduledEventAttributes: &historypb.WorkflowTaskScheduledEventAttributes{
						TaskQueue: &taskqueuepb.TaskQueue{Name: taskQueue},
					},
				},
			},
			{
				EventId:   21,
				EventTime: ts(),
				EventType: enumspb.EVENT_TYPE_WORKFLOW_TASK_STARTED,
				Attributes: &historypb.HistoryEvent_WorkflowTaskStartedEventAttributes{
					WorkflowTaskStartedEventAttributes: &historypb.WorkflowTaskStartedEventAttributes{
						ScheduledEventId: 20,
					},
				},
			},
			{
				EventId:   22,
				EventTime: ts(),
				EventType: enumspb.EVENT_TYPE_WORKFLOW_TASK_COMPLETED,
				Attributes: &historypb.HistoryEvent_WorkflowTaskCompletedEventAttributes{
					WorkflowTaskCompletedEventAttributes: &historypb.WorkflowTaskCompletedEventAttributes{
						ScheduledEventId: 20,
						StartedEventId:   21,
					},
				},
			},
			{
				EventId:   23,
				EventTime: ts(),
				EventType: enumspb.EVENT_TYPE_WORKFLOW_EXECUTION_COMPLETED,
				Attributes: &historypb.HistoryEvent_WorkflowExecutionCompletedEventAttributes{
					WorkflowExecutionCompletedEventAttributes: &historypb.WorkflowExecutionCompletedEventAttributes{
						WorkflowTaskCompletedEventId: 22,
						Result:                       encodePayloads(TransferResponse{TransferID: "tx-replay", Status: "COMPLETED"}),
					},
				},
			},
		},
	}

	replayer := worker.NewWorkflowReplayer()
	replayer.RegisterWorkflow(DurableTransferWorkflow)

	err := replayer.ReplayWorkflowHistory(nil, history)
	require.NoError(t, err)
}

// Test_ReplayHistoryFromFile replays workflow history from an exported JSON file.
// To generate the fixture from a running cluster:
//
//	temporal workflow show --workflow-id <id> --output json > internal/bank/temporal/testdata/transfer_history.json
func Test_ReplayHistoryFromFile(t *testing.T) {
	f, err := os.Open("testdata/transfer_history.json")
	if os.IsNotExist(err) {
		t.Skip("No history fixture — export one from a running Temporal cluster")
	}
	require.NoError(t, err)
	defer f.Close()

	history, err := client.HistoryFromJSON(f, client.HistoryJSONOptions{})
	require.NoError(t, err)

	replayer := worker.NewWorkflowReplayer()
	replayer.RegisterWorkflow(DurableTransferWorkflow)

	err = replayer.ReplayWorkflowHistory(nil, history)
	require.NoError(t, err)
}
