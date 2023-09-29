package main

import (
	"github.com/Mrpye/go-action-workflow/actions/action_store"
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

	//*****************
	//Add custom action
	//*****************
	wf.ActionList["store"] = action_store.Action_Store //add the action for Storing data

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
	err = wf.RunJob("store-example", nil, nil)
	if err != nil {
		println(err.Error())
	}
}
