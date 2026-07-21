package fetcher

import (
	"fmt"
	"time"

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

func ToQuotation(d decimal.Decimal) api.Quotation {
	units := d.Truncate(0)
	nanoDecimal := d.Sub(units).Mul(decimal.NewFromInt(1_000_000_000))

	return api.Quotation{
		Units: units.String(),
		Nano:  int(nanoDecimal.IntPart()),
	}
}

func ToMoneyValue(m domain.Money) api.MoneyValue {
	q := ToQuotation(m.Amount)

	return api.MoneyValue{
		Currency: m.Currency,
		Units:    q.Units,
		Nano:     q.Nano,
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
		AccountID:             raw.AccountID,
		Total:                 toMoney(raw.TotalAmountPortfolio),
		TotalAmountCurrencies: toMoney(raw.TotalAmountCurrencies),
	}

	portfolio.Positions = make([]domain.Position, len(raw.Positions))
	for i, pos := range raw.Positions {
		portfolio.Positions[i] = domain.Position{
			Figi:                 pos.Figi,
			InstrumentType:       pos.InstrumentType,
			Quantity:             toDecimal(pos.Quantity.Units, pos.Quantity.Nano),
			AveragePositionPrice: toMoney(pos.AveragePositionPrice),
			CurrentPrice:         toMoney(pos.CurrentPrice),
			PositionUID:          pos.PositionUID,
			InstrumentUID:        pos.InstrumentUID,
			Ticker:               pos.Ticker,
			ClassCode:            pos.ClassCode,
		}
	}

	return portfolio
}

func convertInstruments(raw api.Instruments) []domain.Instrument {
	instruments := make([]domain.Instrument, len(raw.Instruments))
	for i, inst := range raw.Instruments {
		instruments[i] = domain.Instrument{
			ClassCode:             inst.ClassCode,
			Ticker:                inst.Ticker,
			InstrumentType:        inst.InstrumentType,
			PositionUID:           inst.PositionUID,
			Figi:                  inst.Figi,
			APITradeAvailableFlag: inst.APITradeAvailableFlag,
			Lot:                   inst.Lot,
			UID:                   inst.UID,
			Name:                  inst.Name,
		}
	}
	return instruments
}

func convertLastPrices(raw api.LastPrices) []domain.LastPrice {
	lastPrices := make([]domain.LastPrice, len(raw.LastPrices))
	for i, lp := range raw.LastPrices {
		lastPrices[i] = domain.LastPrice{
			ClassCode:     lp.ClassCode,
			Ticker:        lp.Ticker,
			Price:         toDecimal(lp.Price.Units, lp.Price.Nano),
			InstrumentUID: lp.InstrumentUID,
			Figi:          lp.Figi,
			Time:          lp.Time,
		}
	}
	return lastPrices
}

func ConvertPostOrderResponse(raw api.PostOrderResponse) domain.PostOrderResponse {
	return domain.PostOrderResponse{
		ClassCode:            raw.ClassCode,
		Ticker:               raw.Ticker,
		OrderID:              raw.OrderID,
		Figi:                 raw.Figi,
		InitialOrderPrice:    toMoney(raw.InitialOrderPrice),
		InitialCommission:    toMoney(raw.InitialCommission),
		Message:              raw.Message,
		LotsExecuted:         raw.LotsExecuted,
		TotalOrderAmount:     toMoney(raw.TotalOrderAmount),
		LotsRequested:        raw.LotsRequested,
		InstrumentUID:        raw.InstrumentUID,
		OrderRequestID:       raw.OrderRequestID,
		ExecutedOrderPrice:   toMoney(raw.ExecutedOrderPrice),
		ExecutedCommission:   toMoney(raw.ExecutedCommission),
		InitialSecurityPrice: toMoney(raw.InitialSecurityPrice),
		ResponseMetadata: struct {
			ServerTime time.Time "json:\"serverTime\""
			TrackingID string    "json:\"trackingId\""
		}{
			ServerTime: raw.ResponseMetadata.ServerTime,
			TrackingID: raw.ResponseMetadata.TrackingID,
		},
	}
}
