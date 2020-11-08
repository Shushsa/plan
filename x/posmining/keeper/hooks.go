package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/plan-crypto/node/x/paramining/types"
)

// When the keeper charges paramining
type ParaminingChargedHook = func(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Int)

// Adds new paramining charged hook
func (k *Keeper) AddParaminingChargedHook(hook ParaminingChargedHook) {
	k.paraminingChargedHooks = append(k.paraminingChargedHooks, hook)
}

// Cal it after charging paramining
func (k Keeper) afterParaminingCharged(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Int) {
	for _, hook := range k.paraminingChargedHooks {
		hook(ctx, addr, amt)
	}
}

// Generates a hook that would be called before moving the coins from one address to another
func (k Keeper) GenerateBeforeTransferHook() func(ctx sdk.Context, from sdk.AccAddress, to sdk.AccAddress, amn sdk.Coins) {
	return func(ctx sdk.Context, sender sdk.AccAddress, receiver sdk.AccAddress, amn sdk.Coins) {
		// The sender should get paramining on his balance
		k.ChargeParamining(ctx, sender, false)

		// The received should save the paramined amount since his percent may change
		k.SaveParamined(ctx, receiver)
	}
}

// Generates a hook that would be called after moving the coins from one address to another
func (k Keeper) GenerateAfterTransferHook() func(ctx sdk.Context, from sdk.AccAddress, to sdk.AccAddress, amn sdk.Coins) {
	return func(ctx sdk.Context, sender sdk.AccAddress, receiver sdk.AccAddress, amn sdk.Coins) {
		// Check if we should change the daily percents
		k.UpdateDailyPercent(ctx, sender)

		k.UpdateDailyPercent(ctx, receiver)
	}
}


// Generates a hook that will be called when somebody changes the structure balance
func (k Keeper) GenerateStructureChangedHook() func(ctx sdk.Context, addr sdk.AccAddress, currentBalance sdk.Int, previousBalance sdk.Int) {
	return func(ctx sdk.Context, addr sdk.AccAddress, currentBalance sdk.Int, previousBalance sdk.Int) {
		// To avoid extra fetching of the paramining record
		currentStructureCoff := types.GetStructureCoff(currentBalance)

		if !currentStructureCoff.Equal(types.GetStructureCoff(previousBalance))  {
			// First save already paramined since the formula will change
			k.SaveParamined(ctx, addr)

			// Update paramining record
			paramining := k.GetParamining(ctx, addr)
			paramining.StructureCoff = currentStructureCoff

			k.SetParamining(ctx, paramining)
		}
	}
}

//_________________________________________________________________________________________

// Slashing hooks
type Hooks struct {
	k Keeper
}

var _ stakingtypes.StakingHooks = Hooks{}

// Create new distribution hooks
func (k Keeper) SlashingHooks() Hooks { return Hooks{k} }

// We should save paramining of every delegator before validator gets slashed
func (h Hooks) BeforeValidatorSlashed(ctx sdk.Context, valAddr sdk.ValAddress, fraction sdk.Dec) {
	h.k.UpdateDelegatorsBeforeSlashing(ctx, valAddr)
}


// nolint - unused hooks
func (h Hooks) BeforeValidatorModified(_ sdk.Context, _ sdk.ValAddress)                         {}
func (h Hooks) AfterValidatorBonded(_ sdk.Context, _ sdk.ConsAddress, _ sdk.ValAddress)         {}
func (h Hooks) AfterValidatorBeginUnbonding(_ sdk.Context, _ sdk.ConsAddress, _ sdk.ValAddress) {}
func (h Hooks) BeforeDelegationRemoved(_ sdk.Context, _ sdk.AccAddress, _ sdk.ValAddress)       {}
func (h Hooks) AfterValidatorCreated(ctx sdk.Context, valAddr sdk.ValAddress) {}
func (h Hooks) AfterValidatorRemoved(ctx sdk.Context, _ sdk.ConsAddress, valAddr sdk.ValAddress) { }
func (h Hooks) BeforeDelegationCreated(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) {}
func (h Hooks) BeforeDelegationSharesModified(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) {}
func (h Hooks) AfterDelegationModified(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) { }
