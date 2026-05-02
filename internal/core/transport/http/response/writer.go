package response

import "net/http"

const StatusCodeUninitialized = -1

// ResponseWriter is a custom decorator around the standard http.ResponseWriter.
// It intercepts and stores the HTTP status code so it can be read later (e.g., by logging middleware).
type ResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

// NewResponseWriter creates a new instance of ResponseWriter.
func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{
		ResponseWriter: w,
		statusCode:     StatusCodeUninitialized,
	}
}

// WriteHeader overrides the underlying WriteHeader method to capture an store
// the status code in memory before sending it to the client.
func (rw *ResponseWriter) WriteHeader(statusCode int) {
	rw.ResponseWriter.WriteHeader(statusCode)
	rw.statusCode = statusCode
}

// StatusCode retrieves the captured HTTP status code.
// It sets the 200 status code if a response was sent
// without explicitly setting a status code.
func (rw *ResponseWriter) StatusCode() int {
	if rw.statusCode == StatusCodeUninitialized {
		panic("no status code set")
	}
	return rw.statusCode
}
