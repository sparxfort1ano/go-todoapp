package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sparxfort1ano/go-todoapp/internal/core/logger"
	"github.com/sparxfort1ano/go-todoapp/internal/core/repository/postgres/pgxpool"
	"github.com/sparxfort1ano/go-todoapp/internal/core/transport/http/middleware"
	"github.com/sparxfort1ano/go-todoapp/internal/core/transport/http/server"
	tasksPostgres "github.com/sparxfort1ano/go-todoapp/internal/features/tasks/repository/postgres"
	tasksService "github.com/sparxfort1ano/go-todoapp/internal/features/tasks/service"
	tasksHTTP "github.com/sparxfort1ano/go-todoapp/internal/features/tasks/transport/http"
	usersPostgres "github.com/sparxfort1ano/go-todoapp/internal/features/users/repository/postgres"
	usersService "github.com/sparxfort1ano/go-todoapp/internal/features/users/service"
	usersHTTP "github.com/sparxfort1ano/go-todoapp/internal/features/users/transport/http"
	"go.uber.org/zap"
)

var (
	timeZone = time.UTC
)

func main() {
	time.Local = timeZone

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	logger, err := logger.NewLogger(logger.NewConfigMust())
	if err != nil {
		fmt.Println("failed to init app logger:", err)
		os.Exit(1)
	}
	defer logger.Close()

	logger.Debug("application time zone", zap.Any("zone", timeZone))

	logger.Debug("initializing postgres conection pool")
	pool, err := pgxpool.NewPool(
		ctx,
		pgxpool.NewConfigMust(),
	)
	if err != nil {
		logger.Fatal("failed to initialize postgres connection pool", zap.Error(err))
	}
	defer pool.Close()

	logger.Debug("initializing feature", zap.String("feature", "users"))
	usersRepository := usersPostgres.NewUsersRepository(pool)
	usersService := usersService.NewUsersService(usersRepository)
	usersHTTPHandler := usersHTTP.NewUsersHTTPHandler(usersService)

	logger.Debug("initializing feature", zap.String("feature", "tasks"))
	tasksRepository := tasksPostgres.NewTasksRepository(pool)
	tasksService := tasksService.NewTaskService(tasksRepository)
	tasksHTTPHandler := tasksHTTP.NewTaskHTTPHandler(tasksService)

	logger.Debug("initializing HTTP server")
	httpServer := server.NewHTTPServer(
		server.NewConfigMust(),
		logger,
		middleware.RequestID(),
		middleware.Logger(logger),
		middleware.Trace(),
		middleware.Panic(),
	)
	apiVersionRouterV1 := server.NewAPIVersionRouter(server.APIVersion1)
	apiVersionRouterV1.RegisterRoutes(usersHTTPHandler.Routes()...)
	apiVersionRouterV1.RegisterRoutes(tasksHTTPHandler.Routes()...)

	httpServer.RegisterAPIRouters(
		apiVersionRouterV1,
	)

	if err := httpServer.Run(ctx); err != nil {
		logger.Error("HTTP server run error")
	}
}
