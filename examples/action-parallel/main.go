package main

import (
	"github.com/Mrpye/go-action-workflow/actions/action_api"
	"github.com/Mrpye/go-action-workflow/actions/action_condition"
	"github.com/Mrpye/go-action-workflow/actions/action_parallel_workflow"
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
	wf.ActionList["parallel"] = action_parallel_workflow.Action_Parallel
	wf.ActionList["sub-workflow"] = action_sub_workflow.Action_SubWorkflow
	wf.ActionList["store"] = action_store.Action_Store
	wf.ActionList["api"] = action_api.Action_CallApi               //add the action for calling APIs
	wf.ActionList["condition"] = action_condition.Action_Condition //add the action for calling APIs

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
	err = wf.RunJob("parallel-example", nil, nil)
	if err != nil {
		println(err.Error())
	}
}
