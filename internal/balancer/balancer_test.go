package balancer

import (
	"testing"

	"github.com/CatSprite-dev/proporcia/internal/domain"
	"github.com/shopspring/decimal"
)

func TestDeficits(t *testing.T) {
	portfolio := domain.Portfolio{
		Total: domain.Money{
			Amount:   decimal.NewFromInt(35000),
			Currency: "RUB",
		},
		Positions: []domain.Position{
			{
				InstrumentUID: "GAZP",
				Quantity:      decimal.NewFromInt(100),
				CurrentPrice:  domain.Money{Amount: decimal.NewFromInt(100), Currency: "RUB"},
			},
			{
				InstrumentUID: "LKOH",
				Quantity:      decimal.NewFromInt(10),
				CurrentPrice:  domain.Money{Amount: decimal.NewFromInt(500), Currency: "RUB"},
			},
		},
	}

	weights := map[string]decimal.Decimal{
		"SBER": decimal.NewFromFloat(0.4), // 0 в портфеле -> дефицит 14000
		"GAZP": decimal.NewFromFloat(0.4), // 10000 в портфеле -> дефицит 4000
		"LKOH": decimal.NewFromFloat(0.2), // 5000 в портфеле -> дефицит 2000
	}

	result := Deficits(portfolio, weights)

	if len(result) != 3 {
		t.Fatalf("expected 3 deficits, got %d", len(result))
	}

	checkDeficit(t, result, "SBER", decimal.NewFromInt(14000))
	checkDeficit(t, result, "GAZP", decimal.NewFromInt(4000))
	checkDeficit(t, result, "LKOH", decimal.NewFromInt(2000))
}

func TestDeficits_OverweightPositionSkipped(t *testing.T) {
	portfolio := domain.Portfolio{
		Total: domain.Money{
			Amount:   decimal.NewFromInt(10000),
			Currency: "RUB",
		},
		Positions: []domain.Position{
			{
				InstrumentUID: "SBER",
				Quantity:      decimal.NewFromInt(100),
				CurrentPrice:  domain.Money{Amount: decimal.NewFromInt(100), Currency: "RUB"},
			},
		},
	}

	weights := map[string]decimal.Decimal{
		"SBER": decimal.NewFromFloat(0.5), // цель 5000, сейчас 10000 -> перевес, дефицита нет
	}

	result := Deficits(portfolio, weights)

	if len(result) != 0 {
		t.Fatalf("expected no deficits for overweight position, got %v", result)
	}
}

func TestDeficits_EmptyPortfolio(t *testing.T) {
	portfolio := domain.Portfolio{
		Total: domain.Money{
			Amount:   decimal.Zero,
			Currency: "RUB",
		},
	}

	weights := map[string]decimal.Decimal{
		"SBER": decimal.NewFromFloat(1.0),
	}

	result := Deficits(portfolio, weights)

	if len(result) != 0 {
		t.Fatalf("expected no deficits when total is zero, got %v", result)
	}
}

func checkDeficit(t *testing.T, result map[string]domain.Money, InstrumentUID string, expected decimal.Decimal) {
	t.Helper()
	got, ok := result[InstrumentUID]
	if !ok {
		t.Fatalf("expected deficit for %s, got none", InstrumentUID)
	}
	if !got.Amount.Equal(expected) {
		t.Errorf("%s: expected deficit %s, got %s", InstrumentUID, expected, got.Amount)
	}
}

func TestLotsToBuy(t *testing.T) {
	deficits := map[string]domain.Money{
		"SBER": {Amount: decimal.NewFromInt(14000), Currency: "RUB"},
		"LKOH": {Amount: decimal.NewFromInt(2000), Currency: "RUB"},
	}

	prices := map[string]domain.PriceInfo{
		"SBER": {Price: decimal.NewFromInt(280), LotSize: 10}, // лот 2800
		"LKOH": {Price: decimal.NewFromInt(500), LotSize: 1},  // лот 500
	}

	result := LotsToBuy(deficits, prices)

	if got := result["SBER"]; got != 5 {
		t.Errorf("SBER: expected 5 lots, got %d", got)
	}
	if got := result["LKOH"]; got != 4 {
		t.Errorf("LKOH: expected 4 lots, got %d", got)
	}
}

func TestLotsToBuy_RoundsDown(t *testing.T) {
	deficits := map[string]domain.Money{
		"SBER": {Amount: decimal.NewFromInt(15000), Currency: "RUB"}, // 15000/2800 = 5.35
	}

	prices := map[string]domain.PriceInfo{
		"SBER": {Price: decimal.NewFromInt(280), LotSize: 10},
	}

	result := LotsToBuy(deficits, prices)

	if got := result["SBER"]; got != 5 {
		t.Errorf("expected floor rounding to 5 lots, got %d", got)
	}
}

func TestLotsToBuy_MissingPriceSkipped(t *testing.T) {
	deficits := map[string]domain.Money{
		"UNKNOWN": {Amount: decimal.NewFromInt(1000), Currency: "RUB"},
	}

	prices := map[string]domain.PriceInfo{}

	result := LotsToBuy(deficits, prices)

	if _, ok := result["UNKNOWN"]; ok {
		t.Errorf("expected no entry for InstrumentUID without price info")
	}
}

func TestAllocateCash(t *testing.T) {
	deficits := map[string]domain.Money{
		"LKOH": {Amount: decimal.NewFromInt(4000), Currency: "RUB"},
		"PLZL": {Amount: decimal.NewFromInt(1000), Currency: "RUB"},
	}
	cash := domain.Money{Amount: decimal.NewFromInt(1000), Currency: "RUB"}

	result := AllocateCash(deficits, cash)

	// LKOH доля 4000/5000 = 0.8 -> 800
	// PLZL доля 1000/5000 = 0.2 -> 200
	checkDeficit(t, result, "LKOH", decimal.NewFromInt(800))
	checkDeficit(t, result, "PLZL", decimal.NewFromInt(200))
}

func TestAllocateCash_EqualDeficits(t *testing.T) {
	deficits := map[string]domain.Money{
		"A": {Amount: decimal.NewFromInt(4000), Currency: "RUB"},
		"B": {Amount: decimal.NewFromInt(4000), Currency: "RUB"},
	}
	cash := domain.Money{Amount: decimal.NewFromInt(20000), Currency: "RUB"}

	result := AllocateCash(deficits, cash)

	checkDeficit(t, result, "A", decimal.NewFromInt(10000))
	checkDeficit(t, result, "B", decimal.NewFromInt(10000))
}

func TestAllocateCash_ZeroTotalDeficit(t *testing.T) {
	deficits := map[string]domain.Money{}
	cash := domain.Money{Amount: decimal.NewFromInt(5000), Currency: "RUB"}

	result := AllocateCash(deficits, cash)

	if len(result) != 0 {
		t.Fatalf("expected empty result when there are no deficits, got %v", result)
	}
}

func TestAllocateCash_ExceedsCash(t *testing.T) {
	deficits := map[string]domain.Money{
		"LKOH": {Amount: decimal.NewFromInt(100000), Currency: "RUB"},
	}
	cash := domain.Money{Amount: decimal.NewFromInt(500), Currency: "RUB"}

	result := AllocateCash(deficits, cash)

	// единственный тикер получает весь кэш целиком
	checkDeficit(t, result, "LKOH", decimal.NewFromInt(500))
}
