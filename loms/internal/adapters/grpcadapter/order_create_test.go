package grpcadapter

import (
	"context"
	"ecom/loms/internal/model"
	grpc_loms "ecom/loms/pkg/api/loms/v1"
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
)

func TestOrderCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	orderManager := NewMockorderManager(ctrl)
	adapter := New(orderManager)

	testCases := []struct {
		name            string
		orderCreateReq  *grpc_loms.OrderCreateRequest
		mockExpectation func()
		expectedErr     error
		expectedOrderID int64
	}{
		{
			name: "success",
			orderCreateReq: &grpc_loms.OrderCreateRequest{
				UserId: 1,
				Items: []*grpc_loms.Item{
					{
						Sku:   1,
						Count: 1,
					},
				},
			},
			mockExpectation: func() {
				orderManager.EXPECT().OrderCreate(gomock.Any(), gomock.Any(), model.User{ID: 1}).Return(model.OrderID(1), nil)
			},
			expectedOrderID: 1,
		},
		{
			name: "error",
			orderCreateReq: &grpc_loms.OrderCreateRequest{
				UserId: 1,
				Items: []*grpc_loms.Item{
					{
						Sku:   1,
						Count: 1,
					},
				},
			},
			mockExpectation: func() {
				orderManager.EXPECT().OrderCreate(gomock.Any(), gomock.Any(), model.User{ID: 1}).Return(model.OrderID(0), model.ErrOrderIDNotFound)
			},
			expectedErr: model.ErrOrderIDNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockExpectation()
			ctx := context.Background()
			orderID, err := adapter.OrderCreate(ctx, tc.orderCreateReq)
			if !errors.Is(err, tc.expectedErr) {
				t.Errorf("name: %s. err expected %v, got %v", tc.name, tc.expectedErr, err)
			}
			if tc.expectedErr == nil {
				if orderID.OrderId != tc.expectedOrderID {
					t.Errorf("name: %s. orderID expected %v, got %v", tc.name, tc.expectedOrderID, orderID)
				}
			}
		})
	}
}
