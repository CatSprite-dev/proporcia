package fetcher

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/CatSprite-dev/proporcia/internal/api"
	"github.com/CatSprite-dev/proporcia/internal/domain"
	"github.com/shopspring/decimal"
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

func (f *Fetcher) ResolvePrices(ctx context.Context, token string, positions []domain.Position, targets []domain.Target) (map[string]domain.PriceInfo, error) {
	prices := make(map[string]domain.PriceInfo, len(targets))
	var missingTargets []domain.Target

	for _, target := range targets {
		found := false
		for _, position := range positions {
			if position.InstrumentUID == target.InstrumentUID {
				prices[target.InstrumentUID] = domain.PriceInfo{
					Price:   position.CurrentPrice.Amount,
					LotSize: target.Lot,
				}
				found = true
				break
			}
		}
		if !found {
			missingTargets = append(missingTargets, target)
		}
	}

	if len(missingTargets) == 0 {
		return prices, nil
	}

	instrumentIDs := make([]string, 0, len(missingTargets))
	for _, target := range missingTargets {
		if target.InstrumentUID != "" {
			instrumentIDs = append(instrumentIDs, target.InstrumentUID)
		}
	}

	if len(instrumentIDs) == 0 {
		return nil, fmt.Errorf("missing targets have no valid UIDs")
	}
	lastPricesResp, err := f.apiClient.GetLastPrices(
		ctx,
		token,
		instrumentIDs,
		api.LastPriceUnspecified,
		api.InstrumentStatusBase,
	)
	if err != nil {
		return nil, fmt.Errorf("get last prices batch: %w", err)
	}

	convertedPrices := convertLastPrices(lastPricesResp)

	priceMap := make(map[string]domain.LastPrice, len(convertedPrices))
	for _, p := range convertedPrices {
		if p.InstrumentUID != "" {
			priceMap[p.InstrumentUID] = p
		}
	}

	for _, target := range missingTargets {
		lastPrice, ok := priceMap[target.InstrumentUID]
		if !ok {
			return nil, fmt.Errorf("no price returned for UID %s (%s)", target.InstrumentUID, target.Ticker)
		}

		if lastPrice.Price.IsZero() || lastPrice.Price.LessThanOrEqual(decimal.Zero) {
			return nil, fmt.Errorf("zero or negative price returned for %s (%s)", target.Ticker, target.InstrumentUID)
		}

		prices[target.InstrumentUID] = domain.PriceInfo{
			Price:   lastPrice.Price,
			LotSize: target.Lot,
		}
	}

	return prices, nil
}
