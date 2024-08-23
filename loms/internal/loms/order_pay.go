package loms

import (
	"context"
	"fmt"

	"ecom/loms/internal/model"
)

func (m *Manager) OrderPay(ctx context.Context, orderPay model.OrderID) error {
	const op = `manager.OrderPay`

	orderInfo, err := m.orderRepo.OrderInfoByOrderID(ctx, orderPay)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if orderInfo.StatusCode != model.OrderStatusAwaitingPayment {
		return fmt.Errorf("%s: %w", op, model.ErrOrderIsNotWaitingForPayment)
	}

	if err := m.orderRepo.OrderPay(ctx, orderPay); err != nil {
		return fmt.Errorf("%s orderPay err: %w", op, err)
	}

	m.producer.SendOrderStatusChanged(ctx, model.OrderChangedMessage{OrderID: int(orderPay), NewStatus: m.statusMappings[model.OrderStatusPayed]})
	return nil
}
