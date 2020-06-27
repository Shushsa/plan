package keeper

import (
	"math"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	planTypes "github.com/cosmos/sdk-tutorials/nameservice/x/nameservice/types"
	"github.com/cosmos/sdk-tutorials/nameservice/x/paramining/types"
)

// Calculates how many tokens has been paramined
func (k Keeper) CalculateParamined(ctx sdk.Context, paramining types.Paramining, coins sdk.Coins) sdk.Int {
	periods := k.GetParaminingPeriods(ctx, paramining)

	totalCoins := sdk.NewInt(0)

	for _, value := range periods {
		timeDiff := k.GetTimeDiffFromSeconds(ctx, value.Total.Int64())

		perTime := k.CalculateCoinsPerTime(ctx, paramining, value.SavingsCoff, coins)

		totalCoins = totalCoins.Add(timeDiff.Seconds.Mul(
			perTime.Second).Add(
			timeDiff.Minutes.Mul(perTime.Minute)).Add(
			timeDiff.Hours.Mul(perTime.Hour)).Add(
			timeDiff.Days.Mul(perTime.Day)))
	}

	paramined := totalCoins.Add(paramining.Paramined)

	coinsAmount := coins.AmountOf(planTypes.PLAN)
	paraminingThreshold := planTypes.GetParaminingThreshold()

	// If the balance is greater that our hardcap (2m plans) - let's just charge the difference
	if paramined.IsPositive() && coinsAmount.Add(paramined).GTE(paraminingThreshold) {
		if coinsAmount.GTE(paraminingThreshold) {
			paramined = sdk.NewInt(0)
		} else {
			paramined = paraminingThreshold.Sub(coinsAmount)
		}
	}

	return paramined
}

// Calculates how many tokens we should charge during every "savings" period - every period is 30 days long
func (k Keeper) GetParaminingPeriods(ctx sdk.Context, paramining types.Paramining) []types.ParaminingPeriod {
	chunks := k.GetSavingsChunks(ctx, int64(ctx.BlockHeader().Time.Sub(paramining.LastTransaction).Seconds()))
	chargeDiff := int64(ctx.BlockHeader().Time.Sub(paramining.LastCharged).Seconds())

	var result []types.ParaminingPeriod

	i := len(chunks) - 1

	var duration int64

	for chargeDiff > 0 {
		if chargeDiff > chunks[i] {
			duration = chunks[i]
		} else {
			duration = chargeDiff
		}

		chargeDiff -= duration

		var savingsCoff int64

		if len(types.Savings) <= i {
			savingsCoff = types.Savings[len(types.Savings)-1]
		} else {
			savingsCoff = types.Savings[i]
		}

		period := types.ParaminingPeriod{
			Total:       sdk.NewInt(duration),
			SavingsCoff: sdk.NewInt(savingsCoff),
		}

		result = append([]types.ParaminingPeriod{period}, result...)

		i -= 1
	}

	return result
}

// Calculates savings chunks - splits all the time passed by 30 days periods and returns them
func (k Keeper) GetSavingsChunks(ctx sdk.Context, seconds int64) []int64 {
	var daysSeparator int64 = 2592000

	if seconds < daysSeparator {
		return []int64{seconds}
	}

	periods := seconds / daysSeparator
	mod := int64(math.Mod(float64(seconds), float64(daysSeparator)))

	var result []int64
	var i int64 = 0

	for i < periods {
		result = append(result, daysSeparator)

		i += 1
	}

	// What's left
	if mod > 0 {
		result = append(result, mod)
	}

	return result
}

// Calculates how many tokens we get per second, minute, hour and day
func (k Keeper) CalculateCoinsPerTime(ctx sdk.Context, paramining types.Paramining, savingsCoff sdk.Int, coins sdk.Coins) types.CoinsPerTime {
	result := types.NewCoinsPerTime()

	if paramining.DailyPercent.IsZero() || k.emissionKeeper.IsThresholdReached(ctx) {
		return result
	}

	toQuo := sdk.NewInt(10000)
	actualPercent := paramining.DailyPercent

	if paramining.StructureCoff.IsZero() == false {
		actualPercent = actualPercent.Mul(paramining.StructureCoff)
		toQuo = toQuo.MulRaw(100) // Добавляем 00 в конец т.к. перемножаем с кофф
	}

	if savingsCoff.IsZero() == false {
		actualPercent = actualPercent.Mul(savingsCoff)
		toQuo = toQuo.MulRaw(100) // Добавляем 00 в конец т.к. перемножаем с кофф
	}

	actualCoins := coins.AmountOf(planTypes.PLAN)

	result.Day = actualCoins.Mul(actualPercent).Quo(toQuo)
	result.Hour = result.Day.QuoRaw(24)
	result.Minute = result.Hour.QuoRaw(60)
	result.Second = result.Minute.QuoRaw(60)

	return result
}

// Get time difference in days, hours, minutes and seconds
func (k Keeper) GetTimeDiffFromSeconds(ctx sdk.Context, seconds int64) types.TimeDifference {
	duration := time.Duration(seconds) * time.Second

	difference := types.NewTimeDifference()
	difference.Total = sdk.NewInt(int64(duration.Seconds()))

	// Меньше минуты
	if duration.Seconds() < 60.0 {
		difference.Seconds = sdk.NewInt(int64(duration.Seconds()))

		return difference
	}

	// Меньше часа
	if duration.Minutes() < 60.0 {
		difference.Minutes = sdk.NewInt(int64(duration.Minutes()))
		difference.Seconds = sdk.NewInt(int64(math.Mod(duration.Seconds(), 60)))

		return difference
	}

	// Меньше дня
	if duration.Hours() < 24.0 {
		difference.Hours = sdk.NewInt(int64(duration.Hours()))
		difference.Minutes = sdk.NewInt(int64(math.Mod(duration.Minutes(), 60)))
		difference.Seconds = sdk.NewInt(int64(math.Mod(duration.Seconds(), 60)))

		return difference
	}

	difference.Days = sdk.NewInt(int64(duration.Hours() / 24))
	difference.Hours = sdk.NewInt(int64(math.Mod(duration.Hours(), 24)))
	difference.Minutes = sdk.NewInt(int64(math.Mod(duration.Minutes(), 60)))
	difference.Seconds = sdk.NewInt(int64(math.Mod(duration.Seconds(), 60)))

	return difference
}

// Returns current savings coff depends on how many days passed since the latest tx
func (k Keeper) GetCurrentSavingsCoff(ctx sdk.Context, paramining types.Paramining) sdk.Int {
	chunks := k.GetSavingsChunks(ctx, int64(ctx.BlockHeader().Time.Sub(paramining.LastTransaction).Seconds()))

	if len(types.Savings) >= len(chunks) {
		return sdk.NewInt(types.Savings[len(chunks)-1])
	}

	return sdk.NewInt(types.Savings[len(types.Savings)-1])
}
