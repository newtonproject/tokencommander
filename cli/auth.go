// Copyright 2016 The go-ethereum Authors
// Copyright 2019 The Newton Foundation
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package cli

import (
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

// NewTransactor is a utility method to easily create a transaction signer from
// an encrypted json key stream and the associated passphrase.
func NewTransactor(keyin io.Reader, passphrase string, networkID *big.Int) (*bind.TransactOpts, error) {
	json, err := ioutil.ReadAll(keyin)
	if err != nil {
		return nil, err
	}
	key, err := keystore.DecryptKey(json, passphrase)
	if err != nil {
		return nil, err
	}
	return NewKeyedTransactor(key.PrivateKey, networkID), nil
}

// NewKeyedTransactor is a utility method to easily create a transaction signer
// from a single private key.
func NewKeyedTransactor(key *ecdsa.PrivateKey, networkID *big.Int) *bind.TransactOpts {
	keyAddr := crypto.PubkeyToAddress(key.PublicKey)
	return &bind.TransactOpts{
		From: keyAddr,
		Signer: func(address common.Address, tx *types.Transaction) (*types.Transaction, error) {
			// force use ChainID
			signer := types.NewEIP155Signer(networkID)
			if address != keyAddr {
				return nil, errors.New("not authorized to sign this account")
			}
			signature, err := crypto.Sign(signer.Hash(tx).Bytes(), key)
			if err != nil {
				return nil, err
			}
			return tx.WithSignature(signer, signature)
		},
	}
}

func NewKeyedTransactorByAccount(wallet *keystore.KeyStore, account accounts.Account, passphrase string, networkID *big.Int) *bind.TransactOpts {
	return &bind.TransactOpts{
		From: account.Address,
		Signer: func(address common.Address, tx *types.Transaction) (*types.Transaction, error) {
			fmt.Println("The tx is as follow: ")
			fmt.Println("\tFrom:", account.Address.String())
			if tx.To() == nil {
				fmt.Println("\tTo: ContractCreate")
			} else {
				fmt.Println("\tTo:", tx.To().String())
			}
			fmt.Println("\tValue:", getWeiAmountTextByUnit(tx.Value(), UnitETH))
			fmt.Println("\tData:", hex.EncodeToString(tx.Data()))
			fmt.Println("\tNonce:", tx.Nonce())
			fmt.Println("\tGasPrice:", getWeiAmountTextByUnit(tx.GasPrice(), UnitETH))
			fmt.Println("\tGasLimit:", tx.Gas())
			fmt.Println("\tGasFee:", getWeiAmountTextByUnit(big.NewInt(0).Mul(tx.GasPrice(), big.NewInt(0).SetUint64(tx.Gas())), UnitETH))

			for trials := 0; trials <= 1; trials++ {
				err := wallet.Unlock(account, passphrase)
				if err == nil {
					break
				}
				if trials >= 1 {
					return nil, fmt.Errorf("failed to unlock account %s (%v)", account.Address.String(), err)

				}
				prompt := fmt.Sprintf("Unlocking account %s", account.Address.String())
				passphrase, _ = getPassPhrase(prompt, false)
			}

			return wallet.SignTx(account, tx, networkID)
		},
	}
}

func NewBatchKeyedTransactorByAccount(wallet *keystore.KeyStore, account accounts.Account, networkID *big.Int) *bind.TransactOpts {
	return &bind.TransactOpts{
		From: account.Address,
		Signer: func(address common.Address, tx *types.Transaction) (*types.Transaction, error) {
			return wallet.SignTx(account, tx, networkID)
		},
	}
}
