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
func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case QueryGet:
			return queryGet(ctx, path[1:], req, k)
		default:
			return nil, sdk.ErrUnknownRequest("Unknown paramining query endpoint")
		}
	}
}

// Returns the paramining record
func queryGet(ctx sdk.Context, path []string, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	addr, err := sdk.AccAddressFromBech32(path[0])

	if err != nil {
		return []byte{}, sdk.ErrUnknownRequest("Wrong address")
	}

	res, err := codec.MarshalJSONIndent(k.cdc, k.GetParaminingResolve(ctx, addr))

	if err != nil {
		panic("could not marshal result to JSON")
	}

	return res, nil
}