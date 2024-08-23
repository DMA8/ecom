package grpcadapter

import (
	"context"
	"ecom/loms/internal/model"
	grpc_loms "ecom/loms/pkg/api/loms/v1"
)

func (g *GrpcAdapter) OrderCreate(ctx context.Context, in *grpc_loms.OrderCreateRequest) (*grpc_loms.OrderCreateResponse, error) {
	orderToCreate := model.OrderCreate{
		Items: convertToItems(in.Items),
	}

	orderID, err := g.lomsManager.OrderCreate(ctx, orderToCreate, model.User{ID: in.UserId})
	if err != nil {
		return nil, err
	}
	return &grpc_loms.OrderCreateResponse{OrderId: int64(orderID)}, nil
}

func convertToItems(in []*grpc_loms.Item) []model.Item {
	items := make([]model.Item, 0, len(in))
	for _, item := range in {
		items = append(items, model.Item{
			SKU:   model.SKU(item.Sku),
			Count: uint16(item.Count),
		})
	}
	return items
}
