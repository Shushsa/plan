package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/Shushsa/plan/x/bank"                        // ouro bank folder not found
	"github.com/Shushsa/plan/x/emission/keeper"             // ouro emission folder not found
	planTypes "github.com/Shushsa/plan/x/nameservice/types" // nameservice folder
	"github.com/Shushsa/plan/x/paramining/types"
)

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	storeKey sdk.StoreKey // Unexposed key to access store from sdk.Context

	coinKeeper     bank.Keeper
	stakingKeeper  staking.Keeper
	emissionKeeper keeper.Keeper

	cdc *codec.Codec // The wire codec for binary encoding/decoding.

	paraminingChargedHooks []ParaminingChargedHook
}

func NewKeeper(cdc *codec.Codec, storeKey sdk.StoreKey, coinKeeper bank.Keeper, stakingKeeper staking.Keeper, emissionKeeper keeper.Keeper) Keeper {
	return Keeper{
		storeKey:               storeKey,
		coinKeeper:             coinKeeper,
		stakingKeeper:          stakingKeeper,
		emissionKeeper:         emissionKeeper,
		cdc:                    cdc,
		paraminingChargedHooks: make([]ParaminingChargedHook, 0),
	}
}

// Updates daily percent based on the paraminable balance
func (k Keeper) UpdateDailyPercent(ctx sdk.Context, addr sdk.AccAddress) {
	paraminableBalance := k.coinKeeper.GetParaminableBalance(ctx, addr)

	paramining := k.GetParamining(ctx, addr)

	newDailyPercent := types.GetDailyPercent(paraminableBalance.AmountOf(planTypes.PLAN))

	if !paramining.DailyPercent.Equal(newDailyPercent) {
		paramining.DailyPercent = newDailyPercent

		k.SetParamining(ctx, paramining)
	}
}

// Charges the paramining and resets the "lastCharged" field
func (k Keeper) ChargeParamining(ctx sdk.Context, addr sdk.AccAddress, isReinvest bool) sdk.Int {
	paraminableBalance := k.coinKeeper.GetParaminableBalance(ctx, addr)

	paramining := k.GetParamining(ctx, addr)

	coinsParamined := k.CalculateParamined(ctx, paramining, paraminableBalance)

	// Reset the fields
	paramining.Paramined = sdk.NewInt(0)

	// Since reinvest doesn't reset "last transaction" field
	if !isReinvest {
		paramining.LastTransaction = ctx.BlockHeader().Time
	}

	paramining.LastCharged = ctx.BlockHeader().Time

	k.SetParamining(ctx, paramining)

	// If we charged at least 0.000001 plan
	if coinsParamined.IsPositive() {
		_, err := k.coinKeeper.AddCoins(ctx, addr, planTypes.NewCoins(coinsParamined))

		if err != nil {
			panic(err)
		}

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeParaminingCharged,
				sdk.NewAttribute(sdk.AttributeKeySender, addr.String()),
				sdk.NewAttribute(AttributeKeyAmount, coinsParamined.String()),
			),
		)

		k.afterParaminingCharged(ctx, addr, coinsParamined)
	}

	return coinsParamined
}

// Saves the paramined without charing it
func (k Keeper) SaveParamined(ctx sdk.Context, addr sdk.AccAddress) sdk.Int {
	paraminableBalance := k.coinKeeper.GetParaminableBalance(ctx, addr)
	paramining := k.GetParamining(ctx, addr)
	paramined := k.CalculateParamined(ctx, paramining, paraminableBalance)

	paramining.Paramined = paramined
	paramining.LastCharged = ctx.BlockHeader().Time

	k.SetParamining(ctx, paramining)

	return paramined
}

// We need to save delegators paramining before they get slashed
func (k Keeper) UpdateDelegatorsBeforeSlashing(ctx sdk.Context, valAddr sdk.ValAddress) {
	delegations := k.stakingKeeper.GetValidatorDelegations(ctx, valAddr)

	for _, delegation := range delegations {
		k.SaveParamined(ctx, delegation.DelegatorAddress)
	}
}

// Resolves paramining, so we can get that data via API
func (k Keeper) GetParaminingResolve(ctx sdk.Context, owner sdk.AccAddress) types.ParaminingResolve {
	paraminableBalance := k.coinKeeper.GetParaminableBalance(ctx, owner)

	paramining := k.GetParamining(ctx, owner)

	return types.ParaminingResolve{
		Paramining:   paramining,
		Paramined:    k.CalculateParamined(ctx, paramining, paraminableBalance),
		CoinsPerTime: k.CalculateCoinsPerTime(ctx, paramining, k.GetCurrentSavingsCoff(ctx, paramining), paraminableBalance),
	}
}
