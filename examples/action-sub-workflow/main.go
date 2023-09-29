package main

import (
	"github.com/Mrpye/go-action-workflow/actions/action_store"
	"github.com/Mrpye/go-action-workflow/actions/action_sub_workflow"
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
	wf.ActionList["sub-workflow"] = action_sub_workflow.Action_SubWorkflow
	wf.ActionList["store"] = action_store.Action_Store

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
	err = wf.RunJob("main-workflow-example", nil, nil)
	if err != nil {
		println(err.Error())
	}
}
