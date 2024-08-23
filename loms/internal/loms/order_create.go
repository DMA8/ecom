package loms

import (
	"context"
	"fmt"

	"ecom/loms/internal/model"
)

func (m *Manager) OrderCreate(ctx context.Context, order model.OrderCreate, user model.User) (model.OrderID, error) {
	const op = `stock.OrderCreate`

	order.Status = model.OrderStatusNew
	orderID, err := m.orderRepo.OrderCreate(ctx, order, user)
	if err != nil {
		return 0, fmt.Errorf("%s %w", op, err)
	}

	m.producer.SendOrderStatusChanged(ctx, model.OrderChangedMessage{OrderID: int(orderID), NewStatus: m.statusMappings[model.OrderStatusAwaitingPayment]})
	return orderID, nil
}
