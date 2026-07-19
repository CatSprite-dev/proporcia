package balancer

import (
	"github.com/CatSprite-dev/proporcia/internal/domain"
	"github.com/shopspring/decimal"
)

type Balancer struct {
	Fetcher domain.APIClient
}

func NewBalancer(fetcher domain.APIClient) *Balancer {
	return &Balancer{Fetcher: fetcher}
}

func Deficits(portfolio domain.Portfolio, weights map[string]decimal.Decimal) map[string]domain.Money {
	total := portfolio.Total.Amount

	targets := make(map[string]domain.Money)
	for ticker, weight := range weights {
		targets[ticker] = domain.Money{
			Amount:   total.Mul(weight),
			Currency: portfolio.Total.Currency,
		}
	}

	current := make(map[string]domain.Money)
	for ticker, target := range targets {
		found := false
		for _, position := range portfolio.Positions {
			if position.Ticker == ticker {
				current[ticker] = domain.Money{
					Amount:   position.CurrentPrice.Amount.Mul(position.Quantity),
					Currency: position.CurrentPrice.Currency,
				}
				found = true
				break
			}
		}
		if !found {
			current[ticker] = domain.Money{
				Amount:   decimal.Zero,
				Currency: target.Currency,
			}
		}
	}

	result := make(map[string]domain.Money)
	for ticker, target := range targets {
		deficit := target.Amount.Sub(current[ticker].Amount)
		if deficit.LessThanOrEqual(decimal.Zero) {
			continue
		}
		result[ticker] = domain.Money{
			Amount:   deficit,
			Currency: target.Currency,
		}
	}

	return result
}
