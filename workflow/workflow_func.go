package workflow

import (
	"errors"
	"fmt"
	"path"
	"strings"
	"text/template"

	"github.com/Mrpye/golib/log"
)

// CreateSubWorkflowEngine creates a new workflow engine from the current workflow engine
// This is useful for creating sub workflows
func (m *Workflow) CreateSubWorkflowEngine() *Workflow {
	wf := CreateWorkflow()
	wf.Manifest = m.Manifest
	wf.ActionList = m.ActionList
	wf.dataBucket = m.dataBucket
	wf.LogLevel = m.LogLevel
	wf.templateFuncMap = wf.GetInbuiltTemplateFuncMap()
	wf.InitFunc = m.InitFunc
	return wf
}

// SetBasePath will set the base path of the manifest file
// base_path - the base path of the manifest file
func (m *Workflow) SetBasePath(base_path string) {
	m.base_path = base_path
}

// BuildPath will build a path from the base path and the file path
// file_path - the file path to build
// returns the full path
func (m *Workflow) BuildPath(file_path string) string {
	return path.Join(m.base_path, file_path)
}

func (m *Workflow) CountArray(v []interface{}) int {
	return len(v)
}

// GetCurrentJob returns the current job
func (m *Workflow) GetCurrentJob() *Job {
	return m.current_job
}

// MapValuesToInput will map the values in the map to the input of the current job
// map_data - the map of values to map to the input
func (m *Workflow) MapValuesToInput(map_data map[string]interface{}) error {
	for k, v := range map_data {
		err := m.current_job.SetInputAnswer(k, v)
		if err != nil {
			return err
		}
	}
	return nil
}

// GetInputValue will return the value of the job input
// key - the key of the parameter to return
// returns the value of the parameter or nil if the parameter does not exist
func (m *Workflow) GetInputValue(key string) interface{} {
	//*************
	//Get the param
	//*************
	input_param := m.current_job.GetInput(key)

	if input_param == nil {
		return nil
	}
	//*******************
	//Get the value as is
	//*******************
	value := input_param.GetValue()
	switch data_val := value.(type) {
	case string:
		parsed_str, _ := m.ParseToken(m.model, string(data_val))
		if m.LogLevel == LOG_VERBOSE {
			log.LogVerbose(fmt.Sprintf("ParseToken Value(%v) Result(%v)\n", data_val, parsed_str))
		}
		return parsed_str
	default:
		return data_val
	}
}

// SetCurrentActionIndex will set the current action index
// index - the index to set
// returns an error if the index is out of range or if the index is greater than the next action
func (m *Workflow) SetCurrentActionIndex(index int) error {
	//***********************************
	// Check if the index is out of range
	//***********************************
	if index < 0 || index >= len(m.current_job.Actions) {
		return errors.New("index out of range")
	}

	//******************************
	// Take a peek at the loop stack
	//******************************
	loop, _ := m.stack.Peek()

	// *************************************************
	//check if the index is greater than the next action
	//if it is then we need to throw an error
	// *************************************************
	if loop != nil {
		temp_index := -1
		for i := loop.Index; i < len(m.current_job.Actions); i++ {
			if m.current_job.Actions[i].Action == "next" {
				temp_index = i
			}
		}
		if index > temp_index {
			return errors.New("cannot set index to a value greater than the next action")
		}
	}
	//**********************
	// Set the current index
	//**********************
	m.current_index = index
	return nil
}

// GetCurrentActionIndex will get the current action index
// returns the current action index
func (m *Workflow) GetCurrentActionIndex() int {
	return m.current_index
}

// SetTemplateFuncMap sets the template function map
func (m *Workflow) SetTemplateFuncMap(f template.FuncMap) {
	m.templateFuncMap = f
}

func (m *Workflow) GetTemplateFuncMap() template.FuncMap {
	return m.templateFuncMap
}

// splitActionParams will split the action into parts
// - action  - the action to split
// - returns - the action parts
// - returns - error
func (m *Workflow) splitActionParams(action string, process_condition_only bool) ([]string, error) {

	//****************************
	//See if to process the tokens
	//****************************
	action_parts := strings.Split(fmt.Sprintf("%v", action), ";")
	parse_tokens := false
	switch strings.ToLower(action_parts[0]) {
	case "for", "next":
		parse_tokens = true
	default:
		if !process_condition_only {
			parse_tokens = true
		}
	}

	//********************************
	//Parse any variable in the action
	//********************************
	lowercase_action := ""
	if parse_tokens {
		parsed_action, err := m.ParseToken(m.model, action)
		if err != nil {
			return nil, err
		}
		lowercase_action = parsed_action
	}
	//***************************
	//Split the action into parts
	//***************************
	action_parts = strings.Split(fmt.Sprintf("%v", lowercase_action), ";")
	action_parts[0] = strings.ToLower(action_parts[0])

	return action_parts, nil
}

// Get Runtime vars from the workflow using key
// key - the key of the parameter to return
// returns the value of the parameter or nil if the parameter does not exist
func (m *Workflow) GetRuntimeVar(key string) (interface{}, error) {

	//***************
	//check the vars
	//***************
	if m.runtime_vars == nil {
		return nil, fmt.Errorf("runtime_vars is nil")
	}

	//********************
	//check the key exists
	//********************
	if m.runtime_vars[key] == nil {
		return nil, fmt.Errorf("key %v does not exist", key)
	}

	//*******************
	//Get the value as is
	//*******************
	value := m.runtime_vars[key]
	return value, nil
}

// GetConfigValue will return the value of the config with the given key
// key - the key of the config to return
// data_type - the type of the config to return
// custom - custom values to pass to the config function
// returns the value of the config or nil if the config does not exist
// returns an error if the config function is not set
func (w *Workflow) GetConfigValue(config_target string, key string, data_type string, custom ...string) (interface{}, error) {
	var tmp_config ReadConfigFunc

	//************************************************
	//See if the config function is set for the target
	//************************************************
	if w.ReadConfigFunc[config_target] != nil {
		tmp_config = w.ReadConfigFunc[config_target]
	} else if config_target != "" {
		return nil, fmt.Errorf("config function not set for target %v", config_target)
	}

	//********************************************************************************
	//If the config target is not set or is default then get the first item in the map
	//********************************************************************************
	if config_target == "default" || config_target == "" {
		//get the first item in the map
		for _, v := range w.ReadConfigFunc {
			tmp_config = v
			break
		}
	}

	//*******************
	//Get the value as is
	//*******************
	result, err := tmp_config(key, data_type, custom...)
	if w.LogLevel == LOG_VERBOSE {
		log.LogVerbose(fmt.Sprintf("GetConfigValue Key(%v) Type(%v) Value(%v)\n", key, data_type, result))
	}
	return result, err
}

// MapTargetConfigValue will call the target config mapper function
// m - the template data
// target - the target to map
// returns the mapped target
// returns an error if the mapper function is not set
func (w *Workflow) MapTargetConfigValue(m interface{}, target interface{}) (interface{}, error) {
	//************************
	//Call the mapper function
	//************************
	return w.TargetMapFunc(w, m, target)
}

// SetAnswers will set the answers for the workflow
// answers - the answers to set
func (m *Workflow) SetAnswers(answers map[string]interface{}) {
	for key, value := range m.Manifest.Parameters {
		if answers[value.Key] != nil {
			m.Manifest.Parameters[key].SetAnswer(answers[value.Key])
		}
	}
}

// AddActionSchema adds an action and target schema to the client
// and adds the action to the workflow
func (m *Workflow) AddActionSchema(sch SchemaEndpoint) {

	//********************************
	//Add the actions to the workflow
	//********************************
	actions := sch.GetActionSchema()
	for key := range actions {
		action_schema := actions[key]
		if action_schema.Action != nil {
			m.ActionList[key] = action_schema.Action
		}
	}

	//*********************************
	//Add function maps to the workflow
	//*********************************
	fm := m.GetTemplateFuncMap()
	if fm == nil {
		if fm == nil {
			m.SetTemplateFuncMap(m.GetInbuiltTemplateFuncMap())

		}
	}
	if fm == nil {
		fm = m.GetTemplateFuncMap()
	}
	//**************************************************
	//Add the functions to the workflow from the library
	//**************************************************
	if sch.GetFunctionMap() != nil {
		functions := sch.GetFunctionMap()
		for key := range functions {
			function_schema := functions[key]
			fm[key] = function_schema.Function
		}
	}

}
