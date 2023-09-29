package workflow

// Job is a workflow job
type Job struct {
	Key           string               `json:"key" yaml:"key" flag:"key k" desc:"key for the job, must be unique within the workflow, used as the job name"`
	Title         string               `json:"title" yaml:"title" flag:"title t" desc:"title for the job"`
	Description   string               `json:"description" yaml:"description" flag:"desc d" desc:"description for the job"`
	Actions       []Action             `json:"actions" yaml:"actions" flag:"actions a" desc:"list of actions for the job"`
	IsSubWorkflow bool                 `json:"is_sub_workflow" yaml:"is_sub_workflow" flag:"sub s" desc:"set to true if this is a sub workflow"`
	Inputs        map[string]Parameter `json:"inputs,omitempty" yaml:"inputs,omitempty" flag:"inputs i" desc:"inputs for the job"`
}

// JobOption is a function that sets an option on the job
type JobOption func(*Job)

// OptionJobKey sets the key on the job
// v is the key
func OptionJobKey(v string) JobOption {
	return func(h *Job) {
		h.Key = v
	}
}

// OptionJobTitle sets the title on the job
// v is the title
func OptionJobTitle(v string) JobOption {
	return func(h *Job) {
		h.Title = v
	}
}

//	OptionJobDescription sets the description on the job
// v is the description
func OptionJobDescription(v string) JobOption {
	return func(h *Job) {
		h.Description = v
	}
}

// OptionJobActions sets the actions on the job
// v is the actions as a slice of Action
func OptionJobActions(v []Action) JobOption {
	return func(h *Job) {
		h.Actions = v
	}
}

// CreateJob creates a new job
// opts are the options to set on the job
// returns the job
func CreateJob(opts ...JobOption) *Job {
	job := &Job{}
	for _, opt := range opts {
		opt(job)
	}
	return job
}
