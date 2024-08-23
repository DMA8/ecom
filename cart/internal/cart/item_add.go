package cart

import (
	"context"
	"fmt"

	"ecom/cart/internal/model"
)

func (c *Cart) AddItem(ctx context.Context, item model.Item, user model.User) error {
	const op = `cart.AddItem`

	itemInfo, err := c.productServiceAPIClient.GetProduct(ctx, item.SKU)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if err := c.IsSKUValid(itemInfo); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	quantity, err := c.lomsServiceAPIClient.StockInfo(ctx, item.SKU)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if quantity < uint64(item.Count) {
		return fmt.Errorf("%s: %w", op, model.ErrNotEnoughStock)
	}

	if err := c.repo.ItemAdd(ctx, item, user); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
