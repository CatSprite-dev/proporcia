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
