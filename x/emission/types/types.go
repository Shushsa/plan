package types

import 	(
	sdk "github.com/cosmos/cosmos-sdk/types"
	plnTypes "github.com/plan-crypto/node/x/plan/types"
)


type Emission struct {
	Current sdk.Int `json:"current"` // Текущая эмиссия
	Threshold sdk.Int `json:"threshold"` // Порог, после которого парамайнинг перестает работать
}

// Достигнут ли порог
func(e Emission) IsThresholdReached() bool {
	return e.Current.GTE(e.Threshold)
}

// Starting emission
func NewEmission() Emission {
	return Emission{
		Current: sdk.NewIntWithDecimal(plnTypes.INITIAL, plnTypes.POINTS),
		Threshold: sdk.NewIntWithDecimal(7000000000, plnTypes.POINTS),
	}
}

func (e Emission) String() string {
	return e.Current.String()
}
