package plan

import (
	"fmt"
	"github.com/plan-crypto/node/x/plan/types"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	slashingTypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
)


func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case types.MsgAmnesty:
			return handleAmnesty(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized plan Msg type: %v", msg.Type())

			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleAmnesty(ctx sdk.Context, keeper Keeper, msg types.MsgAmnesty) sdk.Result {
	if msg.Owner.String() != types.GenesisWallet {
		return sdk.Result{}
	}

	keeper.slashingKeeper.IterateValidatorSigningInfos(ctx,
		func(address sdk.ConsAddress, info slashingTypes.ValidatorSigningInfo) bool {
			info.JailedUntil = time.Unix(0, 0)
			info.Tombstoned = false

			keeper.slashingKeeper.SetValidatorSigningInfo(ctx, info.Address, info)

			ctx.EventManager().EmitEvent(
				sdk.NewEvent(
					sdk.EventTypeMessage,
					sdk.NewAttribute(sdk.AttributeKeySender, msg.Owner.String()),
					sdk.NewAttribute("amnesty", info.Address.String()),
					sdk.NewAttribute("success", "1"),

				),
			)

			return false
		},
	)

	return sdk.Result{Events: ctx.EventManager().Events()}
}
