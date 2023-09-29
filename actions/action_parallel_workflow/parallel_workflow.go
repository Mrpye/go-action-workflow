package action_parallel_workflow

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/Mrpye/go-action-workflow/workflow"
	"github.com/Mrpye/golib/log"
)

type ParallelSchema struct {
}

// GetFunctionMap returns the function map for this schema
func (s ParallelSchema) GetFunctionMap() map[string]workflow.FunctionSchema {
	return map[string]workflow.FunctionSchema{}
}

// GetTargetSchema returns the target schema for this action
func (s ParallelSchema) GetTargetSchema() map[string]workflow.TargetSchema {
	//Build a test schema
	return map[string]workflow.TargetSchema{}
}

// GetActions returns the actions for this schema
func (s ParallelSchema) GetActionSchema() map[string]workflow.ActionSchema {

	//no short d f h
	return map[string]workflow.ActionSchema{
		"parallel": {
			Action:         Action_Parallel,
			Short:          "Store a value in the data bucket",
			Long:           "Store a value in the data bucket",
			ProcessResults: false,
			ConfigSchema: map[string]*workflow.Schema{
				"actions": {
					Type:        workflow.TypeList,
					Partial:     false,
					Required:    true,
					Description: "A list of actions to run in parallel",
					Short:       "a",
					Default:     []string{},
				},
			},
		},
	}
}

// GetAction returns the action for the target
func GetSchema() workflow.SchemaEndpoint {
	return ParallelSchema{}
}

func Action_Parallel(w *workflow.Workflow, m *workflow.TemplateData) error {

	//**********************************
	//Get a string value from the config
	//**********************************
	inputs, err := w.GetConfigTokenInterface("actions", m, true)
	if err != nil {
		return err
	}

	//*********************************
	//Map the actions to the job struct
	//*********************************
	var actions []workflow.Action
	bs, _ := json.MarshalIndent(inputs, "", "  ")
	err = json.Unmarshal(bs, &actions)
	if err != nil {
		return err
	}

	//*********************
	// Setup the wait group
	// And error store
	//*********************
	var wg sync.WaitGroup
	wg.Add(len(actions))
	err_store := make([]error, len(actions))

	//****************
	// Run the actions
	//****************
	for i := range actions {
		go func(i int) {
			defer wg.Done()
			err_store[i] = w.RunAction(&actions[i])
			if w.LogLevel > workflow.LOG_INFO {
				log.LogVerbose(fmt.Sprintf("Action `%s` completed with error: %v", actions[i].Action, err_store[i]))
			}
		}(i)
	}
	//***************************
	//Wait for then all to finish
	//***************************
	wg.Wait()

	//*******************************
	//Check for errors
	//And combine them into one error
	//*******************************
	err = nil
	for i := range err_store {
		if err_store[i] != nil {
			if err == nil {
				err = fmt.Errorf("action:`%s` error %s", actions[i].Action, err_store[i].Error())
			} else {
				err = fmt.Errorf("%w; action:`%s` error %s;", err, actions[i].Action, err_store[i].Error())
			}
		}
	}

	return err
}
