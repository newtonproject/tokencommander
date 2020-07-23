package cli

import (
	"fmt"
	"github.com/spf13/viper"

	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
)

func (cli *CLI) buildAddCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "add <contractAddress> [symbol]",
		Short:                 "Add custom contract",
		Args:                  cobra.MinimumNArgs(1),
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {

			if !common.IsHexAddress(args[0]) {
				fmt.Println("Contract address invalid")
				return
			}
			contractAddress := common.HexToAddress(args[0])
			cli.contractAddress = contractAddress.String()

			var symbol string
			if len(args) > 1 {
				symbol = args[1]
			}

			if symbol == "" {
				SimpleToken, err := cli.GetSimpleToken()
				if err != nil {
					fmt.Println("GetSimpleToken Error: ", err)
					fmt.Println(cmd.UsageString())
					return
				}
				symbol, err = SimpleToken.Symbol(nil)
				if err != nil {
					fmt.Println("Get Symbol error, please enter the symbol")
					return
				}
			}

			// re-check
			if symbol == "" {
				fmt.Println("Symbol is empty")
				return
			}

			addressStr := viper.GetString(fmt.Sprintf("Contracts.%s", symbol))
			if addressStr != "" {
				if common.HexToAddress(addressStr) != contractAddress {
					fmt.Printf("This symbol %s has been used by the contract address %s, please choose a new symbol\n", symbol, addressStr)
					return
				}
				fmt.Println("This contract address has been added before.")
				return
			}
			viper.Set(fmt.Sprintf("Contracts.%s", symbol), contractAddress.String())
			err := viper.WriteConfigAs(cli.config)
			if err != nil {
				fmt.Println("WriteConfig:", err)
				return
			}

			fmt.Printf("The contract %s(%s) has been added.\n", symbol, contractAddress.String())

			cli.buildInfoCmd().Run(cmd, args)

			return
		},
	}

	return cmd
}
