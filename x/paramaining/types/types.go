package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"time"
)

var Savings = [...]int64{
	0, // 0-30 days
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

// Структура хранения данных парамайнинга
type Paramining struct {
	Owner sdk.AccAddress `json:"owner"` // Владелец

	DailyPercent sdk.Int `json:"daily_percent"` // Дневной процент начисления парамайнинга
	StructureCoff sdk.Int `json:"structure_coff"` // Коэффициент структуры

	Paramined    sdk.Int `json:"paramined"` // Как много токенов уже напарамайнено, но не снято - юзаем для расчета при изменении условий

	LastTransaction time.Time `json:"last_transaction"` // Когда последний раз была исходящая транзакция
	LastCharged time.Time `json:"last_charged"` // Когда последний раз был charge (начисление парамайнинга)
}


// Возвращает новый Paramining
func NewParamining(owner sdk.AccAddress) Paramining {
	return Paramining{
		Owner: owner,
		Paramined: sdk.NewInt(0),
		DailyPercent: sdk.NewInt(0),
		StructureCoff: sdk.NewInt(0),
	}
}

// Для подсчета начисляемых токенов за какое-то время (сутки, час, минута, секунда)
type CoinsPerTime struct {
	Day sdk.Int `json:"day"`
	Hour sdk.Int `json:"hour"`
	Minute sdk.Int `json:"minute"`
	Second sdk.Int `json:"second"`
}

// Возвращает новый CoinsPerTime
func NewCoinsPerTime() CoinsPerTime {
	return CoinsPerTime{
		Day: sdk.NewInt(0),
		Hour: sdk.NewInt(0),
		Minute: sdk.NewInt(0),
		Second: sdk.NewInt(0),
	}
}

// Для подсчета логики с savings
type ParaminingPeriod struct {
	Total sdk.Int `json:"total"` // общее время
	SavingsCoff sdk.Int `json:"savings_coff"` // кофф накопления
}

// Разница во времени
type TimeDifference struct {
	Days sdk.Int `json:"days"` // Кол-во days
	Hours sdk.Int `json:"hours"` // Кол-во часов
	Minutes sdk.Int `json:"minutes"` // Кол-во минут
	Seconds sdk.Int `json:"seconds"` // Кол-во секунд

	Total sdk.Int `json:"total"` // Общее время в секундах
}

func NewTimeDifference() TimeDifference {
	return TimeDifference{
		Days: sdk.NewInt(0),
		Hours: sdk.NewInt(0),
		Minutes: sdk.NewInt(0),
		Seconds: sdk.NewInt(0),
		Total: sdk.NewInt(0),
	}
}


func InBetween(i sdk.Int, minRaw, maxRaw int64) bool {
	min := sdk.NewInt(minRaw)
	max := sdk.NewInt(maxRaw)

	if i.GTE(min) && i.LTE(max) {
		return true
	} else {
		return false
	}
}

// Возвращает коэффициент структуры в зависимости от баланса
func GetStructureCoff(balance sdk.Int) sdk.Int {
	if balance.LT(sdk.NewInt(1000000000)) {
		return sdk.NewInt(0)
	}

	if InBetween(balance, 1000000000, 9999999999) {
		return sdk.NewInt(218)
	}

	if InBetween(balance, 10000000000, 99999999999) {
		return sdk.NewInt(236)
	}

	if InBetween(balance, 100000000000, 999999999999) {
		return sdk.NewInt(277)
	}

	if InBetween(balance, 1000000000000, 9999999999999) {
		return sdk.NewInt(305)
	}

	if InBetween(balance, 10000000000000, 99999999999999) {
		return sdk.NewInt(336)
	}

	if InBetween(balance, 100000000000000, 999999999999999) {
		return sdk.NewInt(388)
	}

	return sdk.NewInt(437)
}

// Возвращает дневной процент в зависимости от баланса
func GetDailyPercent(balance sdk.Int) sdk.Int {
	if balance.LT(sdk.NewInt(10000)) {
		return sdk.NewInt(0)
	}

	if InBetween(balance, 10000, 99999999) {
		return sdk.NewInt(6)
	}

	if InBetween(balance, 100000000, 999999999) {
		return sdk.NewInt(7)
	}

	if InBetween(balance, 1000000000, 9999999999) {
		return sdk.NewInt(9)
	}

	if InBetween(balance, 10000000000, 49999999999) {
		return sdk.NewInt(10)
	}

	if InBetween(balance, 50000000000, 99999999999) {
		return sdk.NewInt(12)
	}

	if InBetween(balance, 100000000000, 499999999999) {
		return sdk.NewInt(14)
	}

	return sdk.NewInt(16)
}
