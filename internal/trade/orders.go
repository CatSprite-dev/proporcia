package trade

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/CatSprite-dev/proporcia/internal/api"
	"github.com/CatSprite-dev/proporcia/internal/domain"
	"github.com/CatSprite-dev/proporcia/internal/fetcher"
	"github.com/CatSprite-dev/proporcia/internal/storage"
	"github.com/shopspring/decimal"
)

type OrderService struct {
	apiClient *api.Client
	storage   *storage.Storage
	logger    *slog.Logger
}

func NewOrderService(apiClient *api.Client, storage *storage.Storage, logger *slog.Logger) *OrderService {
	return &OrderService{apiClient: apiClient, storage: storage, logger: logger}
}

func (os *OrderService) PostOrder(
	ctx context.Context,
	token string,
	quantityLots string,
	price decimal.Decimal,
	accountID string,
	orderID string,
	instrumentID string,
) (domain.PostOrderResponse, error) {
	raw, err := os.apiClient.PostOrder(
		ctx,
		token,
		quantityLots,
		fetcher.ToQuotation(price),
		api.OrderDirectionBuy,
		accountID,
		api.OrderTypeMarket,
		orderID,
		instrumentID,
		api.TimeInForceUnspecified,
		api.PriceTypeUnspecified,
		false,
	)

	if err != nil {
		return domain.PostOrderResponse{}, fmt.Errorf("post order: %w", err)
	}

	postOrderResp := fetcher.ConvertPostOrderResponse(raw)
	return postOrderResp, nil
}

func (os *OrderService) Buy(ctx context.Context, token string, accountID string, lotsToBuy map[string]int) ([]domain.PostOrderResponse, error) {
	targets, err := os.storage.GetTargets(ctx)
	if err != nil {
		return []domain.PostOrderResponse{}, fmt.Errorf("get targets: %w", err)
	}

	orders := make([]domain.PostOrderResponse, 0, len(lotsToBuy))
	for _, target := range targets {
		if lots, ok := lotsToBuy[target.InstrumentUID]; ok {
			if lots <= 0 {
				continue
			}
			postOrderResp, err := os.PostOrder(ctx, token, strconv.Itoa(lots), decimal.Zero, accountID, "", target.InstrumentUID)
			if err != nil {
				os.logger.Error("buy target failed", "ticker", target.Ticker, "error", err)
				continue
			}
			orders = append(orders, postOrderResp)
		}
	}
	return orders, nil
}

func (os *OrderService) SaveOrders(ctx context.Context, orders []domain.PostOrderResponse) error {
	err := os.storage.SaveOrders(ctx, orders)
	if err != nil {
		return fmt.Errorf("save orders: %w", err)
	}
	return nil
}
