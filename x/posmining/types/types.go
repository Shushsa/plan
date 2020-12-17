package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var Savings = [...]int64{
	0,   // 0-30 days
	160, // 1.60 or 50%, 30-60 days
	160, // 1.60 or 51%, 60-90 days
	180, // 1.80 or 52%, 90-120 days
	180, // 1.80 or 53%, 120-150 days
	180, // 1.80 or 54%, 150-180 days
	200, // 2.00 or 55%, 180-210 days
	200, // 2.00 or 55%, 210-240 days
	200, // 2.00 or 55%, 240-270 days
	220, // 2.20 or 55%, 270-300 days
	220, // 2.20 or 55%, 300-330 days
	220, // 2.20 or 55%, 330-360 days
	220, // 2.20 or 100%, >360 days
}

func GetSavingCoff(i int) sdk.Int {
	if len(Savings) > i {
		return sdk.NewInt(Savings[i])
	}

	return sdk.NewInt(Savings[len(Savings)-1])
}

// Структура хранения данных POS-майнинга
type Posmining struct {
	Owner           sdk.AccAddress `json:"owner"`            // Владелец
	DailyPercent    sdk.Int        `json:"daily_percent"`    // Дневной процент начисления майнинга
	StructureCoff   sdk.Int        `json:"structure_coff"`   // Коэффициент структуры
	Posmined        sdk.Int        `json:"posmined"`         // Сколько уже намайнено, но не снято - для расчета при изменении условий
	LastTransaction time.Time      `json:"last_transaction"` // Когда последний раз была исходящая транзакция
	LastCharged     time.Time      `json:"last_charged"`     // Когда последний раз был charge (начисление майнинга)
}

// Возвращает новый Posmining
func NewPosmining(owner sdk.AccAddress) Posmining {
	return Posmining{
		Owner:         owner,
		Posmined:      sdk.NewInt(0),
		DailyPercent:  sdk.NewInt(0),
		StructureCoff: sdk.NewInt(0),
	}
}

// Current Correction
type Correction struct {
	StartDate           time.Time            `json:"start_date"`           // datetime of the updated coff
	OpeningPrice        sdk.Int              `json:"opening_price"`        // the market price being used
	PreviousCorrections []PreviousCorrection `json:"previous_corrections"` // previous regulation periods
}

// Updates the regulation when we get new market data
func (t *Correction) Update(current time.Time, price sdk.Int, coff sdk.Int) {
	prev := PreviousCorrection{
		StartDate:    t.StartDate,
		EndDate:      current,
		OpeningPrice: t.OpeningPrice,
	}

	t.PreviousCorrections = append([]PreviousCorrection{prev}, t.PreviousCorrections...)

	t.StartDate = current
	t.OpeningPrice = price
}

type PreviousCorrection struct {
	StartDate    time.Time `json:"start_date"`    // дата и время начала регуляции
	EndDate      time.Time `json:"end_date"`      // дата и время конца регуляции
	OpeningPrice sdk.Int   `json:"opening_price"` // цена, при которой поменялась регуляция
}
