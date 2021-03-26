package cli

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/newtonproject/tokencommander/contracts/ERC721"
	"github.com/spf13/cobra"
)

var MinterRole = crypto.Keccak256([]byte("MINTER_ROLE")) // 0x9f2df0fed2c77648de5860a4cc508cd0818c85b8b8a1ab4ceeef8d981c8956a6}

func (cli *CLI) buildMintCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "mint <address> [--uri <tokenUri>]",
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
			erc721Token := simpleToken.(*ERC721.NRC7Full)

			parsed, err := abi.JSON(strings.NewReader(ERC721.NRC7FullABI))
			if err != nil {
				fmt.Println(err)
				return
			}
			contractAddress := common.HexToAddress(cli.contractAddress)
			contract := bind.NewBoundContract(contractAddress, parsed, cli.client, cli.client, cli.client)

			if cli.address == "" || !common.IsHexAddress(cli.address) {
				fmt.Println("Error: not set from address of owner or from address illegal")
				return
			}

			var MinterRole32 [32]byte
			copy(MinterRole32[:], MinterRole[:])
			isMinter, err := erc721Token.HasRole(nil, MinterRole32, common.HexToAddress(cli.address))
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
				tx, err = erc721Token.Mint(opts, toAddress)
				if err != nil {
					fmt.Printf("Error: mint error(%s)\n", err)
					return
				}
			} else {
				tx, err = erc721Token.MintWithTokenURI(opts, toAddress, tokenUri)
				if err != nil {
					fmt.Printf("Error: mint error(%s)\n", err)
					return
				}
			}

			fmt.Printf("Succeed mint token for address %s, TxID %s.\n", toAddress.String(), tx.Hash().String())
			fmt.Println("Waiting for transaction to be mined...")
			txr, err := bind.WaitMined(ctx, cli.client, tx)
			if err != nil {
				fmt.Println(err)
				return
			}
			// txr.Logs

			if len(txr.Logs) > 0 {
				log := *(txr.Logs[0])
				var transferLog ERC721.NRC7FullTransfer
				err = contract.UnpackLog(&transferLog, "Transfer", log)
				if err != nil {
					fmt.Println("Unpack Log error: ", err)
					return
				}

				if transferLog.From != (common.Address{}) {
					fmt.Println("From address from mint log is not zero")
					return
				}

				fmt.Println("The tokenID is: ", transferLog.TokenId.String())
			}

			return
		},
	}

	cmd.Flags().String("url", "", "mint with token url")

	return cmd
}
