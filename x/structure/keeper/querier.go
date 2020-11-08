package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

const (
	QueryGet = "get"
	QueryUpper = "upper"
)

// NewQuerier is the module level router for state queries
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case QueryGet:
			return queryGet(ctx, path[1:], req, keeper)
		case QueryUpper:
			return queryUpper(ctx, path[1:], req, keeper)
		default:
			return nil, sdk.ErrUnknownRequest("Unknown structure query endpoint")
		}
	}
}
func queryGet(ctx sdk.Context, path []string, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	addr, err := sdk.AccAddressFromBech32(path[0])

	if err != nil {
		return []byte{}, sdk.ErrUnknownRequest("Wrong address")
	}

	requestedStructure := k.GetStructure(ctx, addr)

	res, err := codec.MarshalJSONIndent(k.cdc, requestedStructure)

	if err != nil {
		panic("could not marshal result to JSON")
	}

	return res, nil
}
func queryUpper(ctx sdk.Context, path []string, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	addr, err := sdk.AccAddressFromBech32(path[0])

	if err != nil {
		return []byte{}, sdk.ErrUnknownRequest("Wrong address")
	}

	requestedStructure := k.GetUpperStructure(ctx, addr)

	res, err := codec.MarshalJSONIndent(k.cdc, requestedStructure)

	if err != nil {
		panic("could not marshal result to JSON")
	}

	return res, nil
}