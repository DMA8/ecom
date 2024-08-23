package grpcadapter

import (
	"context"
	"ecom/loms/internal/model"
	grpc_loms "ecom/loms/pkg/api/loms/v1"
)

func (g *GrpcAdapter) OrderInfo(ctx context.Context, in *grpc_loms.OrderInfoRequest) (*grpc_loms.OrderInfoResponse, error) {
	order, err := g.lomsManager.OrderInfoByOrderID(ctx, model.OrderID(in.OrderId))
	if err != nil {
		return nil, err
	}
	return modelOrderInfoToDTO(order), nil
}

func modelOrderInfoToDTO(in model.OrderInfo) *grpc_loms.OrderInfoResponse {
	return &grpc_loms.OrderInfoResponse{
		UserId: int64(in.User.ID),
		Status: in.Status,
		Items:  itemsDTO(in.Items),
	}
}

func itemsDTO(in []model.Item) []*grpc_loms.Item {
	res := make([]*grpc_loms.Item, 0, len(in))
	for _, item := range in {
		res = append(res, &grpc_loms.Item{
			Sku:   uint32(item.SKU),
			Count: uint32(item.Count),
		})
	}
	return res
}
