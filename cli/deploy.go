package cli

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func (cli *CLI) buildDeployCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deploy <--name tokenname> <--symbol tokensymbol> <--total totalSupplyAmount> [--decimals decimal]",
		Short: fmt.Sprintf("Deploy %s contract", cli.blockchain.String()),
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {

			save, _ := cmd.Flags().GetBool("save")

			fromAddress := viper.GetString("from")
			if fromAddress == "" || !common.IsHexAddress(fromAddress) {
				fmt.Println("Error: not set from address of owner")
				fmt.Println(cmd.UsageString())
				return
			}

			name, _ := cmd.Flags().GetString("name")
			if name == "" {
				fmt.Println("Error: not set name")
				fmt.Println(cmd.UsageString())
				return
			}

			symbol, _ := cmd.Flags().GetString("symbol")
			if symbol == "" {
				fmt.Println("Error: not set symbol")
				fmt.Println(cmd.UsageString())
				return
			}

			var decimals uint8
			var totalSupply *big.Int
			if cli.mode != ModeERC721 {
				decimals, _ = cmd.Flags().GetUint8("decimals")
				if decimals < 0 || decimals > 18 {
					fmt.Println("Error: not set decimals or decimals invalid")
					fmt.Println(cmd.UsageString())
					return
				}

				totalSupplyStr, _ := cmd.Flags().GetString("total")
				if !IsDecimalString(totalSupplyStr) {
					fmt.Printf("totalSupply(%v) illegal\n", totalSupplyStr)
					return
				}
				var ok bool
				totalSupply, ok = getWeiAmountWeiByStringWithDecimals(totalSupplyStr, 10, decimals)
				if !ok {
					fmt.Println("Error: totalSupply invalid")
					fmt.Println(cmd.UsageString())
					return
				}
			}

			if cli.contractAddress == "" {
				save = true
			}
			cli.Deploy(fromAddress, name, symbol, decimals, totalSupply)

			if save {
				viper.WriteConfigAs(cli.config)
			}
		},
	}

	cmd.Flags().StringP("name", "n", "", "the name of the token")
	cmd.Flags().Uint8P("decimals", "d", 18, "the decimals of the token, 0~18")
	cmd.Flags().StringP("total", "t", "", "the total supply of the token")

	cmd.Flags().Bool("save", false, "save contract address to config file")

	cmd.MarkFlagRequired("name")
	cmd.MarkFlagRequired("symbol")
	// cmd.MarkFlagRequired("total")

	return cmd
}
