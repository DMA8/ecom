package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"ecom/cart/pkg/logger"
	"ecom/loms/internal/adapters/grpcadapter"
	"ecom/loms/internal/adapters/kafka/producer"
	"ecom/loms/internal/config"
	"ecom/loms/internal/loms"
	"ecom/loms/internal/repository/order"
	"ecom/loms/internal/repository/stock"
	"ecom/loms/pkg/migrations"

	"github.com/IBM/sarama"
	"github.com/jackc/pgx/v5"
)

/*
1) Развернуть kafka кластер в docker-compose.
2) Интегрировать LOMS с кафкой: LOMS пишет в кафку изменения статуса заказа (создание/отмена).
3) Создать сервис нотификаций(просто main с консьюмер группой достаточно).
	Сервис нотификаций должен вычитывать сообщения о статусе заказа и отправлять нотификации
	об изменениях статуса заказа (просто писать в лог сообщение с номером заказа и статусом).
4) Обеспечить гарантию отправки сообщения о статусе заказа (отказоустойчевый кластер/реплики, ack и вот это все).
5) Обеспечить упорядоченную обработку статусов одного заказа (добиться упорядоченной отправки событий по заказу и чтение событий по заказу).
6) Обеспечbть семантику exectly once (рассказывал как)
7) 💎 Применить паттерн Outbox.

*/

/*
(Logistics and Order Management System)
Сервис отвечает за учет заказов и стоки по товарам.
*/
func main() {
	var exitCode int
	wg := &sync.WaitGroup{}

	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	logger.SetDefaultLogger(cfg.LogLevel)

	if err := migrations.Migrate(cfg.PSQLConnStr); err != nil {
		log.Println(err)
	}
	slog.Info("migrations ok!")

	ctx, cancelGlobal := context.WithCancel(context.Background())

	psqlConnection, err := initDB(ctx, cfg)
	if err != nil {
		slog.Error(fmt.Sprintf("couln't connect to psql err: %s", err.Error()))
		os.Exit(1)
	}

	repoStock := stock.New(psqlConnection)
	repoOrder := order.New(psqlConnection, repoStock)

	producer, err := producer.New(
		strings.Split(cfg.KafkaBrokers, ";"),
		producer.WithRequiredAcks(sarama.NoResponse),
		producer.WithProducerPartitioner(sarama.NewHashPartitioner),
		producer.WithMaxOpenRequests(5),
		producer.WithMaxRetries(5),
		producer.WithRetryBackoff(10*time.Millisecond),
		producer.WithProducerFlushMessages(3),
		producer.WithProducerFlushFrequency(5*time.Second),
	)
	if err != nil {
		slog.Error(fmt.Sprintf("kafka is not ready: %s", err.Error()))
		os.Exit(1)
	}

	manager := loms.New(repoOrder, repoStock, producer)

	grpcServer := grpcadapter.New(manager)
	servgrpc := grpcServer.Server()

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer func() {
			if r := recover(); r != nil {
				slog.Info("panic recover servgrpc.Serve()", "err", r)
				exitCode = 1
				cancelGlobal()
			}
		}()
		grpcLis, err := net.Listen(cfg.GrpcNetwork, cfg.ServiceURLGrpc)
		if err != nil {
			slog.Error(fmt.Sprintf("couln't create grpc listener err: %s", err.Error()))
			exitCode = 1
			cancelGlobal()
			return
		}
		slog.Info("launching grpc serv", "port", cfg.ServiceURLGrpc)
		if err := servgrpc.Serve(grpcLis); err != nil {
			slog.Error(fmt.Sprintf("couln't run grpc serv err: %s", err.Error()))
			exitCode = 1
		}
		cancelGlobal()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

		select {
		case <-ctx.Done():
			slog.Info("context done, shutting down")
		case <-signals:
			slog.Info("got signal, shutting down")
		}
		servgrpc.GracefulStop()
	}()
	wg.Wait()
	slog.Info("shutdown ok!", "exitCode", exitCode)
	os.Exit(exitCode)
}

func initDB(ctx context.Context, cfg *config.Config) (*pgx.Conn, error) {
	const op = "main.initDB"

	psqlConnection, err := pgx.Connect(ctx, cfg.PSQLConnStr)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if err := psqlConnection.Ping(ctx); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return psqlConnection, nil
}
