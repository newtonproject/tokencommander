package cli

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/newtonproject/tokencommander/contracts/ERC20"
	"github.com/newtonproject/tokencommander/contracts/ERC721"
)

func (cli *CLI) balanceOf(address common.Address) *big.Int {

	simpleToken, err := cli.GetSimpleToken()
	if err != nil {
		fmt.Printf("Balance: GetSimpleToken Error(%v)\n", err)
		return big.NewInt(0)
	}
	balance, err := simpleToken.BalanceOf(nil, address)
	if err != nil {
		fmt.Printf("Balance: BalanceAt Error(%v)\n", err)
		return big.NewInt(0)
	}

	return balance
}

func (cli *CLI) balanceOfText(address common.Address) string {

	simpleToken, err := cli.GetSimpleToken()
	if err != nil {
		return fmt.Sprintf("GetSimpleToken Error(%v)", err)
	}
	balance, err := simpleToken.BalanceOf(nil, address)
	if err != nil {
		return fmt.Sprintf("BalanceOf Error(%v)", err)
	}
	if cli.mode == ModeERC721 {
		return balance.String()
	}
	decimals, err := simpleToken.(*ERC20.BaseToken).Decimals(nil)
	if err != nil {
		return fmt.Sprintf("Decimals: Get Decimals Error(%v)\n", err)
	}
	symbol, err := simpleToken.Symbol(nil)
	if err != nil {
		return fmt.Sprintf("Symbol: Get Symbol Error(%v)\n", err)
	}

	return getAmountTextByWeiWithDecimals(balance, decimals) + " " + symbol
}

func (cli *CLI) getTokensOfOwner(address common.Address) []*big.Int {
	if cli.mode != ModeERC721 {
		return nil
	}
	simpleToken, err := cli.GetSimpleToken()
	if err != nil {
		return nil // fmt.Sprintf("GetSimpleToken Error(%v)", err)
	}

	lenBig, err := simpleToken.(*ERC721.NRC7Full).BalanceOf(nil, address)
	if err != nil {
		return nil
	}
	tokens := make([]*big.Int, lenBig.Uint64())
	for i := uint64(0); i < lenBig.Uint64(); i++ {
		token, err := simpleToken.(*ERC721.NRC7Full).TokenOfOwnerByIndex(nil, address, big.NewInt(0).SetUint64(i))
		if err != nil {
			return nil // fmt.Sprintf("Decimals: Get Decimals Error(%v)\n", err)
		}
		tokens[i] = big.NewInt(0).Set(token)
	}

	return tokens
}
