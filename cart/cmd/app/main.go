package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"ecom/cart/internal/adapters/grpcadapter"
	"ecom/cart/internal/cart"
	"ecom/cart/internal/clients/lomsgrpc"
	"ecom/cart/internal/clients/productservicegrpc"
	"ecom/cart/internal/config"
	"ecom/cart/internal/repository"
	product "ecom/cart/pkg/api/productService"
	"ecom/cart/pkg/logger"
	grpc_loms "ecom/loms/pkg/api/loms/v1"
	"ecom/loms/pkg/migrations"

	"github.com/jackc/pgx/v5"
	"google.golang.org/grpc"
)

/*
Сервис отвечает за пользовательскую корзину и позволяет оформить заказ.
*/
func main() {
	var exitCode int
	wg := &sync.WaitGroup{}

	ctx, cancelGlobal := context.WithCancel(context.Background())

	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	if err := migrations.Migrate(cfg.PSQLConnStr); err != nil {
		log.Println(err)
	}
	slog.Info("migrations ok!")
	psqlConnection, err := initDB(ctx, cfg)
	if err != nil {
		slog.Error(fmt.Sprintf("couln't connect to psql err: %s", err.Error()))
		os.Exit(1)
	}

	logger.SetDefaultLogger(cfg.LogLevel)
	repo := repository.New(psqlConnection)

	ctx = context.WithValue(ctx, config.ConfigKey, cfg)

	connLoms, err := grpc.DialContext(
		ctx,
		cfg.LomsURLGrpc,
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatal(err)
	}
	rateLim := productservicegrpc.NewRateLimInterceptor(cfg.ProductServiceRateLimRPS)

	connProductService, err := grpc.DialContext(
		ctx,
		cfg.ProductServiceGrpcURL,
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(rateLim.RateLimiterInterceptor),
	)
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		rateLim.ResetCounterWorker(ctx)
	}()

	productServiceGrpcCli := product.NewProductServiceClient(connProductService)
	productServiceFacade := productservicegrpc.New(productServiceGrpcCli)
	clientLoms := grpc_loms.NewLOMSServiceClient(connLoms)
	lomsFacadeGRPC := lomsgrpc.New(clientLoms)

	cart := cart.New(repo, productServiceFacade, lomsFacadeGRPC)

	grpcServer := grpcadapter.New(cart)
	grpcServ := grpcServer.Server()

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer func() {
			if r := recover(); r != nil {
				slog.Info("panic recover srv.ListenAndServe()", "err", r)
				exitCode = 1
				cancelGlobal()
			}
		}()
		grpcListen, err := net.Listen(cfg.GrpcNetwork, cfg.ServiceURLGrpc)
		if err != nil {
			slog.Error(fmt.Sprintf("couln't run grpc serv err: %s", err.Error()))
			exitCode = 1
			cancelGlobal()
			return
		}
		slog.Info("launching grpc serv", "port", cfg.ServiceURLGrpc)
		if err := grpcServ.Serve(grpcListen); err != nil {
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
		grpcServ.GracefulStop()
	}()
	wg.Wait()
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
