package in

// CreateTaskParams includes the parameters of the incoming port when creating a task.
type CreateTaskParams struct {
	Title        string
	Description  *string
	AuthorUserID int
}

func NewCreateTaskParams(
	title string,
	description *string,
	authorUsedID int,
) CreateTaskParams {
	return CreateTaskParams{
		Title:        title,
		Description:  description,
		AuthorUserID: authorUsedID,
	}
}

// CreateTaskResult includes the result returned from the service after creating the task.
type CreateTaskResult struct {
	Task Task
}

func NewCreateTaskResult(
	task Task,
) CreateTaskResult {
	return CreateTaskResult{
		Task: task,
	}
}
