// Package domain provides the core business models (e.g. domain, nullable, pagination).
// It sits at the center of architecture and has no dependencies on other layers.
package domain

var (
	UninializedID      = -1
	UninializedVersion = -1
)

// NewUserUninialized creates a new User entity before it is persisted to storage.
// The ID and Version are set to placeholder values until the database assigns them.
func NewUserUninialized(fullName string, phoneNumber *string) User {
	return NewUser(
		UninializedID,
		UninializedVersion,
		fullName,
		phoneNumber,
	)
}
