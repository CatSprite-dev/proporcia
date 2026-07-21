package orders

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/CatSprite-dev/proporcia/internal/api"
	"github.com/CatSprite-dev/proporcia/internal/domain"
	"github.com/CatSprite-dev/proporcia/internal/fetcher"
	"github.com/shopspring/decimal"
)

type OrderService struct {
	apiClient *api.Client
	logger    *slog.Logger
}

func NewOrderService(apiClient *api.Client, logger *slog.Logger) *OrderService {
	return &OrderService{apiClient: apiClient, logger: logger}
}

func (os *OrderService) PostOrder(
	ctx context.Context,
	token string,
	quantity string,
	price decimal.Decimal,
	accountID string,
	orderID string,
	instrumentID string,
) (domain.PostOrderResponse, error) {
	raw, err := os.apiClient.PostOrder(
		ctx,
		token,
		quantity,
		fetcher.ToQuotation(price),
		api.OrderDirectionBuy,
		accountID,
		api.OrderTypeBestPrice,
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
