package in

// GetTasksParams includes the parameters of the incoming port when getting tasks.
type GetTasksParams struct {
	Pagination Pagination
	UserID     *int
}

func NewGetTasksParams(
	pagination Pagination,
	userID *int,
) GetTasksParams {
	return GetTasksParams{
		Pagination: pagination,
		UserID:     userID,
	}
}

// GetTasksResult includes the result returned from the service after getting the task.
type GetTasksResult struct {
	Tasks []Task
}

func NewGetTasksResult(
	tasks []Task,
) GetTasksResult {
	return GetTasksResult{
		Tasks: tasks,
	}
}
