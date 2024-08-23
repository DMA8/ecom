package lomsgrpc

import (
	"context"
	"ecom/cart/internal/model"
	grpc_loms "ecom/loms/pkg/api/loms/v1"

	"google.golang.org/grpc"
)

type loms interface {
	StockInfo(ctx context.Context, req *grpc_loms.StockInfoRequest, opt ...grpc.CallOption) (*grpc_loms.StockInfoResponse, error)
	OrderCreate(ctx context.Context, req *grpc_loms.OrderCreateRequest, opt ...grpc.CallOption) (*grpc_loms.OrderCreateResponse, error)
}

type LOMSCliGrpc struct {
	l loms
}

func New(l loms) *LOMSCliGrpc {
	return &LOMSCliGrpc{
		l: l,
	}
}

func (l *LOMSCliGrpc) StockInfo(ctx context.Context, sku model.SKU) (uint64, error) {
	resp, err := l.l.StockInfo(ctx, &grpc_loms.StockInfoRequest{
		Sku: uint32(sku),
	})
	if err != nil {
		return 0, err
	}

	return uint64(resp.Count), nil
}

func (l *LOMSCliGrpc) OrderCreate(ctx context.Context, items []model.Item, user model.User) (model.OrderID, error) {
	itemsCreate := make([]*grpc_loms.Item, 0, len(items))
	for _, item := range items {
		itemsCreate = append(itemsCreate, &grpc_loms.Item{
			Sku:   uint32(item.SKU),
			Count: uint32(item.Count),
		})
	}

	resp, err := l.l.OrderCreate(ctx, &grpc_loms.OrderCreateRequest{
		Items:  itemsCreate,
		UserId: user.ID,
	})
	if err != nil {
		return 0, err
	}

	return model.OrderID(resp.OrderId), nil
}
