package paramining

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/Shushsa/plan/x/paramining/keeper"
	"github.com/Shushsa/plan/x/paramining/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

type GenesisState struct {
	ParaminingRecords []types.Paramining `json:"paramining_records"`
}

// Genesis initialization
func NewGenesisState() GenesisState {
	return GenesisState{
		ParaminingRecords: []types.Paramining{},
	}
}

// Genesis validation
func ValidateGenesis(data GenesisState) error {
	return nil
}

// Genesis default state
func DefaultGenesisState() GenesisState {
	return GenesisState{
		ParaminingRecords: []types.Paramining{},
	}
}

// Genesis init
func InitGenesis(ctx sdk.Context, k keeper.Keeper, data GenesisState) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}

// Genesis export
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) GenesisState {
	var paraminingRecords []types.Paramining

	iterator := k.GetParaminingIterator(ctx)

	for ; iterator.Valid(); iterator.Next() {
		addr := sdk.AccAddress(iterator.Key())

		paramining := k.GetParamining(ctx, addr)
		paraminingRecords = append(paraminingRecords, paramining)
	}

	return GenesisState{ParaminingRecords: paraminingRecords}
}
