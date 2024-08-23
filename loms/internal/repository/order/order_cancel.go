package order

import (
	"context"
	"ecom/loms/internal/model"
	"fmt"
	"log/slog"
)

func (r *Repository) OrderCancel(ctx context.Context, order model.OrderID) error {
	const op = `repositoryOrder.OrderCancel`
	var err error

	tx, err := r.psqlConnection.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
				slog.Error(fmt.Sprintf("err '%s' while rollbacking transaction", rollbackErr))
			}
		}
	}()

	orderInfo, err := r.OrderInfoByOrderIDTx(ctx, tx, order)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if orderInfo.StatusCode == model.OrderStatusCancelled {
		err = fmt.Errorf("%s: %w", op, model.ErrOrderAlreadyCanceled)
		return err
	}

	skuReservedByOrder, err := r.stockRepo.SkuReservedByOrderTx(ctx, tx, order)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if err = r.stockRepo.IncreaseStockTx(ctx, tx, skuReservedByOrder); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if err = r.stockRepo.ReserveRemoveByOrderIDTx(ctx, tx, order); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if err = r.orderSetStatusTx(ctx, tx, order, model.OrderStatusCancelled); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return tx.Commit(ctx)
}
