package cli

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
)

func (cli *CLI) buildBalanceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "balance [address1] [address2] [address3]...",
		Short:                 "Balance of address on Token",
		Args:                  cobra.MinimumNArgs(0),
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {

			var addressList []common.Address
			if len(args) > 0 {
				for i := 0; i < len(args); i++ {
					if common.IsHexAddress(args[i]) {
						addressList = append(addressList, common.HexToAddress(args[i]))
					}
				}
			} else {
				err := cli.buildWallet()
				if err != nil {
					fmt.Println("BuildClient error: ", err)
					return
				}
				for _, account := range cli.wallet.Accounts() {
					addressList = append(addressList, account.Address)
				}
			}

			for _, address := range addressList {
				if cli.mode == ModeERC721 {
					tokens := cli.getTokensOfOwner(address)
					if tokens != nil && len(tokens) > 0 {
						fmt.Println(address.String(), cli.balanceOfText(address), tokens)
						continue
					}
				}
				fmt.Println(address.String(), cli.balanceOfText(address))
			}

			return
		},
	}

	return cmd
}
