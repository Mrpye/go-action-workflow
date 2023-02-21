package tests

import (
	"testing"

	"github.com/Mrpye/go-workflow/actions/store"
	"github.com/Mrpye/go-workflow/actions/tests"
	"github.com/Mrpye/go-workflow/workflow"
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
	wf.ActionList["ActionStore"] = store.ActionStore
	wf.ActionList["ActionTest"] = tests.ActionTest
	wf.ActionList["ActionFailTest"] = tests.ActionFailTest
	wf.ActionList["ActionJSAndMap"] = tests.ActionJSAndMap

	//*************************
	//load the workflow manifest
	//*************************
	err := wf.LoadManifest("../examples/full-test-example/workflow.yaml")
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

	if err != nil {
		t.Error(err)
	}

}
