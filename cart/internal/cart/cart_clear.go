package cart

import (
	"context"
	"fmt"

	"ecom/cart/internal/model"
)

func (c *Cart) CartClear(ctx context.Context, user model.User) error {
	const op = `cart.CartClear`

	if err := c.repo.DeleteItemsByUser(ctx, user); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
