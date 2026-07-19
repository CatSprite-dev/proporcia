package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/CatSprite-dev/proporcia/internal/api"
	"github.com/CatSprite-dev/proporcia/internal/balancer"
	"github.com/CatSprite-dev/proporcia/internal/config"
	"github.com/CatSprite-dev/proporcia/internal/fetcher"
	"github.com/CatSprite-dev/proporcia/internal/storage"
	"github.com/CatSprite-dev/proporcia/internal/targets"
	"github.com/shopspring/decimal"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.NewConfig()
	if err != nil {
		logger.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	client := api.NewClient(cfg.BaseURL, logger)

	fetcher := fetcher.NewFetcher(client, logger)

	db, err := storage.NewStorage(cfg.DBPath, logger)
	if err != nil {
		logger.Error("failed to initialize storage", "error", err)
		os.Exit(1)
	}

	if err := db.Init(ctx); err != nil {
		logger.Error("failed to initialize database", "error", err)
		os.Exit(1)
	}

	err = targets.Sync(ctx, *cfg, db, fetcher, logger)
	if err != nil {
		logger.Error("failed to sync targets", "error", err)
		os.Exit(1)
	}

	accounts, err := fetcher.GetAccounts(ctx, cfg.Token)
	if err != nil {
		logger.Error("failed to get accounts", "error", err)
		os.Exit(1)
	}

	if len(accounts) == 0 {
		logger.Error("no accounts found")
		os.Exit(1)
	}

	portfolio, err := fetcher.GetPortfolio(ctx, cfg.Token, accounts[0].ID)
	if err != nil {
		logger.Error("failed to get portfolio", "error", err)
		os.Exit(1)
	}

	logger.Info("portfolio loaded", "positions", len(portfolio.Positions))

	dbTargets, err := db.GetTargets(ctx)
	if err != nil {
		logger.Error("failed to get targets from database", "error", err)
		os.Exit(1)
	}
	weights := make(map[string]decimal.Decimal)
	for _, target := range dbTargets {
		weights[target.Ticker] = target.Weight
	}

	deficits := balancer.Deficits(portfolio, weights)
	prices, err := fetcher.ResolvePrices(ctx, cfg.Token, portfolio.Positions, dbTargets)
	if err != nil {
		logger.Error("failed to resolve prices", "error", err)
		os.Exit(1)
	}

	realDeficits := balancer.AllocateCash(deficits, portfolio.TotalAmountCurrencies)

	lotsToBuy := balancer.LotsToBuy(realDeficits, prices)

	for ticker, lots := range lotsToBuy {
		logger.Info("lots to buy", "ticker", ticker, "lots", lots)
	}
}
