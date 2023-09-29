package main

import (
	"github.com/Mrpye/go-action-workflow/actions/action_file"
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
	wf.ActionList["copy"] = action_file.Action_Copy     //add the action for calling APIs
	wf.ActionList["delete"] = action_file.Action_Delete //add the action for calling APIs
	wf.ActionList["rename"] = action_file.Action_Rename //add the action for calling APIs
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
