package main

import (
	"github.com/Mrpye/go-workflow/actions/api"
	"github.com/Mrpye/go-workflow/actions/condition"
	"github.com/Mrpye/go-workflow/actions/parallel_workflow"
	"github.com/Mrpye/go-workflow/actions/store"
	"github.com/Mrpye/go-workflow/actions/sub_workflow"
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
	wf.ActionList["parallel"] = parallel_workflow.Action_Parallel
	wf.ActionList["sub-workflow"] = sub_workflow.Action_SubWorkflow
	wf.ActionList["store"] = store.Action_Store
	wf.ActionList["api"] = api.Action_CallApi               //add the action for calling APIs
	wf.ActionList["condition"] = condition.Action_Condition //add the action for calling APIs

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
	err = wf.RunJob("parallel-example")
	if err != nil {
		println(err.Error())
	}
}
