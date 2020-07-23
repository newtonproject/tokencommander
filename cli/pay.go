package cli

import (
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func (cli *CLI) buildPayCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "pay <amount|tokenID|all> <--to toAddress> [--from fromAddress]",
		Aliases: []string{"transfer"},
		Short:   "Command about transaction",
		Args:    cobra.MinimumNArgs(1),
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {

			amountStr := args[0]

			fromAddressStr := viper.GetString("from")
			if fromAddressStr == "" || !common.IsHexAddress(fromAddressStr) {
				fmt.Println("Error: not set from address of owner or from address illegal")
				fmt.Fprint(os.Stderr, cmd.UsageString())
				return
			}
			fromAddress := common.HexToAddress(fromAddressStr)

			toAddressStr, err := cmd.Flags().GetString("to")
			if err != nil {
				fmt.Println("Error: required flag(s) \"to\" not set")
				fmt.Fprint(os.Stderr, cmd.UsageString())
				return
			}
			if !common.IsHexAddress(toAddressStr) {
				fmt.Println("Error: illegal to address", toAddressStr)
				return
			}
			toAddress := common.HexToAddress(toAddressStr)

			nowait, _ := cmd.Flags().GetBool("nowait")
			cli.pay(fromAddress, toAddress, amountStr, nowait)

			return
		},
	}

	cmd.Flags().StringP("to", "t", "", "the address pay to")
	cmd.MarkFlagRequired("to")
	cmd.Flags().Bool("nowait", false, "do not wait for tx to be mined")

	return cmd
}
