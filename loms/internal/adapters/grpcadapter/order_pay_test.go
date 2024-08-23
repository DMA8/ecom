package grpcadapter

import (
	"context"
	"ecom/loms/internal/model"
	grpc_loms "ecom/loms/pkg/api/loms/v1"
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
)

func TestOrderPay(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	orderManager := NewMockorderManager(ctrl)
	adapter := New(orderManager)

	testCases := []struct {
		name            string
		orderPayReq     *grpc_loms.OrderPayRequest
		mockExpectation func()
		expectedErr     error
		expectedOrderID int64
	}{
		{
			name: "success",
			orderPayReq: &grpc_loms.OrderPayRequest{
				OrderId: 1,
			},
			mockExpectation: func() {
				orderManager.EXPECT().OrderPay(gomock.Any(), model.OrderID(1)).Return(nil)
			},
			expectedOrderID: 1,
		},
		{
			name: "error",
			orderPayReq: &grpc_loms.OrderPayRequest{
				OrderId: 1,
			},
			mockExpectation: func() {
				orderManager.EXPECT().OrderPay(gomock.Any(), model.OrderID(1)).Return(model.ErrOrderIDNotFound)
			},
			expectedErr: model.ErrOrderIDNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockExpectation()
			ctx := context.Background()
			_, err := adapter.OrderPay(ctx, tc.orderPayReq)
			if !errors.Is(err, tc.expectedErr) {
				t.Errorf("name: %s. err expected %v, got %v", tc.name, tc.expectedErr, err)
			}
		})
	}
}
