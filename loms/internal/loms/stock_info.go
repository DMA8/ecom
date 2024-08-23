package loms

import (
	"context"
	"fmt"

	"ecom/loms/internal/model"
)

func (m *Manager) StockInfo(ctx context.Context, sku model.SKU) (int64, error) {
	const op = `stock.Stock.IsEnoughStock`

	q, err := m.stockRepo.StockQuantity(ctx, sku)
	if err != nil {
		return 0, fmt.Errorf("%s %w", op, err)
	}

	return q, nil
}
