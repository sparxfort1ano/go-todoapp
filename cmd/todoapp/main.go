package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/sparxfort1ano/go-todoapp/internal/core/logger"
	"github.com/sparxfort1ano/go-todoapp/internal/core/repository/postgres"
	"github.com/sparxfort1ano/go-todoapp/internal/core/transport/http/middleware"
	"github.com/sparxfort1ano/go-todoapp/internal/core/transport/http/server"
	userPostgres "github.com/sparxfort1ano/go-todoapp/internal/features/users/repository/postgres"
	userService "github.com/sparxfort1ano/go-todoapp/internal/features/users/service"
	userHTTP "github.com/sparxfort1ano/go-todoapp/internal/features/users/transport/http"
	"go.uber.org/zap"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	logger, err := logger.NewLogger(logger.NewConfigMust())
	if err != nil {
		fmt.Println("failed to init app logger:", err)
		os.Exit(1)
	}
	defer logger.Close()

	logger.Debug("initializing postgres conection pool")
	pool, err := postgres.NewConnectionPool(
		ctx,
		postgres.NewConfigMust(),
	)
	if err != nil {
		logger.Fatal("failed to initialize postgres connection pool", zap.Error(err))
	}
	defer pool.Close()

	logger.Debug("initializing feature", zap.String("feature", "users"))
	usersRepository := userPostgres.NewUsersRepository(pool)
	usersService := userService.NewUsersService(usersRepository)
	usersHTTPHandler := userHTTP.NewUsersHTTPHandler(usersService)

	logger.Debug("initializing HTTP server")
	httpServer := server.NewHTTPServer(
		server.NewConfigMust(),
		logger,
		middleware.RequestID(),
		middleware.Logger(logger),
		middleware.Panic(),
		middleware.Trace(),
	)
	apiVersionRouter := server.NewAPIVersionRouter(server.APIVersion1)
	apiVersionRouter.RegisterRoutes(usersHTTPHandler.Routes()...)
	httpServer.RegisterAPIRouters(apiVersionRouter)

	if err := httpServer.Run(ctx); err != nil {
		logger.Error("HTTP server run error")
	}
}
