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

func (f *Fetcher) FindInstrument(ctx context.Context, token string, ticker string) (domain.Instrument, error) {
	raw, err := f.apiClient.FindInstruments(ctx, token, ticker, api.InstrumentTypeUnspecified, true)
	if err != nil {
		return domain.Instrument{}, fmt.Errorf("find instrument: %w", err)
	}

	instruments := convertInstruments(raw)

	if len(instruments) == 0 {
		return domain.Instrument{}, fmt.Errorf("instrument not found: %s", ticker)
	}

	var result domain.Instrument
	found := false
	for _, instrument := range instruments {
		if instrument.Ticker == ticker && instrument.ClassCode == "TQBR" {
			result = instrument
			found = true
			break
		}
	}

	if !found {
		return domain.Instrument{}, fmt.Errorf("instrument not found: %s", ticker)
	}

	return result, nil
}
