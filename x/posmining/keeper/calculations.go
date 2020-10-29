package keeper

import (
	"math"
	"time"

	"github.com/Shushsa/plan/x/posmining/types"

	"github.com/Shushsa/plan/x/coins"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Returns a list of saving posmining periods
func (k Keeper) GetSavingPeriods(ctx sdk.Context, posmining types.Posmining) []types.PosminingPeriod {
	var daysSeparator int64 = 2592000

	lastTx := posmining.LastTransaction
	secondsDiff := int64(ctx.BlockTime().Sub(lastTx).Seconds())

	if secondsDiff < daysSeparator {
		return []types.PosminingPeriod{types.NewPosminingPeriod(lastTx, ctx.BlockTime(), sdk.NewInt(0))}
	}

	periods := secondsDiff / daysSeparator
	mod := int64(math.Mod(float64(secondsDiff), float64(daysSeparator)))

	var result []types.PosminingPeriod
	var i int64 = 0

	for i < periods {
		result = append(result, types.NewPosminingPeriod(
			lastTx.Add(time.Duration(daysSeparator*i)*time.Second),
			lastTx.Add(time.Duration(daysSeparator*(i+1))*time.Second),
			sdk.NewInt(0),
		))

		i += 1
	}

	// What's left
	if mod > 0 {
		latestPeriod := lastTx.Add(time.Duration(daysSeparator*periods) * time.Second)

		result = append(result, types.NewPosminingPeriod(
			latestPeriod,
			latestPeriod.Add(time.Duration(mod)*time.Second),
			sdk.NewInt(0),
		))
	}

	return result
}

// Calculates and returns a group of posmining periods
func (k Keeper) GetPosminingGroup(ctx sdk.Context, posmining types.Posmining, coin coins.Coin, balance sdk.Int) types.PosminingGroup {
	group := types.NewPosminingGroup(posmining, balance)

	// For the custom coins, we just have to apply the usual percents during the whole time
	if !coin.Default {
		group.Add(types.NewPosminingPeriod(posmining.LastCharged, ctx.BlockTime(), sdk.NewInt(0)))

		//return group
	}

	return group
}

// Calculates how many tokens has been posmined
func (k Keeper) CalculatePosmined(ctx sdk.Context, posmining types.Posmining, coin coins.Coin, coinsAmount sdk.Int) sdk.Int {
	// If we have a threshold set and it's already has been reach, we should always return 0
	if coin.PosminingThreshold.IsPositive() && coinsAmount.GTE(coin.PosminingThreshold) {
		return sdk.NewInt(0)
	}

	posmined := posmining.Paramined.Add(k.GetPosminingGroup(ctx, posmining, coin, coinsAmount).Paramined)

	if coin.PosminingThreshold.IsPositive() && posmined.IsPositive() && coinsAmount.Add(posmined).GTE(coin.PosminingThreshold) {
		posmined = coin.PosminingThreshold.Sub(coinsAmount)
	}

	return posmined
}
