package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	core_config "github.com/D1skord/todo/internal/core/config"
	core_logger "github.com/D1skord/todo/internal/core/logger"
	core_pgx_pool "github.com/D1skord/todo/internal/core/repository/postgres/pool/pgx"
	core_http_middleware "github.com/D1skord/todo/internal/core/transport/http/middleware"
	core_http_server "github.com/D1skord/todo/internal/core/transport/http/server"
	statistics_postgres_repository "github.com/D1skord/todo/internal/features/statistics/repository/postgres"
	statistics_service "github.com/D1skord/todo/internal/features/statistics/service"
	statistics_transport_http "github.com/D1skord/todo/internal/features/statistics/transport/http"
	tasks_postgres_repository "github.com/D1skord/todo/internal/features/tasks/repository"
	tasks_service "github.com/D1skord/todo/internal/features/tasks/service"
	tasks_transport_http "github.com/D1skord/todo/internal/features/tasks/transport/http"
	users_postgres_repository "github.com/D1skord/todo/internal/features/users/repository/postgres"
	users_service "github.com/D1skord/todo/internal/features/users/service"
	users_transport_http "github.com/D1skord/todo/internal/features/users/transport/http"
	"go.uber.org/zap"

	_ "github.com/D1skord/todo/docs"
)

// @title 			Golang ToDo API
// @version 		1.0
// @description		Todo Application REST-API scheme
// @host 			127.0.0.1:5050
// @BasePath 		/api/v1
func main() {
	cfg := core_config.NewConfigMust()
	time.Local = cfg.TimeZone

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

	logger.Debug("application time zone", zap.Any("time_zone", time.Local))

	logger.Debug("Starting ToDo application...")

	logger.Debug("initialize connection pool")

	pool, err := core_pgx_pool.NewPool(
		ctx,
		core_pgx_pool.NewConfigMust(),
	)

	if err != nil {
		logger.Fatal("failed to init postgres connection pool", zap.Error(err))
	}
	defer pool.Close()

	logger.Debug("initializing feature", zap.String("feature", "users"))

	usersRepository := users_postgres_repository.NewUsersRepository(pool)
	usersService := users_service.NewUserService(usersRepository)
	usersTransportHTTP := users_transport_http.NewUsersHTTPHandler(usersService)

	logger.Debug("initializing feature", zap.String("feature", "tasks"))
	tasksRepository := tasks_postgres_repository.NewTasksRepository(pool)
	tasksService := tasks_service.NewTasksService(tasksRepository)
	tasksTransportHTTP := tasks_transport_http.NewTasksHTTPHandler(tasksService)

	logger.Debug("initializing feature", zap.String("feature", "statistics"))
	statisticsRepository := statistics_postgres_repository.NewStatisticsRepository(pool)
	statisticsService := statistics_service.NewStatisticsService(statisticsRepository)
	statisticsTransportHTTP := statistics_transport_http.NewStatisticsHTTPHandler(statisticsService)

	httpConfig := core_http_server.NewConfigMust()
	logger.Debug("initializing HTTP server")
	httpServer := core_http_server.NewHTTPServer(
		core_http_server.NewConfigMust(),
		logger,
		core_http_middleware.CORS(httpConfig.AllowedOrigins),
		core_http_middleware.RequestId(),
		core_http_middleware.Logger(logger),
		core_http_middleware.Trace(logger),
		core_http_middleware.Panic(),
	)

	apiVersionRouterV1 := core_http_server.NewAPIVersionRouter(core_http_server.APIVersion1)
	apiVersionRouterV1.RegisterRoutes(usersTransportHTTP.Routes()...)
	apiVersionRouterV1.RegisterRoutes(tasksTransportHTTP.Routes()...)
	apiVersionRouterV1.RegisterRoutes(statisticsTransportHTTP.Routes()...)

	/*
		Example of usage apiVersionRouterV2 with separate MiddleWares
		//apiVersionRouterV2 := core_http_server.NewAPIVersionRouter(
		//	core_http_server.APIVersion2,
		//	core_http_middleware.Dummy("api v2 middleware"),
		//)
		//apiVersionRouterV2.RegisterRoutes(usersTransportHTTP.Routes()...)

		*
	*/
	httpServer.RegisterAPIRouters(
		apiVersionRouterV1,
		//apiVersionRouterV2,
	)

	httpServer.RegisterSwagger()

	if err := httpServer.Run(ctx); err != nil {
		logger.Error("HTTP server run error", zap.Error(err))
	}
}
