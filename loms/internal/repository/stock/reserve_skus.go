package stock

import (
	"context"
	"fmt"
	"log/slog"

	"ecom/loms/internal/model"

	"github.com/jackc/pgx/v5"
	"github.com/lib/pq"
)

func (r *RepositoryStock) ReserveSKUsTx(ctx context.Context, tx pgx.Tx, orderID model.OrderID, items []model.Item) error {
	const op = `repositoryStock.ReserveSKUsTx`

	const query = `
	SELECT sku, quantity
	FROM stocks
	WHERE sku = ANY($1)
	FOR UPDATE`

	skusNeedToReserve := make(map[uint32]uint32, len(items))
	SKUs := make([]int32, 0, len(items))
	for _, item := range items {
		SKUs = append(SKUs, int32(item.SKU))
		skusNeedToReserve[uint32(item.SKU)] = uint32(item.Count)
	}

	rows, err2 := tx.Query(ctx, query, pq.Array(SKUs))
	if err2 != nil {
		slog.Error(fmt.Sprintf("%s: %s", op, err2))
		return fmt.Errorf("%s: %w", op, err2)
	}
	defer rows.Close()

	var skuReaded int
	for rows.Next() {
		var sku uint32
		var totalStockQuantitySKU uint32

		if err := rows.Scan(&sku, &totalStockQuantitySKU); err != nil {
			return fmt.Errorf("%s!: %w", op, err)
		}
		if totalStockQuantitySKU < skusNeedToReserve[sku] {
			return model.ErrNotEnoughStock
		}
		skuReaded++
	}
	if rows.Err() != nil {
		return fmt.Errorf("%s: %w", op, rows.Err())
	}

	if skuReaded != len(items) {
		return model.ErrNotEnoughStock
	}

	if err := r.createReservationTx(ctx, tx, int64(orderID), items); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if err := r.updateStocksTx(ctx, tx, items); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *RepositoryStock) createReservationTx(ctx context.Context, tx pgx.Tx, orderID int64, items []model.Item) error {
	const op = `repositoryStock.CreateReservation`

	const query = `
	INSERT INTO stock_reserve (order_id, sku, quantity)
	VALUES ($1, $2, $3)`

	for _, item := range items {
		if _, err := tx.Exec(ctx, query, orderID, item.SKU, item.Count); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}
	return nil
}

func (r *RepositoryStock) updateStocksTx(ctx context.Context, tx pgx.Tx, items []model.Item) error {
	const op = `repositoryStock.UpdateStocks`

	const query = `
	UPDATE stocks
	SET quantity = quantity - $1
	WHERE sku = $2`

	for _, item := range items {
		if _, err := tx.Exec(ctx, query, item.Count, item.SKU); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}
	return nil
}
