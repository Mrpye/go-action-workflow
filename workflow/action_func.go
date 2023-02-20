package workflow

import (
	"strconv"
	"strings"
)

// GetConfigString returns the Action config as a string
// Key is the name of the config
func (m *Action) GetConfigString(key string) string {
	if value, ok := m.Config[key]; ok {
		switch val := value.(type) {
		case bool:
			return strconv.FormatBool(val)
		case string:
			return val
		case int:
			return strconv.Itoa(val)
		case int64:
			return strconv.FormatInt(val, 10)
		case float64:
			return strconv.FormatFloat(val, 'f', -1, 64)
		default:
			return ""
		}
	}
	return ""
}

// GetConfigBool returns the Action config as a bool
// Key is the name of the config
func (m *Action) GetConfigBool(key string) bool {
	if value, ok := m.Config[key]; ok {
		switch val := value.(type) {
		case bool:
			return val
		case string:
			if strings.ToLower(val) == "t" {
				return true
			} else if strings.ToLower(val) == "f" {
				return false
			} else if strings.ToLower(val) == "yes" {
				return true
			} else if strings.ToLower(val) == "no" {
				return false
			}
			b, _ := strconv.ParseBool(val)
			return b
		case int:
			return val > 0
		default:
			return false
		}
	}
	return false
}

// GetConfig returns the Action config as an interface
// Key is the name of the config
func (m *Action) GetConfig(key string) interface{} {
	if val, ok := m.Config[key]; ok {
		return val
	}
	return ""
}

// GetConfigInt returns the Action config as an int
// Key is the name of the config
func (m *Action) GetConfigInt(key string) int {
	if value, ok := m.Config[key]; ok {
		switch val := value.(type) {
		case string:
			int_val, err := strconv.Atoi(val)
			if err != nil {
				return 0
			}
			return int_val
		case int:
			return val
		case float64:
			return int(val)
		default:
			// User defined types work as well
			return 0
		}
	}
	return 0
}
