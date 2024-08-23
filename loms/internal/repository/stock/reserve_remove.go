package stock

import (
	"context"
	"fmt"

	"ecom/loms/internal/model"

	"github.com/jackc/pgx/v5"
)

func (r *RepositoryStock) ReserveRemoveByOrderIDTx(ctx context.Context, tx pgx.Tx, orderID model.OrderID) error {
	const op = `repositoryStock.ReserveRemoveByOrderIDTx`

	const query = `
	DELETE FROM stock_reserve
	WHERE order_id = $1`

	commTag, err := tx.Exec(ctx, query, orderID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if commTag.RowsAffected() < 1 {
		return fmt.Errorf("%s: %w", op, fmt.Errorf("no reserve for orderID: %d", orderID))
	}

	return nil
}
