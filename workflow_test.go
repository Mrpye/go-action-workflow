package tests

import (
	"testing"

	"github.com/Mrpye/go-action-workflow/actions/action_condition"
	"github.com/Mrpye/go-action-workflow/actions/action_parallel_workflow"
	"github.com/Mrpye/go-action-workflow/actions/action_store"
	"github.com/Mrpye/go-action-workflow/actions/action_sub_workflow"
	"github.com/Mrpye/go-action-workflow/actions/action_tests"
	"github.com/Mrpye/go-action-workflow/workflow"
)

func TestWorkflow(t *testing.T) {

	//*****************
	//create a workflow
	//*****************
	wf := workflow.CreateWorkflow()

	//**********************************
	//Only show errors and print actions
	//**********************************
	wf.LogLevel = workflow.LOG_QUIET

	//*******************
	//Add a custom action
	//*******************
	wf.ActionList["ActionStore"] = action_store.Action_Store
	wf.ActionList["ActionTest"] = action_tests.ActionTest
	wf.ActionList["ActionFailTest"] = action_tests.ActionFailTest
	wf.ActionList["ActionJSAndMap"] = action_tests.ActionJSAndMap
	wf.ActionList["parallel"] = action_parallel_workflow.Action_Parallel
	wf.ActionList["sub-workflow"] = action_sub_workflow.Action_SubWorkflow
	wf.ActionList["store"] = action_store.Action_Store
	wf.ActionList["condition"] = action_condition.Action_Condition //add the action for calling APIs

	//*************************
	//load the workflow manifest
	//*************************
	err := wf.LoadManifest("examples/full-test-example/workflow.yaml")
	if err != nil {
		t.Error(err)
		return
	}

	//********************
	//Run the workflow job
	//********************
	err = wf.RunJob("test-example", nil, nil)
	if err != nil {
		t.Error(err)
		return
	}

	if err != nil {
		t.Error(err)
	}

}
