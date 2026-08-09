package in

// GetTaskParams includes the parameters of the incoming port when getting a task.
type GetTaskParams struct {
	ID int
}

func NewGetTaskParams(
	id int,
) GetTaskParams {
	return GetTaskParams{
		ID: id,
	}
}

// GetTaskResult includes the result returned from the service after getting the task.
type GetTaskResult struct {
	Task Task
}

func NewGetTaskResult(
	task Task,
) GetTaskResult {
	return GetTaskResult{
		Task: task,
	}
}
