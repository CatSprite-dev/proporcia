package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
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

	url := c.baseURL + "/rest/tinkoff.public.invest.api.contract.v1.InstrumentsService/FindInstruments"

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
