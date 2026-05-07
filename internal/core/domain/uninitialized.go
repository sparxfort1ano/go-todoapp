// Package domain provides the core business models (e.g. domain, nullable, pagination).
// It sits at the center of architecture and has no dependencies on other layers.
package domain

import "time"

var (
	UninializedID      = -1
	UninializedVersion = -1
)

// NewUserUninitialized creates a new User entity before it is persisted to storage.
// The ID and Version are set to placeholder values until the database assigns them.
func NewUserUninitialized(fullName string, phoneNumber *string) User {
	return NewUser(
		UninializedID,
		UninializedVersion,
		fullName,
		phoneNumber,
	)
}

// NewTaskUninitialized creates a new Task entity before it is persisted to storage.
// The ID and Version are set to placeholder values until the database assigns them.
func NewTaskUninitialized(
	title string,
	description *string,
	authorUserID int,
) Task {
	return NewTask(
		UninializedID,
		UninializedVersion,
		title,
		description,
		false,
		time.Now(),
		nil,
		authorUserID,
	)
}
