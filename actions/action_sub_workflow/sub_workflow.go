package action_sub_workflow

import (
	"encoding/json"

	"github.com/Mrpye/go-action-workflow/workflow"
	"github.com/Mrpye/golib/log"
)

type SubWorkflowSchema struct {
}

// GetFunctionMap returns the function map for this schema
func (s SubWorkflowSchema) GetFunctionMap() map[string]workflow.FunctionSchema {
	return map[string]workflow.FunctionSchema{}
}

// GetTargetSchema returns the target schema for this action
func (s SubWorkflowSchema) GetTargetSchema() map[string]workflow.TargetSchema {
	//Build a test schema
	return map[string]workflow.TargetSchema{}
}

// GetActions returns the actions for this schema
func (s SubWorkflowSchema) GetActionSchema() map[string]workflow.ActionSchema {

	//no short d f h
	return map[string]workflow.ActionSchema{
		"sub_workflow": {
			Action:         Action_SubWorkflow,
			Short:          "Run a sub workflow",
			Long:           "Run a sub workflow",
			ProcessResults: false,
			ConfigSchema: map[string]*workflow.Schema{
				"inputs": {
					Type:        workflow.TypeMap,
					Partial:     false,
					Required:    false,
					Description: "The inputs for the sub workflow",
					ConfigKey:   "inputs",
					Short:       "i",
					Default: []string{
						"test=test",
					},
				},
				"job": {
					Type:        workflow.TypeString,
					Partial:     false,
					Required:    true,
					Description: "The name of the sub workflow to run",
					Short:       "j",
					Default:     "",
				},
			},
		},
	}
}

// GetAction returns the action for the target
func GetSchema() workflow.SchemaEndpoint {
	return SubWorkflowSchema{}
}

func Action_SubWorkflow(w *workflow.Workflow, m *workflow.TemplateData) error {

	//**********************************
	//Get a string value from the config
	//**********************************
	inputs, err := w.GetConfigTokenMap("inputs", m, true)
	if err != nil {
		return err
	}
	job, err := w.GetConfigTokenString("job", m, true)
	if err != nil {
		return err
	}

	if w.LogLevel >= workflow.LOG_VERBOSE {
		bs, _ := json.MarshalIndent(inputs, "", "  ")
		log.LogVerbose("Running Sub Workflow:" + job + " with inputs:\n " + string(bs))
	}
	//Map the values to the workflow inputs
	wf := w.CreateSubWorkflowEngine()
	err = wf.RunSubWorkflow(job, inputs)
	return err
}
