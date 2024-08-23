package cart

import (
	"context"
	"fmt"

	"ecom/cart/internal/model"
)

func (c *Cart) CartCheckout(ctx context.Context, user model.User) (model.OrderID, error) {
	const op = `cart.CartCheckout`

	allItems, err := c.repo.ItemsByUser(ctx, user)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	if len(allItems) == 0 {
		return 0, fmt.Errorf("%s: %w", op, model.ErrCheckoutEmptyCart)
	}

	orderID, err := c.lomsServiceAPIClient.OrderCreate(ctx, allItems, user)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	if err := c.repo.DeleteItemsByUser(ctx, user); err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return orderID, nil
}
