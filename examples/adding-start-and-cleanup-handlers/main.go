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

	//*************************************
	//Add the startup and cleanup functions
	//*************************************
	wf.InitFunc = Startup
	wf.CleanFunc = Clean

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
	err = wf.RunJob("adding-start-and-cleanup-handlers")
	if err != nil {
		println(err.Error())
	}

}

func AddHttp(val string) string {
	return fmt.Sprintf("http://%s", val)
}

func Startup(w *workflow.Workflow) error {
	//****************************************************
	//Save some values to the data bucket
	//You could read values from a config file or database
	//and save them to the data bucket
	//****************************************************
	w.SetValueToDataBucket("target", "host", "localhost")
	w.SetValueToDataBucket("target", "port", 8080)

	//**************************************************************
	// Add a custom function to the template function map
	//You don't necessarily have to do this in the startup function
	//You could do it in the main function
	//**************************************************************
	template_functions := w.GetTemplateFuncMap()
	template_functions["add_http"] = AddHttp

	//***********************************
	//Save it back to the workflow engine
	//***********************************
	w.SetTemplateFuncMap(template_functions)

	return nil
}

func Clean(w *workflow.Workflow) error {
	fmt.Println("Cleaning up")
	return nil
}
