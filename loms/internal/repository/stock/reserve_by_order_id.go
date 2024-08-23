package stock

import (
	"context"
	"ecom/loms/internal/model"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func (r *RepositoryStock) SkuReservedByOrderTx(ctx context.Context, tx pgx.Tx, orderID model.OrderID) ([]model.Item, error) {
	const op = `repositoryStock.ReserveRemoveByOrderIDTx`

	const query = `
	SELECT sku, quantity
	FROM stock_reserve
	WHERE order_id = $1
	FOR UPDATE`

	rows, err := tx.Query(ctx, query, orderID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	reserves := make([]model.Item, 0, 5)

	for rows.Next() {
		var sku uint32
		var quantity uint32

		if err := rows.Scan(&sku, &quantity); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		reserves = append(reserves, model.Item{
			SKU:   model.SKU(sku),
			Count: uint16(quantity),
		})
	}
	return reserves, nil
}
