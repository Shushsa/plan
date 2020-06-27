package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	ouroTypes "github.com/ouroboros-crypto/node/x/ouroboros/types"
)

// This hook should be called when we change the structure balance
type StructureChangedHook = func(ctx sdk.Context, addr sdk.AccAddress, currentBalance sdk.Int, previousBalance sdk.Int)

func (k *Keeper) AddStructureChangedHook(hook StructureChangedHook) {
	k.structureChangedHooks = append(k.structureChangedHooks, hook)
}

// Generates a hook that would be called after moving the coins from one address to another
func (k Keeper) GenerateAfterTransferHook() func(ctx sdk.Context, from sdk.AccAddress, to sdk.AccAddress, amn sdk.Coins) {
	return func(ctx sdk.Context, sender sdk.AccAddress, receiver sdk.AccAddress, amn sdk.Coins) {
		// If the receive is already in someone's structure, we have to get through the whole structure
		coinsAmount := amn.AmountOf(ouroTypes.OURO)

		// If the receiver isn't genesis and the receiver is already in some structure
		if receiver.String() != ouroTypes.GenesisWallet && !k.AddToStructure(ctx, sender, receiver, coinsAmount) {
			// Taking the coins out of the sender's upper structure
			k.DecreaseStructureBalance(ctx, sender, coinsAmount)

			// And add them to the receiver's upper structure
			k.IncreaseStructureBalance(ctx, receiver, coinsAmount)
		}
	}
}

// Generates a hook that would be called after charging paramined coins
func (k Keeper) GenerateParaminingChargedHook() func(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Int) {
	return func(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Int) {
		// Increase structure balance
		k.IncreaseStructureBalance(ctx, addr, amt)
	}
}
