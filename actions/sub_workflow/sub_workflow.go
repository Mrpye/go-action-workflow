package sub_workflow

import (
	"encoding/json"

	"github.com/Mrpye/go-action-workflow/workflow"
	"github.com/Mrpye/golib/lib"
)

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

	if w.Verbose >= workflow.LOG_VERBOSE {
		bs, _ := json.MarshalIndent(inputs, "", "  ")
		lib.LogVerbose("Running Sub Workflow:" + job + " with inputs:\n " + string(bs))
	}
	//Map the values to the workflow inputs
	wf := w.CreateSubWorkflowEngine()
	err = wf.RunSubWorkflow(job, inputs)
	return err
}
