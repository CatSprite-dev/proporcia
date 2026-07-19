package domain

import (
	"github.com/shopspring/decimal"
)

type Money struct {
	Amount   decimal.Decimal
	Currency string
}

type Target struct {
	ID        int64           `json:"id"`
	Name      string          `json:"name"`
	Ticker    string          `json:"ticker"`
	Weight    decimal.Decimal `json:"weight"`
	UID       string          `json:"uid"`
	ClassCode string          `json:"class_code"`
	Lot       int             `json:"lot"`
	Type      string          `json:"type"`
}

type Account struct {
	ID string
}

type Portfolio struct {
	AccountID string `json:"accountId"`
	Positions []Position
	Total     Money `json:"total"`
}

type Position struct {
	Figi                 string          `json:"figi"`
	InstrumentType       string          `json:"instrumentType"`
	Quantity             decimal.Decimal `json:"quantity"`
	AveragePositionPrice Money           `json:"averagePositionPrice"`
	CurrentPrice         Money           `json:"currentPrice"`
	PositionUID          string          `json:"positionUid"`
	InstrumentUID        string          `json:"instrumentUid"`
	Ticker               string          `json:"ticker"`
	ClassCode            string          `json:"classCode"`
}

type Instrument struct {
	ClassCode             string `json:"classCode"`
	Ticker                string `json:"ticker"`
	InstrumentType        string `json:"instrumentType"`
	PositionUID           string `json:"positionUid"`
	Figi                  string `json:"figi"`
	APITradeAvailableFlag bool   `json:"apiTradeAvailableFlag"`
	Lot                   int    `json:"lot"`
	UID                   string `json:"uid"`
	Name                  string `json:"name"`
}
