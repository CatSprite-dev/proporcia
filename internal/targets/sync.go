package targets

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"github.com/CatSprite-dev/proporcia/internal/config"
	"github.com/CatSprite-dev/proporcia/internal/domain"
	"github.com/CatSprite-dev/proporcia/internal/fetcher"
	"github.com/CatSprite-dev/proporcia/internal/storage"
	"github.com/shopspring/decimal"
)

// allowedInstrumentTypes ограничивает бота только акциями и облигациями
var allowedInstrumentTypes = map[string]bool{
	"share": true,
	"bond":  true,
}

func isAllowedInstrumentType(instrumentType string) bool {
	return allowedInstrumentTypes[instrumentType]
}

type discTarget struct {
	Ticker string          `json:"ticker"`
	Weight decimal.Decimal `json:"weight"`
}

func Sync(ctx context.Context, cfg config.Config, storage *storage.Storage, fetcher *fetcher.Fetcher, logger *slog.Logger) error {
	targetsDisk, err := os.ReadFile(cfg.TargetsPath)
	if err != nil {
		return fmt.Errorf("read targets file: %w", err)
	}
	var discTargets []discTarget
	if err := json.Unmarshal(targetsDisk, &discTargets); err != nil {
		return fmt.Errorf("unmarshal targets: %w", err)
	}

	dbTargets, err := storage.GetTargets(ctx)
	if err != nil {
		return fmt.Errorf("get targets from db: %w", err)
	}

	dbByTicker := make(map[string]domain.Target, len(dbTargets))
	for _, t := range dbTargets {
		dbByTicker[t.Ticker] = t
	}

	diskTickers := make(map[string]struct{}, len(discTargets))

	added := 0
	updated := 0
	for _, t := range discTargets {
		diskTickers[t.Ticker] = struct{}{}

		target, exists := dbByTicker[t.Ticker]
		if !exists {
			instrument, err := fetcher.FindInstrument(ctx, cfg.Token, t.Ticker)
			if err != nil {
				return fmt.Errorf("find instrument for ticker %s: %w", t.Ticker, err)
			}
			if !isAllowedInstrumentType(instrument.InstrumentType) {
				return fmt.Errorf("instrument type not allowed for %s: %s", t.Ticker, instrument.InstrumentType)
			}
			target = domain.Target{
				Name:          instrument.Name,
				Ticker:        t.Ticker,
				InstrumentUID: instrument.InstrumentUID,
				ClassCode:     instrument.ClassCode,
				Lot:           instrument.Lot,
				Type:          instrument.InstrumentType,
			}
			added++
		} else {
			updated++
		}
		target.Weight = t.Weight

		if err := storage.UpsertTarget(ctx, target); err != nil {
			return fmt.Errorf("upsert target %s: %w", t.Ticker, err)
		}

	}

	deleted := 0
	for _, t := range dbTargets {
		if _, exists := diskTickers[t.Ticker]; !exists {
			if err := storage.DeleteTarget(ctx, t.Ticker); err != nil {
				return fmt.Errorf("delete target %s: %w", t.Ticker, err)
			}
			deleted++
		}
	}

	logger.Info("sync completed", "added", added, "deleted", deleted, "updated", updated)

	return nil
}
