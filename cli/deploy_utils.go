//go:generate abigen --sol contracts/contracts/contracts/NRC20/BaseToken.sol --pkg ERC20 --out contracts/ERC20/SimpleToken.go
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
	"github.com/spf13/viper"
)

// Deploy deploy contract
func (cli *CLI) Deploy(address, name, symbol, baseTokenURI string, decimals uint8, totalSupply *big.Int) {
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
		contractAddress, tx, _, err = ERC721.DeployNRC7Full(opts, client, name, symbol, baseTokenURI)
	} else {
		contractAddress, tx, _, err = ERC20.DeployBaseToken(opts, client, name, symbol, decimals,
			totalSupply, totalSupply, true, true)
	}
	if err != nil {
		fmt.Println("DeployContract error: ", err)
		return
	}

	fmt.Printf("Contract %s deploy at address %s\n", cli.mode, contractAddress.String())
	fmt.Printf("Transaction waiting to be mined: 0x%x\n", tx.Hash())
	cli.contractAddress = contractAddress.String()
	viper.Set("contractaddress", cli.contractAddress)
	_, err = bind.WaitDeployed(opts.Context, client, tx)
	if err != nil {
		fmt.Println("WaitDeployed error: ", err)
		return
	}

	fmt.Printf("Contract %s deploy success\n", cli.mode)
}
