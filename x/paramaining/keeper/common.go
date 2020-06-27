package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/Shushsa/plan/x/paramining/types"
)

// Fetches a paramining record by the owner - if one doesn't exist, it'll create a new one
func (k Keeper) GetParamining(ctx sdk.Context, owner sdk.AccAddress) types.Paramining {
	store := ctx.KVStore(k.storeKey)

	if !store.Has(owner.Bytes()) {
		newParamining := types.NewParamining(owner)

		newParamining.LastTransaction = ctx.BlockHeader().Time
		newParamining.LastCharged = ctx.BlockHeader().Time

		return newParamining
	}

	var upperStructure types.Paramining

	k.cdc.MustUnmarshalBinaryBare(store.Get(owner.Bytes()), &upperStructure)

	return upperStructure
}

// Saves the paramining record
func (k Keeper) SetParamining(ctx sdk.Context, paramining types.Paramining) {
	store := ctx.KVStore(k.storeKey)

	store.Set(paramining.Owner.Bytes(), k.cdc.MustMarshalBinaryBare(paramining))
}

// Returns an iterator that allows to iterate over the records
func (k Keeper) GetParaminingIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)

	return sdk.KVStorePrefixIterator(store, nil)
}
