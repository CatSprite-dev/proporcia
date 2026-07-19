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
				Ticker:       "GAZP",
				Quantity:     decimal.NewFromInt(100),
				CurrentPrice: domain.Money{Amount: decimal.NewFromInt(100), Currency: "RUB"},
			},
			{
				Ticker:       "LKOH",
				Quantity:     decimal.NewFromInt(10),
				CurrentPrice: domain.Money{Amount: decimal.NewFromInt(500), Currency: "RUB"},
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
				Ticker:       "SBER",
				Quantity:     decimal.NewFromInt(100),
				CurrentPrice: domain.Money{Amount: decimal.NewFromInt(100), Currency: "RUB"},
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

func checkDeficit(t *testing.T, result map[string]domain.Money, ticker string, expected decimal.Decimal) {
	t.Helper()
	got, ok := result[ticker]
	if !ok {
		t.Fatalf("expected deficit for %s, got none", ticker)
	}
	if !got.Amount.Equal(expected) {
		t.Errorf("%s: expected deficit %s, got %s", ticker, expected, got.Amount)
	}
}
