package domain

import "github.com/shopspring/decimal"

type Money struct {
	Amount   decimal.Decimal
	Currency string
}

type Account struct {
	ID string
}

type Portfolio struct {
	AccountID string `json:"accountId"`
	Positions []Position
}

type Position struct {
	Figi                 string          `json:"figi"`
	InstrumentType       string          `json:"instrumentType"`
	Quantity             decimal.Decimal `json:"quantity"`
	AveragePositionPrice Money           `json:"averagePositionPrice"`
	CurrentPrice         Money           `json:"currentPrice"`
	QuantityLots         decimal.Decimal `json:"quantityLots"`
	PositionUID          string          `json:"positionUid"`
	InstrumentUID        string          `json:"instrumentUid"`
	Ticker               string          `json:"ticker"`
	ClassCode            string          `json:"classCode"`
}
