package main

import (
	"context"
	"fmt"
	"os"

	"github.com/CatSprite-dev/proporcia/internal/api"
	"github.com/CatSprite-dev/proporcia/internal/config"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		os.Exit(1)
	}

	client := api.NewClient(cfg.BaseURL)

	token := os.Getenv("TOKEN")

	accounts, err := client.GetAccounts(context.Background(), token, api.AccountStatusUnspecified)
	if err != nil {
		fmt.Printf("get account error: %v", err)
	}

	portfolio, err := client.GetPortfolio(context.Background(), token, accounts.Accounts[0].ID)
	if err != nil {
		fmt.Printf("get portfolio error: %v", err)
	}

	fmt.Printf("Portfolio: %+v\n", portfolio)
}
