package repository

// DeleteTaskParams includes the parameters of the outcoming port when deleting a task.
type DeleteTaskParams struct {
	ID int
}

func NewDeleteTaskParams(
	id int,
) DeleteTaskParams {
	return DeleteTaskParams{
		ID: id,
	}
}

// DeleteTaskResult includes the result returned from the repository after deleting the task.
type DeleteTaskResult struct{}

func NewDeleteTaskResult() DeleteTaskResult {
	return DeleteTaskResult{}
}
