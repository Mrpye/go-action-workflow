package main

import (
	"github.com/Mrpye/go-action-workflow/actions/condition"
	"github.com/Mrpye/go-action-workflow/actions/parallel_workflow"
	"github.com/Mrpye/go-action-workflow/actions/store"
	"github.com/Mrpye/go-action-workflow/actions/sub_workflow"
	"github.com/Mrpye/go-action-workflow/actions/tests"
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
	wf.Verbose = workflow.LOG_INFO

	//*******************
	//Add a custom action
	//*******************
	wf.ActionList["ActionStore"] = store.Action_Store
	wf.ActionList["ActionTest"] = tests.ActionTest
	wf.ActionList["ActionFailTest"] = tests.ActionFailTest
	wf.ActionList["ActionJSAndMap"] = tests.ActionJSAndMap
	wf.ActionList["parallel"] = parallel_workflow.Action_Parallel
	wf.ActionList["sub-workflow"] = sub_workflow.Action_SubWorkflow
	wf.ActionList["store"] = store.Action_Store
	wf.ActionList["condition"] = condition.Action_Condition //add the action for calling APIs

	//*************************
	//load the workflow manifest
	//*************************
	err := wf.LoadManifest("./workflow.yaml")
	if err != nil {
		println(err.Error())
		return
	}

	//********************
	//Run the workflow job
	//********************
	err = wf.RunJob("test-example")
	if err != nil {
		println(err.Error())
		return
	}

	println("Test Passed")

}
