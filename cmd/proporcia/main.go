package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/CatSprite-dev/proporcia/internal/app"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	logger.Info("starting application...")

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	application, err := app.New(ctx, logger)
	if err != nil {
		logger.Error("failed to init app", "error", err)
		os.Exit(1)
	}
	defer application.Close()

	logger.Info("running portfolio rebalance...")
	err = application.Run(ctx)
	if err != nil {
		logger.Error("application run failed", "error", err)
		os.Exit(1)
	}

	logger.Info("application finished successfully")
}
