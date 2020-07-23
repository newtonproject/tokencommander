package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func (cli *CLI) buildRootCmd() {

	if cli.rootCmd != nil {
		cli.rootCmd.ResetFlags()
		cli.rootCmd.ResetCommands()
	}

	short := fmt.Sprintf("%s is a commandline client on %s for users to interact with the %s contract",
		cli.Name, cli.blockchain.String(), strings.Join(ModeERCList, "/"))
	rootCmd := &cobra.Command{
		Use:              cli.Name, // "TokenCommander",
		Short:            short,
		Run:              cli.help,
		PersistentPreRun: cli.setup,
	}
	cli.rootCmd = rootCmd

	// Global flags
	rootCmd.PersistentFlags().StringVarP(&cli.config, "config", "c", defaultConfigFile, "The `path` to config file")
	rootCmd.PersistentFlags().StringP("walletPath", "w", defaultWalletPath, "Wallet storage `directory`")
	rootCmd.PersistentFlags().StringP("rpcURL", "i", defaultRPCURL, fmt.Sprintf("%s json rpc or ipc `url`", cli.blockchain.String()))
	rootCmd.PersistentFlags().StringP("contractAddress", "a", defaultContractAddress, "Contract `address`")
	rootCmd.PersistentFlags().StringP("from", "f", "", "the from `address` who pay gas")

	rootCmd.PersistentFlags().String("mode", ModeERC20, fmt.Sprintf(`use %s token`, strings.Join(ModeERCList, "|")))
	rootCmd.PersistentFlags().StringP("symbol", "s", "", "the symbol of the contract, this'll overwrite the `--contractAddress` when load token")

	// Basic commands
	rootCmd.AddCommand(cli.buildInitCmd())    // init
	rootCmd.AddCommand(cli.buildVersionCmd()) // version

	// deploy
	rootCmd.AddCommand(cli.buildDeployCmd())

	// account
	rootCmd.AddCommand(cli.buildAccountCmd())

	// info
	rootCmd.AddCommand(cli.buildInfoCmd())

	// pay
	rootCmd.AddCommand(cli.buildPayCmd())
	rootCmd.AddCommand(cli.buildBatchPayCmd()) // batch pay

	// balance
	rootCmd.AddCommand(cli.buildBalanceCmd())

	// ERC721
	rootCmd.AddCommand(cli.buildMintCmd()) // mint

	// add
	rootCmd.AddCommand(cli.buildAddCmd())

}
