package stock

import (
	"context"
	"ecom/loms/internal/model"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func (r *RepositoryStock) IncreaseStockTx(ctx context.Context, tx pgx.Tx, items []model.Item) error {
	const op = `repositoryStock.IncreaseStockTx`

	const query = `
	UPDATE stocks
	SET quantity = quantity + $1
	WHERE sku = $2`

	for _, item := range items {
		if _, err := tx.Exec(ctx, query, item.Count, item.SKU); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	return nil
}
