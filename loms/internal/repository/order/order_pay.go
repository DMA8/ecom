package order

import (
	"context"
	"ecom/loms/internal/model"
	"fmt"
	"log/slog"
)

// update status
func (r *Repository) OrderPay(ctx context.Context, orderPay model.OrderID) error {
	const op = `repositoryOrder.OrderPay`
	var err error

	tx, err := r.psqlConnection.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(ctx); rbErr != nil {
				slog.Error(fmt.Sprintf("err '%s' while rollbacking transaction", rbErr))
			}
		}
	}()

	if err = r.orderSetStatusPaidTx(ctx, tx, orderPay, model.OrderStatusPayed); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if err = r.stockRepo.ReserveRemoveByOrderIDTx(ctx, tx, orderPay); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return tx.Commit(ctx)
}
