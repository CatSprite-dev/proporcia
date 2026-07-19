package fetcher

import (
	"fmt"

	"github.com/CatSprite-dev/proporcia/internal/api"
	"github.com/CatSprite-dev/proporcia/internal/domain"
	"github.com/shopspring/decimal"
)

func toDecimal(units string, nano int) decimal.Decimal {
	u, err := decimal.NewFromString(units)
	if err != nil {
		panic(fmt.Sprintf("invalid decimal units from broker: %q: %v", units, err))
	}
	n := decimal.New(int64(nano), -9)
	return u.Add(n)
}

func toMoney(mv api.MoneyValue) domain.Money {
	return domain.Money{
		Amount:   toDecimal(mv.Units, mv.Nano),
		Currency: mv.Currency,
	}
}

func convertAccounts(raw api.Accounts) []domain.Account {
	accounts := make([]domain.Account, len(raw.Accounts))
	for i, a := range raw.Accounts {
		accounts[i] = domain.Account{ID: a.ID}
	}
	return accounts
}

func convertPortfolio(raw api.Portfolio) domain.Portfolio {
	portfolio := domain.Portfolio{
		AccountID: raw.AccountID,
	}

	portfolio.Positions = make([]domain.Position, len(raw.Positions))
	for i, pos := range raw.Positions {
		portfolio.Positions[i] = domain.Position{
			Figi:                 pos.Figi,
			InstrumentType:       pos.InstrumentType,
			Quantity:             toDecimal(pos.Quantity.Units, pos.Quantity.Nano),
			AveragePositionPrice: toMoney(pos.AveragePositionPrice),
			CurrentPrice:         toMoney(pos.CurrentPrice),
			QuantityLots:         toDecimal(pos.QuantityLots.Units, pos.QuantityLots.Nano),
			PositionUID:          pos.PositionUID,
			InstrumentUID:        pos.InstrumentUID,
			Ticker:               pos.Ticker,
			ClassCode:            pos.ClassCode,
		}
	}

	return portfolio
}
