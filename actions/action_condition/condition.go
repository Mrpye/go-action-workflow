// This is the workflow action for calling an api
package action_condition

import (
	"fmt"

	"github.com/Knetic/govaluate"
	"github.com/Mrpye/go-action-workflow/workflow"
	"github.com/Mrpye/golib/log"
)

type ConditionSchema struct {
}

// GetFunctionMap returns the function map for this schema
func (s ConditionSchema) GetFunctionMap() map[string]workflow.FunctionSchema {
	return map[string]workflow.FunctionSchema{}
}

// GetTargetSchema returns the target schema for this action
func (s ConditionSchema) GetTargetSchema() map[string]workflow.TargetSchema {
	//Build a test schema
	return map[string]workflow.TargetSchema{}
}

// GetActions returns the actions for this schema
func (s ConditionSchema) GetActionSchema() map[string]workflow.ActionSchema {

	//no short d f h
	return map[string]workflow.ActionSchema{
		"condition": {
			Action:         Action_Condition,
			Short:          "Run an action based on a condition",
			Long:           "Run an action based on a condition",
			ProcessResults: true,
			ConfigSchema: map[string]*workflow.Schema{
				"condition": {
					Type:        workflow.TypeString,
					Partial:     false,
					Required:    true,
					Description: "The condition to evaluate must evaluate to a boolean",
					Short:       "c",
					Default:     "",
				},
				"pass": {
					Type:        workflow.TypeString,
					Partial:     false,
					Required:    true,
					Description: "Action or Action Key to run if the condition passes",
					Short:       "p",
					Default:     "",
				},
				"fail": {
					Type:        workflow.TypeString,
					Partial:     false,
					Required:    true,
					Description: "Action or Action Key to run if the condition failed",
					Short:       "f",
					Default:     "",
				},
			},
		},
	}
}

// GetAction returns the action for the target
func GetSchema() workflow.SchemaEndpoint {
	return ConditionSchema{}
}

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
		if w.LogLevel >= workflow.LOG_INFO {
			log.PrintlnOK(fmt.Sprintf("condition passed going to action %s", pass))
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
		if w.LogLevel >= workflow.LOG_INFO {
			log.PrintlnFail(fmt.Sprintf("condition failed going to action %s", fail))

		}
	}

	return nil
}
