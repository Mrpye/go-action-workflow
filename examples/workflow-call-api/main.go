package main

import (
	"github.com/Mrpye/go-workflow/actions/api"
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

	//*****************
	//Add custom action
	//*****************
	wf.ActionList["api"] = api.CallApi //add the action for calling APIs

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
	err = wf.RunJob("simple-example")
	if err != nil {
		println(err.Error())
	}
}
