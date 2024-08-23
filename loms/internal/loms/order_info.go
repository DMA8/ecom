package loms

import (
	"context"
	"fmt"

	"ecom/loms/internal/model"
)

func (m *Manager) OrderInfoByOrderID(ctx context.Context, orderID model.OrderID) (model.OrderInfo, error) {
	const op = `OrderInfo`

	orderInfo, err := m.orderRepo.OrderInfoByOrderID(ctx, orderID)
	if err != nil {
		return model.OrderInfo{}, fmt.Errorf("%s %w", op, err)
	}

	orderInfo.Status, err = m.ConvertStatus(orderInfo.StatusCode)
	if err != nil {
		return model.OrderInfo{}, fmt.Errorf("%s %w", op, err)
	}

	if len(orderInfo.Items) == 0 {
		return model.OrderInfo{}, fmt.Errorf("%s %w", op, model.ErrOrderIDNotFound)
	}

	return orderInfo, nil
}
