// Package service acts as the service layer for the Web feature.
package service

// WebService encapsulates the core business logic for web management.
type WebService struct {
	webRepository WebRepository
}

// WebRepository defines the contract that decouples the service layer
// from the underlying repository logic.
type WebRepository interface {
	// GetFile extracts the index.html file from the filesystem.
	GetFile(filePath string) ([]byte, error)
}

// NewWebService creates a new instance of WebService.
func NewWebService(
	webRepository WebRepository,
) *WebService {
	return &WebService{
		webRepository: webRepository,
	}
}
