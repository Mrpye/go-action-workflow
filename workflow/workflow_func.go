package workflow

import (
	"io/ioutil"
	"strings"

	"gopkg.in/yaml.v3"
)

//LoadManifest loads the manifest
//package_path: the path to the manifest file
func (m *Workflow) LoadManifest(package_path string) error {
	//****************
	//Load the package
	//****************
	manifest := CreateManifest()
	file, _ := ioutil.ReadFile(package_path)
	err := yaml.Unmarshal(file, &manifest)
	if err != nil {
		return err
	}
	m.Manifest = *manifest

	return nil
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
	return p.GetValue()

}
