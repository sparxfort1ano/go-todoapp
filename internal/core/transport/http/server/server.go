// Package server provides utilities to configure, run and gracefully shutdown
// the main HTTP server, along with API versioning and routing mechanisms.
package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/sparxfort1ano/go-todoapp/internal/core/logger"
	"github.com/sparxfort1ano/go-todoapp/internal/core/transport/http/middleware"
	"go.uber.org/zap"
)

// HTTPServer encapsulates the main HTTP multiplexer, global middleware chain,
// server configuration and the application logger.
type HTTPServer struct {
	mux        *http.ServeMux
	cfg        config
	log        *logger.Logger
	middleware []middleware.Middleware
}

// NewHTTPServer creates a new instance of HTTPServer.
func NewHTTPServer(cfg config, log *logger.Logger, middleware ...middleware.Middleware) *HTTPServer {
	return &HTTPServer{
		mux:        http.NewServeMux(),
		cfg:        cfg,
		log:        log,
		middleware: middleware,
	}
}

// RegisterAPIRouters mounts version-specific sub-routers (e.g. v1, v2) onto the main HTTP server.
// It automatically wraps each router with the version's specific middleware.
func (s *HTTPServer) RegisterAPIRouters(routers ...*APIVersionRouter) {
	for _, router := range routers {
		prefix := "/api/" + string(router.apiVersion)

		s.mux.Handle(
			prefix+"/",
			http.StripPrefix(prefix, router.withMiddleware()),
		)
	}
}

// Run starts the HTTP server, supporting its graceful shutdown.
func (s *HTTPServer) Run(ctx context.Context) error {
	mux := middleware.ChainMiddleware(s.mux, s.middleware...)

	server := &http.Server{
		Addr:    s.cfg.Addr,
		Handler: mux,
	}

	ch := make(chan error, 1)

	go func() {
		defer close(ch)

		s.log.Warn("start HTTP server", zap.String("addr", s.cfg.Addr))

		err := server.ListenAndServe()
		if !errors.Is(err, http.ErrServerClosed) {
			ch <- err
		}
	}()

	select {
	case err := <-ch:
		if err != nil {
			return fmt.Errorf("listen and serve HTTP: %w", err)
		}
	case <-ctx.Done():
		s.log.Warn("shutdown HTTP server...")

		shutdownCtx, cancel := context.WithTimeout(
			context.Background(),
			s.cfg.ShutdownTimeout,
		)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			server.Close()

			return fmt.Errorf("shutdown HTTP server: %w", err)
		}

		s.log.Warn("HTTP server stopped")
	}

	return nil
}
