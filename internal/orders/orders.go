package orders

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

func (os *OrderService) Buy(ctx context.Context, token string, accountID string, lotsToBuy map[string]int) ([]domain.PostOrderResponse, error) {
	targets, err := os.storage.GetTargets(ctx)
	if err != nil {
		return []domain.PostOrderResponse{}, fmt.Errorf("get targets: %w", err)
	}

	responses := make([]domain.PostOrderResponse, 0, len(lotsToBuy))
	for _, target := range targets {
		if lots, ok := lotsToBuy[target.Ticker]; ok {
			if lots <= 0 {
				continue
			}
			postOrderResp, err := os.PostOrder(ctx, token, strconv.Itoa(lots), decimal.Zero, accountID, "", target.UID)
			if err != nil {
				return []domain.PostOrderResponse{}, fmt.Errorf("buy target: %w", err)
			}
			responses = append(responses, postOrderResp)
		}
	}
	return responses, nil
}
