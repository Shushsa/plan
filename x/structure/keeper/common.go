package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ouroboros-crypto/node/x/structure/types"
)

// Returns "upper structure" that's just a pointer to the next structure above
func (k Keeper) GetUpperStructure(ctx sdk.Context, address sdk.AccAddress) types.UpperStructure {
	store := ctx.KVStore(k.fastAccessKey)

	if !store.Has(address.Bytes()) {
		return types.NewUpperStructure(address)
	}

	var upperStructure types.UpperStructure

	k.cdc.MustUnmarshalBinaryBare(store.Get(address.Bytes()), &upperStructure)

	return upperStructure
}

// Saves pointer to the upper account in the sturcture
func (k Keeper) SetUpperStructure(ctx sdk.Context, address sdk.AccAddress, upperStructure types.UpperStructure) {
	store := ctx.KVStore(k.fastAccessKey)

	store.Set(address.Bytes(), k.cdc.MustMarshalBinaryBare(upperStructure))
}

// Checks if the user in any structure yet
func (k Keeper) HasUpperStructure(ctx sdk.Context, address sdk.AccAddress) bool {
	return !k.GetUpperStructure(ctx, address).Owner.Empty()
}

// Iterator
func (k Keeper) GetUpperStructureIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.fastAccessKey)

	return sdk.KVStorePrefixIterator(store, nil)
}

// Returns the structure record by its owner
func (k Keeper) GetStructure(ctx sdk.Context, owner sdk.AccAddress) types.Structure {
	store := ctx.KVStore(k.storeKey)

	if !store.Has(owner.Bytes()) {
		return types.NewStructure(owner)
	}

	var structure types.Structure

	k.cdc.MustUnmarshalBinaryBare(store.Get(owner.Bytes()), &structure)

	return structure
}

// Saves the sturcture record
func (k Keeper) SetStructure(ctx sdk.Context, structure types.Structure) {
	store := ctx.KVStore(k.storeKey)

	store.Set(structure.Owner.Bytes(), k.cdc.MustMarshalBinaryBare(structure))
}

// Iterator
func (k Keeper) GetStructureIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)

	return sdk.KVStorePrefixIterator(store, nil)
}
