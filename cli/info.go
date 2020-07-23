package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/newtonproject/tokencommander/contract/ERC20"
)

func (cli *CLI) buildInfoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "info [-a contractAddress] [-s contractSymbol]",
		Short:                 "Show contract basic info",
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {
			SimpleToken, err := cli.GetSimpleToken()
			if err != nil {
				fmt.Println("GetSimpleToken Error: ", err)
				fmt.Println(cmd.UsageString())
				return
			}

			fmt.Printf("The contract address(%s) basic information is as follows:\n", cli.contractAddress)

			name, err := SimpleToken.Name(nil)
			if err != nil {
				fmt.Printf("Name: Get name Error(%v)\n", err)
			} else {
				fmt.Println("Name: ", name)
			}

			symbol, err := SimpleToken.Symbol(nil)
			if err != nil {
				fmt.Printf("Symbol: Get symbol Error(%v)\n", err)
			} else {
				fmt.Println("Symbol: ", symbol)
			}

			if cli.mode == ModeERC20 {
				decimals, err := SimpleToken.(*ERC20.SimpleToken).Decimals(nil)
				if err != nil {
					fmt.Printf("Decimals: Get decimals Error(%v)\n", err)
				} else {
					fmt.Println("Decimals: ", decimals)
				}

				totalSupply, err := SimpleToken.TotalSupply(nil)
				if err != nil {
					fmt.Printf("TotalSupply: Get totalSupply Error(%v)\n", err)
				} else {
					fmt.Println("TotalSupply: ", getAmountTextByWeiWithDecimals(totalSupply, decimals), symbol)
				}
			} else if cli.mode == ModeERC721 {
				totalSupply, err := SimpleToken.TotalSupply(nil)
				if err != nil {
					fmt.Printf("TotalSupply: Get totalSupply Error(%v)\n", err)
				} else {
					fmt.Println("TotalSupply: ", totalSupply.String())
				}
			}

			return
		},
	}

	return cmd
}
