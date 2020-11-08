package types

import sdk "github.com/cosmos/cosmos-sdk/types"

const (
	QueryGetPosmining = "get"
)

// Для возвращения скомпанованного ответа
type PosminingResolve struct {
	Coin         string       `json:"coin"`
	Posmined     sdk.Int      `json:"posmined"`
	Posmining    Posmining    `json:"posmining"`
	CoinsPerTime CoinsPerTime `json:"coins_per_time"`
	SavingsCoff  sdk.Int      `json:"savings_coff"`
}

func (r PosminingResolve) String() string {
	return r.Posmined.String()
}
