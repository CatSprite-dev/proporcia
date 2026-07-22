package domain

import (
	"time"

	"github.com/shopspring/decimal"
)

type Money struct {
	Amount   decimal.Decimal
	Currency string
}

type PriceInfo struct {
	Price   decimal.Decimal
	LotSize int
}

type Target struct {
	ID            int64           `json:"id"`
	Name          string          `json:"name"`
	Ticker        string          `json:"ticker"`
	Weight        decimal.Decimal `json:"weight"`
	InstrumentUID string          `json:"instrumentUid"`
	ClassCode     string          `json:"class_code"`
	Lot           int             `json:"lot"`
	Type          string          `json:"type"`
}

type Account struct {
	ID string
}

type Portfolio struct {
	AccountID             string `json:"accountId"`
	Positions             []Position
	Total                 Money `json:"total"`
	TotalAmountCurrencies Money `json:"totalAmountCurrencies"`
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
	InstrumentUID         string `json:"instrumentUid"`
	Name                  string `json:"name"`
}

type LastPrice struct {
	ClassCode     string          `json:"classCode"`
	Ticker        string          `json:"ticker"`
	Price         decimal.Decimal `json:"price"`
	InstrumentUID string          `json:"instrumentUid"`
	Figi          string          `json:"figi"`
	Time          time.Time       `json:"time"`
}

type PostOrderResponse struct {
	ClassCode string `json:"classCode"`
	Ticker    string `json:"ticker"`
	OrderID   string `json:"orderId"`
	Figi      string `json:"figi"`

	// InitialOrderPrice — начальная цена заявки (произведение количества запрошенных лотов на цену).
	InitialOrderPrice Money `json:"initialOrderPrice"`

	// InitialCommission — начальная комиссия, рассчитанная при выставлении заявки.
	InitialCommission Money `json:"initialCommission"`

	// Message — дополнительные данные об исполнении заявки.
	Message string `json:"message"`

	// LotsExecuted — количество исполненных лотов.
	LotsExecuted string `json:"lotsExecuted"`

	// TotalOrderAmount — итоговая стоимость заявки, включающая все комиссии.
	TotalOrderAmount Money `json:"totalOrderAmount"`

	// LotsRequested — количество запрошенных лотов.
	LotsRequested string `json:"lotsRequested"`

	// InstrumentUID — UID-идентификатор инструмента.
	InstrumentUID string `json:"instrumentUid"`

	// OrderRequestID — идентификатор ключа идемпотентности, переданный клиентом (в формате UID, до 36 символов).
	OrderRequestID string `json:"orderRequestId"`

	// ExecutedOrderPrice — исполненная средняя цена одного инструмента в заявке.
	ExecutedOrderPrice Money `json:"executedOrderPrice"`

	// ExecutedCommission — фактическая комиссия по итогам исполнения заявки.
	ExecutedCommission Money `json:"executedCommission"`

	// InitialSecurityPrice — начальная цена за 1 инструмент (для получения стоимости лота требуется умножить на лотность).
	InitialSecurityPrice Money `json:"initialSecurityPrice"`

	// ResponseMetadata — метаданные ответа от сервера (время сервера, tracking_id).
	ResponseMetadata struct {
		ServerTime time.Time `json:"serverTime"`
		TrackingID string    `json:"trackingId"`
	} `json:"responseMetadata"`
}
