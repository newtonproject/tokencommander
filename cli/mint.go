package cli

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/newtonproject/tokencommander/contract/ERC721"
	"github.com/spf13/cobra"
)

func (cli *CLI) buildMintCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "mint <address> [tokenID] [--uri <tokenUri>]",
		Short:                 fmt.Sprintf("Command to mint tokenID for address, only for %s", ModeERC721),
		Args:                  cobra.MinimumNArgs(1),
		Aliases:               []string{"mine"},
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {

			if cli.mode != ModeERC721 {
				fmt.Println(errOnlyERC721)
				return
			}

			simpleToken, err := cli.GetSimpleToken()
			if err != nil {
				fmt.Println(err)
				return
			}
			erc721Token := simpleToken.(*ERC721.SimpleToken)

			if cli.address == "" || !common.IsHexAddress(cli.address) {
				fmt.Println("Error: not set from address of owner or from address illegal")
				return
			}
			isMinter, err := erc721Token.IsMinter(nil, common.HexToAddress(cli.address))
			if err != nil {
				fmt.Printf("Error: check minter error(%s)\n", err)
				return
			}
			if !isMinter {
				fmt.Printf("The from address(%s) is not minter\n", cli.address)
				return
			}

			toAddressStr := args[0]
			if toAddressStr == "" || !common.IsHexAddress(toAddressStr) {
				fmt.Println("Error: the address of token owner illegal")
				fmt.Fprint(os.Stderr, cmd.UsageString())
				return
			}
			toAddress := common.HexToAddress(toAddressStr)

			tokenID := big.NewInt(0)
			if len(args) > 1 {
				var ok bool
				tokenID, ok = tokenID.SetString(args[1], 10)
				if !ok {
					fmt.Printf("convert %s to tokenID error\n", args[1])
					return
				}
			} else {
				tokenID, err = erc721Token.TotalSupply(nil)
				if err != nil {
					fmt.Printf("Error: get totalSupply error(%s)\n", err)
					return
				}
			}

			exists, err := erc721Token.Exists(nil, tokenID)
			if err != nil {
				fmt.Println("Exists: ", err)
				return
			}
			if exists {
				fmt.Printf("Error: the token(%s) is exists, please specify another ID\n", tokenID.String())
				return
			}

			var tokenUri string
			if cmd.Flags().Changed("uri") {
				tokenUri, err = cmd.Flags().GetString("uri")
				if err != nil {
					fmt.Println("Get token url error: ", err)
					return
				}
			}

			opts, err := cli.getTransactOpts(cli.address)
			if err != nil {
				fmt.Println("GetTransactOpts: ", err)
				return
			}
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
			defer cancel()
			opts.Context = ctx

			var tx *types.Transaction
			if tokenUri == "" {
				tx, err = erc721Token.Mint(opts, toAddress, tokenID)
				if err != nil {
					fmt.Printf("Error: mint error(%s)\n", err)
					return
				}
			} else {
				tx, err = erc721Token.MintWithTokenURI(opts, toAddress, tokenID, tokenUri)
				if err != nil {
					fmt.Printf("Error: mint error(%s)\n", err)
					return
				}
			}

			fmt.Printf("Succeed mint token %s for address %s, TxID %s.\n", tokenID.String(), toAddress.String(), tx.Hash().String())
			fmt.Println("Waiting for transaction to be mined...")
			if _, err := bind.WaitMined(ctx, cli.client, tx); err != nil {
				fmt.Println(err)
				return
			}
			showTransactionReceipt(cli.rpcURL, tx.Hash().String())

			return
		},
	}

	cmd.Flags().String("uri", "", "mint with token uri")

	return cmd
}
