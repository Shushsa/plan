package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const CorrectionUpdatePeriod = 24

var Savings = [...]int64{
	0,   // 0-30 days
	150, // 1.50 or 50%, 30-60 days
	151, // 1.51 or 51%, 60-90 days
	152, // 1.52 or 52%, 90-120 days
	153, // 1.53 or 53%, 120-150 days
	154, // 1.54 or 54%, 150-180 days
	155, // 1.55 or 55%, 180-210 days
	155, // 1.55 or 55%, 210-240 days
	155, // 1.55 or 55%, 240-270 days
	155, // 1.55 or 55%, 270-300 days
	155, // 1.55 or 55%, 300-330 days
	155, // 1.55 or 55%, 330-360 days
	200, // 2.00 or 100%, >360 days
}

func GetSavingCoff(i int) sdk.Int {
	if len(Savings) > i {
		return sdk.NewInt(Savings[i])
	}

	return sdk.NewInt(Savings[len(Savings)-1])
}

// Структура хранения данных парамайнинга
type Posmining struct {
	Owner sdk.AccAddress `json:"owner"` // Владелец

	DailyPercent  sdk.Int `json:"daily_percent"`  // Дневной процент начисления парамайнинга
	StructureCoff sdk.Int `json:"structure_coff"` // Коэффициент структуры

	Paramined sdk.Int `json:"paramined"` // Как много токенов уже напарамайнено, но не снято - юзаем для расчета при изменении условий

	LastTransaction time.Time `json:"last_transaction"` // Когда последний раз была исходящая транзакция
	LastCharged     time.Time `json:"last_charged"`     // Когда последний раз был charge (начисление парамайнинга)
}

// Возвращает новый Posmining
func NewPosmining(owner sdk.AccAddress) Posmining {
	return Posmining{
		Owner:         owner,
		Paramined:     sdk.NewInt(0),
		DailyPercent:  sdk.NewInt(0),
		StructureCoff: sdk.NewInt(0),
	}
}
