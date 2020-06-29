package keeper

import (
	"fmt"

	"github.com/Shushsa/plan/x/paramining/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		switch msg := msg.(type) {
		case types.MsgReinvest:
			return handleReinvest(ctx, k, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized paramining Msg type: %v", msg.Type())
			// return sdk.ErrUnknownRequest(errMsg).Result()
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}

func handleReinvest(ctx sdk.Context, k Keeper, msg types.MsgReinvest) (*sdk.Result, error) {
	reinvested := k.ChargeParamining(ctx, msg.Owner, true)

	event := sdk.NewEvent(
		sdk.EventTypeMessage,
		sdk.NewAttribute(sdk.AttributeKeySender, msg.Owner.String()),
		sdk.NewAttribute(AttributeKeyAmount, reinvested.String()),
	)

	ctx.EventManager().EmitEvent(event)

	return &sdk.Result{Events: []sdk.Event{event}}, nil
}
