package emission

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/Shushsa/plan/x/emission/keeper"
	"github.com/Shushsa/plan/x/emission/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

type GenesisState struct {
}

// Init
func NewGenesisState() GenesisState {
	return GenesisState{}
}

// Validate
func ValidateGenesis(data GenesisState) error {
	return nil
}

// Default state
func DefaultGenesisState() GenesisState {
	return GenesisState{}
}

// Init from state - just set the emission
func InitGenesis(ctx sdk.Context, k keeper.Keeper, data GenesisState) []abci.ValidatorUpdate {
	// Стартовая эмиссия
	k.SetEmission(ctx, types.NewEmission())

	return []abci.ValidatorUpdate{}
}

// Export
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) GenesisState {
	return GenesisState{}
}
