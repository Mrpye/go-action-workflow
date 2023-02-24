package main

import (
	"encoding/json"
	"fmt"
	"log"

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

	//*******************
	//Add a custom action
	//*******************
	wf.ActionList["MultiPrint"] = MultiPrint

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
	err = wf.RunJob("add-custom-actions-example")
	if err != nil {
		println(err.Error())
	}

}

//**************************
//print will print a message
//**************************
func MultiPrint(w *workflow.Workflow, m *workflow.TemplateData) error {

	//**********************************
	//Get a string value from the config
	//**********************************
	string_value, err := w.GetConfigTokenString("string_value", m, true)
	if err != nil {
		return err
	}
	//*******************************
	//Get a int value from the config
	//*******************************
	int_value, err := w.GetConfigTokenInt("int_value", m, true)
	if err != nil {
		return err
	}
	//********************************
	//Get a bool value from the config
	//********************************
	bool_value, err := w.GetConfigTokenBool("bool_value", m, true)
	if err != nil {
		return err
	}
	//*******************************
	//Get a map value from the config
	//*******************************
	map_value, err := w.GetConfigTokenMap("map_value", m, true)
	if err != nil {
		return err
	}

	//****************
	//Print the values
	//****************
	println(string_value)
	println(int_value)
	println(bool_value)
	b, err := json.Marshal(map_value)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(b))

	return nil
}
