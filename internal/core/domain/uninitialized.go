// Package domain provides the core business models (e.g. domain, nullable, pagination, statistics).
// It sits at the center of architecture and has no dependencies on other layers.
package domain

import "time"

var (
	UninitalizedID      = -1
	UninitalizedVersion = -1
)

// NewUserUninitialized creates a new User entity before it is persisted to storage.
// The ID and Version are set to placeholder values until the database assigns them.
func NewUserUninitialized(fullName string, phoneNumber *string) User {
	return NewUser(
		UninitalizedID,
		UninitalizedVersion,
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
		UninitalizedID,
		UninitalizedVersion,
		title,
		description,
		false,
		time.Now(),
		nil,
		authorUserID,
	)
}
