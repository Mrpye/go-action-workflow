package action_exec

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/Mrpye/go-action-workflow/workflow"
)

type ExecSchema struct {
}

// GetFunctionMap returns the function map for this schema
func (s ExecSchema) GetFunctionMap() map[string]workflow.FunctionSchema {
	return map[string]workflow.FunctionSchema{}
}

// GetTargetSchema returns the target schema for this action
func (s ExecSchema) GetTargetSchema() map[string]workflow.TargetSchema {
	//Build a test schema
	return map[string]workflow.TargetSchema{}
}

// GetActions returns the actions for this schema
func (s ExecSchema) GetActionSchema() map[string]workflow.ActionSchema {

	//no short d f h
	return map[string]workflow.ActionSchema{
		"exec": {
			Action:         Action_Exec,
			Short:          "Store a value in the data bucket",
			Long:           "Store a value in the data bucket",
			ProcessResults: false,
			ConfigSchema: map[string]*workflow.Schema{
				"app": {
					Type:        workflow.TypeString,
					Partial:     false,
					Required:    true,
					Description: "The application to execute",
					Short:       "a",
					Default:     "",
				},
				"args": {
					Type:        workflow.TypeString,
					Partial:     false,
					Required:    true,
					Description: "the arguments to pass",
					Short:       "i",
					Default:     "",
				},
			},
		},
	}
}

// GetAction returns the action for the target
func GetSchema() workflow.SchemaEndpoint {
	return ExecSchema{}
}

func Action_Exec(w *workflow.Workflow, m *workflow.TemplateData) error {

	//**********************************
	//Get a string value from the config
	//**********************************
	app_path, err := w.GetConfigTokenString("app", m, true)
	if err != nil {
		return err
	}
	app_path = w.BuildPath(app_path)

	//********************************
	//Get a args value from the config
	//********************************
	command, err := w.GetConfigTokenString("args", m, true)
	if err != nil {
		return err
	}

	quoted := false
	args := strings.FieldsFunc(command, func(r rune) bool {
		if r == '"' {
			quoted = !quoted
		}
		return !quoted && r == ' '
	})

	for x := range args {
		args[x] = strings.ReplaceAll(args[x], "\"", "")
	}

	// ********************************
	// execute the command
	// ********************************
	output, err := exec.Command(app_path, args...).Output()
	if err != nil {
		if w.LogLevel > workflow.LOG_QUIET {
			fmt.Println(string(output))
		}

		return fmt.Errorf("error executing command: %v %v err: %s", app_path, command, err.Error())
	}

	//******************
	//Display the result
	//******************
	if w.LogLevel > workflow.LOG_QUIET {
		fmt.Println(string(output))
	}

	//*******************
	//Process the results
	//*******************
	err = w.ActionProcessResults(m, output)
	if err != nil {
		return err
	}

	//**************
	//Return the err
	//**************
	return nil
}
