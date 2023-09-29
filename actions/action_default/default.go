// This is the workflow action for calling an api
package action_default

import (
	"fmt"

	"github.com/Mrpye/go-action-workflow/workflow"
)

type DefaultSchema struct {
}

// GetTargetSchema returns the target workflow for this action
func (s DefaultSchema) GetTargetSchema() map[string]workflow.TargetSchema {
	//Build a test workflow
	return map[string]workflow.TargetSchema{}
}

// GetFunctionMap returns the function map for this workflow
func (s DefaultSchema) GetFunctionMap() map[string]workflow.FunctionSchema {
	return map[string]workflow.FunctionSchema{
		"image_name": {
			Cmd:         "image_name",
			Description: "gets the name of the image from the image path",
			Target:      "action_docker",
			ParameterSchema: map[string]*workflow.Schema{
				"image": {
					Type:        workflow.TypeString,
					Required:    true,
					Description: "docker.io/circleci/[slim-base]:latest",
				},
			},
		},
	}
}

// GetActions returns the actions for this workflow
func (s DefaultSchema) GetActionSchema() map[string]workflow.ActionSchema {

	//no short d f h
	return map[string]workflow.ActionSchema{
		"end": {
			Short: "End the workflow",
			Long:  "End the workflow",
		},
		"for": {
			Short:        "for loop",
			Long:         "for loop",
			InlineParams: true,
			InlineFormat: func(cfg map[string]interface{}) string {
				return fmt.Sprintf("for;%s;%s;%s", cfg["variable"], cfg["from"], cfg["to"]) //;v;x;z
			},
			ConfigSchema: map[string]*workflow.Schema{
				"variable": {
					Type:        workflow.TypeString,
					Partial:     false,
					Required:    true,
					Description: "Variable for the loop",
					Short:       "v",
					Default:     "",
				},
				"from": {
					Type:        workflow.TypeString,
					Partial:     false,
					Required:    true,
					Description: "From Value",
					Short:       "f",
					Default:     "",
				},
				"to": {
					Type:        workflow.TypeString,
					Partial:     false,
					Required:    true,
					Description: "To Value",
					Short:       "t",
					Default:     "",
				},
			},
		},
		"next": {
			Short: "next loop",
			Long:  "next loop",
		},
		"fail": {
			Short: "Fail the workflow",
			Long:  "Fail the workflow",
		},
		"goto": {
			Short:        "goto",
			Long:         "goto",
			InlineParams: true,
			InlineFormat: func(cfg map[string]interface{}) string {
				return fmt.Sprintf("goto;%s", cfg["key"]) //;v;x;z
			},
			ConfigSchema: map[string]*workflow.Schema{
				"key": {
					Type:        workflow.TypeString,
					Partial:     false,
					Required:    true,
					Description: "Action to goto",
					Short:       "i",
					Default:     "",
				},
			},
		},
		"print": {
			Short:        "Print a value",
			Long:         "Print a value",
			InlineParams: true,
			InlineFormat: func(cfg map[string]interface{}) string {
				return fmt.Sprintf("print;%s", cfg["message"]) //;v;x;z
			},
			ConfigSchema: map[string]*workflow.Schema{
				"message": {
					Type:        workflow.TypeString,
					Partial:     false,
					Required:    true,
					Description: "Message to print",
					Short:       "m",
					Default:     "",
				},
			},
		},
		"action": {
			Short:        "Runs a Global Action",
			Long:         "Runs a Global Action",
			InlineParams: true,
			InlineFormat: func(cfg map[string]interface{}) string {
				return fmt.Sprintf("action;%s", cfg["action_key"]) //;v;x;z
			},
			ConfigSchema: map[string]*workflow.Schema{
				"action_key": {
					Type:        workflow.TypeString,
					Partial:     false,
					Required:    true,
					Description: "Message to print",
					Short:       "a",
					Default:     "",
				},
			},
		},
	}
}

// GetAction returns the action for the target
func GetSchema() workflow.SchemaEndpoint {
	return DefaultSchema{}
}
