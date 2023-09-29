package workflow

import (
	"fmt"
	"io/ioutil"

	golib_str "github.com/Mrpye/golib/str"
	"gopkg.in/yaml.v3"
)

type Answers struct {
	Name    string   `json:"package_name" yaml:"package_name"`
	Version string   `json:"package_version" yaml:"package_version"`
	Answers []Answer `json:"answers" yaml:"answers"`
}

type Answer struct {
	Key         string      `json:"key" yaml:"key" flag:"key k" desc:"key for the parameter, must be unique within the workflow, used as the parameter name"`
	Title       string      `json:"title,omitempty" yaml:"title,omitempty" flag:"title t" desc:"title for the parameter, used as the parameter title"`
	Description string      `json:"description,omitempty" yaml:"description,omitempty" flag:"desc d" desc:"description for the parameter, used as the parameter description"`
	InputType   string      `json:"type,omitempty" yaml:"type,omitempty" flag:"type y" desc:"type of the parameter, must be one of string, int, float, bool"`
	Value       interface{} `json:"value" yaml:"value" flag:"value v" desc:"value of the parameter, must be a string, int, float, or bool"`
}

func (m *Answers) GetPackageProjectName() string {
	return golib_str.CreateKey(fmt.Sprintf("%s-%s", m.Name, m.Version), "-")
}

func (m *Answers) Save(file_name string) error {
	file, _ := yaml.Marshal(m)
	err := ioutil.WriteFile(file_name, file, 0644)
	return err
}
