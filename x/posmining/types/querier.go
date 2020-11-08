package types

import sdk "github.com/cosmos/cosmos-sdk/types"

// Для возвращения скомпанованного ответа
type ParaminingResolve struct {
	Paramined sdk.Int `json:"paramined"`
	SavingsCoff sdk.Int `json:"savings_coff"`
	Paramining Paramining `json:"paramining"`
	CoinsPerTime CoinsPerTime `json:"coins_per_time"`
}


func (r ParaminingResolve) String() string {
	return r.Paramined.String()
}