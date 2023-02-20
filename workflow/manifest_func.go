package workflow

import (
	"strings"
)

func (m *Manifest) ParameterExists(key string, app_profile string) bool {
	for _, o := range m.Parameters {
		if o.Key == strings.ToLower(key) {
			return true
		}
	}
	return false
}

//get the parameter
func (m *Manifest) GetParameter(key string, app_profile string) *Parameter {
	for _, o := range m.Parameters {
		if o.Key == strings.ToLower(key) {
			return &o
		}
	}
	return nil
}

func (m *Manifest) GetJob(key string) *Job {
	for i, o := range m.Jobs {
		if o.Key == strings.ToLower(key) {
			return &m.Jobs[i]
		}
	}
	return nil
}
