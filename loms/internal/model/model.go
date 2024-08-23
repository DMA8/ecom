package model

import (
	"fmt"
	"time"
)

type OrderStatus int

var TraceIDKey = struct{}{}

const (
	OrderStatusCancelled       OrderStatus = 0
	OrderStatusNew             OrderStatus = 1
	OrderStatusAwaitingPayment OrderStatus = 2
	OrderStatusPayed           OrderStatus = 3
	OrderStatusFailed          OrderStatus = 4
)

var (
	ErrOrderIDNotFound             = fmt.Errorf("orderID not found")
	ErrNotEnoughStock              = fmt.Errorf("not enough stock")
	ErrSKUNotValid                 = fmt.Errorf("sku not valid")
	ErrNoTokenForProductService    = fmt.Errorf("no token for product service")
	ErrOrderIsNotWaitingForPayment = fmt.Errorf("order is not waiting for payment")
	ErrOrderAlreadyCanceled        = fmt.Errorf("order is already canceled")
	ErrSKUNotFound                 = fmt.Errorf("sku not found")
)

type ProductServiceTokenKeyT struct{}

var ProductServiceTokenKey = ProductServiceTokenKeyT{}

type (
	SKU     uint32
	OrderID int64
)

type Item struct {
	SKU   SKU    `json:"sku"`
	Count uint16 `json:"count"`
}

type User struct {
	ID int64
}

type OrderCreate struct {
	Items  []Item
	Status OrderStatus
}

type OrderInfo struct {
	ID         int64
	CreatedAt  time.Time
	UpdatedAt  time.Time
	StatusCode OrderStatus
	Status     string
	Items      []Item
	User       User
}

type OrderChangedMessage struct {
	OrderID   int
	Ts        time.Time
	NewStatus string
}
