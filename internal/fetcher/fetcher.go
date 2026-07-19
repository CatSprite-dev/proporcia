package fetcher

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/CatSprite-dev/proporcia/internal/api"
	"github.com/CatSprite-dev/proporcia/internal/domain"
)

type Fetcher struct {
	apiClient *api.Client
	logger    *slog.Logger
}

func NewFetcher(apiClient *api.Client, logger *slog.Logger) *Fetcher {
	return &Fetcher{
		apiClient: apiClient,
		logger:    logger,
	}
}

func (f *Fetcher) GetAccounts(ctx context.Context, token string) ([]domain.Account, error) {
	raw, err := f.apiClient.GetAccounts(ctx, token, api.AccountStatusUnspecified)
	if err != nil {
		return nil, fmt.Errorf("get accounts: %w", err)
	}

	accounts := convertAccounts(raw)

	return accounts, nil
}

func (f *Fetcher) GetPortfolio(ctx context.Context, token string, accountID string) (domain.Portfolio, error) {
	raw, err := f.apiClient.GetPortfolio(ctx, token, accountID)
	if err != nil {
		return domain.Portfolio{}, fmt.Errorf("get portfolio: %w", err)
	}

	portfolio := convertPortfolio(raw)

	return portfolio, nil
}
