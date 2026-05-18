package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	logger "github.com/sparxfort1ano/go-todoapp/internal/core/logger"
	httpresponse "github.com/sparxfort1ano/go-todoapp/internal/core/transport/http/response"
	"go.uber.org/zap"
)

// CORS adds necessary headers to allow cross-origin requests from
// trusted domains.
func CORS(allowedOriginsList []string) Middleware {
	allowedOrigins := make(map[string]struct{})
	for _, origin := range allowedOriginsList {
		allowedOrigins[origin] = struct{}{}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			if _, ok := allowedOrigins[origin]; ok {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Accept")
			}

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

const requestIDHeader = "X-Request-ID"

// RequestID ensures every request has a unique identifier.
// It reads the X-Request-ID head from the client or generates a new UUID.
func RequestID() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get(requestIDHeader)
			if requestID == "" {
				requestID = uuid.NewString()
			}

			r.Header.Set(requestIDHeader, requestID)
			w.Header().Set(requestIDHeader, requestID)

			next.ServeHTTP(w, r)
		})
	}
}

// Logger injects a context-aware zap.Logger into the request context.
// It binds the request_id and URL to the logger for consistent structured logging.
func Logger(log *logger.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get(requestIDHeader)

			l := log.With(
				zap.String("request_id", requestID),
				zap.String("url", r.URL.String()),
			)

			ctx := logger.IntoContext(r.Context(), l)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// Trace logs the start and completion of an HTTP request handling.
// It prevents the server from crashing and returns a graceful 500 response.
func Trace() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// EXCLUSION: bypass Swagger routes without wrapping the ResponseWriter.
			// Swagger uses http.FileServer to serve large static UI files (JS/CSS).
			// FileServer relies on hidden interfaces like http.Flusher and io.ReaderFrom
			// for chunked transfer encoding. Our custom ResponseWriter wrapper hides
			// these interfaces, which would break the file streaming.
			// Furthermore, tracing static infrastructure files creates unnecessary log noise.
			if strings.HasPrefix(r.URL.Path, "/swagger/") {
				next.ServeHTTP(w, r)
				return
			}

			ctx := r.Context()
			log := logger.FromContext(ctx)
			rw := httpresponse.NewResponseWriter(w)

			before := time.Now()
			log.Debug(
				">>> incoming HTTP request",
				zap.String("http_method", r.Method),
				zap.Time("time", before.UTC()),
			)

			next.ServeHTTP(rw, r)

			log.Debug(
				"<<< done HTTP request",
				zap.Int("status_code", rw.StatusCode()),
				zap.Duration("latency", time.Since(before)),
			)
		})
	}
}

// Panic recovers from unexpected panics during HTTP request handling.
// It prevents the server from crashing and returns a graceful 500 response.
func Panic() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := logger.FromContext(ctx)
			responseHandler := httpresponse.NewHTTPResponseHandler(log, w)

			defer func() {
				if p := recover(); p != nil {
					responseHandler.PanicResponse(
						p,
						"during handle HTTP request got unexpected panic",
					)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
