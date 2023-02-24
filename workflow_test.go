package tests

import (
	"testing"

	"github.com/Mrpye/go-action-workflow/actions/condition"
	"github.com/Mrpye/go-action-workflow/actions/parallel_workflow"
	"github.com/Mrpye/go-action-workflow/actions/store"
	"github.com/Mrpye/go-action-workflow/actions/sub_workflow"
	"github.com/Mrpye/go-action-workflow/actions/tests"
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
	wf.Verbose = workflow.LOG_QUIET

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
	err := wf.LoadManifest("examples/full-test-example/workflow.yaml")
	if err != nil {
		t.Error(err)
		return
	}

	//********************
	//Run the workflow job
	//********************
	err = wf.RunJob("test-example")
	if err != nil {
		t.Error(err)
		return
	}

	if err != nil {
		t.Error(err)
	}

}
