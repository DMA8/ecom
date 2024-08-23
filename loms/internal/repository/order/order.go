package order

import (
	"context"
	"ecom/loms/internal/model"

	"github.com/jackc/pgx/v5"
)

type StockRepository interface {
	ReserveSKUsTx(ctx context.Context, tx pgx.Tx, orderID model.OrderID, items []model.Item) error
	ReserveRemoveByOrderIDTx(ctx context.Context, tx pgx.Tx, orderID model.OrderID) error
	SkuReservedByOrderTx(ctx context.Context, tx pgx.Tx, orderID model.OrderID) ([]model.Item, error)
	IncreaseStockTx(ctx context.Context, tx pgx.Tx, items []model.Item) error
}

type Repository struct {
	psqlConnection *pgx.Conn
	stockRepo      StockRepository
}

func New(psqlConnection *pgx.Conn, stock StockRepository) *Repository {
	return &Repository{
		psqlConnection: psqlConnection,
		stockRepo:      stock,
	}
}
