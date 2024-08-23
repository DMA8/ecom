package grpcadapter

import (
	"context"
	"ecom/loms/internal/model"
	grpc_loms "ecom/loms/pkg/api/loms/v1"
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
)

func TestOrderCancel(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	orderManager := NewMockorderManager(ctrl)
	adapter := New(orderManager)

	testCases := []struct {
		name            string
		orderCancelReq  *grpc_loms.OrderCancelRequest
		mockExpectation func()
		expectedErr     error
	}{
		{
			name: "success",
			orderCancelReq: &grpc_loms.OrderCancelRequest{
				OrderId: 1,
			},
			mockExpectation: func() {
				orderManager.EXPECT().OrderCancel(gomock.Any(), model.OrderID(1)).Return(nil)
			},
		},
		{
			name: "error",
			orderCancelReq: &grpc_loms.OrderCancelRequest{
				OrderId: 1,
			},
			mockExpectation: func() {
				orderManager.EXPECT().OrderCancel(gomock.Any(), model.OrderID(1)).Return(model.ErrOrderIDNotFound)
			},
			expectedErr: model.ErrOrderIDNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockExpectation()
			ctx := context.Background()
			_, err := adapter.OrderCancel(ctx, tc.orderCancelReq)
			if !errors.Is(err, tc.expectedErr) {
				t.Errorf("name: %s. err expected %v, got %v", tc.name, tc.expectedErr, err)
			}
		})
	}
}
