package keeper

import (
	"github.com/Shushsa/plan/x/coins"
	"github.com/Shushsa/plan/x/posmining/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Saves the regulation record
func (k Keeper) SetCorrection(ctx sdk.Context, regulation types.Correction) {
	store := ctx.KVStore(k.storeKey)

	store.Set([]byte("correction"), k.Cdc.MustMarshalBinaryBare(regulation))
}

// Gets the regulation record
func (k Keeper) GetCorrection(ctx sdk.Context) types.Correction {
	store := ctx.KVStore(k.storeKey)

	var regulation types.Correction

	k.Cdc.MustUnmarshalBinaryBare(store.Get([]byte("correction")), &regulation)

	return regulation
}

// Fetches a posmining record by the owner and the coin - if one doesn't exist, it'll create a new one
func (k Keeper) GetPosmining(ctx sdk.Context, owner sdk.AccAddress, coin coins.Coin) types.Posmining {
	store := ctx.KVStore(k.storeKey)
	key := owner.Bytes()

	if !coin.Default {
		key = []byte(coin.Symbol + owner.String())
	}

	if !store.Has(key) {
		newPosmining := types.NewPosmining(owner)

		newPosmining.LastTransaction = ctx.BlockHeader().Time
		newPosmining.LastCharged = ctx.BlockHeader().Time

		return newPosmining
	}

	var posmining types.Posmining

	k.Cdc.MustUnmarshalBinaryBare(store.Get(key), &posmining)

	return posmining
}


// Returns an iterator that allows to iterate over the records
func (k Keeper) GetPosminingIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)

	return sdk.KVStorePrefixIterator(store, nil)
}

// Saves the posmining record
func (k Keeper) SetPosmining(ctx sdk.Context, posmining types.Posmining, coin coins.Coin) {
	store := ctx.KVStore(k.storeKey)

	key := posmining.Owner.Bytes()

	if !coin.Default {
		key = []byte(coin.Symbol + posmining.Owner.String())
	}

	store.Set(key, k.Cdc.MustMarshalBinaryBare(posmining))
}
