package plan

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/plan-crypto/node/x/emission/keeper"
	paraminingKeeper "github.com/plan-crypto/node/x/paramining/keeper"
	structureKeeper "github.com/plan-crypto/node/x/structure/keeper"
)

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	accountKeeper    auth.AccountKeeper
	coinKeeper       bank.Keeper
	structureKeeper  structureKeeper.Keeper
	paraminingKeeper paraminingKeeper.Keeper
	emissionKeeper   keeper.Keeper
	supplyKeeper     supply.Keeper
	slashingKeeper   slashing.Keeper


	cdc *codec.Codec // The wire codec for binary encoding/decoding.
}

func NewKeeper(cdc *codec.Codec, accountKeeper auth.AccountKeeper, coinKeeper bank.Keeper, structureKeeper structureKeeper.Keeper, paraminingKeeper paraminingKeeper.Keeper, emissionKeeper keeper.Keeper, supplyKeeper supply.Keeper, slashingKeeper slashing.Keeper) Keeper {
	return Keeper{
		cdc:              cdc,
		accountKeeper:    accountKeeper,
		coinKeeper:       coinKeeper,
		structureKeeper:  structureKeeper,
		paraminingKeeper: paraminingKeeper,
		emissionKeeper:   emissionKeeper,
		supplyKeeper:     supplyKeeper,
		slashingKeeper:     slashingKeeper,
	}
}
