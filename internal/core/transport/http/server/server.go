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
func (h *HTTPServer) RegisterAPIRouters(routers ...*APIVersionRouter) {
	for _, router := range routers {
		prefix := "/api/" + string(router.apiVersion)

		h.mux.Handle(
			prefix+"/",
			http.StripPrefix(prefix, router),
		)
	}
}

// Run starts the HTTP server, supporting its graceful shutdown.
func (h *HTTPServer) Run(ctx context.Context) error {
	mux := middleware.ChainMiddleware(h.mux, h.middleware...)

	server := &http.Server{
		Addr:    h.cfg.Addr,
		Handler: mux,
	}

	ch := make(chan error, 1)

	go func() {
		defer close(ch)

		h.log.Warn("start HTTP server", zap.String("addr", h.cfg.Addr))

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
		h.log.Warn("shutdown HTTP server...")

		shutdownCtx, cancel := context.WithTimeout(
			context.Background(),
			h.cfg.ShutdownTimeout,
		)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			server.Close()

			return fmt.Errorf("shutdown HTTP server: %w", err)
		}

		h.log.Warn("HTTP server stopped")
	}

	return nil
}
