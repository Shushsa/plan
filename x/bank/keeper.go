package bank

import (
	planTypes "github.com/Shushsa/plan/x/nameservice/types" // nameservice folder
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	sdkbank "github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/staking"
)

type Keeper struct {
	sdkbank.BaseKeeper

	ak            auth.AccountKeeper
	StakingKeeper staking.Keeper
	paramSpace    params.Subspace

	beforeTransferHooks []CoinsTransferHook
	afterTransferHooks  []CoinsTransferHook
}

/*
func NewBankKeeper(ak auth.AccountKeeper,
	paramSpace params.Subspace,
	codespace sdk.CodespaceType) Keeper { // ouro error on this string

	keeper := Keeper{
		BaseKeeper:          sdkbank.NewBaseKeeper(ak, paramSpace, codespace),
		ak:                  ak,
		paramSpace:          paramSpace,
		beforeTransferHooks: make([]CoinsTransferHook, 0),
		afterTransferHooks:  make([]CoinsTransferHook, 0),
	}

	return keeper
} // ouro error
*/

// Returns the balance that should be used during calculations of paramining
func (k Keeper) GetParaminableBalance(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins {
	// return k.GetCoins(ctx, addr).Add(k.GetStackedCoins(ctx, addr)) // ouro error version

	get_coins := k.GetCoins(ctx, addr)
	get_stacked_coins := k.GetStackedCoin(ctx, addr)

	return get_coins.Add(get_stacked_coins)
}

// Returns both stacked and unbounding coins
func (k Keeper) GetStackedCoin(ctx sdk.Context, addr sdk.AccAddress) sdk.Coin {
	result := sdk.NewInt(0)

	// First let's get through the stakes
	stakes := k.StakingKeeper.GetAllDelegatorDelegations(ctx, addr)

	for _, value := range stakes {
		result = result.Add(value.GetShares().TruncateInt())
	}

	// Then let's get through the unbounding coins
	unbounding := k.StakingKeeper.GetAllUnbondingDelegations(ctx, addr)

	for _, value := range unbounding {
		for _, entry := range value.Entries {
			result = result.Add(entry.Balance)
		}
	}

	return planTypes.NewCoin(result)
}

// SendCoins moves coins from one account to another
func (keeper Keeper) SendCoins(
	ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins,
) error {
	keeper.beforeCoinsTransfer(ctx, fromAddr, toAddr, amt)

	err := keeper.BaseKeeper.SendCoins(ctx, fromAddr, toAddr, amt)

	if err == nil {
		keeper.afterCoinsTransfer(ctx, fromAddr, toAddr, amt)
	}

	return err
}
