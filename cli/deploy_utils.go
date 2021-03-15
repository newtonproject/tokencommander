//go:generate abigen --sol contract/ERC20/SimpleToken.sol --pkg ERC20 --out contract/ERC20/SimpleToken.go
//go:generate abigen --sol contract/ERC721/SimpleToken.sol --pkg ERC721 --out contract/ERC721/SimpleToken.go
package cli

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/newtonproject/tokencommander/contract/ERC20"
	"github.com/newtonproject/tokencommander/contract/ERC721"
	"github.com/spf13/viper"
)

// Deploy deploy contract
func (cli *CLI) Deploy(address, name, symbol string, decimals uint8, totalSupply *big.Int) {
	var err error

	opts, err := cli.getTransactOpts(address)
	if err != nil {
		fmt.Println("GetTransactOpts: ", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()
	opts.Context = ctx

	cli.BuildClient()
	client := cli.client
	var contractAddress common.Address
	tx := new(types.Transaction)
	if cli.mode == ModeERC721 {
		contractAddress, tx, _, err = ERC721.DeploySimpleToken(opts, client, name, symbol)
	} else {
		contractAddress, tx, _, err = ERC20.DeploySimpleToken(opts, client, name, symbol, decimals, totalSupply)
	}
	if err != nil {
		fmt.Println("DeployContract error: ", err)
		return
	}

	fmt.Printf("Contract deploy at address %s\n", contractAddress.String())
	fmt.Printf("Transaction waiting to be mined: 0x%x\n", tx.Hash())
	cli.contractAddress = contractAddress.String()
	viper.Set("contractaddress", cli.contractAddress)
	_, err = bind.WaitDeployed(opts.Context, client, tx)
	if err != nil {
		fmt.Println("WaitDeployed error: ", err)
		return
	}

	fmt.Println("Contract deploy success")
}
