package account

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/romangurevitch/go-training/pkg/client/bank"
	"github.com/spf13/cobra"
)

func getBalanceCmd(bankClient bank.Client) *cobra.Command {
	return &cobra.Command{
		Use:   "balance [account-id]",
		Short: "Check the balance of an account",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			accountID := args[0]
			ctx := context.Background()

			res, err := bankClient.GetAccount(ctx, accountID)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}

			b, _ := json.MarshalIndent(res, "", "  ")
			fmt.Println(string(b))
		},
	}
}
