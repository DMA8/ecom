package model

import (
	"fmt"
)

var (
	ErrNotEnoughStock                   = fmt.Errorf("not enough stock")
	ErrInvalidSKU                       = fmt.Errorf("invalid sku")
	ErrUserNotFound                     = fmt.Errorf("user not found")
	ErrCartIsEmpty                      = fmt.Errorf("cart is empty")
	ErrNoTokenForProductService         = fmt.Errorf("no token for product service")
	ErrCheckoutEmptyCart                = fmt.Errorf("checkout empty cart")
	ErrProductServiceSuspiciousResponse = fmt.Errorf("suspicious response from product service")
)

type ProductServiceTokenKeyT struct{}

var ProductServiceTokenKey = ProductServiceTokenKeyT{}

type (
	SKU     uint32
	OrderID int64
)

type Item struct {
	SKU   SKU
	Count uint16
}

type ProductInfo struct {
	Name  string `json:"name"`
	Price uint32 `json:"price"`
}

type User struct {
	ID int64 `json:"user"`
}

type UserCart struct {
	TotalPrice uint32
	Items      []ItemInCart
}

type ItemInCart struct {
	Sku   SKU
	Count uint16
	Name  string
	Price uint32
}
