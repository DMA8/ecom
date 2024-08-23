package loms

import (
	context "context"
	model "ecom/loms/internal/model"
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
)

func TestOrderCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repoStock := NewMockStockRepository(ctrl)
	repoOrder := NewMockOrderRepository(ctrl)
	producer := NewMockProducer(ctrl)
	service := New(repoOrder, repoStock, producer)

	testCases := []struct {
		name            string
		orderID         model.OrderID
		user            model.User
		order           model.OrderCreate
		orderRepoExpect func()
		repoOrderExpect func()
		expectedErr     error
		expectedOrderID model.OrderID
	}{
		{
			name:    "success",
			orderID: 1,
			user:    model.User{ID: 1},
			order:   model.OrderCreate{Items: []model.Item{{SKU: model.SKU(1), Count: 1}}},
			orderRepoExpect: func() {
				repoOrder.EXPECT().OrderCreate(gomock.Any(), gomock.Any(), gomock.Any()).Return(model.OrderID(1), nil).Times(1)
			},
			expectedOrderID: model.OrderID(1),
		},
		{
			name:    "error",
			orderID: 1,
			orderRepoExpect: func() {
				repoOrder.EXPECT().OrderCreate(gomock.Any(), gomock.Any(), gomock.Any()).Return(model.OrderID(0), model.ErrNotEnoughStock).Times(1)
			},
			expectedErr:     model.ErrNotEnoughStock,
			expectedOrderID: model.OrderID(0),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.orderRepoExpect()
			ctx := context.Background()
			id, err := service.OrderCreate(ctx, tc.order, tc.user)
			if !errors.Is(err, tc.expectedErr) {
				t.Errorf("name: %s. err expected %v, got %v", tc.name, tc.expectedErr, err)
			}
			if id != tc.expectedOrderID {
				t.Errorf("name: %s. id expected %v, got %v", tc.name, tc.expectedOrderID, id)
			}
		})
	}
}
