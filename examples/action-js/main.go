package main

import (
	"github.com/Mrpye/go-action-workflow/actions/action_js"
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
	wf.ActionList["js"] = action_js.Action_RunJS //add the action for calling APIs

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
	err = wf.RunJob("call-js-example", nil, nil)
	if err != nil {
		println(err.Error())
	}
}
