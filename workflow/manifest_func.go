package workflow

import (
	"strings"

	go_data_chain "github.com/Mrpye/go-data-chain"
)

// ParameterExists checks if a parameter exists
// key is the key of the parameter
// returns true if the parameter exists
func (m *Manifest) ParameterExists(key string) bool {
	for _, o := range m.Parameters {
		if o.Key == strings.ToLower(key) {
			return true
		}
	}
	return false
}

// GetParameter gets a parameter
// key is the key of the parameter
// returns the parameter or nil if it does not exist
func (m *Manifest) GetParameter(key string) *Parameter {
	for _, o := range m.Parameters {
		if o.Key == strings.ToLower(key) {
			return &o
		}
	}
	return nil
}

// GetJob gets a job from the manifest
// key is the key of the job
// returns the job or nil if it does not exist
func (m *Manifest) GetJob(key string) *Job {
	for i, o := range m.Jobs {
		if o.Key == strings.ToLower(key) {
			return &m.Jobs[i]
		}
	}
	return nil
}

// GetJob gets a job from the manifest
// key is the key of the job
// returns the job or nil if it does not exist
func (m *Manifest) GetGlobalAction(key string) *Action {
	for i, o := range m.Actions {
		if o.Key == strings.ToLower(key) {
			return &m.Actions[i]
		}
	}
	return nil
}

func (m *Manifest) GlobalActionKeyExists(key string) bool {
	if key == "" {
		return false
	}
	for _, o := range m.Actions {
		if o.Key == strings.ToLower(key) {
			return true
		}
	}
	return false
}

// DataModel returns the data model for the manifest
// returns the data model for the manifest as a data chain
func (m *Manifest) DataModel() *go_data_chain.Data {
	if m.Data != nil {
		return go_data_chain.CreateDataChain(m.Data, true)
	}
	return nil
}
