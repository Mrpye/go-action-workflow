package workflow

import (
	"github.com/Mrpye/golib/convert"
)

// GetConfigString returns the Action config as a string
// Key is the name of the config
// returns the config as a string
func (m *Action) GetConfigString(key string) string {
	if value, ok := m.Config[key]; ok {
		return convert.ToString(value)
	}
	return ""
}

// GetConfigBool returns the Action config as a bool
// Key is the name of the config
// returns the config as a bool
func (m *Action) GetConfigBool(key string) bool {
	if value, ok := m.Config[key]; ok {
		return convert.ToBool(value)
	}
	return false
}

// GetConfig returns the Action config as an interface
// Key is the name of the config
// returns the config as an interface
func (m *Action) GetConfig(key string) interface{} {
	if val, ok := m.Config[key]; ok {
		return val
	}
	return ""
}

// GetConfigInt returns the Action config as an int
// Key is the name of the config
// returns the config as an int
func (m *Action) GetConfigInt(key string) int {
	if value, ok := m.Config[key]; ok {
		return convert.ToInt(value)
	}
	return 0
}

// GetConfigInt returns the Action config as an int
// Key is the name of the config
// returns the config as an int
func (m *Action) GetConfigFloat64(key string) float64 {
	if value, ok := m.Config[key]; ok {
		return convert.ToFloat64(value)
	}
	return 0
}
