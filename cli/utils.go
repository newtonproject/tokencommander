package cli

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/btcsuite/btcutil/base58"
	"github.com/ethereum/go-ethereum/common"
	prompt2 "github.com/ethereum/go-ethereum/console/prompt"
)

// IsDecimalString Check whether amount string is legal amount
var IsDecimalString = regexp.MustCompile(`^[1-9]\d*$|^0$|^0\.\d*$|^[1-9](\d)*\.(\d)*$`).MatchString

func showSuccess(msg string, args ...interface{}) {
	fmt.Printf(msg+"\n", args...)
}

// getPassPhrase retrieves the password associated with an account,
// requested interactively from the user.
func getPassPhrase(prompt string, confirmation bool) (string, error) {
	// prompt the user for the password
	if prompt != "" {
		fmt.Println(prompt)
	}
	password, err := prompt2.Stdin.PromptPassword("Enter passphrase (empty for no passphrase): ")
	if err != nil {
		return "", err
	}
	if confirmation {
		confirm, err := prompt2.Stdin.PromptPassword("Enter same passphrase again: ")
		if err != nil {
			return "", err
		}
		if password != confirm {
			return "", fmt.Errorf("Passphrases do not match")
		}
	}
	return password, nil
}

func stringInSlice(str string, list []string) bool {
	for _, v := range list {
		if v == str {
			return true
		}
	}
	return false
}

func getWeiAmountTextUnitByUnit(amount *big.Int, unit string) string {
	if amount == nil {
		return fmt.Sprintf("0 %s", UnitWEI)
	}
	amountStr := amount.String()
	amountStrLen := len(amountStr)
	if unit == "" {
		if amountStrLen <= 18 {
			// show in WEI
			unit = UnitWEI
		} else {
			unit = UnitETH
		}
	}

	return fmt.Sprintf("%s %s", getWeiAmountTextByUnit(amount, unit), unit)
}

func getWeiAmountTextByUnit(amount *big.Int, unit string) string {
	if amount == nil {
		return "0"
	}
	amountStr := amount.String()
	amountStrLen := len(amountStr)

	switch unit {
	case UnitETH:
		var amountStrDec, amountStrInt string
		if amountStrLen <= 18 {
			amountStrDec = strings.Repeat("0", 18-amountStrLen) + amountStr
			amountStrInt = "0"
		} else {
			amountStrDec = amountStr[amountStrLen-18:]
			amountStrInt = amountStr[:amountStrLen-18]
		}
		amountStrDec = strings.TrimRight(amountStrDec, "0")
		if len(amountStrDec) <= 0 {
			return amountStrInt
		}
		return amountStrInt + "." + amountStrDec

	case UnitWEI:
		return amountStr
	}

	return "Illegal Unit"
}

func getWeiAmountWeiByStringWithDecimals(amountStr string, base int, decimals uint8) (*big.Int, bool) {
	amount := new(big.Int)
	var aStr string
	index := strings.Index(amountStr, ".")
	if index < 0 {
		aStr = amountStr + strings.Repeat("0", int(decimals))
		return amount.SetString(aStr, base)
	} else if index == 0 {
		amountStr = "0" + amountStr
		index++
	}

	if index+1 >= len(amountStr) || len(amountStr[index+1:]) > int(decimals) {
		return nil, false
	}

	aStr = amountStr + strings.Repeat("0", int(decimals))
	aStr = aStr[:index] + aStr[index+1:index+1+int(decimals)]
	return amount.SetString(aStr, base)
}

func getAmountTextByWeiWithDecimals(amount *big.Int, decimals uint8) string {
	amountStr := amount.String()
	len := len(amountStr)

	if len <= int(decimals) {
		aStr := strings.TrimRight(amountStr, "0")
		if aStr == "" {
			return "0"
		}
		return "0." + strings.Repeat("0", int(decimals)-len) + aStr
	}

	if decimals == 0 {
		return amountStr[:len-int(decimals)]
	}
	aStr := strings.TrimRight(amountStr[len-int(decimals):], "0")
	if aStr == "" {
		return amountStr[:len-int(decimals)]
	}
	return amountStr[:len-int(decimals)] + "." + aStr
}

// showTransactionReceipt
func showTransactionReceipt(url, txStr string) {
	var jsonStr = []byte(fmt.Sprintf(`{"jsonrpc":"2.0","method":"eth_getTransactionReceipt","params":["%s"],"id":1}`, txStr))
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	clientHttp := &http.Client{}

	resp, err := clientHttp.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var body json.RawMessage
		if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
			fmt.Println(err)
			return
		}

		bodyStr, err := json.MarshalIndent(body, "", "    ")
		if err != nil {
			fmt.Println("JSON marshaling failed: ", err)
			return
		}
		fmt.Printf("%s\n", bodyStr)

		return
	}
}

func getFaucet(rpcURL, address string) {
	url := fmt.Sprintf("%s/faucet?address=%s", rpcURL, address)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Get error: %v\n", err)
		return
	}
	if resp.StatusCode == 200 {
		fmt.Printf("Get faucet for %s\n", address)
	}
}

func addressToNew(chainID []byte, address common.Address) string {
	input := append(chainID, address.Bytes()...)
	return "NEW" + base58.CheckEncode(input, 0)
}

func newToAddress(chainID []byte, newAddress string) (common.Address, error) {
	if newAddress[:3] != "NEW" {
		return common.Address{}, errors.New("not NEW address")
	}

	decoded, version, err := base58.CheckDecode(newAddress[3:])
	if err != nil {
		return common.Address{}, err
	}
	if version != 0 {
		return common.Address{}, errors.New("illegal version")
	}
	if len(decoded) < 20 {
		return common.Address{}, errors.New("illegal decoded length")
	}
	if !bytes.Equal(decoded[:len(decoded)-20], chainID) {
		return common.Address{}, errors.New("illegal ChainID")
	}

	address := common.BytesToAddress(decoded[len(decoded)-20:])

	return address, nil
}
