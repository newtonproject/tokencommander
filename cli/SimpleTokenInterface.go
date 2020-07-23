package cli

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

// SimpleToken simpleToken interface
type SimpleToken interface {
	// Name returns the name of the token
	Name(opts *bind.CallOpts) (string, error)

	// Symbol returns the symbol of the token
	Symbol(opts *bind.CallOpts) (string, error)

	// TotalSupply returns the total token supply
	TotalSupply(opts *bind.CallOpts) (*big.Int, error)

	// BalanceOf returns the account balance of another account with address owner
	BalanceOf(opts *bind.CallOpts, owner common.Address) (*big.Int, error)
}
