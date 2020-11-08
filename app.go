package Node_Github

import (
	"encoding/json"
	"fmt"
	"github.com/plan-crypto/node/x/emission/keeper"
	"os"

	abci "github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"

	bam "github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/genaccounts"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/supply"
	xbank "github.com/plan-crypto/node/x/bank"
	"github.com/plan-crypto/node/x/emission"
	emissionTypes "github.com/plan-crypto/node/x/emission/types"
	pln "github.com/plan-crypto/node/x/plan"
	"github.com/plan-crypto/node/x/paramining"
	paraminingKeeper "github.com/plan-crypto/node/x/paramining/keeper"
	"github.com/plan-crypto/node/x/structure"
	structureKeeper "github.com/plan-crypto/node/x/structure/keeper"
	"runtime"
)

const appName = "plan"
const HaltHeight = 0

var (
	// default home directories for the application CLI
	DefaultCLIHome = getCliPath()

	// DefaultNodeHome sets the folder where the applcation data and configuration will be stored
	DefaultNodeHome = getNodePath()

	// NewBasicManager is in charge of setting up basic module elemnets
	ModuleBasics = module.NewBasicManager(
		genaccounts.AppModuleBasic{},
		genutil.AppModuleBasic{},
		auth.AppModuleBasic{},
		bank.AppModuleBasic{},
		staking.AppModuleBasic{},
		distr.AppModuleBasic{},
		params.AppModuleBasic{},
		slashing.AppModuleBasic{},
		supply.AppModuleBasic{},
		gov.AppModuleBasic{},
		emission.AppModuleBasic{},
		paramining.AppModuleBasic{},
		structure.AppModuleBasic{},
		pln.AppModuleBasic{},
	)

	// account permissions
	maccPerms = map[string][]string{
		auth.FeeCollectorName:     nil,
		distr.ModuleName:          nil,
		gov.ModuleName:            nil,
		staking.BondedPoolName:    []string{supply.Burner, supply.Staking},
		staking.NotBondedPoolName: []string{supply.Burner, supply.Staking},
	}
)

func getCliPath() string {
	if runtime.GOOS == "windows" {
		return os.ExpandEnv("$UserProfile/.plancli")
	}

	return os.ExpandEnv("$HOME/.plancli")
}

func getNodePath() string {
	if runtime.GOOS == "windows" {
		return os.ExpandEnv("$UserProfile/.pland")
	}

	return os.ExpandEnv("$HOME/.pland")
}

// MakeCodec generates the necessary codecs for Amino
func MakeCodec() *codec.Codec {
	var cdc = codec.New()
	ModuleBasics.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	return cdc
}

type planApp struct {
	*bam.BaseApp
	cdc *codec.Codec

	// Keys to access the substores
	keyMain          *sdk.KVStoreKey
	keyAccount       *sdk.KVStoreKey
	keySupply        *sdk.KVStoreKey
	keyStaking       *sdk.KVStoreKey
	tkeyStaking      *sdk.TransientStoreKey
	keyDistr         *sdk.KVStoreKey
	tkeyDistr        *sdk.TransientStoreKey
	keyNS            *sdk.KVStoreKey
	keyParams        *sdk.KVStoreKey
	tkeyParams       *sdk.TransientStoreKey
	keyGov           *sdk.KVStoreKey
	keyEmission      *sdk.KVStoreKey
	keySlashing      *sdk.KVStoreKey
	keyStructure     *sdk.KVStoreKey
	keyStructureFast *sdk.KVStoreKey
	keyParamining    *sdk.KVStoreKey

	// Keepers
	accountKeeper  auth.AccountKeeper
	bankKeeper     xbank.Keeper
	stakingKeeper  staking.Keeper
	slashingKeeper slashing.Keeper
	distrKeeper    distr.Keeper
	supplyKeeper   supply.Keeper
	paramsKeeper   params.Keeper
	govKeeper      gov.Keeper


	emissionKeeper   keeper.Keeper
	structureKeeper  structureKeeper.Keeper
	paraminingKeeper paraminingKeeper.Keeper
	plnKeeper       pln.Keeper

	// Module Manager
	mm *module.Manager
}

func NewPlanApp(logger log.Logger, db dbm.DB, baseAppOptions ...func(*bam.BaseApp)) *planApp {
	// First define the top level codec that will be shared by the different modules
	cdc := MakeCodec()

	// BaseApp handles interactions with Tendermint through the ABCI protocol
	bApp := bam.NewBaseApp(appName, logger, db, auth.DefaultTxDecoder(cdc), baseAppOptions...)

	// Here you initialize your application with the store keys it requires
	var app = &planApp{
		BaseApp: bApp,
		cdc:     cdc,

		keyMain:          sdk.NewKVStoreKey(bam.MainStoreKey),
		keyAccount:       sdk.NewKVStoreKey(auth.StoreKey),
		keySupply:        sdk.NewKVStoreKey(supply.StoreKey),
		keyStaking:       sdk.NewKVStoreKey(staking.StoreKey),
		tkeyStaking:      sdk.NewTransientStoreKey(staking.TStoreKey),
		keyDistr:         sdk.NewKVStoreKey(distr.StoreKey),
		tkeyDistr:        sdk.NewTransientStoreKey(distr.TStoreKey),
		keyParams:        sdk.NewKVStoreKey(params.StoreKey),
		tkeyParams:       sdk.NewTransientStoreKey(params.TStoreKey),
		keyGov:           sdk.NewKVStoreKey(gov.StoreKey),
		keyEmission:      sdk.NewKVStoreKey(emissionTypes.StoreKey),
		keySlashing:      sdk.NewKVStoreKey(slashing.StoreKey),
		keyStructure:     sdk.NewKVStoreKey(structureKeeper.StoreKey),
		keyStructureFast: sdk.NewKVStoreKey(structureKeeper.FastAccessKey),
		keyParamining:    sdk.NewKVStoreKey(paraminingKeeper.StoreKey),
	}

	// The ParamsKeeper handles parameter storage for the application
	app.paramsKeeper = params.NewKeeper(app.cdc, app.keyParams, app.tkeyParams, params.DefaultCodespace)
	// Set specific supspaces
	authSubspace := app.paramsKeeper.Subspace(auth.DefaultParamspace)
	bankSupspace := app.paramsKeeper.Subspace(bank.DefaultParamspace)
	stakingSubspace := app.paramsKeeper.Subspace(staking.DefaultParamspace)
	distrSubspace := app.paramsKeeper.Subspace(distr.DefaultParamspace)
	slashingSubspace := app.paramsKeeper.Subspace(slashing.DefaultParamspace)
	govSubspace := app.paramsKeeper.Subspace(gov.DefaultParamspace)


	// The AccountKeeper handles address -> account lookups
	app.accountKeeper = auth.NewAccountKeeper(
		app.cdc,
		app.keyAccount,
		authSubspace,
		auth.ProtoBaseAccount,
	)

	// The BankKeeper allows you perform sdk.Coins interactions
	app.bankKeeper = xbank.NewBankKeeper(
		app.accountKeeper,
		bankSupspace,
		bank.DefaultCodespace,
	)

	// The SupplyKeeper collects transaction fees and renders them to the fee distribution module
	app.supplyKeeper = supply.NewKeeper(
		app.cdc,
		app.keySupply,
		app.accountKeeper,
		app.bankKeeper,
		supply.DefaultCodespace,
		maccPerms,
	)

	// The staking keeper
	stakingKeeper := staking.NewKeeper(
		app.cdc,
		app.keyStaking,
		app.tkeyStaking,
		app.supplyKeeper,
		stakingSubspace,
		staking.DefaultCodespace,
	)

	app.distrKeeper = distr.NewKeeper(
		app.cdc,
		app.keyDistr,
		distrSubspace,
		&stakingKeeper,
		app.supplyKeeper,
		distr.DefaultCodespace,
		auth.FeeCollectorName,
	)

	app.slashingKeeper = slashing.NewKeeper(
		app.cdc,
		app.keySlashing,
		&stakingKeeper,
		slashingSubspace,
		slashing.DefaultCodespace,
	)

	// register the staking hooks
	// NOTE: stakingKeeper above is passed by reference, so that it will contain these hooks
	app.stakingKeeper = *stakingKeeper.SetHooks(
		staking.NewMultiStakingHooks(
			app.distrKeeper.Hooks(),
			app.slashingKeeper.Hooks(),
			app.paraminingKeeper.SlashingHooks(),
			app.emissionKeeper.SlashingHooks(),
		),
	)

	// Since we cannot do that during the keepers creation
	app.bankKeeper.StakingKeeper = app.stakingKeeper

	// Keeper that handles all the structures related stuff
	app.structureKeeper = structureKeeper.NewKeeper(
		app.cdc,
		app.keyStructure,
		app.keyStructureFast,
	)

	app.emissionKeeper = keeper.NewKeeper(
		app.cdc,
		app.keyEmission,
		app.stakingKeeper,
	)

	// Keeper that handles all the paramining related stuff
	app.paraminingKeeper = paraminingKeeper.NewKeeper(
		app.cdc,
		app.keyParamining,
		app.bankKeeper,
		app.stakingKeeper,
		app.emissionKeeper,
	)

	// Helping keeper that's mostly being used for the API calls
	app.plnKeeper = pln.NewKeeper(
		app.cdc,
		app.accountKeeper,
		app.bankKeeper,
		app.structureKeeper,
		app.paraminingKeeper,
		app.emissionKeeper,
		app.supplyKeeper,
		app.slashingKeeper,
	)

	govRouter := gov.NewRouter()

	govRouter.AddRoute(gov.RouterKey, gov.ProposalHandler).
		AddRoute(params.RouterKey, params.NewParamChangeProposalHandler(app.paramsKeeper)).
		AddRoute(distr.RouterKey, distr.NewCommunityPoolSpendProposalHandler(app.distrKeeper))

	app.govKeeper = gov.NewKeeper(
		app.cdc,
		app.keyGov,
		app.paramsKeeper,
		govSubspace,
		app.supplyKeeper,
		stakingKeeper,
		gov.DefaultCodespace,
		govRouter,
	)

	// Structure changes hooks
	app.structureKeeper.AddStructureChangedHook(app.paraminingKeeper.GenerateStructureChangedHook())

	// Paramining charged hooks
	app.paraminingKeeper.AddParaminingChargedHook(app.structureKeeper.GenerateParaminingChargedHook())
	app.paraminingKeeper.AddParaminingChargedHook(app.emissionKeeper.GenerateParaminingChargedHook())

	// Cons transfer hooks
	app.bankKeeper.AddBeforeHook(app.paraminingKeeper.GenerateBeforeTransferHook())
	app.bankKeeper.AddAfterHook(app.paraminingKeeper.GenerateAfterTransferHook())
	app.bankKeeper.AddAfterHook(app.structureKeeper.GenerateAfterTransferHook())

	app.mm = module.NewManager(
		genaccounts.NewAppModule(app.accountKeeper),
		genutil.NewAppModule(app.accountKeeper, app.stakingKeeper, app.BaseApp.DeliverTx),
		auth.NewAppModule(app.accountKeeper),
		bank.NewAppModule(app.bankKeeper, app.accountKeeper),
		paramining.NewAppModule(app.paraminingKeeper),
		emission.NewAppModule(app.emissionKeeper),
		pln.NewAppModule(app.plnKeeper, app.bankKeeper),
		structure.NewAppModule(app.structureKeeper),
		supply.NewAppModule(app.supplyKeeper, app.accountKeeper),
		distr.NewAppModule(app.distrKeeper, app.supplyKeeper),
		slashing.NewAppModule(app.slashingKeeper, app.stakingKeeper),
		staking.NewAppModule(app.stakingKeeper, app.distrKeeper, app.accountKeeper, app.supplyKeeper),
		gov.NewAppModule(app.govKeeper, app.supplyKeeper),
	)

	app.mm.SetOrderBeginBlockers(distr.ModuleName, slashing.ModuleName)
	app.mm.SetOrderEndBlockers(staking.ModuleName)

	// Sets the order of Genesis - Order matters, genutil is to always come last
	app.mm.SetOrderInitGenesis(
		genaccounts.ModuleName,
		distr.ModuleName,
		staking.ModuleName,
		auth.ModuleName,
		bank.ModuleName,
		supply.ModuleName,
		gov.ModuleName,
		emission.ModuleName,
		paramining.ModuleName,
		structure.ModuleName,
		pln.ModuleName,
		slashing.ModuleName,
		genutil.ModuleName,
	)

	// register all module routes and module queriers
	app.mm.RegisterRoutes(app.Router(), app.QueryRouter())

	// The initChainer handles translating the genesis.json file into initial state for the network
	app.SetInitChainer(app.InitChainer)
	app.SetBeginBlocker(app.BeginBlocker)
	app.SetEndBlocker(app.EndBlocker)

	// The AnteHandler handles signature verification and transaction pre-processing
	app.SetAnteHandler(
		auth.NewAnteHandler(
			app.accountKeeper,
			app.supplyKeeper,
			auth.DefaultSigVerificationGasConsumer,
		),
	)

	app.MountStores(
		app.keyMain,
		app.keyAccount,
		app.keySupply,
		app.keyStaking,
		app.tkeyStaking,
		app.keyDistr,
		app.tkeyDistr,
		app.keySlashing,
		app.keyParams,
		app.tkeyParams,
		app.keyGov,
		app.keyEmission,
		app.keyStructure,
		app.keyStructureFast,
		app.keyParamining,
	)


	err := app.LoadLatestVersion(app.keyMain)

	if err != nil {
		cmn.Exit(err.Error())
	}

	if HaltHeight != 0 && app.LastBlockHeight() == HaltHeight {
		cmn.Exit(fmt.Sprintf("Reached the halt height %d", HaltHeight))
	}

	return app
}

// GenesisState represents chain state at the start of the chain. Any initial state (account balances) are stored here.
type GenesisState map[string]json.RawMessage

func NewDefaultGenesisState() GenesisState {
	return ModuleBasics.DefaultGenesis()
}

func (app *planApp) InitChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	var genesisState GenesisState

	err := app.cdc.UnmarshalJSON(req.AppStateBytes, &genesisState)
	if err != nil {
		panic(err)
	}

	return app.mm.InitGenesis(ctx, genesisState)
}

func (app *planApp) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	return app.mm.BeginBlock(ctx, req)
}
func (app *planApp) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	return app.mm.EndBlock(ctx, req)
}
func (app *planApp) LoadHeight(height int64) error {
	return app.LoadVersion(height, app.keyMain)
}

//_________________________________________________________

func (app *planApp) ExportAppStateAndValidators(forZeroHeight bool, jailWhiteList []string,
) (appState json.RawMessage, validators []tmtypes.GenesisValidator, err error) {

	// as if they could withdraw from the start of the next block
	ctx := app.NewContext(true, abci.Header{Height: app.LastBlockHeight()})

	genState := app.mm.ExportGenesis(ctx)
	appState, err = codec.MarshalJSONIndent(app.cdc, genState)
	if err != nil {
		return nil, nil, err
	}

	validators = staking.WriteValidators(ctx, app.stakingKeeper)

	return appState, validators, nil
}
