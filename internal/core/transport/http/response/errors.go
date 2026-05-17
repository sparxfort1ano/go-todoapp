package response

// ErrorResponse represents the standard JSON body for HTTP error responses.
type ErrorResponse struct {
	Error   string `json:"error" example:"full error text"`
	Message string `json:"message" example:"short human-readable message"`
}
