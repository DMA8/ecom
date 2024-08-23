package order

import (
	"context"
	"ecom/loms/internal/model"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func (r *Repository) orderSetStatusTx(ctx context.Context, tx pgx.Tx, orderID model.OrderID, status model.OrderStatus) error {
	const op = `repositoryOrder.OrderSetStatusTx`

	const query = `
	UPDATE orders
	SET status = $1
	WHERE id = $2`

	if _, err := tx.Exec(ctx, query, status, orderID); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (r *Repository) orderSetStatusPaidTx(ctx context.Context, tx pgx.Tx, orderID model.OrderID, newStatus model.OrderStatus) error {
	const op = `repositoryOrder.OrderSetStatusPaidTx`

	const query = `
	UPDATE orders
	SET status = $1
	WHERE id = $2 AND status = $3`

	ct, err := tx.Exec(ctx, query, newStatus, orderID, model.OrderStatusAwaitingPayment)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if ct.RowsAffected() == 0 {
		return fmt.Errorf("%s: couldn't set status 'payed' for orderID: %d ", op, orderID)
	}

	return nil
}
