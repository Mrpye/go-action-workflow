package workflow

import (
	"fmt"
	"strings"
)

// GetInput returns the Job input as an interface
// - key is the name of the parameter
// - returns the parameter
func (m *Job) GetInput(key string) *Parameter {
	if val, ok := m.Inputs[key]; ok {
		return &val
	}
	return nil
}

// SetInputAnswer sets the answer for the input parameter
// - key is the name of the parameter
// - value is the answer
// - returns an error if the input does not exist
func (m *Job) SetInputAnswer(key string, value interface{}) error {
	if val, ok := m.Inputs[key]; ok {
		val.SetAnswer(value)
		m.Inputs[key] = val
		return nil
	}
	return fmt.Errorf("cannot find input with key %s", key)
}

// GetKeyIndex returns the index of the action with the given key
// - key is the name of the action
// - returns the index of the action
func (m *Job) GetKeyIndex(key string) int {
	for i, o := range m.Actions {
		if strings.EqualFold(o.Key, key) {
			return i
		}
	}
	return -1
}

// ActionExists returns true if the action exists
// - key is the name of the action
// - returns true if the action exists
func (m *Job) ActionExists(key string) bool {
	parts := strings.Split(key, ";")
	switch parts[0] {
	case "end", "print", "goto", "fail", "for", "next", "wait", "wait-seconds", "wait-minutes":
		return true
	default:
		for _, o := range m.Actions {
			if o.Key == strings.ToLower(key) {
				return true
			}
		}
	}
	return false
}

// ActionKeyExists returns true if the action exists
// - key is the name of the action
// - returns true if the action exists
func (m *Job) ActionKeyExists(key string) bool {
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

// GetActionByKey returns the action with the given key
// - key is the name of the action
// - returns the action
func (m *Job) GetActionByKey(key string) *Action {
	for i, o := range m.Actions {
		if o.Key == strings.ToLower(key) {
			return &m.Actions[i]
		}
	}
	return nil
}
