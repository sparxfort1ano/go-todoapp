package in

// PatchTaskParams includes the parameters of the incoming port when patching a task with using nullable pattern.
type PatchTaskParams struct {
	ID    int
	Patch TaskPatch
}

func NewPatchTaskParams(
	id int,
	patch TaskPatch,
) PatchTaskParams {
	return PatchTaskParams{
		ID:    id,
		Patch: patch,
	}
}

// PatchTaskResult includes the result returned from the service after patching the task.
type PatchTaskResult struct {
	Task Task
}

func NewPatchTaskResult(
	task Task,
) PatchTaskResult {
	return PatchTaskResult{
		Task: task,
	}
}
