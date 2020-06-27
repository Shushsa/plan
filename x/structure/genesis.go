package structure

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ouroboros-crypto/node/x/structure/keeper"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/ouroboros-crypto/node/x/structure/types"
)

type GenesisState struct {
	UpperStructureRecords []types.UpperStructure `json:"upper_structure_records"`
	StructureRecords []types.Structure `json:"structure_records"`
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
	return GenesisState{
		UpperStructureRecords: []types.UpperStructure{},
		StructureRecords: []types.Structure{},
	}
}

// Init from state
func InitGenesis(ctx sdk.Context, k  keeper.Keeper, data GenesisState) []abci.ValidatorUpdate {
	for _, record := range data.UpperStructureRecords {
		k.SetUpperStructure(ctx, record.Address, record)
	}

	for _, record := range data.StructureRecords {
		k.SetStructure(ctx, record)
	}

	return []abci.ValidatorUpdate{}
}

// Export
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) GenesisState {
	var structureRecords []types.Structure
	var upperStructureRecords []types.UpperStructure

	iterator := k.GetUpperStructureIterator(ctx)

	for ; iterator.Valid(); iterator.Next() {
		addr := sdk.AccAddress(iterator.Key())

		upperStructure := k.GetUpperStructure(ctx, addr)
		upperStructureRecords = append(upperStructureRecords, upperStructure)
	}

	iterator = k.GetStructureIterator(ctx)

	for ; iterator.Valid(); iterator.Next() {
		addr := sdk.AccAddress(iterator.Key())

		structure := k.GetStructure(ctx, addr)
		structureRecords = append(structureRecords, structure)
	}

	return GenesisState{UpperStructureRecords: upperStructureRecords, StructureRecords: structureRecords}
}
