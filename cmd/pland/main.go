package main

import (
	"encoding/json"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/plan-crypto/node"
	"github.com/spf13/viper"
	"io"

	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/x/genaccounts"
	genaccscli "github.com/cosmos/cosmos-sdk/x/genaccounts/client/cli"
	"github.com/cosmos/cosmos-sdk/x/staking"

	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/libs/log"

	sdk "github.com/cosmos/cosmos-sdk/types"
	genutilcli "github.com/cosmos/cosmos-sdk/x/genutil/client/cli"
	addrs "github.com/plan-crypto/node/x/plan/types"

	abci "github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tendermint/libs/db"
	tmtypes "github.com/tendermint/tendermint/types"
)

func main() {
	cobra.EnableCommandSorting = false

	cdc := Node_Github.MakeCodec()

	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(addrs.AccAddr, addrs.AccPub)
	config.SetBech32PrefixForValidator(addrs.ValAddr, addrs.ValPub)
	config.SetBech32PrefixForConsensusNode(addrs.ConsAddr, addrs.ConsPub)
	config.Seal()

	ctx := server.NewDefaultContext()

	rootCmd := &cobra.Command{
		Use:               "pland",
		Short:             "Plan App Daemon (server)",
		PersistentPreRunE: server.PersistentPreRunEFn(ctx),
	}

	// CLI commands to initialize the chain
	rootCmd.AddCommand(
		genutilcli.InitCmd(ctx, cdc, Node_Github.ModuleBasics, Node_Github.DefaultNodeHome),
		genutilcli.CollectGenTxsCmd(ctx, cdc, genaccounts.AppModuleBasic{}, Node_Github.DefaultNodeHome),
		genutilcli.GenTxCmd(ctx, cdc, Node_Github.ModuleBasics, staking.AppModuleBasic{}, genaccounts.AppModuleBasic{}, Node_Github.DefaultNodeHome, Node_Github.DefaultCLIHome),
		genutilcli.ValidateGenesisCmd(ctx, cdc, Node_Github.ModuleBasics),
		// AddGenesisAccountCmd allows users to add accounts to the genesis file
		genaccscli.AddGenesisAccountCmd(ctx, cdc, Node_Github.DefaultNodeHome, Node_Github.DefaultCLIHome),
	)

	server.AddCommands(ctx, cdc, rootCmd, newApp, exportAppStateAndTMValidators)

	// prepare and add flags
	executor := cli.PrepareBaseCmd(rootCmd, "PLN", Node_Github.DefaultNodeHome)
	err := executor.Execute()
	if err != nil {
		panic(err)
	}
}

func newApp(logger log.Logger, db dbm.DB, traceStore io.Writer) abci.Application {
	if Node_Github.HaltHeight != 0 {
		return Node_Github.NewPlanApp(logger, db, baseapp.SetHaltHeight(uint64(Node_Github.HaltHeight)))
	}

	return Node_Github.NewPlanApp(logger, db, baseapp.SetHaltHeight(uint64(viper.GetInt(server.FlagHaltHeight))))
}

func exportAppStateAndTMValidators(
	logger log.Logger, db dbm.DB, traceStore io.Writer, height int64, forZeroHeight bool, jailWhiteList []string,
) (json.RawMessage, []tmtypes.GenesisValidator, error) {

	if height != -1 {
		plnApp := Node_Github.NewPlanApp(logger, db)
		err := plnApp.LoadHeight(height)
		if err != nil {
			return nil, nil, err
		}
		return plnApp.ExportAppStateAndValidators(forZeroHeight, jailWhiteList)
	}

	plnApp := Node_Github.NewPlanApp(logger, db)

	return plnApp.ExportAppStateAndValidators(forZeroHeight, jailWhiteList)
}