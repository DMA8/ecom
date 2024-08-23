package order

import (
	"context"
	"ecom/loms/internal/model"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
)

// insert
func (r *Repository) OrderCreate(ctx context.Context, order model.OrderCreate, user model.User) (model.OrderID, error) {
	const op = `repositoryOrder.OrderCreate`
	var err error

	tx, err := r.psqlConnection.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
				slog.Error(fmt.Sprintf("%s: %s", op, rollbackErr))
			}
		}
	}()

	orderID, err := r.orderCreateTx(ctx, tx, order, user)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	err = r.stockRepo.ReserveSKUsTx(ctx, tx, orderID, order.Items)
	if err != nil {
		if err = r.orderSetStatusTx(ctx, tx, orderID, model.OrderStatusFailed); err != nil {
			slog.Error(fmt.Sprintf("%s: %s", op, err))
			return 0, fmt.Errorf("%s: %w", op, err)
		}
	}

	err = r.orderSetStatusTx(ctx, tx, orderID, model.OrderStatusAwaitingPayment)
	if err != nil {
		slog.Error(fmt.Sprintf("%s: %s", op, err))
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	if commitErr := tx.Commit(ctx); commitErr != nil {
		slog.Error(fmt.Sprintf("%s: %s", op, err))
	}
	return orderID, nil
}

func (r *Repository) orderCreateTx(ctx context.Context, tx pgx.Tx, order model.OrderCreate, user model.User) (model.OrderID, error) {
	const op = `repositoryOrder.orderCreateTx`
	const query = `
	INSERT INTO orders (user_id, status, order_content)
	VALUES ($1, $2, $3)
	RETURNING id`

	orderContent, err := json.Marshal(order.Items)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	var orderID model.OrderID
	if err := tx.QueryRow(ctx, query, user.ID, model.OrderStatusNew, orderContent).Scan(&orderID); err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return orderID, nil
}
