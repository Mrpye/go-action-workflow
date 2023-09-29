package main

import (
	"github.com/Mrpye/go-action-workflow/actions/action_api"
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
	wf.ActionList["api"] = action_api.Action_CallApi //add the action for calling APIs

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
	err = wf.RunJob("call-api-example", nil, nil)
	if err != nil {
		println(err.Error())
	}
}
