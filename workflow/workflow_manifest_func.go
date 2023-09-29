package workflow

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path"
	"strings"
	"time"

	go_data_chain "github.com/Mrpye/go-data-chain"
	"github.com/Mrpye/golib/dir"
	"github.com/Mrpye/golib/file"
	"gopkg.in/yaml.v3"
)

// AddManifest adds a manifest to the workflow
//func (m *Workflow) AddManifest() {
//	m.Manifest = *CreateManifest()
//}

func (m *Workflow) ManifestGetJob(job_key string) *Job {
	return m.Manifest.GetJob(job_key)
}

// UpdateMeta updates the meta data of the manifest
// meta: the meta data to update
func (m *Workflow) ManifestUpdateMeta(meta *MetaData) {
	if meta.Name != "" {
		m.Manifest.Meta.Name = meta.Name
	}
	if meta.Description != "" {
		m.Manifest.Meta.Description = meta.Description
	}
	if meta.Version != "" {
		m.Manifest.Meta.Version = meta.Version
	}
	if meta.Author != "" {
		m.Manifest.Meta.Author = meta.Author
	}
	if meta.CreatedDate == "" {
		m.Manifest.Meta.CreatedDate = meta.CreatedDate
	}
	if meta.UpdateDate != "" {
		m.Manifest.Meta.UpdateDate = meta.UpdateDate
	}
}

// AddMeta adds the meta data to the manifest overwriting any existing meta data
// meta: the meta data to add
func (m *Workflow) ManifestCreateMeta(meta *MetaData) {
	m.Manifest.Meta = *meta

}

// AddJob adds a job to the manifest
// job: the job to add
// returns an error if the job already exists or if the job key is empty
func (m *Workflow) ManifestAddJob(job *Job) error {
	if job.Key == "" {
		return errors.New("Job key cannot be empty")
	}
	if m.Manifest.GetJob(job.Key) != nil {
		return fmt.Errorf("Job already exists: %s", job.Key)
	}
	m.Manifest.Jobs = append(m.Manifest.Jobs, *job)
	return nil
}

// DeleteJob delete a job from the manifest
// job_key: the key of the job to delete
// returns an error if the job does not exist
func (m *Workflow) ManifestDeleteJob(job_key string) error {
	job := m.Manifest.GetJob(job_key)
	if job == nil {
		return fmt.Errorf("Job does not exist: %s", job_key)
	}

	for i, j := range m.Manifest.Jobs {
		if j.Key == job_key {
			m.Manifest.Jobs = append(m.Manifest.Jobs[:i], m.Manifest.Jobs[i+1:]...)
			return nil
		}
	}
	return nil
}

// AddActionToJob adds an action to a job
// job_key: the key of the job to add the action to
// action_obj: the action to add
func (m *Workflow) ManifestAddActionToJob(job_key string, action_obj *Action) error {
	job := m.Manifest.GetJob(job_key)
	if job == nil {
		return fmt.Errorf("Job does not exist: %s", job_key)
	}

	if job.ActionKeyExists(action_obj.Key) {
		return fmt.Errorf("Action already exists: %s", action_obj.Key)
	}

	job.Actions = append(job.Actions, *action_obj)
	return nil
}

// AddActionToJob adds an action to a job
// job_key: the key of the job to add the action to
// action_obj: the action to add
func (m *Workflow) ManifestAddGlobalAction(action_obj *Action) error {

	if m.Manifest.GlobalActionKeyExists(action_obj.Key) {
		return fmt.Errorf("global Action already exists: %s", action_obj.Key)
	}

	m.Manifest.Actions = append(m.Manifest.Actions, *action_obj)
	return nil
}

// DeleteActionFromJob deletes an action from a job by key
// job_key: the key of the job to delete the action from
// action_key: the key of the action to delete
// returns an error if the job or action does not exist
func (m *Workflow) ManifestDeleteActionFromJob(job_key string, action_key string) error {
	job := m.Manifest.GetJob(job_key)
	if job == nil {
		return fmt.Errorf("Job does not exist: %s", job_key)
	}

	for i, action := range job.Actions {
		if action.Key == action_key {
			job.Actions = append(job.Actions[:i], job.Actions[i+1:]...)
			return nil
		}
	}

	return fmt.Errorf("Action does not exist with key: %s", action_key)
}

// ManifestCreateParameter creates a parameter in the manifest
// param: the parameter to create
// returns an error if the parameter already exists
func (m *Workflow) ManifestCreateParameter(param *Parameter) error {
	if m.Manifest.GetParameter(param.Key) != nil {
		return fmt.Errorf("Parameter already exists: %s", param.Key)
	}
	m.Manifest.Parameters = append(m.Manifest.Parameters, *param)
	return nil
}

// ManifestDeleteParameter deletes a parameter from the manifest
// param_key: the key of the parameter to delete
// returns an error if the parameter does not exist
func (m *Workflow) ManifestDeleteParameter(param_key string) error {
	param := m.Manifest.GetParameter(param_key)
	if param == nil {
		return fmt.Errorf("Parameter does not exist: %s", param_key)
	}

	for i, p := range m.Manifest.Parameters {
		if p.Key == param_key {
			m.Manifest.Parameters = append(m.Manifest.Parameters[:i], m.Manifest.Parameters[i+1:]...)
			return nil
		}
	}
	return nil
}

// LoadManifest loads the manifest
// package_path: the path to the manifest file
// returns an error if the manifest is not valid
func (m *Workflow) LoadManifest(package_path string) error {

	//get the path from the package path
	p := path.Dir(package_path)
	m.SetBasePath(p)

	//****************
	//Load the package
	//****************
	file, _ := ioutil.ReadFile(package_path)
	return m.LoadManifestFromString(string(file))
}

// LoadManifestFromString loads the manifest from a string
// manifest_string: the manifest as a string
// returns an error if the manifest is not valid
func (m *Workflow) LoadManifestFromString(manifest_string string) error {
	manifest := CreateManifest()
	err := yaml.Unmarshal([]byte(manifest_string), &manifest)
	if err != nil {
		return err
	}
	m.Manifest = *manifest
	return nil
}

func getNowAsString() string {
	return time.Now().Format(time.RFC3339)
}

// SaveManifest saves the manifest
// file_name: the name of the file to save
// returns an error if the manifest is not valid
func (m *Workflow) SaveManifest(file_name string) error {

	err := dir.MakeDirAll(file_name)
	if err != nil {
		return err
	}
	//set the date to now as string
	if m.Manifest.Meta.CreatedDate == "" {
		m.Manifest.Meta.CreatedDate = getNowAsString()
	}
	m.Manifest.Meta.UpdateDate = getNowAsString()

	file, err := yaml.Marshal(m.Manifest)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(file_name, file, 0644)
	return err
}

// GetDataItem returns the data item with the given name
// item: the name of the data item to return
// returns the data item or nil if the item does not exist
func (m *Workflow) GetDataItem(item string) *go_data_chain.Data {
	if m.Manifest.Data != nil {
		data := m.Manifest.DataModel()
		if data != nil && data.Err == nil {
			return data.GetMapItem(item)
		}
	}
	return nil
}

// GetParamValue will return the value of the parameter with the given key
// key - the key of the parameter to return
// returns the value of the parameter or nil if the parameter does not exist
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
		return parsed_str //m.convertTo(p.InputType, parsed_str)
	default:
		return data_val
	}
}

// LoadAnswerFile loads the answers from the given file
// answer_file: the file to load the answers from
// returns an error if the answers are not valid
func (m *Workflow) LoadAnswerFile(answer_file string) error {
	//*****************************************
	//Read the answer file and put answers back
	//*****************************************
	file_data, err := file.ReadFileToString(answer_file)
	if err != nil {
		return err
	}
	data := Answers{}
	err = yaml.Unmarshal([]byte(file_data), &data)
	if err != nil {
		return err
	}

	err = m.loadAnswerObj(data)
	if err != nil {
		return err
	}

	return nil
}

// LoadAnswerObj loads the answers from the given object
// data: the answers to load
// returns an error if the answers are not valid
func (m *Workflow) loadAnswerObj(data Answers) error {
	//*****************************************
	//Read the answer file and put answers back
	//*****************************************
	if data.Name != m.Manifest.Meta.Name {
		return fmt.Errorf("this answer file is not for this package! answer file:%s package:%s", data.Name, m.Manifest.Meta.Name)
	}
	if data.Version != m.Manifest.Meta.Version {
		return fmt.Errorf("this answer file is not for this version of the package! answer file:%s package:%s", data.Version, m.Manifest.Meta.Version)
	}
	//*******************
	//Restore the answers
	//*******************
	for _, o := range data.Answers {
		for j, k := range m.Manifest.Parameters {
			if strings.EqualFold(k.Key, o.Key) {
				m.Manifest.Parameters[j].SetAnswer(o.Value)
				break
			}
		}
	}

	return nil
}

func (m *Workflow) BuildAnswerFile() (*Answers, error) {
	//****************************************
	//Read the params and build an answer file
	//****************************************
	answers := &Answers{
		Name:    m.Manifest.Meta.Name,
		Version: m.Manifest.Meta.Version,
	}

	for _, o := range m.Manifest.Parameters {
		//******************************
		//See if the parameter is hidden
		//******************************
		answers.Answers = append(answers.Answers, Answer{Key: o.Key, Title: o.Title, Description: o.Description, Value: o.GetValue(), InputType: o.InputType})
	}

	return answers, nil
}

// Gets the data from the template using a string
func (m *Workflow) GetDataFromString(path string) (interface{}, error) {
	if strings.HasPrefix(path, "$") {
		//The we read the data
		data_path := path[1:]
		parts := strings.Split(data_path, ".")
		data := m.Manifest.DataModel()
		for i, part := range parts {
			if strings.EqualFold(part, "data") && i == 0 {
				continue
			}
			data = data.GetMapItem(part)
			if data == nil {
				return nil, fmt.Errorf("could not find data for path:%s", path)
			}
		}
		if data != nil {
			return data.ToInterface(), nil
		}
	}
	return path, nil
}
