package action_system

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/Mrpye/go-action-workflow/workflow"
	"github.com/Mrpye/golib/str"
)

type SystemSchema struct {
}

/*
Cmd             string
	Description     string
	Target          string //The module that the function is for e.g action_git
	ParameterSchema map[string]*Schema
	Function        any //The function that is called to execute the function
*/

// GetFunctionMap returns the function map for this schema
func (s SystemSchema) GetFunctionMap() map[string]workflow.FunctionSchema {
	return map[string]workflow.FunctionSchema{
		"get_env": {
			Cmd:         "image_name",
			Description: "get the environment variable",
			Target:      "action_system",
			Function:    GetEnv,
			ParameterSchema: map[string]*workflow.Schema{
				"key": {
					Type:        workflow.TypeString,
					Required:    true,
					Description: "The environment variable key",
				},
			},
		},
	}
}

// GetTargetSchema returns the target schema for this action
func (s SystemSchema) GetTargetSchema() map[string]workflow.TargetSchema {
	//Build a test schema
	return map[string]workflow.TargetSchema{}
}

// GetActions returns the actions for this schema
func (s SystemSchema) GetActionSchema() map[string]workflow.ActionSchema {

	//no short d f h
	return map[string]workflow.ActionSchema{
		"sub_workflow": {
			Action:         Action_SetEnv,
			Short:          "Set Environment Variables",
			Long:           "Set Environment Variables",
			ProcessResults: false,
			ConfigSchema: map[string]*workflow.Schema{
				"env_vars": {
					Type:        workflow.TypeMap,
					Partial:     false,
					Required:    true,
					Description: "Map of environment variables to set",
					Short:       "e",
					Default:     []string{},
				},
			},
		},
	}
}

// GetAction returns the action for the target
func GetSchema() workflow.SchemaEndpoint {
	return SystemSchema{}
}

func Action_SetEnv(w *workflow.Workflow, m *workflow.TemplateData) error {

	//**********************************
	//Get a string value from the config
	//**********************************
	envs, err := w.GetConfigTokenMap("env_vars", m, true)
	if err != nil {
		return err
	}

	//*********************
	//loop through the vars
	//*********************
	for k, val := range envs {
		processed_value, err := w.ParseToken(m, val.(string))
		if err != nil {
			return err
		}
		key := strings.ToUpper(str.Clean(k, "_"))
		if key == "" {
			return errors.New("missing key from set env")
		}
		os.Setenv(key, processed_value)
		if w.LogLevel >= workflow.LOG_INFO {
			fmt.Printf("Env Var set %s=%s", key, processed_value)
		}
	}
	return err
}
