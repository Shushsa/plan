package types

import sdk "github.com/cosmos/cosmos-sdk/types"


const (
	PLN = "pln"
	POINTS = 6 // points
	INITIAL = 10000000 // initial emission
	PARAMINING_THRESHOLD = 2000000 // paramining hard cap
)

func NewPlnCoin(amount int64) sdk.Coin {
	return sdk.NewInt64Coin(PLN, amount)
}

func NewCoin(amount sdk.Int) sdk.Coin {
	return sdk.NewCoin(PLN, amount)
}
func NewCoins(amount sdk.Int) sdk.Coins {
	return sdk.NewCoins(sdk.NewCoin(PLN, amount))
}

// Max structure level
func GetMaxLevel() sdk.Int {
	return sdk.NewInt(100)
}

func GetParaminingThreshold() sdk.Int {
	return sdk.NewIntWithDecimal(PARAMINING_THRESHOLD, POINTS)
}