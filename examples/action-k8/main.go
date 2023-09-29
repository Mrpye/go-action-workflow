package main

import (
	"errors"
	"fmt"

	"github.com/Mrpye/go-action-workflow/actions/action_k8"
	"github.com/Mrpye/go-action-workflow/workflow"
	"github.com/spf13/viper"
)

func main() {
	//***********
	//Setup viper
	//***********
	viper.SetConfigFile("config.json")
	viper.ReadInConfig()

	//*****************
	//create a workflow
	//*****************
	wf := workflow.CreateWorkflow()

	//**********************************
	//Only show errors and print actions
	//**********************************
	wf.LogLevel = workflow.LOG_INFO

	//*****************
	//Add custom action
	//*****************
	wf.ActionList["k8_get_ws_items"] = action_k8.Action_K8GetWorkspace   //add the action for Storing data
	wf.ActionList["k8_get_service_ip"] = action_k8.Action_K8GetServiceIP //add the action for Storing data
	wf.ActionList["k8_pod_exec"] = action_k8.Action_K8PodExec            //add the action for Storing data
	wf.ActionList["k8_wait"] = action_k8.Action_K8WaitCompleteStatus     //add the action for Storing data
	//*************************
	//Setup the config function
	//*************************
	wf.ReadConfigFunc["viper"] = ReadViperConfig

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
	//create the runtime vars
	runtime_vars := make(map[string]interface{})
	runtime_vars["env"] = "dev"
	err = wf.RunJob("k8-example", runtime_vars, nil)
	if err != nil {
		println(err.Error())
	}
}

func ReadViperConfig(key string, data_type string, custom ...string) (interface{}, error) {
	if key == "" {
		return nil, errors.New("key is empty")
	}

	//****************************************************************************
	//custom is used to pass other string data to the config function
	//this is useful if you want to pass a config file name or something like that
	//****************************************************************************
	/*for _, v := range custom {
		println(v)
	}*/

	key = fmt.Sprintf("%s.%s.%s", custom[0], custom[1], key)

	switch data_type {
	case "string":
		return viper.GetString(key), nil
	case "int":
		return viper.GetInt(key), nil
	case "bool":
		return viper.GetBool(key), nil
	case "float64":
		return viper.GetFloat64(key), nil
	default:
		return viper.GetString(key), nil
	}
}
