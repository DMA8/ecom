//go:generate mockgen -source $GOFILE -destination mock_contract.go -package ${GOPACKAGE}
package loms

import (
	"context"

	"ecom/loms/internal/model"
)

type StockRepository interface {
	StockQuantity(ctx context.Context, sku model.SKU) (int64, error)
}

type OrderRepository interface {
	OrderCreate(ctx context.Context, order model.OrderCreate, user model.User) (model.OrderID, error)
	OrderInfoByOrderID(ctx context.Context, orderID model.OrderID) (model.OrderInfo, error)
	OrderPay(ctx context.Context, orderPay model.OrderID) error
	OrderCancel(ctx context.Context, orderPay model.OrderID) error
}

type Producer interface {
	SendOrderStatusChanged(ctx context.Context, order model.OrderChangedMessage)
}

type ProductServiceCli interface {
	IsSKUValid(ctx context.Context, sku model.SKU) error
}

type Manager struct {
	orderRepo      OrderRepository
	stockRepo      StockRepository
	producer       Producer
	statusMappings map[model.OrderStatus]string
}

func New(or OrderRepository, sr StockRepository, producer Producer) *Manager {
	return &Manager{
		orderRepo: or,
		stockRepo: sr,
		statusMappings: map[model.OrderStatus]string{
			model.OrderStatusNew:             "new",
			model.OrderStatusAwaitingPayment: "awaiting_payment",
			model.OrderStatusPayed:           "paid",
			model.OrderStatusFailed:          "failed",
			model.OrderStatusCancelled:       "cancelled",
		},
		producer: producer,
	}
}
