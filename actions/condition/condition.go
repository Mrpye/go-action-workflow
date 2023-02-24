//This is the workflow action for calling an api
package condition

import (
	"fmt"

	"github.com/Knetic/govaluate"
	"github.com/Mrpye/go-workflow/workflow"
	"github.com/Mrpye/golib/lib"
)

// CallApi is the main function for the action
func Action_Condition(w *workflow.Workflow, m *workflow.TemplateData) error {

	//*********************
	//Get the config values
	//*********************
	condition, err := w.GetConfigTokenString("condition", m, true)
	if err != nil {
		return err
	}

	//***********************
	//Evaluate the condition
	//***********************
	expression, err := govaluate.NewEvaluableExpression(condition)
	if err != nil {
		return err
	}
	result, err := expression.Evaluate(nil)
	if err != nil {
		return err
	}

	//***************************************
	//Check to see if the result is a boolean
	//***************************************
	_, ok := result.(bool)
	if !ok {
		return fmt.Errorf("must evaluate to a boolean")
	}

	pass, err := w.GetConfigTokenString("pass", m, false)
	if err != nil {
		return err
	}

	fail, err := w.GetConfigTokenString("fail", m, false)
	if err != nil {
		return err
	}

	if pass == "" && fail == "" {
		return fmt.Errorf("at lease one of pass or fail must be set")
	}

	//***************************
	//See what to do if it passes
	//***************************
	if result == true && pass != "" {
		if pass == "end" {
			return workflow.ErrEndWorkflow
		}
		index := w.GetCurrentJob().GetKeyIndex(pass)
		if index == -1 {
			return fmt.Errorf("cannot find label %s", pass)
		}
		err := w.SetCurrentActionIndex(index - 1)
		if err != nil {
			return err
		}
		if w.Verbose >= workflow.LOG_INFO {
			lib.PrintlnOK(fmt.Sprintf("condition passed going to action %s", pass))
		}
	} else if result == false && fail != "" {
		if pass == "end" {
			return workflow.ErrEndWorkflow
		}
		index := w.GetCurrentJob().GetKeyIndex(fail)
		if index == -1 {
			return fmt.Errorf("cannot find label %s", fail)
		}
		err := w.SetCurrentActionIndex(index - 1)
		if err != nil {
			return err
		}
		if w.Verbose >= workflow.LOG_INFO {
			lib.PrintlnFail(fmt.Sprintf("condition failed going to action %s", fail))

		}
	}

	return nil
}
