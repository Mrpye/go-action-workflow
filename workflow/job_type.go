package workflow

type Job struct {
	Key           string               `json:"key" yaml:"key"`
	Title         string               `json:"title" yaml:"title"`
	Description   string               `json:"description" yaml:"description"`
	Actions       []Action             `json:"actions" yaml:"actions"`
	IsSubWorkflow bool                 `json:"is_sub_workflow" yaml:"is_sub_workflow"`
	Inputs        map[string]Parameter `json:"inputs" yaml:"inputs"`
}

type JobOption func(*Job)

func OptionJobKey(v string) JobOption {
	return func(h *Job) {
		h.Key = v
	}
}

func OptionJobTitle(v string) JobOption {
	return func(h *Job) {
		h.Title = v
	}
}

func OptionJobDescription(v string) JobOption {
	return func(h *Job) {
		h.Description = v
	}
}

func OptionJobActions(v []Action) JobOption {
	return func(h *Job) {
		h.Actions = v
	}
}

func CreateJob(opts ...JobOption) *Job {
	job := &Job{}
	for _, opt := range opts {
		opt(job)
	}
	return job
}
