package main

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/Mrpye/go-action-workflow/actions/action_file"
	"github.com/Mrpye/go-action-workflow/actions/action_git"
	"github.com/Mrpye/go-action-workflow/workflow"
	"github.com/Mrpye/golib/log"
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
	wf.ActionList["git_download"] = action_git.Action_GitDownload
	wf.ActionList["file_delete"] = action_file.Action_Delete

	wf.TargetMapFunc = MapConfigToTarget

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
	err = wf.RunJob("example", map[string]interface{}{"env": "dev"}, nil)
	if err != nil {
		println(err.Error())
	}
}

func MapConfigToTarget(w *workflow.Workflow, m interface{}, target interface{}) (interface{}, error) {
	//*****************************
	//Get the env from runtime_vars
	//*****************************
	env, err := w.GetRuntimeVar("env")
	if err != nil {
		return nil, err
	}

	env = env.(string)
	if w.LogLevel == workflow.LOG_VERBOSE {
		log.LogVerbose(fmt.Sprintf("env: %s\n", env))
	}

	//*******************
	//Get the target_name
	//*******************
	target_name, err := w.GetConfigTokenString("target_name", m.(*workflow.TemplateData), false)
	if err != nil {
		return nil, err
	}
	type_name := strings.ReplaceAll(strings.ToLower(reflect.TypeOf(target).String()), "*", "")
	if target_name == "" {
		target_name = type_name
	} else {
		target_name = type_name + "_" + target_name
	}
	if w.LogLevel == workflow.LOG_VERBOSE {
		log.LogVerbose(fmt.Sprintf("Target name: %s\n", target_name))
	}

	//**********************************************
	//Use reflection to map the config to the target
	//**********************************************
	v := reflect.Indirect(reflect.ValueOf(target))
	typeOfS := v.Type()
	for i := 0; i < v.NumField(); i++ {
		c := typeOfS.Field(i).Tag
		g := c.Get("yaml")
		if g == "" {
			continue
		}
		//******************************************
		//Remap the config value to the target filed
		//******************************************
		if v.Field(i).Kind() == reflect.Int {
			temp_config_val, _ := w.GetConfigValue("", target_name+"."+g, "int", "targets", env.(string))
			if w.LogLevel == workflow.LOG_VERBOSE {
				log.LogVerbose(fmt.Sprintf("%s: %s\n", g, temp_config_val))
			}
			v.Field(i).SetInt(temp_config_val.(int64))
		} else if v.Field(i).Kind() == reflect.String {
			temp_config_val, _ := w.GetConfigValue("", target_name+"."+g, "string", "targets", env.(string))
			if w.LogLevel == workflow.LOG_VERBOSE {
				log.LogVerbose(fmt.Sprintf("%s: %v\n", g, temp_config_val))
			}
			v.Field(i).SetString(temp_config_val.(string))
		} else if v.Field(i).Kind() == reflect.Bool {
			temp_config_val, _ := w.GetConfigValue("", target_name+"."+g, "bool", "targets", env.(string))
			if w.LogLevel == workflow.LOG_VERBOSE {
				log.LogVerbose(fmt.Sprintf("%s: %v\n", g, temp_config_val))
			}
			v.Field(i).SetBool(temp_config_val.(bool))
		} else if v.Field(i).Kind() == reflect.Float64 {
			temp_config_val, _ := w.GetConfigValue("", target_name+"."+g, "float", "targets", env.(string))
			if w.LogLevel == workflow.LOG_VERBOSE {
				log.LogVerbose(fmt.Sprintf("%s: %v\n", g, temp_config_val))
			}
			v.Field(i).SetFloat(temp_config_val.(float64))
		}
	}

	return target, nil
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
