package workflow

import (
	"errors"
	"text/template"
)

// Log Level constants
const (
	LOG_QUIET   = 0
	LOG_INFO    = 1
	LOG_VERBOSE = 2
)

var (
	ErrEndWorkflow = errors.New("end workflow")
)

//Workflow is the main struct for the workflow
type Workflow struct {
	Verbose         int
	Manifest        Manifest
	InitFunc        EventFunc
	CleanFunc       EventFunc
	templateFuncMap template.FuncMap
	stack           loopStack
	model           *TemplateData
	ActionList      map[string]ActionFunc
	dataBucket      map[string]map[string]interface{}
	current_index   int
	current_job     *Job
}

//ActionFunc is the function that is called for each action
type ActionFunc func(*Workflow, *TemplateData) error

type EventFunc func(*Workflow) error

//WorkflowOption is a function that sets a workflow option
type WorkflowOption func(*Workflow)

// OptionWorkflowVerbose sets the verbose option
func OptionWorkflowVerbose(v int) WorkflowOption {
	return func(h *Workflow) {
		h.Verbose = v
	}
}

// OptionWorkflowManifest sets the manifest option
func OptionWorkflowManifest(v Manifest) WorkflowOption {
	return func(h *Workflow) {
		h.Manifest = v
	}
}

// createWorkflow creates a new workflow
func CreateWorkflow(opts ...WorkflowOption) *Workflow {
	workflow := &Workflow{}
	workflow.ActionList = make(map[string]ActionFunc)
	workflow.dataBucket = make(map[string]map[string]interface{})
	for _, opt := range opts {
		opt(workflow)
	}
	return workflow
}

func (m *Workflow) UpdateWorkflow(opts ...WorkflowOption) *Workflow {
	for _, opt := range opts {
		opt(m)
	}
	return m
}
