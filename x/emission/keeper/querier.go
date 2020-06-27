package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

const (
	QueryGet = "get"
)

// NewQuerier is the module level router for state queries
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case QueryGet:
			return queryGet(ctx, keeper)
		default:
			return nil, sdk.ErrUnknownRequest("Unknown emission query endpoint")
		}
	}
}

// Получаем текущую эмиссию
func queryGet(ctx sdk.Context, keeper Keeper) ([]byte, sdk.Error) {
	res, err := codec.MarshalJSONIndent(keeper.cdc, keeper.GetEmission(ctx))

	if err != nil {
		panic("could not marshal result to JSON")
	}

	return res, nil
}