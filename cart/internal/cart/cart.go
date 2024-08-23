package cart

import (
	"context"

	"ecom/cart/internal/model"
)

type Repository interface {
	ItemAdd(ctx context.Context, item model.Item, user model.User) error
	ItemDelete(ctx context.Context, sku model.SKU, user model.User) error
	DeleteItemsByUser(ctx context.Context, user model.User) error
	ItemsByUser(ctx context.Context, user model.User) ([]model.Item, error)
}

type (
	ProductService interface {
		GetProduct(ctx context.Context, sku model.SKU) (model.ProductInfo, error)
	}

	LOMSService interface {
		StockInfo(ctx context.Context, sku model.SKU) (uint64, error)
		OrderCreate(ctx context.Context, items []model.Item, user model.User) (model.OrderID, error)
	}
)

type Cart struct {
	productServiceAPIClient ProductService
	lomsServiceAPIClient    LOMSService
	repo                    Repository
}

func New(r Repository, productService ProductService, lomsService LOMSService) *Cart {
	return &Cart{
		repo:                    r,
		productServiceAPIClient: productService,
		lomsServiceAPIClient:    lomsService,
	}
}
