package repository

// UpdateTaskParams includes the parameters of the outcoming port when updating a task.
type UpdateTaskParams struct {
	Task Task
}

func NewUpdateTaskParams(
	task Task,
) UpdateTaskParams {
	return UpdateTaskParams{
		Task: task,
	}
}

// UpdateTaskResult includes the result returned from the repository after updating the task.
type UpdateTaskResult struct {
	Task Task
}

func NewUpdateTaskResult(
	task Task,
) UpdateTaskResult {
	return UpdateTaskResult{
		Task: task,
	}
}
