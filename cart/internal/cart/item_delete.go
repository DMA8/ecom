package cart

import (
	"context"
	"fmt"

	"ecom/cart/internal/model"
)

func (c *Cart) ItemDelete(ctx context.Context, sku model.SKU, user model.User) error {
	const op = `cart.ItemDelete`

	if err := c.repo.ItemDelete(ctx, sku, user); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
