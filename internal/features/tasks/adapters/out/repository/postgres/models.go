package postgres

import (
	"time"

	"github.com/sparxfort1ano/go-todoapp/internal/features/tasks/ports/out/repository"
)

// TaskModel represents the database schema for a user (DAO).
type TaskModel struct {
	ID           int
	Version      int
	Title        string
	Description  *string
	Completed    bool
	CreatedAt    time.Time
	CompletedAt  *time.Time
	AuthorUserID int
}

func modelToRepo(taskModel TaskModel) repository.Task {
	return repository.NewTask(
		taskModel.ID,
		taskModel.Version,
		taskModel.Title,
		taskModel.Description,
		taskModel.Completed,
		taskModel.CreatedAt,
		taskModel.CompletedAt,
		taskModel.AuthorUserID,
	)
}

func modelsToRepo(taskModels []TaskModel) []repository.Task {
	tasks := make([]repository.Task, len(taskModels))

	for i, model := range taskModels {
		tasks[i] = modelToRepo(model)
	}

	return tasks
}
