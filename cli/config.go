package cli

import (
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/viper"
)

const defaultConfigFile = "./config.toml"
const defaultWalletPath = "./wallet/"
const defaultContractAddress = ""

var defaultRPCURL string

const defaultNEWRPCURL = "https://rpc1.newchain.newtonproject.org"
const defaultETHRPCUrl = "https://ethrpc.service.newtonproject.org"

func InitRPCUrl(bc BlockChain) {
	// default RPC Url
	defaultRPCURL = defaultETHRPCUrl
	if bc == NewChain {
		defaultRPCURL = defaultNEWRPCURL
	}
}

func defaultConfig(cli *CLI) {
	viper.BindPFlag("walletPath", cli.rootCmd.PersistentFlags().Lookup("walletPath"))
	viper.BindPFlag("rpcURL", cli.rootCmd.PersistentFlags().Lookup("rpcURL"))
	viper.BindPFlag("contractAddress", cli.rootCmd.PersistentFlags().Lookup("contractAddress"))
	viper.BindPFlag("from", cli.rootCmd.PersistentFlags().Lookup("from"))
	viper.BindPFlag("mode", cli.rootCmd.PersistentFlags().Lookup("mode"))

	viper.SetDefault("walletPath", defaultWalletPath)
	viper.SetDefault("rpcURL", defaultRPCURL)
	viper.SetDefault("contractAddress", defaultContractAddress)
	viper.SetDefault("mode", ModeERC20)
}

func setupConfig(cli *CLI) error {

	// var ret bool
	var err error

	defaultConfig(cli)

	viper.SetConfigName(defaultConfigFile)
	viper.AddConfigPath(".")
	cfgFile := cli.config
	if cfgFile != "" {
		if _, err = os.Stat(cfgFile); err == nil {
			viper.SetConfigFile(cfgFile)
			err = viper.ReadInConfig()
		} else {
			// The default configuration is enabled.
			// fmt.Println(err)
			err = nil
		}
	} else {
		// The default configuration is enabled.
		err = nil
	}

	if rpcURL := viper.GetString("rpcURL"); rpcURL != "" {
		cli.rpcURL = rpcURL
	}
	if walletPath := viper.GetString("walletPath"); walletPath != "" {
		cli.walletPath = walletPath
	}
	if walletPassword := viper.GetString("Password"); walletPassword != "" {
		cli.walletPassword = walletPassword
	}
	if contractAddress := viper.GetString("contractAddress"); contractAddress != "" && common.IsHexAddress(contractAddress) {
		cli.contractAddress = contractAddress
	}
	if address := viper.GetString("from"); address != "" && common.IsHexAddress(address) {
		cli.address = address
	}
	if mode := viper.GetString("mode"); mode != "" {
		if mode != ModeERC20 && mode != ModeERC721 {
			return fmt.Errorf("not support mode %s, only support %s|%s", mode, ModeERC20, ModeERC721)
		}
		cli.mode = mode
	}

	return err
}
