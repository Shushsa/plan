package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	structure "github.com/plan-crypto/node/x/structure/types"
	paramining "github.com/plan-crypto/node/x/paramining/types"
)

// Profile response
type ProfileResolve struct {
	Owner sdk.AccAddress `json:"owner"`

	Balance sdk.Int `json:"balance"`

	Paramining paramining.ParaminingResolve  `json:"paramining"`

	Structure structure.Structure `json:"structure"`
}


func (r ProfileResolve) String() string {
	return r.Balance.String()
}