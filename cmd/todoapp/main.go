package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sparxfort1ano/go-todoapp/internal/core/config"
	"github.com/sparxfort1ano/go-todoapp/internal/core/logger"
	"github.com/sparxfort1ano/go-todoapp/internal/core/repository/postgres/pgxpool"
	redispool "github.com/sparxfort1ano/go-todoapp/internal/core/repository/redis/goredis"
	"github.com/sparxfort1ano/go-todoapp/internal/core/transport/http/middleware"
	"github.com/sparxfort1ano/go-todoapp/internal/core/transport/http/server"
	statsPostgres "github.com/sparxfort1ano/go-todoapp/internal/features/statistics/repository/postgres"
	statsService "github.com/sparxfort1ano/go-todoapp/internal/features/statistics/service"
	statsHTTP "github.com/sparxfort1ano/go-todoapp/internal/features/statistics/transport/http"
	tasksHTTP "github.com/sparxfort1ano/go-todoapp/internal/features/tasks/adapters/in/transport/http"
	tasksCached "github.com/sparxfort1ano/go-todoapp/internal/features/tasks/adapters/out/repository/cached"
	tasksPostgres "github.com/sparxfort1ano/go-todoapp/internal/features/tasks/adapters/out/repository/postgres"
	tasksService "github.com/sparxfort1ano/go-todoapp/internal/features/tasks/service"
	usersPostgres "github.com/sparxfort1ano/go-todoapp/internal/features/users/repository/postgres"
	usersService "github.com/sparxfort1ano/go-todoapp/internal/features/users/service"
	usersHTTP "github.com/sparxfort1ano/go-todoapp/internal/features/users/transport/http"

	"go.uber.org/zap"

	_ "github.com/sparxfort1ano/go-todoapp/docs"
)

// @title 		Golang Todo API
// @version 	1.0
// @description Todo Application REST-API scheme
// @host 		localhost:5050
// @BasePath 	/api/v1
func main() {
	cfg := config.NewConfigMust()
	time.Local = cfg.TimeZone

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	logger, err := logger.NewLogger(logger.NewConfigMust())
	if err != nil {
		fmt.Println("failed to init app logger:", err)
		os.Exit(1)
	}
	defer logger.Close()

	logger.Debug("application time zone", zap.Any("zone", time.Local))

	logger.Debug("initializing postgres connection pool")
	pgxPool, err := pgxpool.NewPool(
		ctx,
		pgxpool.NewConfigMust(),
	)
	if err != nil {
		logger.Fatal("failed to initialize postgres connection pool", zap.Error(err))
	}
	defer pgxPool.Close()

	logger.Debug("initializing redis connection pool")
	redisPool, err := redispool.NewPool(
		ctx,
		redispool.NewConfigMust(),
	)
	if err != nil {
		logger.Fatal("failed to initialize redis client", zap.Error(err))
	}
	defer redisPool.Close()

	logger.Debug("initializing feature", zap.String("feature", "users"))
	usersRepository := usersPostgres.NewUsersRepository(pgxPool)
	usersService := usersService.NewUsersService(usersRepository)
	usersHTTPHandler := usersHTTP.NewUsersHTTPHandler(usersService)

	logger.Debug("initializing feature", zap.String("feature", "tasks"))
	tasksRepository := tasksCached.NewCachedRepository(
		redisPool,
		tasksPostgres.NewTasksRepository(pgxPool),
	)
	tasksService := tasksService.NewTaskService(tasksRepository)
	tasksHTTPHandler := tasksHTTP.NewTaskHTTPHandler(tasksService)

	logger.Debug("initializing feature", zap.String("feature", "statistics"))
	statisticsRepository := statsPostgres.NewStatisticsRepository(pgxPool)
	statisticsService := statsService.NewStatisticsService(statisticsRepository)
	statisticsHTTPHandler := statsHTTP.NewStatisticsHTTPHandler(statisticsService)

	logger.Debug("initializing HTTP server")
	httpConfig := server.NewConfigMust()
	httpServer := server.NewHTTPServer(
		httpConfig,
		logger,
		middleware.CORS(httpConfig.AllowedOrigins),
		middleware.RequestID(),
		middleware.Logger(logger),
		middleware.Trace(),
		middleware.Panic(),
	)
	apiVersionRouterV1 := server.NewAPIVersionRouter(server.APIVersion1)
	apiVersionRouterV1.RegisterRoutes(usersHTTPHandler.Routes()...)
	apiVersionRouterV1.RegisterRoutes(tasksHTTPHandler.Routes()...)
	apiVersionRouterV1.RegisterRoutes(statisticsHTTPHandler.Routes()...)

	httpServer.RegisterAPIRouters(
		apiVersionRouterV1,
	)
	httpServer.RegisterSwagger()

	if err := httpServer.Run(ctx); err != nil {
		logger.Error("HTTP server run error")
	}
}
