package repository

// SaveTaskParams includes the parameters of the outcoming port when saving a new task.
type SaveTaskParams struct {
	Task Task
}

func NewSaveTaskParams(
	task Task,
) SaveTaskParams {
	return SaveTaskParams{
		Task: task,
	}
}

// SaveTaskResult includes the result returned from the repository after saving the new task.
type SaveTaskResult struct {
	Task Task
}

func NewSaveTaskResult(
	task Task,
) SaveTaskResult {
	return SaveTaskResult{
		Task: task,
	}
}
