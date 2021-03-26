package cli

import (
	"bufio"
	"context"
	"fmt"
	"math/big"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/newtonproject/tokencommander/contracts/ERC20"
	"github.com/spf13/cobra"
)

func (cli *CLI) buildBatchPayCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "batchpay <batch.txt>",
		Aliases:               []string{"batch"},
		Short:                 fmt.Sprintf("Batch pay base on file <batch.txt>, only support for %s", ModeERC20),
		Args:                  cobra.MinimumNArgs(1),
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {

			if cli.mode != ModeERC20 {
				fmt.Printf("Only support for %s\n", ModeERC20)
				return
			}

			batchFileName := args[0]
			file, err := os.Open(batchFileName)
			if err != nil {
				fmt.Println(err)
				return
			}
			defer file.Close()

			err = cli.BuildClient()
			if err != nil {
				fmt.Println("BuildClient error: ", err)
				return
			}
			client := cli.client

			simpleToken, err := cli.GetSimpleToken()
			if err != nil {
				fmt.Println("GetSimpleToken Error: ", err)
				return
			}
			erc20, ok := simpleToken.(*ERC20.BaseToken)
			if !ok {
				fmt.Printf("Only support for %s\n", ModeERC20)
				return
			}

			callOpts := new(bind.CallOpts)
			callOpts.Pending = true

			decimals, err := erc20.Decimals(callOpts)
			if err != nil {
				fmt.Printf("Decimals: Get Decimals Error(%v)\n", err)
				return
			}
			symbol, err := erc20.Symbol(nil)
			if err != nil {
				fmt.Printf("Symbol: Get Symbol Error(%v)\n", err)
				return
			}

			// check from
			if cli.address == "" {
				fmt.Println("Not set from address")
				return
			}
			if !common.IsHexAddress(cli.address) {
				fmt.Println("from address not valid hex")
				return
			}
			address := common.HexToAddress(cli.address)
			if address == (common.Address{}) {
				fmt.Println("From address not set")
				return
			}
			wallet := keystore.NewKeyStore(cli.walletPath,
				keystore.StandardScryptN, keystore.StandardScryptP)
			if !wallet.HasAddress(address) {
				fmt.Println("From address not in wallet")
				return
			}

			ctx := context.Background()

			chainID, err := client.NetworkID(ctx)
			if err != nil {
				fmt.Println(err)
				return
			}

			gasPrice := big.NewInt(0)
			if cmd.Flags().Changed("price") {
				price, err := cmd.Flags().GetUint64("price")
				if err != nil {
					fmt.Println("price get error: ", err)
					return
				}
				gasPrice = big.NewInt(0).SetUint64(price)
			} else {
				gasPrice, err = client.SuggestGasPrice(ctx)
				if err != nil {
					fmt.Println("SuggestGasPrice error: ", err)
					return
				}
			}

			nonce := uint64(0)
			if cmd.Flags().Changed("nonce") {
				nonce, err = cmd.Flags().GetUint64("nonce")
				if err != nil {
					fmt.Println("nonce get error: ", err)
					return
				}
			} else {
				nonce, err = client.PendingNonceAt(ctx, address)
				if err != nil {
					fmt.Println("PendingNonceAt error: ", err)
					return
				}
			}

			type pay struct {
				to     common.Address
				amount *big.Int
			}

			batchList := make([]pay, 0)
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				text := scanner.Text()
				l := strings.Split(text, ",")
				if len(l) != 2 {
					fmt.Println("parse error: ", text)
					return
				}

				var to common.Address
				if common.IsHexAddress(l[0]) {
					to = common.HexToAddress(l[0])
				} else {
					if cli.blockchain != NewChain {
						fmt.Println("Convert address error: ", l[0])
						return
					}
					to, err = newToAddress(chainID.Bytes(), l[0])
					if err != nil {
						fmt.Println("NewChain: address is invalid hex address or convert from NEW Address to hex error: ", l[0])
						return
					}
				}
				if to == (common.Address{}) {
					fmt.Println("Warning: to address is zero: ", l[0])
				}

				amount, ok := getWeiAmountWeiByStringWithDecimals(l[1], 10, decimals)
				if !ok {
					fmt.Printf("Amount: convert (%s) from string to amount with decimals(%d) error\n", l[1], decimals)
					return
				}

				batchList = append(batchList, pay{
					to:     to,
					amount: amount})

			}

			fmt.Println("Please confirm the transactions below:")
			totalAmount := big.NewInt(0)
			for _, b := range batchList {
				fmt.Printf("%s,%s\n", b.to.String(),
					getAmountTextByWeiWithDecimals(b.amount, decimals))

				// total
				totalAmount.Add(totalAmount, b.amount)
			}
			fmt.Println("Number of transactions:", len(batchList))

			if totalAmount.Cmp(big.NewInt(0)) <= 0 {
				fmt.Println("Total pay amount is zero")
				return
			}

			balance, err := erc20.BalanceOf(callOpts, address)
			if err != nil {
				fmt.Println(err)
				return
			}

			if balance.Cmp(totalAmount) < 0 {
				fmt.Println("Error: Insufficient funds")
				return
			}

			fmt.Println("Total pay amount:", getAmountTextByWeiWithDecimals(totalAmount, decimals), symbol)

			opts, err := cli.getBatchTransactOpts(address.String())
			if err != nil {
				fmt.Println("GetTransactOpts: ", err)
				return
			}
			opts.Context = ctx
			opts.GasPrice = gasPrice

			wait, _ := cmd.Flags().GetBool("wait")
			gasTotal := big.NewInt(0)
			for _, b := range batchList {
				to := b.to
				amount := b.amount
				opts.Nonce = big.NewInt(0).SetUint64(nonce)
				tx, err := erc20.Transfer(opts, to, amount)
				if err != nil {
					fmt.Println(err)
					return
				}

				fmt.Printf("Succeed broadcast pay %s %s to %s from %s with nonce %d, TxID %s.\n",
					getAmountTextByWeiWithDecimals(amount, decimals), symbol,
					to.String(), address.String(), nonce, tx.Hash().String())

				if wait {
					txr, err := bind.WaitMined(ctx, client, tx)
					if err != nil {
						fmt.Println(err)
						return
					}
					if txr.Status == 1 {
						fmt.Printf("Succeed mined txID %s.\n", txr.TxHash.String())
					} else {
						fmt.Printf("Succeed mined txID %s but status failed.\n", txr.TxHash.String())
					}
					gasTotal.Add(gasTotal, big.NewInt(0).Mul(tx.GasPrice(), big.NewInt(0).SetUint64(txr.GasUsed)))
				} else {
					gasTotal.Add(gasTotal, big.NewInt(0).Mul(tx.GasPrice(), big.NewInt(0).SetUint64(tx.Gas())))
				}

				nonce++
			}

			fmt.Printf("Total Gas is: %s %s\n", getWeiAmountTextByUnit(gasTotal, UnitETH), UnitETH)
		},
	}

	cmd.Flags().Uint64P("price", "p", 1, fmt.Sprintf("the gasPrice used for each paid gas (unit in %s)", UnitWEI))
	cmd.Flags().Uint64P("nonce", "n", 0, "the number of nonce to start")
	cmd.Flags().Bool("wait", false, "wait for transaction to mined")

	return cmd
}
