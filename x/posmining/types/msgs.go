package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const RouterKey = ModuleName // this was defined in your key.go file

// Реинвест пара
type MsgReinvest struct {
	Owner sdk.AccAddress `json:"owner"`
}

// NewMsgSetName is a constructor function for MsgSetName
func NewMsgReinvest(owner sdk.AccAddress) MsgReinvest {
	return MsgReinvest{
		Owner: owner,
	}
}

func (msg MsgReinvest) Route() string { return RouterKey }

func (msg MsgReinvest) Type() string { return "reinvest" }

func (msg MsgReinvest) ValidateBasic() sdk.Error {
	if msg.Owner.Empty() {
		return sdk.ErrInvalidAddress(msg.Owner.String())
	}
	return nil
}

func (msg MsgReinvest) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgReinvest) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}