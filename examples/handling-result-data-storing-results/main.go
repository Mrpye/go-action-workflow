package main

import (
	"encoding/json"

	"github.com/Mrpye/go-workflow/workflow"
)

func main() {
	//*****************
	//create a workflow
	//*****************
	wf := workflow.CreateWorkflow()

	//**********************************
	//Only show errors and print actions
	//**********************************
	wf.Verbose = workflow.LOG_QUIET

	//*******************
	//Add a custom action
	//*******************
	wf.ActionList["MultiPrint"] = MultiPrint

	//*************************
	//load the workflow manifest
	//*************************
	err := wf.LoadManifest("./workflow.yaml")
	if err != nil {
		println(err.Error())
	}

	//********************
	//Run the workflow job
	//********************
	err = wf.RunJob("handling-result-data-storing-results")
	if err != nil {
		println(err.Error())
	}

}

//**************************
//print will print a message
//**************************
func MultiPrint(w *workflow.Workflow) error {
	//*******************************
	//Get a map value from the config
	//*******************************
	map_value, err := w.GetConfigTokenMap("map_value", w.Model, true)
	if err != nil {
		return err
	}

	//***************************************
	//Convert to json and process the results
	//***************************************
	b, err := json.Marshal(map_value)
	if err != nil {
		return err
	}

	//***********************************
	//This function processes the results
	//***********************************
	err = w.ActionProcessResults(string(b))
	if err != nil {
		return err
	}
	return nil
}
