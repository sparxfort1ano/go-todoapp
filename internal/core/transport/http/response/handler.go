// Package response provides utilities for formatiing HTTP responses,
// mapping domain errors to appropriate HTTP status codes and logging.
package response

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	errs "github.com/sparxfort1ano/go-todoapp/internal/core/errors"
	"github.com/sparxfort1ano/go-todoapp/internal/core/logger"
	"go.uber.org/zap"
)

// HTTPResponseHandler wraps an http.ResponseWriter to provide standardized
// JSON formatting and automated error logging capabilities.
type HTTPResponseHandler struct {
	log *logger.Logger
	rw  http.ResponseWriter
}

// NewHTTPResponseHandler creates a new instance of HTTPResponseHandler.
func NewHTTPResponseHandler(log *logger.Logger, rw http.ResponseWriter) *HTTPResponseHandler {
	return &HTTPResponseHandler{
		log: log,
		rw:  rw,
	}
}

// JSONResponse serializes the response body to JSON, sets the HTTP status code
// and logs an error if the encoding process fails.
func (h *HTTPResponseHandler) JSONResponse(
	responseBody any,
	statusCode int,
) {
	h.rw.WriteHeader(statusCode)

	if err := json.NewEncoder(h.rw).Encode(responseBody); err != nil {
		h.log.Error("write HTTP response", zap.Error(err))
	}
}

func (h *HTTPResponseHandler) NoContentResponse() {
	h.rw.WriteHeader(http.StatusNoContent)
}

// errorResponse writes a standard JSON error structure to the client
// and logs the internal error details.
func (h *HTTPResponseHandler) errorResponse(
	statusCode int,
	err error,
	msg string,
) {
	response := map[string]string{
		"message": msg,
		"error":   err.Error(),
	}
	h.JSONResponse(
		response,
		statusCode,
	)
}

// ErrorResponse maps sentinel errors to the correct
// HTTP status codes and logging level, ensuring uniform error handling across the app.
func (h HTTPResponseHandler) ErrorResponse(err error, msg string) {
	var (
		statusCode int
		logFunc    func(string, ...zap.Field)
	)

	switch {
	case errors.Is(err, errs.ErrInvalidArgument):
		statusCode = http.StatusBadRequest
		logFunc = h.log.Warn
	case errors.Is(err, errs.ErrNotFound):
		statusCode = http.StatusNotFound
		logFunc = h.log.Debug
	case errors.Is(err, errs.ErrConflict):
		statusCode = http.StatusConflict
		logFunc = h.log.Warn
	default:
		statusCode = http.StatusInternalServerError
		logFunc = h.log.Error
	}

	logFunc(msg, zap.Error(err))

	h.errorResponse(
		statusCode,
		err,
		msg,
	)
}

// PanicResponse handles unexpected panics by logging the stack trace
// and returning Internal Server Error to the client.
func (h *HTTPResponseHandler) PanicResponse(p any, msg string) {
	err := fmt.Errorf("unexpected panic: %v", p)
	h.log.Error(msg, zap.Error(err))

	h.errorResponse(
		http.StatusInternalServerError,
		err,
		msg,
	)
}
