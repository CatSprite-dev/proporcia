package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func (client *Client) GetAccounts(ctx context.Context, token string, accountStatus AccountStatus) (Accounts, error) {
	type AccountsRequest struct {
		Status AccountStatus `json:"status,omitempty"`
	}

	url := client.baseURL + "/rest/tinkoff.public.invest.api.contract.v1.UsersService/GetAccounts"

	payload := AccountsRequest{Status: accountStatus}

	data, err := client.DoRequest(ctx, url, http.MethodPost, token, payload)
	if err != nil {
		return Accounts{}, fmt.Errorf("do request error (api.GetAccounts): %w", err)
	}

	var accounts Accounts
	err = json.Unmarshal(data, &accounts)
	if err != nil {
		return Accounts{}, fmt.Errorf("unmarshal error (api.GetAccounts): %w", err)
	}
	return accounts, nil
}

func (client *Client) GetPortfolio(ctx context.Context, token string, accountID string) (Portfolio, error) {
	type PortfolioRequest struct {
		AccountID string `json:"accountId"`
	}

	url := client.baseURL + "/rest/tinkoff.public.invest.api.contract.v1.OperationsService/GetPortfolio"

	payload := PortfolioRequest{AccountID: accountID}

	data, err := client.DoRequest(ctx, url, http.MethodPost, token, payload)
	if err != nil {
		return Portfolio{}, fmt.Errorf("do request error (GetPortfolio): %w", err)
	}

	var userPortfolio Portfolio
	err = json.Unmarshal(data, &userPortfolio)
	if err != nil {
		return Portfolio{}, fmt.Errorf("unmarshal error (GetPortfolio): %w", err)
	}

	return userPortfolio, nil
}
