package domain

import (
	"fmt"
	"time"
	"unicode/utf8"

	errs "github.com/sparxfort1ano/go-todoapp/internal/core/errors"
)

// Task represents the core business entity of a task in the system.
// It contains all the essential data and business logic tied to a task.
type Task struct {
	ID      int
	Version int

	Title       string
	Description *string
	Completed   bool
	CreatedAt   time.Time
	CompletedAt *time.Time

	AuthorUserID int
}

// NewTask reconstitutes an existing Task entity from storage
// with a known ID and Version.
func NewTask(
	id int,
	version int,
	title string,
	description *string,
	completed bool,
	createdAt time.Time,
	completedAt *time.Time,
	authorUserID int,
) Task {
	return Task{
		ID:           id,
		Version:      version,
		Title:        title,
		Description:  description,
		Completed:    completed,
		CreatedAt:    createdAt,
		CompletedAt:  completedAt,
		AuthorUserID: authorUserID,
	}
}

// Validate checks whether the business rules for the Task entity are met.
// It returns error if the data contradicts the rules
// such as unappropriate length or bad logic
// for the task `created` and `completed` statuses.
func (t *Task) Validate() error {
	titleLen := utf8.RuneCountInString(t.Title)
	if titleLen < 1 || titleLen > 100 {
		return fmt.Errorf(
			"invalid `Title` len: %d: %w",
			titleLen,
			errs.ErrInvalidArgument,
		)
	}

	if t.Description != nil {
		descriptionLen := utf8.RuneCountInString(*t.Description)
		if descriptionLen < 1 || descriptionLen > 1000 {
			return fmt.Errorf(
				"invalid `Description` len: %d: %w",
				descriptionLen,
				errs.ErrInvalidArgument,
			)
		}
	}

	if t.Completed {
		if t.CompletedAt == nil {
			return fmt.Errorf(
				"`CompletedAt` can't be `nil` if `Completed`==`true`: %w",
				errs.ErrInvalidArgument,
			)
		}

		if t.CompletedAt.Before(t.CreatedAt) {
			return fmt.Errorf(
				"`CompletedAt` can't be before `CreatedAt`: %w",
				errs.ErrInvalidArgument,
			)
		}
	} else {
		if t.CompletedAt != nil {
			return fmt.Errorf(
				"`CompletedAt` must be `nil` if `Completed`==`false`: %w",
				errs.ErrInvalidArgument,
			)
		}
	}

	return nil
}

// TaskPatch represents the data used to partitially update an existing Task.
// Only the fields with Set=true are applied during the patch operation.
type TaskPatch struct {
	Title       Nullable[string]
	Description Nullable[string]
	Completed   Nullable[bool]
}

// NewTaskPatch creates a new instance of TaskPatch.
func NewTaskPatch(
	title Nullable[string],
	description Nullable[string],
	completed Nullable[bool],
) TaskPatch {
	return TaskPatch{
		Title:       title,
		Description: description,
		Completed:   completed,
	}
}

// Validate checks whether the TaskPatch data corresponds to
// the required fields of Task entity.
func (p *TaskPatch) Validate() error {
	if p.Title.Set && p.Title.Value == nil {
		return fmt.Errorf(
			"`Title` can't be patched to NULL: %w",
			errs.ErrInvalidArgument,
		)
	}

	if p.Completed.Set && p.Completed.Value == nil {
		return fmt.Errorf(
			"`Completed` can't be patched NULL: %w",
			errs.ErrInvalidArgument,
		)
	}

	return nil
}

// ApplyPatch modifies the Task entity using the provided TaskPatch data.
// It validates both the patch data
// and the resulting task state before applying changes.
func (t *Task) ApplyPatch(patch TaskPatch) error {
	if err := patch.Validate(); err != nil {
		return fmt.Errorf("validate task patch: %w", err)
	}

	tmp := *t
	if patch.Title.Set {
		tmp.Title = *patch.Title.Value
	}

	if patch.Description.Set {
		tmp.Description = patch.Description.Value
	}

	if patch.Completed.Set {
		tmp.Completed = *patch.Completed.Value

		if tmp.Completed {
			completedAt := time.Now()
			tmp.CompletedAt = &completedAt
		} else {
			tmp.CompletedAt = nil
		}
	}

	if err := tmp.Validate(); err != nil {
		return fmt.Errorf("validate patched task: %w", err)
	}

	*t = tmp

	return nil
}
