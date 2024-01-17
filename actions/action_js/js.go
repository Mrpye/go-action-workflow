package action_js

import (
	"strings"

	"github.com/Mrpye/go-action-workflow/workflow"
	"github.com/Mrpye/golib/file"
)

type JSSchema struct {
}

// GetTargetSchema returns the target schema for this action
func (s JSSchema) GetTargetSchema() map[string]workflow.TargetSchema {
	//Build a test schema
	return map[string]workflow.TargetSchema{}
}

// GetFunctionMap returns the function map for this schema
func (s JSSchema) GetFunctionMap() map[string]workflow.FunctionSchema {
	return map[string]workflow.FunctionSchema{}
}

// GetActions returns the actions for this schema
func (s JSSchema) GetActionSchema() map[string]workflow.ActionSchema {

	//no short d f h
	return map[string]workflow.ActionSchema{
		"js": {
			Action:         Action_RunJS,
			Short:          "Run a javascript file",
			Long:           "Run a javascript file",
			ProcessResults: true,
			ConfigSchema: map[string]*workflow.Schema{
				"js": {
					Type:        workflow.TypeString,
					Partial:     false,
					Required:    false,
					Description: "javascript code to run",
					Short:       "j",
					Default:     "",
				},
				"js_file": {
					Type:        workflow.TypeString,
					Partial:     false,
					Required:    true,
					Description: "javascript file to run",
					Short:       "f",
					Default:     "",
				},
			},
		},
	}
}

// GetAction returns the action for the target
func GetSchema() workflow.SchemaEndpoint {
	return JSSchema{}
}

func Action_RunJS(w *workflow.Workflow, m *workflow.TemplateData) error {

	code, err := w.GetConfigTokenString("js", m, false)
	if err != nil {
		return err
	}

	if code == "" {
		//*********************
		//Get the config values
		//*********************
		js_file, err := w.GetConfigTokenString("js_file", m, true)
		if err != nil {
			return err
		}
		js_file = w.BuildPath(js_file)

		//*****************************
		//See if we have multiple files
		//*****************************
		files := strings.Split(js_file, ";")

		for _, o := range files {
			file_part := o
			file_data, err := file.ReadFileToString(file_part)
			if err != nil {
				return err
			}
			code = code + file_data + "\n"
		}
	}
	//************
	//Run our code
	//************
	vm := w.CreateJSEngine()
	vm.Set("model", m)
	_, err = vm.RunString(code)
	if err != nil {
		return err
	}
	return nil
}
