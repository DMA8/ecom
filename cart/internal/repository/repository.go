package repository

import (
	"context"
	"fmt"

	"ecom/cart/internal/model"

	"github.com/jackc/pgx/v5"
)

type Repository struct {
	conn *pgx.Conn
}

func New(conn *pgx.Conn) *Repository {
	return &Repository{
		conn: conn,
	}
}

func (r *Repository) ItemAdd(ctx context.Context, item model.Item, user model.User) error {
	const op = `repository.ItemAdd`

	const query = `
		INSERT INTO items (user_id, sku, count)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, sku)
		DO UPDATE SET count = items.count + $3`

	_, err := r.conn.Exec(ctx, query, user.ID, item.SKU, item.Count)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (r *Repository) ItemDelete(ctx context.Context, sku model.SKU, user model.User) error {
	const op = `repository.ItemDelete`

	const query = `
		DELETE FROM items
		WHERE user_id = $1 AND sku = $2`

	_, err := r.conn.Exec(ctx, query, user.ID, sku)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (r *Repository) DeleteItemsByUser(ctx context.Context, user model.User) error {
	const op = `repository.DeleteItemsByUser`

	const query = `
		DELETE FROM items
		WHERE user_id = $1`

	_, err := r.conn.Exec(ctx, query, user.ID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (r *Repository) ItemsByUser(ctx context.Context, user model.User) ([]model.Item, error) {
	const op = `repository.ItemsByUser`

	const query = `
		SELECT sku, count
		FROM items
		WHERE user_id = $1`

	rows, err := r.conn.Query(ctx, query, user.ID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	result := make([]model.Item, 0, 5)
	for rows.Next() {
		var sku uint32
		var count uint32

		if err := rows.Scan(&sku, &count); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		result = append(result, model.Item{
			SKU:   model.SKU(sku),
			Count: uint16(count),
		})
	}

	return result, rows.Err()
}
