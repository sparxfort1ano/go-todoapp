package http

import (
	"time"

	"github.com/sparxfort1ano/go-todoapp/internal/features/tasks/ports/in"
)

// TaskDTOResponse represents the outgoing JSON body for multiple task features.
type TaskDTOResponse struct {
	ID           int        `json:"id" example:"15"`
	Version      int        `json:"version" example:"3"`
	Title        string     `json:"title" example:"Домашка"`
	Description  *string    `json:"description" example:"Сделать до четверга математику"`
	Completed    bool       `json:"completed" example:"false"`
	CreatedAt    time.Time  `json:"created_at" example:"2026-02-26T10:30:00Z"`
	CompletedAt  *time.Time `json:"completed_at" example:"null"`
	AuthorUserID int        `json:"author_user_id" example:"5"`
}

func taskDTOFromIn(task in.Task) TaskDTOResponse {
	return TaskDTOResponse{
		ID:           task.ID,
		Version:      task.Version,
		Title:        task.Title,
		Description:  task.Description,
		Completed:    task.Completed,
		CreatedAt:    task.CreatedAt,
		CompletedAt:  task.CompletedAt,
		AuthorUserID: task.AuthorUserID,
	}
}

func tasksDTOFromIn(tasks []in.Task) []TaskDTOResponse {
	tasksDTO := make([]TaskDTOResponse, len(tasks))

	for i, task := range tasks {
		tasksDTO[i] = taskDTOFromIn(task)
	}

	return tasksDTO
}
