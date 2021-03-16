package cli

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/newtonproject/tokencommander/contract/ERC20"
	"github.com/newtonproject/tokencommander/contract/ERC721"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	buildCommit string
	buildDate   string
)

// CLI represents a command-line interface. This class is
// not threadsafe.
type CLI struct {
	Name       string
	rootCmd    *cobra.Command
	version    string
	walletPath string
	rpcURL     string
	config     string

	contractAddress string
	localSymbol     string
	client          *ethclient.Client
	wallet          *keystore.KeyStore
	account         accounts.Account
	SimpleToken     SimpleToken
	walletPassword  string
	address         string
	mode            string

	blockchain BlockChain
}

// NewCLI returns an initialized CLI
func NewCLI() *CLI {
	bc, err := getBlockChain()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	version := "v0.5.2"
	if buildCommit != "" {
		version = fmt.Sprintf("%s-%s", version, buildCommit)
	}
	if buildDate != "" {
		version = fmt.Sprintf("%s-%s", version, buildDate)
	}
	version = fmt.Sprintf("%s-%s", version, bc.String())

	// init unit
	bc.Init()

	cli := &CLI{
		Name:       filepath.Base(os.Args[0]), // "TokenCommander"
		rootCmd:    nil,
		version:    version,
		walletPath: "",
		rpcURL:     "",
		//	testing:         false,
		config:          "",
		contractAddress: "",
		client:          nil,
		SimpleToken:     nil,
		walletPassword:  "",
		mode:            ModeERC20,
		blockchain:      bc,
	}

	cli.buildRootCmd()
	return cli
}

// BuildClient BuildClient
func (cli *CLI) BuildClient() error {
	var err error
	if cli.client == nil {
		cli.client, err = ethclient.Dial(cli.rpcURL)
		if err != nil {
			return fmt.Errorf("Failed to connect to the NewChain client: %v", err)
		}
	}
	return nil
}

// BuildSimpleToken BuildClient
func (cli *CLI) buildSimpleToken() (SimpleToken, error) {
	var err error
	if cli.client == nil {
		cli.BuildClient()
	}

	symbol := cli.localSymbol
	if symbol != "" {
		addressStr := viper.GetString(fmt.Sprintf("Contracts.%s", symbol))
		if addressStr == "" {
			return nil, fmt.Errorf("contract address of symbol %s not set", symbol)
		}
		if !common.IsHexAddress(addressStr) {
			return nil, fmt.Errorf("contract address from symbol %s invalid", symbol)
		}
		cli.contractAddress = addressStr
	}

	if !common.IsHexAddress(cli.contractAddress) {
		return nil, fmt.Errorf("contract address is invalid")
	}

	if cli.mode == ModeERC721 {
		cli.SimpleToken, err = ERC721.NewSimpleToken(common.HexToAddress(cli.contractAddress), cli.client)
	} else {
		cli.SimpleToken, err = ERC20.NewSimpleToken(common.HexToAddress(cli.contractAddress), cli.client)
	}
	if err != nil {
		return nil, fmt.Errorf("NewSimpleToken Error(%v)", err)
	}
	return cli.SimpleToken, nil
}

// GetSimpleToken GetSimpleToken
func (cli *CLI) GetSimpleToken() (SimpleToken, error) {
	if cli.SimpleToken == nil {
		return cli.buildSimpleToken()
	}
	return cli.SimpleToken, nil
}

func (cli *CLI) buildWallet() error {
	if cli.wallet == nil {
		cli.wallet = keystore.NewKeyStore(cli.walletPath,
			keystore.LightScryptN, keystore.LightScryptP)
		if len(cli.wallet.Accounts()) == 0 {
			return fmt.Errorf("Empty wallet, create account first")
		}
	}

	return nil
}

func (cli *CLI) buildAccount(address string) error {

	err := cli.buildWallet()
	if err != nil {
		return err
	}

	if !common.IsHexAddress(address) {
		if common.IsHexAddress(cli.address) {
			address = cli.address
		} else {
			return fmt.Errorf("Error: address(%s) invalid", address)
		}
	}
	cli.account, err = cli.wallet.Find(accounts.Account{Address: common.HexToAddress(address)})
	if err != nil {
		return fmt.Errorf("Error: Can not get the keystore file of address %s", address)
	}
	cli.address = address

	return nil
}

func (cli *CLI) getTransactOpts(address string) (*bind.TransactOpts, error) {
	err := cli.buildAccount(address)
	if err != nil {
		return nil, err
	}

	cli.BuildClient()
	networkID, err := cli.client.NetworkID(context.Background())
	if err != nil {
		fmt.Println("NetworkID Error: ", err)
		return nil, err
	}

	opts := NewKeyedTransactorByAccount(cli.wallet, cli.account, cli.walletPassword, networkID)
	return opts, nil
}

func (cli *CLI) getBatchTransactOpts(address string) (*bind.TransactOpts, error) {
	err := cli.buildAccount(address)
	if err != nil {
		return nil, err
	}

	cli.BuildClient()
	networkID, err := cli.client.NetworkID(context.Background())
	if err != nil {
		fmt.Println("NetworkID Error: ", err)
		return nil, err
	}

	wallet := cli.wallet
	passphrase := cli.walletPassword
	account := cli.account
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

	opts := NewBatchKeyedTransactorByAccount(wallet, account, networkID)

	return opts, nil
}

// Execute parses the command line and processes it.
func (cli *CLI) Execute() {
	cli.rootCmd.Execute()
}

// setup turns up the CLI environment, and gets called by Cobra before
// a command is executed.
func (cli *CLI) setup(cmd *cobra.Command, args []string) {
	err := setupConfig(cli)
	if err != nil {
		fmt.Println(err)
		fmt.Fprint(os.Stderr, cmd.UsageString())
		os.Exit(1)
	}
	if cmd.Flags().Changed("symbol") {
		if symbol, _ := cmd.Flags().GetString("symbol"); symbol != "" {
			cli.localSymbol = symbol
		}
	}
}

func (cli *CLI) help(cmd *cobra.Command, args []string) {
	fmt.Fprint(os.Stderr, cmd.UsageString())

	os.Exit(-1)

}

// TestCommand test command
func (cli *CLI) TestCommand(command string) string {
	// cli.testing = true
	result := cli.Run(strings.Fields(command)...)
	//	cli.testing = false
	return result
}

// Run executes CLI with the given arguments. Used for testing. Not thread safe.
func (cli *CLI) Run(args ...string) string {
	oldStdout := os.Stdout

	r, w, _ := os.Pipe()

	os.Stdout = w

	cli.rootCmd.SetArgs(args)
	cli.rootCmd.Execute()
	cli.buildRootCmd()

	w.Close()

	os.Stdout = oldStdout

	var stdOut bytes.Buffer
	io.Copy(&stdOut, r)
	return stdOut.String()
}

// Embeddable returns a CLI that you can embed into your own Go programs. This
// is not thread-safe.
func (cli *CLI) Embeddable() *CLI {

	return cli
}

// SetPassword SetPassword
func (cli *CLI) SetPassword(_passPhrase string) *CLI {
	cli.walletPassword = _passPhrase
	return cli
}
