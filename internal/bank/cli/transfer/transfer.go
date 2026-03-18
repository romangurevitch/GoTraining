package transfer

import (
	"context"
	"fmt"
	"strconv"

	"github.com/romangurevitch/go-training/internal/pkg/json"
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

	return cmd
}

func getCreateTransferCmd(bankClient bank.Client) *cobra.Command {
	return &cobra.Command{
		Use:   "create [from-id] [to-id] [amount]",
		Short: "Create a new transfer",
		Args:  cobra.ExactArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			amount, err := strconv.ParseInt(args[2], 10, 64)
			cobra.CheckErr(err)

			req := &api.CreateTransferRequest{
				FromAccountID: args[0],
				ToAccountID:   args[1],
				Amount:        amount,
			}

			res, err := bankClient.Transfer(context.Background(), req)
			cobra.CheckErr(err)
			fmt.Println(json.ToJSONString(res))
		},
	}
}
