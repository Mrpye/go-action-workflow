package parallel_workflow

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/Mrpye/go-workflow/workflow"
	"github.com/Mrpye/golib/lib"
)

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
			if w.Verbose > workflow.LOG_INFO {
				lib.LogVerbose(fmt.Sprintf("Action `%s` completed with error: %v", actions[i].Action, err_store[i]))
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
