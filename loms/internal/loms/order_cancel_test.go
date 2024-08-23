package loms

import (
	context "context"
	model "ecom/loms/internal/model"
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
)

func TestOrderCancel(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repoStock := NewMockStockRepository(ctrl)
	repoOrder := NewMockOrderRepository(ctrl)
	producer := NewMockProducer(ctrl)
	service := New(repoOrder, repoStock, producer)

	testCases := []struct {
		name            string
		orderID         model.OrderID
		orderRepoExpect func()
		repoOrderExpect func()
		expectedErr     error
	}{
		{
			name:    "success",
			orderID: 1,
			orderRepoExpect: func() {
				repoOrder.EXPECT().OrderCancel(gomock.Any(), gomock.Eq(model.OrderID(1))).Return(nil)
			},
		},
		{
			name:    "error",
			orderID: 1,
			orderRepoExpect: func() {
				repoOrder.EXPECT().OrderCancel(gomock.Any(), gomock.Eq(model.OrderID(1))).Return(model.ErrOrderIDNotFound)
			},
			expectedErr: model.ErrOrderIDNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.orderRepoExpect()
			ctx := context.Background()
			err := service.OrderCancel(ctx, tc.orderID)
			if !errors.Is(err, tc.expectedErr) {
				t.Errorf("expected %v, got %v", tc.expectedErr, err)
			}
		})
	}
}
