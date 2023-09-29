package main

import (
	"github.com/Mrpye/go-action-workflow/actions/action_condition"
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
	wf.LogLevel = workflow.LOG_INFO
	wf.ActionList["condition"] = action_condition.Action_Condition

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
	err = wf.RunJob("condition-example", nil, nil)
	if err != nil {
		println(err.Error())
	}
}
