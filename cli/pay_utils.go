package cli

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/newtonproject/tokencommander/contracts/ERC20"
	"github.com/newtonproject/tokencommander/contracts/ERC721"
)

// GasFail GasFail info from SuggestGasPrice
var GasFail = "failed to estimate gas needed: gas required exceeds allowance or always failing transaction"

// TxFailAlways Replacement information GasFail
var TxFailAlways = "This is a transaction that will always fail. Please check contract and parameters again."

// SubmitTransaction SubmitTransaction
func (cli *CLI) pay(fromAddress, toAddress common.Address, amountStr string, nowait bool) {
	var err error

	cli.BuildClient()
	client := cli.client

	simpleToken, err := cli.GetSimpleToken()
	if err != nil {
		fmt.Println("GetSimpleToken Error: ", err)
		return
	}

	callOpts := new(bind.CallOpts)
	callOpts.Pending = true

	symbol, err := simpleToken.Symbol(callOpts)
	if err != nil {
		fmt.Printf("Symbol: Get Symbol Error(%v)\n", err)
		return
	}

	opts, err := cli.getTransactOpts(fromAddress.String())
	if err != nil {
		fmt.Println("GetTransactOpts: ", err)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()
	opts.Context = ctx

	var tx *types.Transaction
	if cli.mode == ModeERC721 {
		if !IsDecimalString(amountStr) {
			fmt.Printf("amount(%v) illegal\n", amountStr)
			return
		}

		tokenID, ok := getWeiAmountWeiByStringWithDecimals(amountStr, 10, 0)
		if !ok {
			fmt.Println("amount invalid: ", amountStr)
			return
		}
		tokenOwner, err := simpleToken.(*ERC721.NRC7Full).OwnerOf(callOpts, tokenID)
		if err != nil {
			fmt.Printf("OwnerOf: OwnerOf Error(%v)\n", err)
			return
		}
		if tokenOwner != fromAddress {
			fmt.Printf("The owner of tokenID(%s) is %s not %s\n", tokenID.String(), tokenOwner.String(), fromAddress.String())
		}

		fmt.Printf("Try to transfer tokenID %s to %s from %s ...\n",
			tokenID, toAddress.String(), fromAddress.String())
		tx, err = simpleToken.(*ERC721.NRC7Full).TransferFrom(opts, fromAddress, toAddress, tokenID)
		if err != nil {
			if GasFail == err.Error() {
				fmt.Println("SubmitTransaction error: ", TxFailAlways)
				return
			}
			fmt.Println("SubmitTransaction error: ", err)
			return
		}

		fmt.Printf("Succeed transfer tokenID %s to %s from %s, TxID %s.\n", tokenID, toAddress.String(), fromAddress.String(), tx.Hash().String())

	} else {
		decimals, err := simpleToken.(*ERC20.BaseToken).Decimals(callOpts)
		if err != nil {
			fmt.Printf("Decimals: Get decimals Error(%v)\n", err)
			return
		}

		balance, err := simpleToken.BalanceOf(callOpts, fromAddress)
		if err != nil {
			fmt.Printf("Balance: BalanceOf Error(%v)\n", err)
			return
		}

		amount := big.NewInt(0)
		if amountStr == "all" {
			amount = big.NewInt(0).Set(balance)
		} else {
			if !IsDecimalString(amountStr) {
				fmt.Printf("amount(%v) illegal\n", amountStr)
				return
			}
			var ok bool
			amount, ok = getWeiAmountWeiByStringWithDecimals(amountStr, 10, decimals)
			if !ok {
				fmt.Println("amount invalid: ", amountStr)
				return
			}
			if balance.Cmp(amount) < 0 {
				fmt.Printf("There is not enough balance(%s) to pay the amount(%s) of current transactions.\n",
					getAmountTextByWeiWithDecimals(balance, decimals), getAmountTextByWeiWithDecimals(amount, decimals))
				return
			}
		}

		fmt.Printf("Try to pay %s %s to %s from %s ...\n",
			getAmountTextByWeiWithDecimals(amount, decimals),
			symbol, toAddress.String(), fromAddress.String())
		tx, err = simpleToken.(*ERC20.BaseToken).Transfer(opts, toAddress, amount)
		if err != nil {
			if GasFail == err.Error() {
				fmt.Println("SubmitTransaction error: ", TxFailAlways)
				return
			}
			fmt.Println("SubmitTransaction error: ", err)
			return
		}

		fmt.Printf("Succeed submit pay %s %s to %s from %s, TxID %s.\n", getAmountTextByWeiWithDecimals(amount, decimals),
			symbol, toAddress.String(), fromAddress.String(), tx.Hash().String())
	}

	if !nowait {
		fmt.Println("Waiting for transaction to be mined...")
		txr, err := bind.WaitMined(ctx, client, tx)
		if err != nil {
			fmt.Println("WaitMined error: ", err)
			return
		}
		showTransactionReceipt(cli.rpcURL, tx.Hash().String())

		txStatus := "success"
		if txr.Status != types.ReceiptStatusSuccessful {
			txStatus = "failed"
		}
		fmt.Printf("The tx %s is confirmed and status is %s, with GasFee(%s) = GasPrice(%s) x GasUsed(%d)\n",
			tx.Hash().String(),
			txStatus,
			getWeiAmountTextByUnit(big.NewInt(0).Mul(tx.GasPrice(), big.NewInt(0).SetUint64(txr.GasUsed)), UnitETH),
			getWeiAmountTextByUnit(tx.GasPrice(), UnitETH),
			txr.GasUsed)
	}

}
