package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"strconv"

	"github.com/newtonproject/tokencommander/contract/ERC20"
	"github.com/newtonproject/tokencommander/contract/ERC721"
	"github.com/spf13/cobra"
)

func (cli *CLI) buildInfoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "info [-a contractAddress] [-s contractSymbol] [TokenID]",
		Short:                 "Show contract basic info",
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {
			SimpleToken, err := cli.GetSimpleToken()
			if err != nil {
				fmt.Println("GetSimpleToken Error: ", err)
				fmt.Println(cmd.UsageString())
				return
			}

			fmt.Printf("The contract address(%s) basic information is as follows:\n", cli.contractAddress)

			name, err := SimpleToken.Name(nil)
			if err != nil {
				fmt.Printf("Name: Get name Error(%v)\n", err)
			} else {
				fmt.Println("Name: ", name)
			}

			symbol, err := SimpleToken.Symbol(nil)
			if err != nil {
				fmt.Printf("Symbol: Get symbol Error(%v)\n", err)
			} else {
				fmt.Println("Symbol: ", symbol)
			}

			if cli.mode == ModeERC20 {
				decimals, err := SimpleToken.(*ERC20.SimpleToken).Decimals(nil)
				if err != nil {
					fmt.Printf("Decimals: Get decimals Error(%v)\n", err)
				} else {
					fmt.Println("Decimals: ", decimals)
				}

				totalSupply, err := SimpleToken.TotalSupply(nil)
				if err != nil {
					fmt.Printf("TotalSupply: Get totalSupply Error(%v)\n", err)
				} else {
					fmt.Println("TotalSupply: ", getAmountTextByWeiWithDecimals(totalSupply, decimals), symbol)
				}
			} else if cli.mode == ModeERC721 {
				totalSupply, err := SimpleToken.TotalSupply(nil)
				if err != nil {
					fmt.Printf("TotalSupply: Get totalSupply Error(%v)\n", err)
				} else {
					fmt.Println("TotalSupply: ", totalSupply.String())
				}

				if len(args) > 0 {
					id, err := strconv.Atoi(args[0])
					if err != nil {
						fmt.Printf("TokenID: Get token ID error(%v)\n", err)
						return
					}
					fmt.Printf("The info of token ID %d is as follows: \n", id)

					idBig := big.NewInt(int64(id))

					exists, err := SimpleToken.(*ERC721.SimpleToken).Exists(nil, idBig)
					if err != nil {
						fmt.Println("\tCheck token ID exists error: ", err)
						return
					}
					if !exists {
						fmt.Printf("\tToken ID %d not exists\n", id)
						return
					}

					// owner
					owner, err := SimpleToken.(*ERC721.SimpleToken).OwnerOf(nil, idBig)
					if err != nil {
						fmt.Printf("\tOwner: Get owner error(%v)\n", err)
					} else {
						fmt.Printf("\tOwner: %s\n", owner.String())
					}

					uri, err := SimpleToken.(*ERC721.SimpleToken).TokenURI(nil, idBig)
					if err != nil {
						fmt.Printf("\tTokenURI: Get token uri error(%v)\n", err)
					} else {
						fmt.Printf("\tTokenURI: %s\n", uri)

						// try to get metadata
						raw, err := getJson(uri)
						if err != nil {
							// fmt.Println("\tMetadata: get metadata error", err) // ignore
						} else {
							rawStr, err := json.MarshalIndent(raw, "", "\t")
							if err != nil {
								fmt.Println("\tMetadata: ", err)
							} else {
								fmt.Println("\tMetadata:", string(rawStr))
							}

						}
					}

				}

			}

			return
		},
	}

	return cmd
}

func getJson(uri string) (json.RawMessage, error) {
	resp, err := http.Get(uri)
	if err != nil {
		return nil, err
	}

	if resp.Body == nil {
		return nil, errors.New("body nil")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var raw json.RawMessage
	err = json.Unmarshal(body, &raw)
	if err != nil {
		return nil, err
	}

	return raw, nil
}
