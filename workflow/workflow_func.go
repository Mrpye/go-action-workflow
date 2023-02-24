package workflow

import (
	"errors"
	"fmt"
	"io/ioutil"
	"text/template"

	go_data_chain "github.com/Mrpye/go-data-chain"
	"github.com/Mrpye/golib/lib"
	"gopkg.in/yaml.v3"
)

// SetTemplateFuncMap sets the template function map
func (m *Workflow) CreateSubWorkflowEngine() *Workflow {
	wf := CreateWorkflow()
	wf.Manifest = m.Manifest
	wf.ActionList = m.ActionList
	wf.dataBucket = m.dataBucket
	wf.Verbose = m.Verbose
	wf.templateFuncMap = wf.GetTemplateFuncMap()
	wf.InitFunc = m.InitFunc
	return wf
}

// SetTemplateFuncMap sets the template function map
func (m *Workflow) GetCurrentJob() *Job {
	return m.current_job
}

//LoadManifest loads the manifest from a string
//package_path: the path to the manifest file
func (m *Workflow) LoadManifestFromString(manifest_string string) error {
	manifest := CreateManifest()
	err := yaml.Unmarshal([]byte(manifest_string), &manifest)
	if err != nil {
		return err
	}
	m.Manifest = *manifest

	return nil
}

func (m *Workflow) GetDataItem(item string) *go_data_chain.Data {
	if m.Manifest.Data != nil {
		data := m.Manifest.DataModel().GetMapItem("data")
		if data != nil {
			return data.GetMapItem(item)
		}

	}
	return nil
}

//LoadManifest loads the manifest
//package_path: the path to the manifest file
func (m *Workflow) LoadManifest(package_path string) error {
	//****************
	//Load the package
	//****************
	file, _ := ioutil.ReadFile(package_path)
	return m.LoadManifestFromString(string(file))
}

//SaveManifest saves the manifest
//file_name: the name of the file to save
func (m *Workflow) SaveManifest(file_name string) error {
	file, _ := yaml.Marshal(m.Manifest)
	//file_name := fmt.Sprintf("answer_%s_%s.yaml", m.Package.Meta.PackageName, m.currentAppProfile)
	err := ioutil.WriteFile(file_name, file, 0644)
	return err
}

// GetParamValue will return the value of the parameter with the given key
// key - the key of the parameter to return
func (m *Workflow) GetParamValue(key string) interface{} {
	//*************
	//Get the param
	//*************
	p := m.Manifest.GetParameter(key)

	if p == nil {
		return nil
	}
	//*******************
	//Get the value as is
	//*******************
	value := p.GetValue()
	switch data_val := value.(type) {
	case string:
		parsed_str, _ := m.ParseToken(m.model, string(data_val))
		return parsed_str
	default:
		return data_val
	}
}

// GetInputValue will return the value of the job input
// key - the key of the parameter to return
func (m *Workflow) MapValuesToInput(key map[string]interface{}) error {
	for k, v := range key {
		err := m.current_job.SetInputAnswer(k, v)
		if err != nil {
			return err
		}
	}
	return nil
}

// GetInputValue will return the value of the job input
// key - the key of the parameter to return
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
		if m.Verbose == LOG_VERBOSE {
			lib.LogVerbose(fmt.Sprintf("GetTokenString Value(%v) Result(%v)\n", data_val, parsed_str))
		}
		return parsed_str
	default:
		return data_val
	}
}

// SetCurrentActionIndex will set the current action index
// index - the index to set
func (m *Workflow) SetCurrentActionIndex(index int) error {
	if index < 0 || index >= len(m.current_job.Actions) {
		return errors.New("index out of range")
	}

	loop, _ := m.stack.Peek()

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
