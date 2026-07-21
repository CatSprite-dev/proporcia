package balancer

import (
	"github.com/CatSprite-dev/proporcia/internal/domain"
	"github.com/shopspring/decimal"
)

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

func LotsToBuy(deficits map[string]domain.Money, prices map[string]domain.PriceInfo) map[string]int {
	result := make(map[string]int)
	for ticker, deficit := range deficits {
		info, ok := prices[ticker]
		if !ok {
			continue
		}
		lotPrice := info.Price.Mul(decimal.NewFromInt(int64(info.LotSize)))
		lots := deficit.Amount.Div(lotPrice).Floor().IntPart()
		result[ticker] = int(lots)
	}
	return result
}

func AllocateCash(deficits map[string]domain.Money, cash domain.Money) map[string]domain.Money {
	totalDeficit := decimal.Zero
	for _, deficit := range deficits {
		totalDeficit = totalDeficit.Add(deficit.Amount)
	}

	result := make(map[string]domain.Money)
	if totalDeficit.LessThanOrEqual(decimal.Zero) {
		return result
	}

	for ticker, deficit := range deficits {
		allocated := deficit.Amount.Div(totalDeficit).Mul(cash.Amount)
		result[ticker] = domain.Money{Amount: allocated, Currency: deficit.Currency}
	}
	return result
}

func BuyPlan(portfolio domain.Portfolio, weights map[string]decimal.Decimal, prices map[string]domain.PriceInfo, cash domain.Money) map[string]int {
	deficits := Deficits(portfolio, weights)
	realDeficits := AllocateCash(deficits, cash)
	lotsToBuy := LotsToBuy(realDeficits, prices)
	return lotsToBuy
}
