package cli

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func (cli *CLI) buildAccountCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "account [new|list|balance]",
		Short: fmt.Sprintf("Manage %s accounts", cli.blockchain.String()),
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			return
		},
	}

	cmd.AddCommand(cli.buildAccountNewCmd())
	cmd.AddCommand(cli.buildAccountListCmd())
	cmd.AddCommand(cli.buildAccountBalanceCmd())

	return cmd
}

func (cli *CLI) buildAccountNewCmd() *cobra.Command {
	accountNewCmd := &cobra.Command{
		Use:                   "new [--faucet] [--numOfNew amount]",
		Short:                 "create a new account",
		Args:                  cobra.MinimumNArgs(0),
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			walletPath := cli.walletPath
			wallet := keystore.NewKeyStore(walletPath,
				keystore.LightScryptN, keystore.LightScryptP)

			if cli.walletPassword == "" {
				cli.walletPassword, err = getPassPhrase("Your new account is locked with a password. Please give a password. Do not forget this password.", true)
				if err != nil {
					fmt.Println("Error: ", err)
					return
				}
			}

			numOfNew, err := cmd.Flags().GetInt("numOfNew")
			if err != nil {
				numOfNew = viper.GetInt("account.numOfNew")
			}
			if numOfNew <= 0 {
				fmt.Printf("number[%d] of new account less then 1\n", numOfNew)
				numOfNew = 1
			}

			faucet, _ := cmd.Flags().GetBool("faucet")

			for i := 0; i < numOfNew; i++ {
				account, err := wallet.NewAccount(cli.walletPassword)
				if err != nil {
					fmt.Println("Account error:", err)
					return
				}
				if faucet {
					getFaucet(cli.rpcURL, account.Address.String())
				}
				fmt.Println(account.Address.Hex())
				if cli.address == "" {
					cli.address = account.Address.String()
				}
			}
		},
	}

	accountNewCmd.Flags().IntP("numOfNew", "n", 1, "number of the new account")
	accountNewCmd.Flags().Bool("faucet", false, "get faucet for new account")
	return accountNewCmd
}

func (cli *CLI) buildAccountListCmd() *cobra.Command {
	accountListCmd := &cobra.Command{
		Use:                   "list",
		Short:                 "list all accounts in the wallet path",
		Args:                  cobra.MinimumNArgs(0),
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {
			walletPath := cli.walletPath
			wallet := keystore.NewKeyStore(walletPath,
				keystore.LightScryptN, keystore.LightScryptP)
			if len(wallet.Accounts()) == 0 {
				fmt.Println("Empty wallet, create account first.")
				return
			}

			for _, account := range wallet.Accounts() {
				fmt.Println(account.Address.Hex())
			}
		},
	}

	return accountListCmd
}

func (cli *CLI) buildAccountBalanceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   fmt.Sprintf("balance [-u %s] [-n pending] [-s] [address1] [address2]...", strings.Join(UnitList, "|")),
		Short:                 "Get balance of address",
		Args:                  cobra.MinimumNArgs(0),
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {
			cli.showBalance(cmd, args, false)
		},
	}

	cmd.Flags().StringP("unit", "u", "", fmt.Sprintf("unit for balance. %s.", UnitString))
	cmd.Flags().Bool("safe", false, "enable safe mode to check balance (force use the block 3 block heights less than the latest)")
	cmd.Flags().StringP("number", "n", "latest", `the integer block number, or the string "latest", "earliest" or "pending"`)

	return cmd
}

func (cli *CLI) showBalance(cmd *cobra.Command, args []string, showSum bool) {
	var err error

	unit, _ := cmd.Flags().GetString("unit")
	if unit != "" && !stringInSlice(unit, UnitList) {
		fmt.Printf("Unit(%s) for invalid. %s.\n", unit, UnitString)
		fmt.Fprint(os.Stderr, cmd.UsageString())
		return
	}

	safe, _ := cmd.Flags().GetBool("safe")

	pending := false
	latest := true
	number := big.NewInt(0)
	if !safe {
		if cmd.Flags().Changed("number") {
			numStr, err := cmd.Flags().GetString("number")
			if err != nil {
				fmt.Println("Error: arg number get error: ", err)
				return
			}
			switch numStr {
			case "pending":
				pending = true
				latest = false
			case "earliest":
				pending = false
				latest = false
				number = number.SetUint64(0)
			case "latest":
				latest = true
				pending = false
			default:
				pending = false
				latest = false
				var ok bool
				number, ok = number.SetString(numStr, 10)
				if !ok {
					fmt.Println("Error: arg number convert to big int error")
					return
				}
				if number.Cmp(big.NewInt(0)) < 0 {
					fmt.Println("Error: arg number is less than 0")
					return
				}
			}
		}
	}

	var addressList []common.Address

	if len(args) <= 0 {
		if err := cli.buildWallet(); err != nil {
			fmt.Println(err)
			return
		}

		for _, account := range cli.wallet.Accounts() {
			addressList = append(addressList, account.Address)
		}

	} else {
		for _, addressStr := range args {
			addressList = append(addressList, common.HexToAddress(addressStr))
		}
	}

	if err := cli.BuildClient(); err != nil {
		fmt.Println(err)
		return
	}
	ctx := context.Background()

	var blockNumber *big.Int
	if safe {
		latestHeader, err := cli.client.HeaderByNumber(ctx, nil)
		if err != nil {
			fmt.Println("HeaderByBlock error: ", err)
			return
		}
		if latestHeader == nil {
			fmt.Println("HeaderByBlock return nil")
			return
		}
		blockNumber = big.NewInt(0).Sub(latestHeader.Number, big.NewInt(3))
		fmt.Printf("Safe mode enable, check balance at block height %s while the latest is %s\n", blockNumber.String(), latestHeader.Number.String())
	} else if !latest && !pending {
		blockNumber = number
	}

	balanceSum := big.NewInt(0)
	for _, address := range addressList {
		var balance *big.Int
		if pending {
			balance, err = cli.client.PendingBalanceAt(ctx, address)
		} else {
			balance, err = cli.client.BalanceAt(ctx, address, blockNumber)
		}
		balanceSum.Add(balanceSum, balance)
		if err != nil {
			fmt.Println("Balance error:", err)
			return
		}
		fmt.Printf("Address[%s] Balance[%s]\n", address.Hex(), getWeiAmountTextUnitByUnit(balance, unit))
	}

	if showSum {
		fmt.Println("Number Of Accounts:", len(addressList))
		fmt.Println("Total Balance:", getWeiAmountTextUnitByUnit(balanceSum, unit))
	}

	return
}
