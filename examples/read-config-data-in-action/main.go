package main

import (
	"errors"
	"strings"

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
	wf.ActionList["config"] = Action_ReadConfig //add the action for Storing data

	//*************************
	//Setup the config function
	//*************************
	wf.ReadConfigFunc["viper"] = ReadViperConfig
	wf.ReadConfigFunc["other"] = ReadOtherConfig

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
	err = wf.RunJob("read-config-example", nil, nil)
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
	for _, v := range custom {
		println(v)
	}

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

func ReadOtherConfig(key string, data_type string, custom ...string) (interface{}, error) {
	if key == "" {
		return nil, errors.New("key is empty")
	}

	//****************************************************************************
	//custom is used to pass other string data to the config function
	//this is useful if you want to pass a config file name or something like that
	//****************************************************************************
	for _, v := range custom {
		println(v)
	}

	switch strings.ToUpper(key) {
	case "A":
		return "config data 1", nil
	case "B":
		return "config data 2", nil
	case "C":
		return 12345, nil
	case "D":
		return true, nil
	default:
		return "default", nil
	}
}

func Action_ReadConfig(w *workflow.Workflow, m *workflow.TemplateData) error {

	value, err := w.GetConfigValue("viper", "targets.git.host", "string", "custom data1", "custom data2", "custom data3")
	if err != nil {
		return err
	}
	println(value.(string))

	value, err = w.GetConfigValue("other", "A", "", "custom data4", "custom data5", "custom data6")
	if err != nil {
		return err
	}
	println(value.(string))

	value, err = w.GetConfigValue("other", "C", "")
	if err != nil {
		return err
	}
	println(value.(int))

	return nil
}
