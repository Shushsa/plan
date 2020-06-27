package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
)


type Keeper struct {
	stakingKeeper staking.Keeper

	storeKey      sdk.StoreKey

	cdc *codec.Codec
}

func NewKeeper(cdc *codec.Codec, storeKey sdk.StoreKey, stakingKeeper staking.Keeper) Keeper {
	return Keeper{
		storeKey:      storeKey,
		cdc:           cdc,
		stakingKeeper:           stakingKeeper,
	}
}


// Checks if the threshold has been reached - in that case, we won't do paramining
func (k Keeper) IsThresholdReached(ctx sdk.Context) bool {
	return k.GetEmission(ctx).IsThresholdReached()
}

// Adding new coins to emission
func (k Keeper) Add(ctx sdk.Context, amount sdk.Int) {
	emission := k.GetEmission(ctx)
	emission.Current = emission.Current.Add(amount)

	k.SetEmission(ctx, emission)
}

// Remove coins from emission
func (k Keeper) Sub(ctx sdk.Context, amount sdk.Int) {
	emission := k.GetEmission(ctx)
	emission.Current = emission.Current.Sub(amount)

	k.SetEmission(ctx, emission)
}

// Remove coins from emission before slashing
func (k Keeper) UpdateBeforeSlashing(ctx sdk.Context, valAddr sdk.ValAddress, fraction sdk.Dec) {
	validator, found := k.stakingKeeper.GetValidator(ctx, valAddr)

	if !found {
		panic("Validator does not exist but got slashed")
	}

	amount := validator.GetTokens()

	slashAmountDec := sdk.NewInt(amount.ToDec().Mul(fraction).Int64())

	k.Sub(ctx, slashAmountDec)
}