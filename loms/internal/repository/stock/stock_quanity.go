package stock

import (
	"context"
	"errors"
	"fmt"

	"ecom/loms/internal/model"

	"github.com/jackc/pgx/v5"
)

// select quantity
func (r *RepositoryStock) StockQuantity(ctx context.Context, sku model.SKU) (int64, error) {
	const op = `repositoryStock.StockQuantity`

	const query = `
	SELECT quantity
	FROM stocks
	WHERE sku = $1`

	var quantity int64

	row := r.psqlConnection.QueryRow(ctx, query, sku)
	if err := row.Scan(&quantity); errors.Is(err, pgx.ErrNoRows) {
		return 0, fmt.Errorf("%s: %w", op, model.ErrSKUNotFound)
	} else if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return quantity, nil
}
