package service

import (
	"github.com/sparxfort1ano/go-todoapp/internal/core/domain"
	"github.com/sparxfort1ano/go-todoapp/internal/features/tasks/ports/in"
	"github.com/sparxfort1ano/go-todoapp/internal/features/tasks/ports/out/repository"
)

func domainToRepo(task domain.Task) repository.Task {
	return repository.NewTask(
		task.ID,
		task.Version,
		task.Title,
		task.Description,
		task.Completed,
		task.CreatedAt,
		task.CompletedAt,
		task.AuthorUserID,
	)
}

func repoToDomain(task repository.Task) domain.Task {
	return domain.NewTask(
		task.ID,
		task.Version,
		task.Title,
		task.Description,
		task.Completed,
		task.CreatedAt,
		task.CompletedAt,
		task.AuthorUserID,
	)
}

func repoTaskToIn(task repository.Task) in.Task {
	return in.NewTask(
		task.ID,
		task.Version,
		task.Title,
		task.Description,
		task.Completed,
		task.CreatedAt,
		task.CompletedAt,
		task.AuthorUserID,
	)
}

func outTasksToIn(tasks []repository.Task) []in.Task {
	tasksIn := make([]in.Task, len(tasks))

	for i, task := range tasks {
		tasksIn[i] = repoTaskToIn(task)
	}

	return tasksIn
}
