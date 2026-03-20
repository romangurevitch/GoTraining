package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/google/uuid"
	clientconfig "github.com/romangurevitch/go-training/cmd/temporal/client/config"
	"github.com/romangurevitch/go-training/internal/temporal/encryption"
	"github.com/romangurevitch/go-training/internal/temporal/order"
	"github.com/romangurevitch/go-training/internal/temporal/workflows"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/converter"
	"go.temporal.io/sdk/workflow"
)

// Config holds the CLI configuration
type Config struct {
	ConfigPath   string
	OrderPayload string
	WorkflowName string
	SignalName   string
	WorkflowID   string
	Wait         bool
}

func main() {
	// Initialize structured logging to stderr
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})))

	cfg := Config{}
	flag.StringVar(&cfg.ConfigPath, "config", "./config/temporal/client/local/config.yaml", "path to config file")
	flag.StringVar(&cfg.OrderPayload, "order", "", "json order payload (required to start workflow)")
	flag.StringVar(&cfg.WorkflowName, "workflow", "ProcessOrder", "workflow to run: ProcessOrder or AutoProcessOrder")
	flag.StringVar(&cfg.SignalName, "signal", "", "signal to send (e.g., pickOrder, shipOrder, markOrderAsDelivered, cancelOrder)")
	flag.StringVar(&cfg.WorkflowID, "workflow-id", "", "ID of the workflow to signal or start")
	flag.BoolVar(&cfg.Wait, "wait", true, "wait for workflow to complete")
	flag.Parse()

	// Create a cancellable context that responds to OS signals
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := run(ctx, cfg); err != nil {
		slog.Error("Application failed", "error", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, cfg Config) error {
	appCfg, err := clientconfig.LoadConfig(cfg.ConfigPath)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	// Connect to Temporal
	c, err := client.Dial(client.Options{
		HostPort: fmt.Sprintf("%s:%d", appCfg.Temporal.Host, appCfg.Temporal.Port),
		Logger:   slog.Default(),
		DataConverter: encryption.NewEncryptionDataConverter(
			converter.GetDefaultDataConverter(),
			encryption.DataConverterOptions{Compress: true},
		),
	})
	if err != nil {
		return fmt.Errorf("dial temporal: %w", err)
	}
	defer c.Close()

	// If signal is provided, handle signaling and return
	if cfg.SignalName != "" {
		return handleSignal(ctx, c, cfg.WorkflowID, cfg.SignalName)
	}

	// Otherwise, handle starting a new workflow
	return handleStart(ctx, c, appCfg.Temporal.TaskQueueName, cfg)
}

func handleSignal(ctx context.Context, c client.Client, workflowID, signalName string) error {
	if workflowID == "" {
		return fmt.Errorf("workflow-id is required when signaling")
	}

	// Use a shorter timeout for signaling, but still respect the parent context
	signalCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	slog.Info("Sending signal", "signal", signalName, "workflow_id", workflowID)
	err := c.SignalWorkflow(signalCtx, workflowID, "", signalName, nil)
	if err != nil {
		return fmt.Errorf("signal workflow %s: %w", workflowID, err)
	}

	slog.Info("Signal sent successfully")
	return nil
}

func handleStart(ctx context.Context, c client.Client, taskQueue string, cfg Config) error {
	if cfg.OrderPayload == "" {
		return fmt.Errorf("json order payload is required to start a workflow")
	}

	workflowID := cfg.WorkflowID
	if workflowID == "" {
		workflowID = "order-" + uuid.New().String()
	}

	var o order.Order
	if err := json.Unmarshal([]byte(cfg.OrderPayload), &o); err != nil {
		return fmt.Errorf("unmarshal order: %w", err)
	}

	// Use a concrete function type instead of interface{} for better type safety
	var wfFunc func(workflow.Context, workflows.Params) (order.OrderStatus, error)

	switch cfg.WorkflowName {
	case "AutoProcessOrder":
		wfFunc = workflows.AutoProcessOrder
	case "ProcessOrder":
		wfFunc = workflows.ProcessOrder
	default:
		return fmt.Errorf("unknown workflow name: %s", cfg.WorkflowName)
	}

	options := client.StartWorkflowOptions{
		ID:        workflowID,
		TaskQueue: taskQueue,
	}

	slog.Info("Starting workflow", "workflow", cfg.WorkflowName, "workflow_id", workflowID)

	// Temporal ExecuteWorkflow accepts the function itself
	workflowRun, err := c.ExecuteWorkflow(ctx, options, wfFunc, workflows.Params{Order: o})
	if err != nil {
		return fmt.Errorf("execute workflow: %w", err)
	}

	if cfg.Wait {
		if cfg.WorkflowName == "ProcessOrder" {
			fmt.Printf("\n⏳ Waiting for signal-driven workflow: %s\n", workflowID)
			fmt.Printf("This workflow will pause at each stage until it receives a signal.\n")
			fmt.Printf("Keep this terminal open to see the final result, but use a NEW terminal to send signals.\n")
			fmt.Printf("\nIn your signaling terminal, run:\nexport ID=%s\n\n", workflowID)
		}

		var result order.OrderStatus
		if err := workflowRun.Get(ctx, &result); err != nil {
			return fmt.Errorf("workflow execution failed: %w", err)
		}
		slog.Info("Workflow completed", "status", result)
	} else {
		fmt.Printf("\n🚀 Workflow started: %s\n", workflowID)
		fmt.Printf("To signal this workflow using 'make', run:\nexport ID=%s\n\n", workflowID)
	}

	return nil
}
