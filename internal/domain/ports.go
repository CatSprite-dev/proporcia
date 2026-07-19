package domain

import "context"

type APIClient interface {
	GetAccounts(ctx context.Context, token string) ([]Account, error)
	GetPortfolio(ctx context.Context, token string, accountID string) (Portfolio, error)
}
