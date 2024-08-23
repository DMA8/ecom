package loms

import (
	"context"
	"fmt"

	"ecom/loms/internal/model"
)

func (m *Manager) OrderCancel(ctx context.Context, orderCancel model.OrderID) error {
	const op = `manager.OrderCancel`

	err := m.orderRepo.OrderCancel(ctx, orderCancel)
	if err != nil {
		return fmt.Errorf("%s %w", op, err)
	}

	m.producer.SendOrderStatusChanged(ctx, model.OrderChangedMessage{OrderID: int(orderCancel), NewStatus: m.statusMappings[model.OrderStatusCancelled]})
	return nil
}
