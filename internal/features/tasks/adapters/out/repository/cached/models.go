package cached

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/sparxfort1ano/go-todoapp/internal/features/tasks/ports/out/repository"
)

// TaskModel represents a task to be serialized in JSON.
// Stored under the key "task:<id>" with a TTL set by the pool.
type TaskModel struct {
	ID           int        `json:"id"`
	Version      int        `json:"version"`
	Title        string     `json:"title"`
	Description  *string    `json:"description"`
	Completed    bool       `json:"completed"`
	CreatedAt    time.Time  `json:"created_at"`
	CompletedAt  *time.Time `json:"completed_at"`
	AuthorUserID int        `json:"author_user_id"`
}

func repoToModel(task repository.Task) TaskModel {
	return TaskModel{
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

func modelToRepo(model TaskModel) repository.Task {
	return repository.NewTask(
		model.ID,
		model.Version,
		model.Title,
		model.Description,
		model.Completed,
		model.CreatedAt,
		model.CompletedAt,
		model.AuthorUserID,
	)
}

func taskKey(id int) string {
	return fmt.Sprintf("task:%d", id)
}

func (m *TaskModel) Serialize() ([]byte, error) {
	bytes, err := json.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("serialize task: %w", err)
	}

	return bytes, nil
}

func (m *TaskModel) Deserialize(bytes []byte) error {
	if err := json.Unmarshal(bytes, m); err != nil {
		return fmt.Errorf("deserialize task: %w", err)
	}

	return nil
}

// TaskListModel represents a tasks list to be serialized in JSON.
// Stored in hash under the key "tasks:<userID>" or "tasks:<all>".
// A hash field is defined by tasksListField and encodes pagination.
type TaskListModel []TaskModel

func repoToModels(tasks []repository.Task) TaskListModel {
	taskList := make(TaskListModel, len(tasks))

	for i, task := range tasks {
		taskList[i] = repoToModel(task)
	}

	return taskList
}

func modelsToRepo(taskList TaskListModel) []repository.Task {
	tasks := make([]repository.Task, len(taskList))

	for i, task := range taskList {
		tasks[i] = modelToRepo(task)
	}

	return tasks
}

func tasksListKey(userID *int) string {
	if userID == nil {
		return "tasks:nil"
	}

	return fmt.Sprintf("tasks:%d", *userID)
}

func tasksListField(page repository.Pagination) string {
	getPageParamStr := func(v *int) string {
		if v == nil {
			return "nil"
		}

		return strconv.Itoa(*v)
	}

	return fmt.Sprintf("%s:%s", getPageParamStr(page.Limit), getPageParamStr(page.Offset))
}

func (m *TaskListModel) Serialize() ([]byte, error) {
	result, err := json.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("serialize task list: %w", err)
	}

	return result, nil
}

func (m *TaskListModel) Deserialize(bytes []byte) error {
	if err := json.Unmarshal(bytes, m); err != nil {
		return fmt.Errorf("deserialize task list: %w", err)
	}

	return nil
}
