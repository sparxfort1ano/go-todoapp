// Package filesystem acts as the repository layer for the Web feature.
// It interacts with the filesystem to extract necessary files.
package filesystem

type WebRepository struct{}

// NewWebRepository creates a new instance of WebRepository.
func NewWebRepository() *WebRepository {
	return &WebRepository{}
}
