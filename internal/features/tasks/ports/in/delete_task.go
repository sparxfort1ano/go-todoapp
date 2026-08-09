package in

// DeleteTaskParams includes the parameters of the incoming port when deleting a task.
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

// DeleteTaskResult includes the result returned from the service after deleting the task.
type DeleteTaskResult struct{}

func NewDeleteTaskResult() DeleteTaskResult {
	return DeleteTaskResult{}
}
