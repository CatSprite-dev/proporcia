package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/CatSprite-dev/proporcia/internal/api"
	"github.com/CatSprite-dev/proporcia/internal/config"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	cfg, err := config.NewConfig()
	if err != nil {
		logger.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	client := api.NewClient(cfg.BaseURL, logger)

	accounts, err := client.GetAccounts(context.Background(), cfg.Token, api.AccountStatusUnspecified)
	if err != nil {
		logger.Error("failed to get accounts", "error", err)
		os.Exit(1)
	}

	if len(accounts.Accounts) == 0 {
		logger.Error("no accounts found")
		os.Exit(1)
	}

	portfolio, err := client.GetPortfolio(context.Background(), cfg.Token, accounts.Accounts[0].ID)
	if err != nil {
		logger.Error("failed to get portfolio", "error", err)
		os.Exit(1)
	}

	logger.Info("portfolio loaded", "positions", len(portfolio.Positions))
}
