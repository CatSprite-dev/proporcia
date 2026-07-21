package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

func (c *Client) GetAccounts(ctx context.Context, token string, accountStatus AccountStatus) (Accounts, error) {
	type AccountsRequest struct {
		Status AccountStatus `json:"status,omitempty"`
	}

	url := c.baseURL + "/rest/tinkoff.public.invest.api.contract.v1.UsersService/GetAccounts"

	payload := AccountsRequest{Status: accountStatus}

	data, err := c.DoRequest(ctx, url, http.MethodPost, token, payload)
	if err != nil {
		return Accounts{}, fmt.Errorf("get accounts: %w", err)
	}

	var accounts Accounts
	if err := json.Unmarshal(data, &accounts); err != nil {
		return Accounts{}, fmt.Errorf("unmarshal accounts: %w", err)
	}

	c.logger.Debug("accounts fetched", "count", len(accounts.Accounts))

	return accounts, nil
}

func (c *Client) GetPortfolio(ctx context.Context, token string, accountID string) (Portfolio, error) {
	type PortfolioRequest struct {
		AccountID string `json:"accountId"`
	}

	url := c.baseURL + "/rest/tinkoff.public.invest.api.contract.v1.OperationsService/GetPortfolio"

	payload := PortfolioRequest{AccountID: accountID}

	data, err := c.DoRequest(ctx, url, http.MethodPost, token, payload)
	if err != nil {
		return Portfolio{}, fmt.Errorf("get portfolio: %w", err)
	}

	var portfolio Portfolio
	if err := json.Unmarshal(data, &portfolio); err != nil {
		return Portfolio{}, fmt.Errorf("unmarshal portfolio: %w", err)
	}

	c.logger.Debug("portfolio fetched", "account_id", accountID, "positions", len(portfolio.Positions))

	return portfolio, nil
}

func (c *Client) FindInstruments(ctx context.Context, token string, query string, instrumentKind InstrumentType, apiTradeAvailableFlag bool) (Instruments, error) {
	type FindInstrumentsRequest struct {
		Query                 string         `json:"query"`
		InstrumentType        InstrumentType `json:"instrumentType,omitempty"`
		ApiTradeAvailableFlag bool           `json:"apiTradeAvailableFlag,omitempty"`
	}

	url := c.baseURL + "/rest/tinkoff.public.invest.api.contract.v1.InstrumentsService/FindInstrument"

	payload := FindInstrumentsRequest{
		Query:                 query,
		InstrumentType:        instrumentKind,
		ApiTradeAvailableFlag: apiTradeAvailableFlag,
	}

	data, err := c.DoRequest(ctx, url, http.MethodPost, token, payload)
	if err != nil {
		return Instruments{}, fmt.Errorf("find instruments: %w", err)
	}

	var instruments Instruments
	if err := json.Unmarshal(data, &instruments); err != nil {
		return Instruments{}, fmt.Errorf("unmarshal instruments: %w", err)
	}

	c.logger.Debug("instruments fetched", "query", query, "count", len(instruments.Instruments))

	return instruments, nil
}

func (c *Client) GetLastPrices(
	ctx context.Context,
	token string,
	instrumentIDs []string,
	priceType LastPriceType,
	instrStatus InstrumentStatus,
) (LastPrices, error) {
	type GetLastPricesRequest struct {
		InstrumentIDs    []string         `json:"instrumentIds"`
		LastPriceType    LastPriceType    `json:"lastPriceType,omitempty"`
		InstrumentStatus InstrumentStatus `json:"instrumentStatus,omitempty"`
	}

	url := c.baseURL + "/rest/tinkoff.public.invest.api.contract.v1.MarketDataService/GetLastPrices"

	payload := GetLastPricesRequest{
		InstrumentIDs:    instrumentIDs,
		LastPriceType:    priceType,
		InstrumentStatus: instrStatus,
	}

	data, err := c.DoRequest(ctx, url, http.MethodPost, token, payload)
	if err != nil {
		return LastPrices{}, fmt.Errorf("get last prices: %w", err)
	}

	var lastPrices LastPrices
	if err := json.Unmarshal(data, &lastPrices); err != nil {
		return LastPrices{}, fmt.Errorf("unmarshal last prices: %w", err)
	}

	c.logger.Debug("last prices fetched", "instrument count", len(lastPrices.LastPrices))

	return lastPrices, nil
}

func (c *Client) PostOrder(
	ctx context.Context,
	token string,
	quantity string,
	price Quotation,
	direction OrderDirection,
	accountID string,
	orderType OrderType,
	orderID string,
	instrumentID string,
	timeInForce TimeInForce,
	priceType PriceType,
	confirmMarginTrade bool,
) (PostOrderResponse, error) {

	if orderID == "" {
		orderID = uuid.NewString()
	}

	type PostOrderRequest struct {
		Quantity           string         `json:"quantity"`
		Price              Quotation      `json:"price,omitempty"`
		Direction          OrderDirection `json:"direction"`
		AccountID          string         `json:"accountId"`
		OrderType          OrderType      `json:"orderType"`
		OrderID            string         `json:"orderId"`
		InstrumentID       string         `json:"instrumentId"`
		TimeInForce        TimeInForce    `json:"timeInForce,omitempty"`
		PriceType          PriceType      `json:"priceType,omitempty"`
		ConfirmMarginTrade bool           `json:"confirmMarginTrade,omitempty"`
	}

	url := c.baseURL + "/rest/tinkoff.public.invest.api.contract.v1.OrdersService/PostOrder"

	payload := PostOrderRequest{
		Quantity:           quantity,
		Price:              price,
		Direction:          direction,
		AccountID:          accountID,
		OrderType:          orderType,
		OrderID:            orderID,
		InstrumentID:       instrumentID,
		TimeInForce:        timeInForce,
		PriceType:          priceType,
		ConfirmMarginTrade: confirmMarginTrade,
	}

	data, err := c.DoRequest(ctx, url, http.MethodPost, token, payload)
	if err != nil {
		return PostOrderResponse{}, fmt.Errorf("post order: %w", err)
	}

	var resp PostOrderResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return PostOrderResponse{}, fmt.Errorf("unmarshal post order response: %w", err)
	}

	c.logger.Debug("order posted", "order_id", resp.OrderID, "ticker", resp.Ticker)

	return resp, nil
}
