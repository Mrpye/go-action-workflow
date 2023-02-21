package workflow

import (
	"io/ioutil"
	"strings"

	"gopkg.in/yaml.v3"
)

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

// GetParam will return the parameter with the given key
// key - the key of the parameter to return
func (m *Workflow) GetParam(key string) *Parameter {
	for i, o := range m.Manifest.Parameters {
		if strings.EqualFold(o.Key, key) {
			return &m.Manifest.Parameters[i]
		}
	}
	return nil
}

// GetParamValue will return the value of the parameter with the given key
// key - the key of the parameter to return
func (m *Workflow) GetParamValue(key string) interface{} {
	//*************
	//Get the param
	//*************
	p := m.GetParam(key)

	if p == nil {
		return nil
	}

	//*******************
	//Get the value as is
	//*******************
	value := p.GetValue()
	switch data_val := value.(type) {
	case string:
		parsed_str, _ := m.ParseToken(m.Model, string(data_val))
		return parsed_str
	case int:
		return data_val
	default:
		return data_val
	}

}
