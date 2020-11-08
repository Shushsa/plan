package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/plan-crypto/node/x/paramining/types"
)

func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case types.MsgReinvest:
			return handleReinvest(ctx, k, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized paramining Msg type: %v", msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleReinvest(ctx sdk.Context, k Keeper, msg types.MsgReinvest) sdk.Result {
	reinvested := k.ChargeParamining(ctx, msg.Owner, true)

	event := sdk.NewEvent(
		sdk.EventTypeMessage,
		sdk.NewAttribute(sdk.AttributeKeySender, msg.Owner.String()),
		sdk.NewAttribute(AttributeKeyAmount, reinvested.String()),
	)

	ctx.EventManager().EmitEvent(event)

	return sdk.Result{Events: []sdk.Event{event}}
}
