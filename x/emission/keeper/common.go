package keeper


import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/plan-crypto/node/x/emission/types"
)

// Returns the emission record
func (k Keeper) GetEmission(ctx sdk.Context) types.Emission {
	store := ctx.KVStore(k.storeKey)

	if !store.Has([]byte(StoreKey)) {
		return types.NewEmission()
	}

	var emission types.Emission

	k.cdc.MustUnmarshalBinaryBare(store.Get([]byte(StoreKey)), &emission)

	return emission
}

// Saves the emission record
func (k Keeper) SetEmission(ctx sdk.Context, emission types.Emission) {
	store := ctx.KVStore(k.storeKey)

	store.Set([]byte(StoreKey), k.cdc.MustMarshalBinaryBare(emission))
}