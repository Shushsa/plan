package keeper_test

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/x/mock"
	"github.com/cosmos/cosmos-sdk/x/staking"
	keeper2 "github.com/Shushsa/plan/x/emission/keeper" // ouro emission folder
	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"time"

	"github.com/cosmos/cosmos-sdk/x/bank"
	xbank "github.com/Shushsa/plan/x/bank"                  // ouro bank folder
	planTypes "github.com/Shushsa/plan/x/nameservice/types" // nameservice folder
	paramining "github.com/Shushsa/plan/x/paramining/keeper"
	"github.com/Shushsa/plan/x/paramining/types"
	"github.com/tendermint/tendermint/crypto"

	"github.com/stretchr/testify/assert"
)

var (
	initCoins = sdk.NewCoins(planTypes.NewPlanCoin(10000000000))
)

type testInput struct {
	mApp     *mock.App
	ctx      sdk.Context
	keeper   paramining.Keeper
	addrs    []sdk.AccAddress
	pubKeys  []crypto.PubKey
	privKeys []crypto.PrivKey
}

// Настраиваем и возвращаем тестовое окружение
func getMockApp(t *testing.T) testInput {
	mApp := mock.NewApp()

	keyParamining := sdk.NewKVStoreKey(paramining.StoreKey)
	keyEmission := sdk.NewKVStoreKey(keeper2.StoreKey)

	bankSupspace := mApp.ParamsKeeper.Subspace(bank.DefaultParamspace)

	bankKeeper := xbank.NewBankKeeper(
		mApp.AccountKeeper,
		bankSupspace,
		bank.DefaultCodespace,
	)

	emissionKeeper := keeper2.NewKeeper(mApp.Cdc, keyEmission, staking.Keeper{})

	keeper := paramining.NewKeeper(
		mApp.Cdc,
		keyParamining,
		bankKeeper,
		staking.Keeper{},
		emissionKeeper,
	)

	mApp.CompleteSetup(keyParamining, keyEmission)

	genAccs, addrs, pubKeys, privKeys := mock.CreateGenAccounts(10, initCoins)

	mock.SetGenesis(mApp, genAccs)

	ctx := mApp.BaseApp.NewContext(true, abci.Header{Time: time.Now()})

	return testInput{
		mApp:     mApp,
		ctx:      ctx,
		keeper:   keeper,
		addrs:    addrs,
		pubKeys:  pubKeys,
		privKeys: privKeys,
	}
}

func TestGetDailyPercent(t *testing.T) {
	// All below 0.01 should be returning 0% daily percent
	assert.True(t, sdk.NewInt(0).Equal(types.GetDailyPercent(sdk.NewInt(0))))
	assert.True(t, sdk.NewInt(0).Equal(types.GetDailyPercent(sdk.NewInt(1000))))

	// 0.01 => 99 - 0.06%
	assert.True(t, sdk.NewInt(6).Equal(types.GetDailyPercent(sdk.NewInt(100000))))                          // 0.1
	assert.True(t, sdk.NewInt(6).Equal(types.GetDailyPercent(sdk.NewIntWithDecimal(67, planTypes.POINTS)))) // 67
	assert.True(t, sdk.NewInt(6).Equal(types.GetDailyPercent(sdk.NewInt(99999999))))                        // 99.999999

	// 100 => 999 - 0.07%
	assert.True(t, sdk.NewInt(7).Equal(types.GetDailyPercent(sdk.NewIntWithDecimal(100, planTypes.POINTS)))) // 100
	assert.True(t, sdk.NewInt(7).Equal(types.GetDailyPercent(sdk.NewIntWithDecimal(764, planTypes.POINTS)))) // 764
	assert.True(t, sdk.NewInt(7).Equal(types.GetDailyPercent(sdk.NewInt(999999999))))                        // 999.999999

	// 1000 => 9999 - 0.09%
	assert.True(t, sdk.NewInt(9).Equal(types.GetDailyPercent(sdk.NewIntWithDecimal(1000, planTypes.POINTS)))) // 1000
	assert.True(t, sdk.NewInt(9).Equal(types.GetDailyPercent(sdk.NewIntWithDecimal(5432, planTypes.POINTS)))) // 5432
	assert.True(t, sdk.NewInt(9).Equal(types.GetDailyPercent(sdk.NewInt(9999999999))))                        // 9999.999999

	// 10.000 => 49.999 - 0.1%
	assert.True(t, sdk.NewInt(10).Equal(types.GetDailyPercent(sdk.NewIntWithDecimal(10000, planTypes.POINTS)))) // 10k
	assert.True(t, sdk.NewInt(10).Equal(types.GetDailyPercent(sdk.NewIntWithDecimal(15873, planTypes.POINTS)))) // 15.873k
	assert.True(t, sdk.NewInt(10).Equal(types.GetDailyPercent(sdk.NewInt(49999999999))))                        // 49.999k

	// 50.000 => 99.999 - 0.12%
	assert.True(t, sdk.NewInt(12).Equal(types.GetDailyPercent(sdk.NewIntWithDecimal(50000, planTypes.POINTS)))) // 50k
	assert.True(t, sdk.NewInt(12).Equal(types.GetDailyPercent(sdk.NewIntWithDecimal(65783, planTypes.POINTS)))) // 65.783k
	assert.True(t, sdk.NewInt(12).Equal(types.GetDailyPercent(sdk.NewInt(99999999999))))                        // 99.999k

	// 100.000 => 499.999 - 0.14%
	assert.True(t, sdk.NewInt(14).Equal(types.GetDailyPercent(sdk.NewIntWithDecimal(100000, planTypes.POINTS)))) // 100k
	assert.True(t, sdk.NewInt(14).Equal(types.GetDailyPercent(sdk.NewIntWithDecimal(356784, planTypes.POINTS)))) // 356.784k
	assert.True(t, sdk.NewInt(14).Equal(types.GetDailyPercent(sdk.NewInt(499999999999))))                        // 499.999k

	// 500.000 => 2.000.000 - 0.16%
	assert.True(t, sdk.NewInt(16).Equal(types.GetDailyPercent(sdk.NewIntWithDecimal(500000, planTypes.POINTS))))  // 500k
	assert.True(t, sdk.NewInt(16).Equal(types.GetDailyPercent(sdk.NewIntWithDecimal(1874659, planTypes.POINTS)))) // 1.874kk
	assert.True(t, sdk.NewInt(16).Equal(types.GetDailyPercent(sdk.NewIntWithDecimal(2000000, planTypes.POINTS)))) // 2kk
}

func TestGetStructureCoff(t *testing.T) {
	// All below 1k - 0
	assert.True(t, sdk.NewInt(0).Equal(types.GetStructureCoff(sdk.NewInt(0))))         // 0
	assert.True(t, sdk.NewInt(0).Equal(types.GetStructureCoff(sdk.NewInt(999999999)))) // 999.999999

	// 1000-9999 - 2.18
	assert.True(t, sdk.NewInt(218).Equal(types.GetStructureCoff(sdk.NewIntWithDecimal(1000, planTypes.POINTS)))) // 1k
	assert.True(t, sdk.NewInt(218).Equal(types.GetStructureCoff(sdk.NewInt(9999999999))))                        // 9.999k

	// 10000-99999 - 2.36
	assert.True(t, sdk.NewInt(236).Equal(types.GetStructureCoff(sdk.NewIntWithDecimal(10000, planTypes.POINTS)))) // 10k
	assert.True(t, sdk.NewInt(236).Equal(types.GetStructureCoff(sdk.NewInt(99999999999))))                        // 99.999k

	// 100000-999999 - 2.77
	assert.True(t, sdk.NewInt(277).Equal(types.GetStructureCoff(sdk.NewIntWithDecimal(100000, planTypes.POINTS)))) // 100k
	assert.True(t, sdk.NewInt(277).Equal(types.GetStructureCoff(sdk.NewInt(999999999999))))                        // 999.999k

	// 1.000.000-9.999.999 - 3.05
	assert.True(t, sdk.NewInt(305).Equal(types.GetStructureCoff(sdk.NewIntWithDecimal(1000000, planTypes.POINTS)))) // 1kkk
	assert.True(t, sdk.NewInt(305).Equal(types.GetStructureCoff(sdk.NewInt(9999999999999))))                        // 9.999kk

	// 10.000.000-99.999.999 - 3.36
	assert.True(t, sdk.NewInt(336).Equal(types.GetStructureCoff(sdk.NewIntWithDecimal(10000000, planTypes.POINTS)))) // 10k
	assert.True(t, sdk.NewInt(336).Equal(types.GetStructureCoff(sdk.NewInt(99999999999999))))                        // 99.999kk

	// 100.000.000-999.999.999 - 3.88
	assert.True(t, sdk.NewInt(388).Equal(types.GetStructureCoff(sdk.NewIntWithDecimal(100000000, planTypes.POINTS)))) // 100kk
	assert.True(t, sdk.NewInt(388).Equal(types.GetStructureCoff(sdk.NewInt(999999999999999))))                        // 999.999k

	// > 1.000.000.000	- 4.37
	assert.True(t, sdk.NewInt(437).Equal(types.GetStructureCoff(sdk.NewIntWithDecimal(1000000000, planTypes.POINTS)))) // 1kkk
	assert.True(t, sdk.NewInt(437).Equal(types.GetStructureCoff(sdk.NewIntWithDecimal(3543000000, planTypes.POINTS)))) // 3.54kkk
}

// Testing the GetTimeDiffFromSeconds method
func TestGetTimeDiffFromSeconds(t *testing.T) {
	app := getMockApp(t)

	keeper, ctx := app.keeper, app.ctx

	secondsPassed := int64(0)

	diff := keeper.GetTimeDiffFromSeconds(ctx, secondsPassed)

	assert.True(t, diff.Total.IsZero())

	// 53 seconds
	secondsPassed = 53

	diff = keeper.GetTimeDiffFromSeconds(ctx, secondsPassed)

	assert.True(t, diff.Total.Equal(sdk.NewInt(secondsPassed)))
	assert.True(t, diff.Seconds.Equal(sdk.NewInt(53)))

	// 3 minutes, 12 seconds
	secondsPassed = 192

	diff = keeper.GetTimeDiffFromSeconds(ctx, secondsPassed)

	assert.True(t, diff.Total.Equal(sdk.NewInt(secondsPassed)))
	assert.True(t, diff.Minutes.Equal(sdk.NewInt(3)))
	assert.True(t, diff.Seconds.Equal(sdk.NewInt(12)))

	// 4 hours, 1 minute and 5 seconds
	secondsPassed = 14400 + 60 + 5

	diff = keeper.GetTimeDiffFromSeconds(ctx, secondsPassed)

	assert.True(t, diff.Total.Equal(sdk.NewInt(secondsPassed)))

	assert.True(t, diff.Seconds.Equal(sdk.NewInt(5)))
	assert.True(t, diff.Minutes.Equal(sdk.NewInt(1)))
	assert.True(t, diff.Hours.Equal(sdk.NewInt(4)))

	// 4 days, 1 hour, 9 minutes and 23 seconds
	secondsPassed = 345600 + 3600 + 540 + 23

	diff = keeper.GetTimeDiffFromSeconds(ctx, secondsPassed)

	assert.True(t, diff.Total.Equal(sdk.NewInt(secondsPassed)))

	assert.True(t, diff.Seconds.Equal(sdk.NewInt(23)))
	assert.True(t, diff.Minutes.Equal(sdk.NewInt(9)))
	assert.True(t, diff.Hours.Equal(sdk.NewInt(1)))
	assert.True(t, diff.Days.Equal(sdk.NewInt(4)))
}

// Testing calculation of coins per time unit (day, hour, minute, second)
func TestCalculateCoinsPerTime(t *testing.T) {
	app := getMockApp(t)

	keeper, ctx := app.keeper, app.ctx

	// 1 PLAN
	coinsAmount := sdk.NewInt(1000000)

	// 0.06% per day = 600
	p := types.Paramining{
		DailyPercent:  sdk.NewInt(6),
		StructureCoff: sdk.NewInt(0),
	}
	result := keeper.CalculateCoinsPerTime(ctx, p, sdk.NewInt(0), sdk.NewCoins(sdk.NewCoin(planTypes.PLAN, coinsAmount)))

	assert.Equal(t, sdk.NewInt(600), result.Day)
	assert.Equal(t, sdk.NewInt(25), result.Hour)
	assert.Equal(t, sdk.NewInt(0), result.Minute)
	assert.Equal(t, sdk.NewInt(0), result.Second)

	// 0.06% per day + 50% savings coff = 900
	p = types.Paramining{
		DailyPercent:  sdk.NewInt(6),
		StructureCoff: sdk.NewInt(0),
	}
	result = keeper.CalculateCoinsPerTime(ctx, p, sdk.NewInt(150), sdk.NewCoins(sdk.NewCoin(planTypes.PLAN, coinsAmount)))

	assert.Equal(t, sdk.NewInt(900), result.Day)
	assert.Equal(t, sdk.NewInt(37), result.Hour)
	assert.Equal(t, sdk.NewInt(0), result.Minute)
	assert.Equal(t, sdk.NewInt(0), result.Second)

	// 0.16% per day * 4,37 structure coff = 6992
	p = types.Paramining{
		DailyPercent:  sdk.NewInt(16),
		StructureCoff: sdk.NewInt(437),
	}

	result = keeper.CalculateCoinsPerTime(ctx, p, sdk.NewInt(0), sdk.NewCoins(sdk.NewCoin(planTypes.PLAN, coinsAmount)))

	assert.Equal(t, sdk.NewInt(6992), result.Day)
	assert.Equal(t, sdk.NewInt(291), result.Hour)
	assert.Equal(t, sdk.NewInt(4), result.Minute)
	assert.Equal(t, sdk.NewInt(0), result.Second)
}

// Testing paramining
func TestCalculateParamined(t *testing.T) {
	app := getMockApp(t)

	keeper, _ := app.keeper, app.ctx

	currentTime := time.Date(2019, 10, 3, 00, 0, 0, 0, time.UTC)
	ctx := app.mApp.BaseApp.NewContext(true, abci.Header{Time: currentTime})

	// 1 PLAN
	coinsAmount := sdk.NewInt(1000000)

	// 600 (0.06%) per day, so it should be 1200 + 25 for 2 days and 1 hour
	p := types.Paramining{
		DailyPercent:    sdk.NewInt(6),
		StructureCoff:   sdk.NewInt(0),
		Paramined:       sdk.NewInt(0),
		LastTransaction: time.Date(2019, 9, 29, 0, 0, 0, 0, time.UTC),
		LastCharged:     time.Date(2019, 9, 30, 23, 0, 0, 0, time.UTC),
	}

	tokens := keeper.CalculateParamined(ctx, p, sdk.NewCoins(sdk.NewCoin(planTypes.PLAN, coinsAmount)))

	assert.Equal(t, tokens, sdk.NewInt(1225))
}

// Testing paramining when the threshold has been reached
func TestCalculateParaminedWithThreshold(t *testing.T) {
	app := getMockApp(t)

	keeper, _ := app.keeper, app.ctx

	currentTime := time.Date(2019, 10, 3, 00, 0, 0, 0, time.UTC)
	ctx := app.mApp.BaseApp.NewContext(true, abci.Header{Time: currentTime})

	p := types.Paramining{
		DailyPercent:    sdk.NewInt(6),
		StructureCoff:   sdk.NewInt(0),
		Paramined:       sdk.NewInt(0),
		LastTransaction: time.Date(2019, 9, 29, 0, 0, 0, 0, time.UTC),
		LastCharged:     time.Date(2019, 9, 30, 23, 0, 0, 0, time.UTC),
	}

	// 1 million 9999... PLAN
	coinsAmount := sdk.NewIntWithDecimal(1999999, 6)

	tokens := keeper.CalculateParamined(ctx, p, sdk.NewCoins(sdk.NewCoin(planTypes.PLAN, coinsAmount)))

	// 1 more token, so it becomes 2kk
	assert.Equal(t, tokens, sdk.NewIntWithDecimal(1, 6))

	// 2kk PLAN
	coinsAmount = sdk.NewIntWithDecimal(2000000, 6)

	tokens = keeper.CalculateParamined(ctx, p, sdk.NewCoins(sdk.NewCoin(planTypes.PLAN, coinsAmount)))

	// Paramining should not working
	assert.Equal(t, tokens, sdk.NewInt(0))

	// 3kk PLAN
	coinsAmount = sdk.NewIntWithDecimal(3000000, 6)

	tokens = keeper.CalculateParamined(ctx, p, sdk.NewCoins(sdk.NewCoin(planTypes.PLAN, coinsAmount)))

	// Still not working
	assert.Equal(t, tokens, sdk.NewInt(0))
}

// Testing the getSavingsChunks method
func TestGetSavingsChunks(t *testing.T) {
	app := getMockApp(t)

	keeper, ctx := app.keeper, app.ctx

	// Should return just a single element since seconds < 3 days
	assert.Equal(t, keeper.GetSavingsChunks(ctx, 15*86400), []int64{15 * 86400})

	// 42 days
	assert.Equal(t, keeper.GetSavingsChunks(ctx, 42*86400), []int64{2592000, 86400 * 12})

	// 137 days
	assert.Equal(t, keeper.GetSavingsChunks(ctx, 137*86400), []int64{2592000, 2592000, 2592000, 2592000, 86400 * 17})
}

// Tests how paramining is being calculated with the savings periods
func TestGetParaminingPeriods(t *testing.T) {
	app := getMockApp(t)

	keeper, _ := app.keeper, app.ctx

	// Today is 3th sempteber
	currentTime := time.Date(2019, 10, 3, 00, 0, 0, 0, time.UTC)

	// The last tx was sent at 30 august, 3 days ago (or 259200 seconds)
	lastTransaction := time.Date(2019, 9, 30, 0, 0, 0, 0, time.UTC)

	// And the last reinvest was 170863 seconds ago
	lastCharge := time.Date(2019, 10, 1, 0, 32, 17, 0, time.UTC)

	ctx := app.mApp.BaseApp.NewContext(true, abci.Header{Time: currentTime})

	// We should get
	expectedPeriod := []types.ParaminingPeriod{{Total: sdk.NewInt(170863), SavingsCoff: sdk.NewInt(0)}}

	assert.Equal(t, keeper.GetParaminingPeriods(ctx, types.Paramining{LastTransaction: lastTransaction, LastCharged: lastCharge}), expectedPeriod)

	// The last tx was sent 33 days ago (or 2851200 seconds)
	lastTransaction = time.Date(2019, 8, 29, 0, 0, 0, 0, time.UTC)

	// The last reinvest was 170863 seconds ago
	lastCharge = time.Date(2019, 10, 1, 0, 32, 17, 0, time.UTC)

	ctx = app.mApp.BaseApp.NewContext(true, abci.Header{Time: currentTime})

	// So we should get savings coff 1.5
	expectedPeriod = []types.ParaminingPeriod{{Total: sdk.NewInt(170863), SavingsCoff: sdk.NewInt(150)}}

	assert.Equal(t, keeper.GetParaminingPeriods(ctx, types.Paramining{LastTransaction: lastTransaction, LastCharged: lastCharge}), expectedPeriod)

	// The last tx was sent 2 years ago
	lastTransaction = time.Date(2017, 8, 29, 0, 0, 0, 0, time.UTC)

	// The last reinvest was 170863 seconds ago
	lastCharge = time.Date(2019, 10, 1, 0, 32, 17, 0, time.UTC)

	ctx = app.mApp.BaseApp.NewContext(true, abci.Header{Time: currentTime})

	// So we should get full 2 savings coff
	expectedPeriod = []types.ParaminingPeriod{types.ParaminingPeriod{Total: sdk.NewInt(170863), SavingsCoff: sdk.NewInt(200)}}

	assert.Equal(t, keeper.GetParaminingPeriods(ctx, types.Paramining{LastTransaction: lastTransaction, LastCharged: lastCharge}), expectedPeriod)

	// Last tx was sent 31 days ago
	lastTransaction = time.Date(2019, 9, 2, 0, 0, 0, 0, time.UTC)

	// Last reinvest was ~2 days ago, so we should calculate one day without the savings coff and the next one with it
	lastCharge = time.Date(2019, 10, 1, 0, 32, 17, 0, time.UTC)
	ctx = app.mApp.BaseApp.NewContext(true, abci.Header{Time: currentTime})

	expectedPeriod = []types.ParaminingPeriod{{Total: sdk.NewInt(84463), SavingsCoff: sdk.NewInt(0)}, {Total: sdk.NewInt(86400), SavingsCoff: sdk.NewInt(150)}}

	assert.Equal(t, keeper.GetParaminingPeriods(ctx, types.Paramining{LastTransaction: lastTransaction, LastCharged: lastCharge}), expectedPeriod)
}
