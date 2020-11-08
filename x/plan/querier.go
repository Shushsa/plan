package plan

import (
	"github.com/cosmos/cosmos-sdk/codec"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/plan-crypto/node/x/plan/types"
)

// query endpoints supported by the plan Querier
const (
	QueryProfile    = "profile"
)

// NewQuerier is the module level router for state queries
func NewQuerier(keeper Keeper, coinKeeper bank.Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case QueryProfile:
			return queryProfile(ctx, path[1:], req, keeper, coinKeeper)
		default:
			return nil, sdk.ErrUnknownRequest("unknown plan query endpoint")
		}
	}
}

func queryProfile(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper, coinKeeper bank.Keeper) ([]byte, sdk.Error) {
	addr, err := sdk.AccAddressFromBech32(path[0])

	if err != nil {
		return []byte{}, sdk.ErrUnknownRequest("Wrong address")
	}

	balance := coinKeeper.GetCoins(ctx, addr)

	res, codecErr := codec.MarshalJSONIndent(keeper.cdc, types.ProfileResolve{
		Owner: addr,
		Balance: balance.AmountOf(types.PLN),
		Paramining: keeper.paraminingKeeper.GetParaminingResolve(ctx, addr),
		Structure: keeper.structureKeeper.GetStructure(ctx, addr),
	})

	if codecErr != nil {
		panic("could not marshal result to JSON")
	}

	return res, nil
}
