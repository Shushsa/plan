package paramining

import (
	"encoding/json"

	"github.com/Shushsa/plan/x/paramining/keeper"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"

	"github.com/Shushsa/plan/x/paramining/client/cli"
	"github.com/Shushsa/plan/x/paramining/client/rest"

	ptypes "github.com/Shushsa/plan/x/paramining/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdk_module "github.com/cosmos/cosmos-sdk/types/module"
	abci "github.com/tendermint/tendermint/abci/types"
)

var (
	_ sdk_module.AppModule      = AppModule{}
	_ sdk_module.AppModuleBasic = AppModuleBasic{}
)

const ModuleName = "paramining"

// app module Basics object
type AppModuleBasic struct{}

func (AppModuleBasic) Name() string {
	return ModuleName
}

func (AppModuleBasic) RegisterCodec(cdc *codec.Codec) {
	ptypes.RegisterCodec(cdc)
}

func (AppModuleBasic) DefaultGenesis() json.RawMessage {
	return ptypes.ModuleCdc.MustMarshalJSON(DefaultGenesisState())
}

// Validation check of the Genesis
func (AppModuleBasic) ValidateGenesis(bz json.RawMessage) error {
	var data GenesisState
	err := ptypes.ModuleCdc.UnmarshalJSON(bz, &data)
	if err != nil {
		return err
	}
	// once json successfully marshalled, passes along to genesis.go
	return ValidateGenesis(data)
}

// Register rest routes
func (AppModuleBasic) RegisterRESTRoutes(ctx context.CLIContext, rtr *mux.Router) {
	rest.RegisterRoutes(ctx, rtr)
}

// Get the root query command of this module
func (AppModuleBasic) GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	return cli.GetQueryCmd(cdc)
}

// Get the root tx command of this module
func (AppModuleBasic) GetTxCmd(cdc *codec.Codec) *cobra.Command {
	return cli.GetTxCmd(cdc)
}

type AppModule struct {
	AppModuleBasic
	keeper keeper.Keeper
}

// NewAppModule creates a new AppModule Object
func NewAppModule(k keeper.Keeper) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{},
		keeper:         k,
	}
}

func (AppModule) Name() string {
	return ModuleName
}

func (am AppModule) RegisterInvariants(ir sdk.InvariantRegistry) {}

func (am AppModule) Route() string {
	return ModuleName
}

func (am AppModule) NewHandler() types.Handler {
	return keeper.NewHandler(am.keeper)
}

func (am AppModule) QuerierRoute() string {
	return ModuleName
}

func (am AppModule) NewQuerierHandler() types.Querier {
	return keeper.NewQuerier(am.keeper)
}

func (am AppModule) BeginBlock(_ sdk.Context, _ abci.RequestBeginBlock) {}

func (am AppModule) EndBlock(types.Context, abci.RequestEndBlock) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}

func (am AppModule) InitGenesis(ctx types.Context, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState GenesisState
	ptypes.ModuleCdc.MustUnmarshalJSON(data, &genesisState)
	return InitGenesis(ctx, am.keeper, genesisState)
}

func (am AppModule) ExportGenesis(ctx sdk.Context) json.RawMessage {
	gs := ExportGenesis(ctx, am.keeper)
	return ptypes.ModuleCdc.MustMarshalJSON(gs)
}
