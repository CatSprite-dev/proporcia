package balancer

import (
	"github.com/CatSprite-dev/proporcia/internal/domain"
	"github.com/shopspring/decimal"
)

func Deficits(portfolio domain.Portfolio, weights map[string]decimal.Decimal) map[string]domain.Money {
	total := portfolio.Total.Amount

	positionsByUID := make(map[string]domain.Position, len(portfolio.Positions))
	for _, pos := range portfolio.Positions {
		positionsByUID[pos.InstrumentUID] = pos
	}

	result := make(map[string]domain.Money)

	for uid, weight := range weights {
		targetAmount := total.Mul(weight)

		currentAmount := decimal.Zero
		if pos, found := positionsByUID[uid]; found {
			currentAmount = pos.CurrentPrice.Amount.Mul(pos.Quantity)
		}

		deficit := targetAmount.Sub(currentAmount)
		if deficit.GreaterThan(decimal.Zero) {
			result[uid] = domain.Money{
				Amount:   deficit,
				Currency: portfolio.Total.Currency,
			}
		}
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

	for UID, deficit := range deficits {
		allocated := deficit.Amount.Div(totalDeficit).Mul(cash.Amount)
		result[UID] = domain.Money{Amount: allocated, Currency: deficit.Currency}
	}
	return result
}

func LotsToBuy(deficits map[string]domain.Money, prices map[string]domain.PriceInfo) map[string]int {
	result := make(map[string]int)
	for UID, deficit := range deficits {
		info, ok := prices[UID]
		if !ok {
			continue
		}
		lotPrice := info.Price.Mul(decimal.NewFromInt(int64(info.LotSize)))
		lots := deficit.Amount.Div(lotPrice).Floor().IntPart()
		result[UID] = int(lots)
	}
	return result
}

func BuyPlan(portfolio domain.Portfolio, weights map[string]decimal.Decimal, prices map[string]domain.PriceInfo, cash domain.Money) map[string]int {
	deficits := Deficits(portfolio, weights)
	realDeficits := AllocateCash(deficits, cash)
	lotsToBuy := LotsToBuy(realDeficits, prices)
	return lotsToBuy
}
