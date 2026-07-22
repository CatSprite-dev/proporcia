package app

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/CatSprite-dev/proporcia/internal/api"
	"github.com/CatSprite-dev/proporcia/internal/balancer"
	"github.com/CatSprite-dev/proporcia/internal/config"
	"github.com/CatSprite-dev/proporcia/internal/fetcher"
	"github.com/CatSprite-dev/proporcia/internal/storage"
	"github.com/CatSprite-dev/proporcia/internal/targets"
	"github.com/CatSprite-dev/proporcia/internal/trade"
	"github.com/shopspring/decimal"
)

type App struct {
	cfg            *config.Config
	client         *api.Client
	fetcherService *fetcher.Fetcher
	orderService   *trade.OrderService
	logger         *slog.Logger
	db             *storage.Storage
}

func New(ctx context.Context, logger *slog.Logger) (*App, error) {
	cfg, err := config.NewConfig()
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	client := api.NewClient(cfg.BaseURL, logger)

	fetcherService := fetcher.NewFetcher(client, logger)

	db, err := storage.NewStorage(cfg.DBPath, logger)
	if err != nil {
		return nil, fmt.Errorf("storage init error: %w", err)
	}

	err = db.Init(ctx)
	if err != nil {
		return nil, fmt.Errorf("db init error: %w", err)
	}

	orderService := trade.NewOrderService(client, db, logger)

	return &App{
		cfg:            cfg,
		client:         client,
		fetcherService: fetcherService,
		orderService:   orderService,
		logger:         logger,
		db:             db,
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	err := targets.Sync(ctx, *a.cfg, a.db, a.fetcherService, a.logger)
	if err != nil {
		return fmt.Errorf("failed to sync targets: %w", err)
	}

	accounts, err := a.fetcherService.GetAccounts(ctx, a.cfg.Token)
	if err != nil || len(accounts) == 0 {
		return fmt.Errorf("failed to get accounts or empty: %w", err)
	}

	portfolio, err := a.fetcherService.GetPortfolio(ctx, a.cfg.Token, accounts[0].ID)
	if err != nil {
		return fmt.Errorf("failed to get portfolio: %w", err)
	}

	dbTargets, err := a.db.GetTargets(ctx)
	if err != nil {
		return fmt.Errorf("failed to get targets: %w", err)
	}

	weights := make(map[string]decimal.Decimal)
	for _, target := range dbTargets {
		weights[target.InstrumentUID] = target.Weight
	}

	prices, err := a.fetcherService.ResolvePrices(ctx, a.cfg.Token, portfolio.Positions, dbTargets)
	if err != nil {
		return fmt.Errorf("failed to resolve prices: %w", err)
	}

	buyPlan := balancer.BuyPlan(portfolio, weights, prices, portfolio.TotalAmountCurrencies)

	orders, err := a.orderService.Buy(ctx, a.cfg.Token, portfolio.AccountID, buyPlan)
	if err != nil {
		return fmt.Errorf("failed to buy targets: %w", err)
	}

	err = a.orderService.SaveOrders(ctx, orders)
	if err != nil {
		a.logger.Error("failed to save orders", "error", err)
	}

	for _, resp := range orders {
		a.logger.Info("order executed",
			"ticker", resp.Ticker,
			"order_id", resp.OrderID,
			"lots_executed", resp.LotsExecuted,
			"total_amount", resp.TotalOrderAmount.Amount,
		)
	}

	return nil
}

func (a *App) Close() {
	if err := a.db.Close(); err != nil {
		a.logger.Error("close storage", "error", err)
	}
}
