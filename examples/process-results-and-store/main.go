package main

import (
	"encoding/json"

	"github.com/Mrpye/go-action-workflow/workflow"
)

func main() {
	//*****************
	//create a workflow
	//*****************
	wf := workflow.CreateWorkflow()

	//**********************************
	//Only show errors and print actions
	//**********************************
	wf.LogLevel = workflow.LOG_QUIET

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
	err = wf.RunJob("handling-result-data-storing-results", nil, nil)
	if err != nil {
		println(err.Error())
	}

}

// **************************
// print will print a message
// **************************
func MultiPrint(w *workflow.Workflow, m *workflow.TemplateData) error {
	//**********************************************
	//Get the model if m is passed then its parallel
	//**********************************************

	//*******************************
	//Get a map value from the config
	//*******************************
	map_value, err := w.GetConfigTokenMap("map_value", m, true)
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
	err = w.ActionProcessResults(m, string(b))
	if err != nil {
		return err
	}
	return nil
}
