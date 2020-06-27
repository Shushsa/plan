package keeper_test

import (
	"github.com/cosmos/cosmos-sdk/x/mock"
	abci "github.com/tendermint/tendermint/abci/types"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"

	ouroTypes "github.com/ouroboros-crypto/node/x/ouroboros/types"
	structure "github.com/ouroboros-crypto/node/x/structure/keeper"
	"github.com/ouroboros-crypto/node/x/structure/types"
	"github.com/stretchr/testify/assert"
	"github.com/tendermint/tendermint/crypto"
)

var (
	initTokens, _ = sdk.NewIntFromString("10000")
	initCoins     = sdk.NewCoins(sdk.NewCoin(ouroTypes.OURO, initTokens))
)

type testInput struct {
	mApp     *mock.App
	ctx      sdk.Context
	keeper   structure.Keeper
	addrs    []sdk.AccAddress
	pubKeys  []crypto.PubKey
	privKeys []crypto.PrivKey
}

// Creates and returns a test app
func getMockApp(t *testing.T) testInput {
	mApp := mock.NewApp()

	keyStructure := sdk.NewKVStoreKey(structure.StoreKey)
	keyStructureFast := sdk.NewKVStoreKey(structure.FastAccessKey)

	keeper := structure.NewKeeper(
		mApp.Cdc,
		keyStructure,
		keyStructureFast,
	)

	mApp.CompleteSetup(keyStructure, keyStructureFast)

	genAccs, addrs, pubKeys, privKeys := mock.CreateGenAccounts(200, initCoins)

	mock.SetGenesis(mApp, genAccs)

	ctx := mApp.BaseApp.NewContext(true, abci.Header{})

	return testInput{
		mApp:     mApp,
		ctx:      ctx,
		keeper:   keeper,
		addrs:    addrs,
		pubKeys:  pubKeys,
		privKeys: privKeys,
	}
}

// Tests common upper structure methods
func TestCommonUpperStructureMethods(t *testing.T) {
	app := getMockApp(t)

	keeper, ctx := app.keeper, app.ctx
	firstAccount, secondAccount := app.addrs[0], app.addrs[1]

	if keeper.HasUpperStructure(app.ctx, firstAccount) == true {
		t.Errorf("%s should not have an upper structure by default", firstAccount)
	}

	// Создаем тестовый поинтер
	keeper.SetUpperStructure(ctx, firstAccount, types.UpperStructure{
		Owner:   secondAccount,
		Address: firstAccount,
	})

	if keeper.HasUpperStructure(ctx, firstAccount) == false {
		t.Errorf("%s should have upper structure now", firstAccount)
	}

	upperStructure := keeper.GetUpperStructure(ctx, firstAccount)

	assert.Equal(t, upperStructure.Address, firstAccount)
	assert.Equal(t, upperStructure.Owner, secondAccount)
}

// Tests common structure methods
func TestCommonStructureMethods(t *testing.T) {
	app := getMockApp(t)

	keeper, ctx := app.keeper, app.ctx
	firstAccount := app.addrs[0]

	zero := sdk.NewInt(0)

	// First structure that doesn't have any followers yet
	firstStructure := keeper.GetStructure(ctx, firstAccount)

	assert.Equal(t, firstStructure.Owner, firstAccount)

	assert.True(t, firstStructure.Balance.Equal(zero))
	assert.True(t, firstStructure.Followers.Equal(zero))
	assert.True(t, firstStructure.MaxLevel.Equal(zero))

	firstStructure.Balance = sdk.NewInt(500)
	firstStructure.Followers = sdk.NewInt(50)
	firstStructure.MaxLevel = sdk.NewInt(10)

	keeper.SetStructure(ctx, firstStructure)

	loadedStructure := keeper.GetStructure(ctx, firstAccount)

	assert.True(t, loadedStructure.Balance.Equal(sdk.NewInt(500)))
	assert.True(t, loadedStructure.Followers.Equal(sdk.NewInt(50)))
	assert.True(t, loadedStructure.MaxLevel.Equal(sdk.NewInt(10)))
}

// Tests adding followers to the first line of the sturcture
func TestAddToStructureFirstLine(t *testing.T) {
	app := getMockApp(t)

	keeper, ctx := app.keeper, app.ctx

	// 0.001 OURO
	coinsAmount := sdk.NewInt(1000)

	firstAccount, secondAccount, thirdAccount := app.addrs[0], app.addrs[1], app.addrs[2]

	// Make sure the account doesn't have an upper structure yet
	assert.False(t, keeper.HasUpperStructure(ctx, secondAccount))

	// Add him to the structure
	assert.True(t, keeper.AddToStructure(ctx, firstAccount, secondAccount, coinsAmount))

	// Make sure the account has an upper structure now
	assert.True(t, keeper.HasUpperStructure(ctx, secondAccount))

	// Make sure we cannot rewrite the structure
	assert.False(t, keeper.AddToStructure(ctx, thirdAccount, secondAccount, coinsAmount))

	// Make sure the upper account got his structure updated
	firstAccountStructure := keeper.GetStructure(ctx, firstAccount)

	assert.True(t, firstAccountStructure.Balance.Equal(coinsAmount))
	assert.True(t, firstAccountStructure.Followers.Equal(sdk.NewInt(1)))
	assert.True(t, firstAccountStructure.MaxLevel.Equal(sdk.NewInt(1)))
}

// Tests adding multiple followers on different levels
func TestAddToBigStructure(t *testing.T) {
	app := getMockApp(t)

	keeper, ctx := app.keeper, app.ctx

	// 0.001 OURO
	coinsAmount := sdk.NewInt(1000)

	// 5 levels
	i := 5

	// 5 levels depth structure - 5 -> 4, 4 -> 3, 3 -> 2, 2 -> 1, 1 -> 0
	for i > 0 {
		firstAccount, secondAccount := app.addrs[i], app.addrs[i-1]

		assert.True(t, keeper.AddToStructure(ctx, firstAccount, secondAccount, coinsAmount))

		i -= 1
	}

	// First get the first structure record
	lastStructure := keeper.GetStructure(ctx, app.addrs[5])

	// That structure should have 5 followers and 5 levels
	assert.True(t, lastStructure.MaxLevel.Equal(sdk.NewInt(5)))
	assert.True(t, lastStructure.Followers.Equal(sdk.NewInt(5)))

	// Adding a few more followers
	assert.True(t, keeper.AddToStructure(ctx, app.addrs[1], app.addrs[6], coinsAmount))
	assert.True(t, keeper.AddToStructure(ctx, app.addrs[2], app.addrs[7], coinsAmount))
	assert.True(t, keeper.AddToStructure(ctx, app.addrs[2], app.addrs[8], coinsAmount))

	// Refreshing it from db
	lastStructure = keeper.GetStructure(ctx, app.addrs[5])

	// Now it should have 3 more followers but the same levels
	assert.True(t, lastStructure.MaxLevel.Equal(sdk.NewInt(5)))
	assert.True(t, lastStructure.Followers.Equal(sdk.NewInt(8)))

	// Checking the second and the third structures
	anotherStructure := keeper.GetStructure(ctx, app.addrs[1])

	assert.True(t, anotherStructure.MaxLevel.Equal(sdk.NewInt(1)))
	assert.True(t, anotherStructure.Followers.Equal(sdk.NewInt(2)))

	anotherStructure = keeper.GetStructure(ctx, app.addrs[2])

	assert.True(t, anotherStructure.MaxLevel.Equal(sdk.NewInt(2)))
	assert.True(t, anotherStructure.Followers.Equal(sdk.NewInt(5)))
}

// Testing the biggest structure
func TestAddToMaxStructure(t *testing.T) {
	app := getMockApp(t)

	keeper, ctx := app.keeper, app.ctx

	// 0.001 OURO
	coinsAmount := sdk.NewInt(1000)

	// its depth will be ~150 levels (but will be reduced to 100 for the first account in line)
	i := 150

	for i > 0 {
		firstAccount, secondAccount := app.addrs[i], app.addrs[i-1]

		assert.True(t, keeper.AddToStructure(ctx, firstAccount, secondAccount, coinsAmount))

		i -= 1
	}

	// Make sure the first account doesn't have more than 100 followers & levels
	lastStructure := keeper.GetStructure(ctx, app.addrs[150])

	assert.True(t, lastStructure.MaxLevel.Equal(sdk.NewInt(100)))
	assert.True(t, lastStructure.Followers.Equal(sdk.NewInt(100)))

	lastStructure.Balance = sdk.NewInt(2000)

	keeper.SetStructure(ctx, lastStructure)

	currentBalance := lastStructure.Balance

	// Make sure that adding follower to ~101 account won't be added to the first account
	assert.True(t, keeper.AddToStructure(ctx, app.addrs[50], app.addrs[151], coinsAmount))

	lastStructure = keeper.GetStructure(ctx, app.addrs[150])
	assert.True(t, lastStructure.MaxLevel.Equal(sdk.NewInt(100)))
	assert.True(t, lastStructure.Followers.Equal(sdk.NewInt(100)))
	assert.True(t, lastStructure.Balance.Equal(currentBalance.Sub(coinsAmount)))
}

// Testing the increaseStructureBalance method
func TestIncreaseStructureBalance(t *testing.T) {
	app := getMockApp(t)

	keeper, ctx := app.keeper, app.ctx

	// 0 OURO
	coinsAmount := sdk.NewInt(0)

	// 5 levels
	i := 5

	// 5 levels depth structure - 5 -> 4, 4 -> 3, 3 -> 2, 2 -> 1, 1 -> 0
	for i > 0 {
		firstAccount, secondAccount := app.addrs[i], app.addrs[i-1]

		assert.True(t, keeper.AddToStructure(ctx, firstAccount, secondAccount, coinsAmount))

		i -= 1
	}

	// The latest account in the structure gets 0.001 OURO
	coinsAmount = sdk.NewInt(1000)
	keeper.IncreaseStructureBalance(ctx, app.addrs[0], coinsAmount)
	assert.True(t, keeper.GetStructure(ctx, app.addrs[0]).Balance.IsZero())

	i = 5

	// Checking the balances
	for i > 1 {
		assert.True(t, keeper.GetStructure(ctx, app.addrs[i]).Balance.Equal(coinsAmount))

		i -= 1
	}
}

func TestDecreaseStructureBalance(t *testing.T) {
	app := getMockApp(t)

	keeper, ctx := app.keeper, app.ctx

	// 0 OURO
	coinsAmount := sdk.NewInt(0)

	// 5 levels
	i := 5

	// 5 levels depth structure - 5 -> 4, 4 -> 3, 3 -> 2, 2 -> 1, 1 -> 0
	for i > 0 {
		firstAccount, secondAccount := app.addrs[i], app.addrs[i-1]

		assert.True(t, keeper.AddToStructure(ctx, firstAccount, secondAccount, coinsAmount))

		i -= 1
	}

	// The latest account gets 0.001 OURO
	coinsAmount = sdk.NewInt(1000)
	keeper.IncreaseStructureBalance(ctx, app.addrs[0], coinsAmount)

	// The latest account transfered 0.0004 OURO to another structure, now it has just 0.0006 OURO in the structure
	transfered := sdk.NewInt(400)
	keeper.DecreaseStructureBalance(ctx, app.addrs[0], transfered)

	i = 5

	// Checking all the balances
	for i > 1 {
		assert.True(t, keeper.GetStructure(ctx, app.addrs[i]).Balance.Equal(coinsAmount.Sub(transfered)))

		i -= 1
	}
}