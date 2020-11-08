package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

)

const RouterKey = ModuleName // this was defined in your key.go file

// Taking the validators out of jail
type MsgAmnesty struct {
	Owner sdk.AccAddress `json:"owner"`
}


// NewMsgSetName is a constructor function for MsgSetName
func NewMsgAmnesty(owner sdk.AccAddress) MsgAmnesty {
	return MsgAmnesty{
		Owner: owner,
	}
}

func (msg MsgAmnesty) Route() string { return RouterKey }

func (msg MsgAmnesty) Type() string { return "amnesty" }

func (msg MsgAmnesty) ValidateBasic() sdk.Error {
	if msg.Owner.Empty() || msg.Owner.String() != GenesisWallet {
		return sdk.ErrInvalidAddress(msg.Owner.String())
	}

	return nil
}

func (msg MsgAmnesty) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgAmnesty) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}