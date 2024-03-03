package dto

type TaskData struct {
	Repo       string
	PipelineID string
	Name       string
	Command    string
	Status     string
	ExitCode   int32
	Output     string
	ExecTime   float64
}
