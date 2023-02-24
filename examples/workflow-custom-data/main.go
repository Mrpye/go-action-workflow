package main

import (
	"fmt"

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
	wf.Verbose = workflow.LOG_QUIET
	wf.ActionList["custom"] = Action_CustomData
	wf.ActionList["print"] = Action_Print

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
	err = wf.RunJob("custom-example")
	if err != nil {
		println(err.Error())
	}
}

func Action_Print(w *workflow.Workflow, m *workflow.TemplateData) error {
	msg, _ := w.GetConfigTokenString("msg", m, true)
	fmt.Println(msg)
	return nil
}

func Action_CustomData(w *workflow.Workflow, m *workflow.TemplateData) error {

	val := m.Manifest.DataModel().GetMapItem("data_test").ToBool()
	fmt.Println(val)
	//*******************
	//Get the custom data
	//*******************
	item := m.Manifest.DataModel().GetMapItem("items").GetArray()

	//**************
	//Create a model
	//**************
	model := w.CreateTemplateData(nil)
	//****************
	//Create an action
	//****************
	action := workflow.CreateAction()
	action.Action = "print"
	model.CurrentAction = action
	//**********************
	//Loop through the items
	//**********************
	for _, o := range item {
		//*****************
		//Update the config
		//*****************
		msg := o.GetMapItem("msg").ToString()
		action.Config = map[string]interface{}{
			"msg": msg,
		}
		//**************
		//Run the action
		//**************
		if val, ok := w.ActionList[action.Action]; ok {
			err := val(w, model)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
