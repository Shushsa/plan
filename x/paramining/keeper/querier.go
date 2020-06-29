package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	abci "github.com/tendermint/tendermint/abci/types"
)

const (
	QueryGet = "get"
)

// NewQuerier is the module level router for state queries
func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err error) {
		switch path[0] {
		case QueryGet:
			return queryGet(ctx, path[1:], req, k)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Unknown paramining query endpoint")
		}
	}
}

// Returns the paramining record
func queryGet(ctx sdk.Context, path []string, req abci.RequestQuery, k Keeper) ([]byte, error) {
	addr, err := sdk.AccAddressFromBech32(path[0])

	if err != nil {
		return []byte{}, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Wrong address")
	}

	res, err := codec.MarshalJSONIndent(k.cdc, k.GetParaminingResolve(ctx, addr))

	if err != nil {
		panic("could not marshal result to JSON")
	}

	return res, nil
}
