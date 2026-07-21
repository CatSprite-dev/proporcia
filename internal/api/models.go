package api

import "time"

type MoneyValue struct {
	Currency string `json:"currency"`
	Units    string `json:"units"`
	Nano     int    `json:"nano"`
}

type Quotation struct {
	Units string `json:"units"`
	Nano  int    `json:"nano"`
}

type Accounts struct {
	Accounts []struct {
		ID          string    `json:"id"`
		Type        string    `json:"type"`
		Name        string    `json:"name"`
		Status      string    `json:"status"`
		OpenedDate  time.Time `json:"openedDate"`
		ClosedDate  time.Time `json:"closedDate"`
		AccessLevel string    `json:"accessLevel"`
	} `json:"accounts"`
}

type Portfolio struct {
	TotalAmountShares     MoneyValue `json:"totalAmountShares"`
	TotalAmountBonds      MoneyValue `json:"totalAmountBonds"`
	TotalAmountEtf        MoneyValue `json:"totalAmountEtf"`
	TotalAmountCurrencies MoneyValue `json:"totalAmountCurrencies"`
	TotalAmountFutures    MoneyValue `json:"totalAmountFutures"`
	ExpectedYield         Quotation  `json:"expectedYield"`
	Positions             []Position `json:"positions"`
	AccountID             string     `json:"accountId"`
	TotalAmountOptions    MoneyValue `json:"totalAmountOptions"`
	TotalAmountSp         MoneyValue `json:"totalAmountSp"`
	TotalAmountPortfolio  MoneyValue `json:"totalAmountPortfolio"`
	VirtualPositions      []any      `json:"virtualPositions"`
	DailyYield            MoneyValue `json:"dailyYield"`
	DailyYieldRelative    Quotation  `json:"dailyYieldRelative"`
}

type Position struct {
	Figi                     string     `json:"figi"`
	InstrumentType           string     `json:"instrumentType"`
	Quantity                 Quotation  `json:"quantity"`
	AveragePositionPrice     MoneyValue `json:"averagePositionPrice"`
	ExpectedYield            Quotation  `json:"expectedYield"`
	AveragePositionPricePt   Quotation  `json:"averagePositionPricePt"`
	CurrentPrice             MoneyValue `json:"currentPrice"`
	AveragePositionPriceFifo MoneyValue `json:"averagePositionPriceFifo"`
	QuantityLots             Quotation  `json:"quantityLots"`
	Blocked                  bool       `json:"blocked"`
	BlockedLots              Quotation  `json:"blockedLots"`
	PositionUID              string     `json:"positionUid"`
	InstrumentUID            string     `json:"instrumentUid"`
	VarMargin                MoneyValue `json:"varMargin"`
	ExpectedYieldFifo        Quotation  `json:"expectedYieldFifo"`
	DailyYield               MoneyValue `json:"dailyYield"`
	Ticker                   string     `json:"ticker"`
	ClassCode                string     `json:"classCode"`
	CurrentNkd               MoneyValue `json:"currentNkd,omitempty"`
}

type Instruments struct {
	Instruments []struct {
		WeekendFlag           bool      `json:"weekendFlag"`
		ClassCode             string    `json:"classCode"`
		Ticker                string    `json:"ticker"`
		InstrumentType        string    `json:"instrumentType"`
		ForQualInvestorFlag   bool      `json:"forQualInvestorFlag"`
		ForIisFlag            bool      `json:"forIisFlag"`
		PositionUID           string    `json:"positionUid"`
		Figi                  string    `json:"figi"`
		APITradeAvailableFlag bool      `json:"apiTradeAvailableFlag"`
		First1MinCandleDate   time.Time `json:"first1minCandleDate"`
		Lot                   int       `json:"lot"`
		UID                   string    `json:"uid"`
		BlockedTcaFlag        bool      `json:"blockedTcaFlag"`
		Name                  string    `json:"name"`
		First1DayCandleDate   time.Time `json:"first1dayCandleDate"`
		Isin                  string    `json:"isin"`
	} `json:"instruments"`
}

type LastPrices struct {
	LastPrices []struct {
		ClassCode     string    `json:"classCode"`
		Ticker        string    `json:"ticker"`
		Price         Quotation `json:"price"`
		InstrumentUID string    `json:"instrumentUid"`
		Figi          string    `json:"figi"`
		Time          time.Time `json:"time"`
	} `json:"lastPrices"`
}

// PostOrderResponse содержит информацию о выставленном торговом поручении.
type PostOrderResponse struct {
	// ClassCode — класс-код (секция торгов).
	ClassCode string `json:"classCode"`

	// Ticker — тикер инструмента.
	Ticker string `json:"ticker"`

	// OrderID — биржевой идентификатор заявки.
	OrderID string `json:"orderId"`

	// Figi — FIGI-идентификатор инструмента.
	Figi string `json:"figi"`

	// InitialOrderPrice — начальная цена заявки (произведение количества запрошенных лотов на цену).
	InitialOrderPrice MoneyValue `json:"initialOrderPrice"`

	// InitialCommission — начальная комиссия, рассчитанная при выставлении заявки.
	InitialCommission MoneyValue `json:"initialCommission"`

	// Message — дополнительные данные об исполнении заявки.
	Message string `json:"message"`

	// LotsExecuted — количество исполненных лотов.
	LotsExecuted string `json:"lotsExecuted"`

	// TotalOrderAmount — итоговая стоимость заявки, включающая все комиссии.
	TotalOrderAmount MoneyValue `json:"totalOrderAmount"`

	// LotsRequested — количество запрошенных лотов.
	LotsRequested string `json:"lotsRequested"`

	// InitialOrderPricePt — начальная цена заявки в пунктах (для фьючерсов).
	InitialOrderPricePt Quotation `json:"initialOrderPricePt"`

	// InstrumentUID — UID-идентификатор инструмента.
	InstrumentUID string `json:"instrumentUid"`

	// OrderRequestID — идентификатор ключа идемпотентности, переданный клиентом (в формате UID, до 36 символов).
	OrderRequestID string `json:"orderRequestId"`

	// ExecutedOrderPrice — исполненная средняя цена одного инструмента в заявке.
	ExecutedOrderPrice MoneyValue `json:"executedOrderPrice"`

	// ExecutedCommission — фактическая комиссия по итогам исполнения заявки.
	ExecutedCommission MoneyValue `json:"executedCommission"`

	// InitialSecurityPrice — начальная цена за 1 инструмент (для получения стоимости лота требуется умножить на лотность).
	InitialSecurityPrice MoneyValue `json:"initialSecurityPrice"`

	// AciValue — значение накопленного купонного дохода (НКД) на дату выставления торговых поручений.
	AciValue MoneyValue `json:"aciValue"`

	// ResponseMetadata — метаданные ответа от сервера (время сервера, tracking_id).
	ResponseMetadata struct {
		ServerTime time.Time `json:"serverTime"`
		TrackingID string    `json:"trackingId"`
	} `json:"responseMetadata"`
}
