package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	core_logger "github.com/D1skord/todo/internal/core/logger"
	core_postgres_pool "github.com/D1skord/todo/internal/core/repository/postgres/pool"
	core_http_middleware "github.com/D1skord/todo/internal/core/transport/http/middleware"
	core_http_server "github.com/D1skord/todo/internal/core/transport/http/server"
	users_postgres_repository "github.com/D1skord/todo/internal/features/users/repository/postgres"
	users_service "github.com/D1skord/todo/internal/features/users/service"
	users_transport_http "github.com/D1skord/todo/internal/features/users/transport/http"
	"go.uber.org/zap"
)

func main() {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer cancel()

	logger, err := core_logger.NewLogger(core_logger.NewConfigMust())
	if err != nil {
		fmt.Println("failed to initialize logger:", err)
		os.Exit(1)
	}
	defer logger.Close()

	logger.Debug("Starting ToDo application...")

	logger.Debug("initialize connection pool")
	pool, err := core_postgres_pool.NewConnectionPool(
		ctx,
		core_postgres_pool.NewConfigMust(),
	)

	if err != nil {
		logger.Fatal("failed to init postgres connection pool", zap.Error(err))
	}
	defer pool.Close()

	logger.Debug("initializing feature", zap.String("feature", "users"))

	usersRepository := users_postgres_repository.NewUsersRepository(pool)
	usersService := users_service.NewUserService(usersRepository)
	usersTransportHTTP := users_transport_http.NewUsersHTTPHandler(usersService)

	logger.Debug("initializing HTTP server")
	httpServer := core_http_server.NewHTTPServer(
		core_http_server.NewConfigMust(),
		logger,
		core_http_middleware.RequestId(),
		core_http_middleware.Logger(logger),
		core_http_middleware.Panic(),
		core_http_middleware.Trace(logger),
	)

	apiVersionRouter := core_http_server.NewAPIVersionRouter(core_http_server.APIVersion1)
	apiVersionRouter.RegisterRoutes(usersTransportHTTP.Routes()...)
	httpServer.RegisterAPIRouters(apiVersionRouter)

	if err := httpServer.Run(ctx); err != nil {
		logger.Error("HTTP server run error", zap.Error(err))
	}
}
