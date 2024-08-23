package order

import (
	"context"
	"ecom/loms/internal/model"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5"
)

// select by order id
func (r *Repository) OrderInfoByOrderID(ctx context.Context, orderID model.OrderID) (model.OrderInfo, error) {
	const op = `repositoryOrder.OrderInfoByOrderID`

	const query = `
	SELECT id, user_id, status, order_content, created_at, updated_at
	FROM orders
	WHERE id = $1
	ORDER BY id`

	var orderInfo model.OrderInfo
	var orderContent []byte

	row := r.psqlConnection.QueryRow(ctx, query, orderID)

	err := row.Scan(
		&orderInfo.ID,
		&orderInfo.User.ID,
		&orderInfo.StatusCode,
		&orderContent,
		&orderInfo.CreatedAt,
		&orderInfo.UpdatedAt,
	)
	if err != nil {
		return model.OrderInfo{}, fmt.Errorf("%s: %w", op, err)
	}

	if err := json.Unmarshal(orderContent, &orderInfo.Items); err != nil {
		return model.OrderInfo{}, fmt.Errorf("%s: %w", op, err)
	}
	return orderInfo, nil
}

func (r *Repository) OrderInfoByOrderIDTx(ctx context.Context, tx pgx.Tx, orderID model.OrderID) (model.OrderInfo, error) {
	const op = `repositoryOrder.OrderInfoByOrderID`

	const query = `
	SELECT id, user_id, status, order_content, created_at, updated_at
	FROM orders
	WHERE id = $1
	FOR UPDATE`

	var orderInfo model.OrderInfo
	var orderContent []byte

	row := tx.QueryRow(ctx, query, orderID)

	err := row.Scan(
		&orderInfo.ID,
		&orderInfo.User.ID,
		&orderInfo.StatusCode,
		&orderContent,
		&orderInfo.CreatedAt,
		&orderInfo.UpdatedAt,
	)
	if err != nil {
		return model.OrderInfo{}, fmt.Errorf("%s: %w", op, err)
	}

	if err := json.Unmarshal(orderContent, &orderInfo.Items); err != nil {
		return model.OrderInfo{}, fmt.Errorf("%s: %w", op, err)
	}
	return orderInfo, nil
}
