package workflow

import (
	"errors"
	"text/template"
)

// Log Level constants
const (
	LOG_QUIET   = 0 //just show print messages
	LOG_INFO    = 1 //show Action messages
	LOG_VERBOSE = 2 //show Action messages and LogLevel messages
)

var (
	//ErrEndWorkflow is the error that is returned when the workflow is ended
	//Use this error to end the workflow but not return an error to the caller
	ErrEndWorkflow = errors.New("end workflow")
)

// Workflow is the main struct for the workflow
type Workflow struct {
	LogLevel        int64
	Manifest        Manifest                  //the manifest that is used for the workflow
	InitFunc        EventFunc                 //the function that is called when the workflow is initialized
	CleanFunc       EventFunc                 //the function that is called when the workflow is cleaned up
	ActionList      map[string]ActionFunc     //the list of actions that are available
	ReadConfigFunc  map[string]ReadConfigFunc //the function that can be called to get config values for an action
	TargetMapFunc   TargetMapFunc             //Called to map config values to a target type
	templateFuncMap template.FuncMap          //the template function map that is used for the template engine

	stack         loopStack                         //the loop stack that stores the loop data when looping
	model         *TemplateData                     //the model that is used for the template engine
	dataBucket    map[string]map[string]interface{} //the data bucket that is used to store data between actions
	current_index int                               //the current index of the action list when running a job
	current_job   *Job                              //the current job that is running
	runtime_vars  map[string]interface{}
	base_path     string //stores the base path of the manifest file
}

// WorkflowOption is a function that sets a workflow option
type WorkflowOption func(*Workflow)

// OptionWorkflowLogLevel sets the LogLevel option
// 0 = quiet
// 1 = info
// 2 = LogLevel
func OptionWorkflowLogLevel(v int64) WorkflowOption {
	return func(h *Workflow) {
		h.LogLevel = v
	}
}

// OptionWorkflowManifest sets the manifest option
// - v: the manifest
func OptionWorkflowManifest(v Manifest) WorkflowOption {
	return func(h *Workflow) {
		h.Manifest = v
	}
}

// createWorkflow creates a new workflow
// - opts: the workflow options
// returns: the workflow
func CreateWorkflow(opts ...WorkflowOption) *Workflow {
	workflow := &Workflow{}
	workflow.ActionList = make(map[string]ActionFunc)
	workflow.ReadConfigFunc = make(map[string]ReadConfigFunc)
	workflow.dataBucket = make(map[string]map[string]interface{})

	//workflow.runtime_vars = make(map[string]interface{})
	for _, opt := range opts {
		opt(workflow)
	}
	return workflow
}

// UpdateWorkflow updates the workflow
// - opts: the workflow options
// returns: the workflow
func (m *Workflow) UpdateWorkflow(opts ...WorkflowOption) *Workflow {
	for _, opt := range opts {
		opt(m)
	}
	return m
}
