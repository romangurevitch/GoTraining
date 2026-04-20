package transfer

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/romangurevitch/go-training/pkg/client/bank"
	"github.com/romangurevitch/go-training/pkg/client/bank/api"
	"github.com/spf13/cobra"
)

// GetTransferCmd returns the 'transfer' command group.
func GetTransferCmd(bankClient bank.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "transfer",
		Short: "Transfer funds between accounts",
	}

	cmd.AddCommand(getCreateTransferCmd(bankClient))
	cmd.AddCommand(getApproveTransferCmd(bankClient))
	cmd.AddCommand(getRejectTransferCmd(bankClient))

	return cmd
}

func getCreateTransferCmd(bankClient bank.Client) *cobra.Command {
	var reference string
	var durable bool

	cmd := &cobra.Command{
		Use:   "create [from-id] [to-id] [amount]",
		Short: "Create a new transfer",
		Args:  cobra.ExactArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			fromID := args[0]
			toID := args[1]
			amount, err := strconv.ParseInt(args[2], 10, 64)
			if err != nil {
				fmt.Printf("Error: invalid amount: %v\n", err)
				os.Exit(1)
			}

			ctx := context.Background()

			if durable {
				req := &api.DurableTransferRequest{
					FromAccountID: fromID,
					ToAccountID:   toID,
					Amount:        amount,
					Reference:     reference,
				}
				res, err := bankClient.StartDurableTransfer(ctx, req)
				if err != nil {
					fmt.Printf("Error: %v\n", err)
					os.Exit(1)
				}
				printJSON(res)
			} else {
				req := &api.CreateTransferRequest{
					FromAccountID: fromID,
					ToAccountID:   toID,
					Amount:        amount,
				}
				res, err := bankClient.Transfer(ctx, req)
				if err != nil {
					fmt.Printf("Error: %v\n", err)
					os.Exit(1)
				}
				printJSON(res)
			}
		},
	}

	cmd.Flags().StringVarP(&reference, "reference", "r", "", "Transfer reference")
	cmd.Flags().BoolVarP(&durable, "durable", "d", false, "Use durable transfer workflow")

	return cmd
}

func getApproveTransferCmd(bankClient bank.Client) *cobra.Command {
	return &cobra.Command{
		Use:   "approve [transfer-id]",
		Short: "Approve a pending durable transfer",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			transferID := args[0]
			ctx := context.Background()
			res, err := bankClient.ApproveTransfer(ctx, transferID)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}
			printJSON(res)
		},
	}
}

func getRejectTransferCmd(bankClient bank.Client) *cobra.Command {
	return &cobra.Command{
		Use:   "reject [transfer-id]",
		Short: "Reject a pending durable transfer",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			transferID := args[0]
			ctx := context.Background()
			res, err := bankClient.RejectTransfer(ctx, transferID)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}
			printJSON(res)
		},
	}
}

func printJSON(v any) {
	b, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(b))
}
