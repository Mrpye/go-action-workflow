package workflow

import (
	"fmt"
	"strings"
)

//Returns input
func (m *Job) GetInput(key string) *Parameter {
	if val, ok := m.Inputs[key]; ok {
		return &val
	}
	return nil
}

func (m *Job) SetInputAnswer(key string, value interface{}) error {
	if val, ok := m.Inputs[key]; ok {
		val.SetAnswer(value)
		m.Inputs[key] = val
		return nil
	}
	return fmt.Errorf("cannot find input with key %s", key)
}

func (m *Job) GetKeyIndex(key string) int {
	for i, o := range m.Actions {
		if strings.EqualFold(o.Key, key) {
			return i
		}
	}
	return -1
}

func (m *Job) ActionExists(key string) bool {
	parts := strings.Split(key, ";")
	switch parts[0] {
	case "end", "print", "goto", "do", "do-end", "loop", "loop-end", "wait", "wait-seconds", "wait-minutes":
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

func (m *Job) GetActionByKey(key string) *Action {
	for i, o := range m.Actions {
		if o.Key == strings.ToLower(key) {
			return &m.Actions[i]
		}
	}
	return nil
}
