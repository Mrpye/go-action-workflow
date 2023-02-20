package workflow

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"text/template"

	"github.com/Mrpye/golib/lib"
)

// TemplateData is the data that is passed to the template engine
type TemplateData struct {
	Meta          *MetaData
	Manifest      *Manifest
	CurrentAction *Action
	//CurrentTarget interface{}
	//Data          map[string]interface{}
}

// SetAction sets the current action
// current_action - the current action
func (m *TemplateData) SetAction(current_action *Action) {
	m.CurrentAction = current_action
}

// CreateTemplateData creates the template data
// current_action - the current action
// target - the current target
// returns the template data
func (m *Workflow) CreateTemplateData(current_action *Action) *TemplateData {
	return &TemplateData{
		Meta:          &m.Manifest.Meta,
		Manifest:      &m.Manifest,
		CurrentAction: current_action,
		//CurrentTarget: target,
	}
}

// GetStackVariable gets the value of a variable from the stack
// variable_name - the name of the variable
// returns the value of the variable
func (m *Workflow) GetStackVariable(variable_name string) int {
	value, err := m.stack.PeekVariable(variable_name)
	if err != nil {
		return -1
	}
	return value.CurrentValue
}

//ParseInterfaceMap will parse the interface map and replace its tokens
// model - the template data
// val - the map to parse
// returns the value of the map with the tokens replaced
// returns an error if the key is not found
func (m *Workflow) ParseInterfaceMap(model *TemplateData, val map[string]interface{}) map[string]interface{} {
	for k, v := range val {
		_ = k
		switch data_val := v.(type) {
		case string:
			parsed_str, _ := m.ParseToken(model, string(data_val))
			val[k] = parsed_str
		case map[string]interface{}:
			m.ParseInterfaceMap(model, data_val)
		}
	}
	return val
}

//GetConfigTokenMap will get the config token map from the action and replaces its tokens
// key - the key of the config
// model - the template data
// required - if the key is required
// returns the value of the map with the tokens replaced
// returns an error if the key is not found
func (m *Workflow) GetConfigTokenMap(key string, model *TemplateData, required bool) (map[string]interface{}, error) {
	if val, ok := model.CurrentAction.Config[key]; ok {
		switch val := val.(type) {
		case map[string]interface{}:
			new_val := lib.CopyMap(val)
			new_val = m.ParseInterfaceMap(model, new_val)
			return new_val, nil
		default:
			// User defined types work as well
			return nil, nil
		}
	}
	if required {
		return nil, fmt.Errorf("required parameter %s missing", key)
	}
	return nil, nil
}

// GetConfigToken will get the config token from the action and replaces its tokens
// key - the key of the config
// model - the template data
// Returns the value of the config
// Returns an error if the key is not found
func (m *Workflow) GetConfigToken(key string, model *TemplateData) (interface{}, error) {
	if val, ok := model.CurrentAction.Config[key]; ok {
		//Need to check the type is a string
		switch val := val.(type) {
		case string:
			new_val, err := m.ParseToken(model, val)
			if err != nil {
				return val, fmt.Errorf("error with key %s, err:%v", key, err)
			}
			return new_val, nil
		default:
			// User defined types work as well
			return val, nil
		}

	}
	return "", nil
}

// GetConfigTokenInterface will get the config token from the action and replaces its tokens and returns it as an interface
// key - the key of the config
// model - the template data
// Returns the value of the config as an interface
// Returns an error if the key is not found
func (m *Workflow) GetConfigTokenInterface(key string, model *TemplateData, required bool) (interface{}, error) {
	if value, ok := model.CurrentAction.Config[key]; ok {
		switch val := value.(type) {
		case string:
			new_val, err := m.ParseToken(model, val)
			if m.Verbose == LOG_VERBOSE {
				lib.LogVerbose(fmt.Sprintf("GetTokenString Value(%v) Result(%v)\n", val, new_val))
			}
			if err != nil {
				return val, fmt.Errorf("error with key %s, err:%v", key, err)
			}
			return new_val, nil
		case []interface{}:
			for i, v := range val {
				switch v := v.(type) {
				case string:
					new_val, err := m.ParseToken(model, v)
					if m.Verbose == LOG_VERBOSE {
						lib.LogVerbose(fmt.Sprintf("GetTokenString Value(%v) Result(%v)\n", v, new_val))
					}
					if err != nil {
						return v, err
					}
					val[i] = new_val
				}
			}
			return val, nil
		default:
			// User defined types work as well
			return val.(string), nil
		}
		//Need to check the type is a string
	}
	if required {
		return "", fmt.Errorf("required parameter %s missing", key)
	}
	return "", nil
}

// GetConfigTokenString will get the config token from the action and replaces its tokens as a string
// key - the key of the config
// model - the template data
// Returns the value of the config as a string
// Returns an error if the key is not found
func (m *Workflow) GetConfigTokenString(key string, model *TemplateData, required bool) (string, error) {
	if val, ok := model.CurrentAction.Config[key]; ok {
		value, err := m.GetTokenString(val, model)
		if err != nil {
			return value, fmt.Errorf("error with key %s, err:%v", key, err)
		}
		if m.Verbose == LOG_VERBOSE {
			lib.LogVerbose(fmt.Sprintf("GetConfigTokenString Param(%s) Value(%s) Result(%s)\n", key, val, value))
		}
		if required && value == "" {
			return value, fmt.Errorf("missing required value for %s", key)
		}
		return value, nil
	}
	if required {
		return "", fmt.Errorf("required parameter %s missing", key)
	}
	return "", nil
}

// GetTokenString will parse the token string
// value - the value to parse
// model - the template data
// Returns the parsed string
// Returns an error if the key is not found
func (m *Workflow) GetTokenString(value interface{}, model *TemplateData) (string, error) {

	//Need to check the type is a string
	switch val := value.(type) {
	case string:
		new_val, err := m.ParseToken(model, val)
		if m.Verbose == LOG_VERBOSE {
			lib.LogVerbose(fmt.Sprintf("GetTokenString Value(%v) Result(%v)\n", val, new_val))
		}
		if err != nil {
			return val, err
		}
		return new_val, nil
	default:
		// User defined types work as well
		return val.(string), nil
	}
}

// GetConfigTokenBool will get the config token from the action and replaces its tokens and converts to a bool
// key - the key of the config
// model - the template data
// Returns the value of the config as a bool
// Returns an error if the key is not found
func (m *Workflow) GetConfigTokenBool(key string, model *TemplateData, required bool) (bool, error) {
	if val, ok := model.CurrentAction.Config[key]; ok {
		//Need to check the type is a string
		value, err := m.GetTokenBool(val, model)
		if err != nil {
			return value, err
		}
		if m.Verbose == LOG_VERBOSE {
			lib.LogVerbose(fmt.Sprintf("GetConfigTokenBool Param(%s) Value(%v) Result(%v)\n", key, val, value))
		}
		return value, nil
	}
	if required {
		return false, fmt.Errorf("required parameter %s missing", key)
	}
	return false, nil
}

// GetTokenBool will parse the token string and convert to a bool
// value - the value to parse
// model - the template data
// Returns the parsed bool
func (m *Workflow) GetTokenBool(value interface{}, model *TemplateData) (bool, error) {
	//Need to check the type is a string
	switch val := value.(type) {
	case string:
		new_val, err := m.ParseToken(model, val)
		if m.Verbose == LOG_VERBOSE {
			lib.LogVerbose(fmt.Sprintf("GetTokenBool Value(%v) Result(%v)\n", val, new_val))
		}
		if err != nil {
			return false, err
		}
		return lib.ConvertToBool(new_val), nil
	case bool:
		return val, nil
	default:
		// User defined types work as well
		return false, nil
	}

}

// GetConfigTokenInt will get the config token from the action and replaces its tokens and converts to an int
// key - the key of the config
// model - the template data
// Returns the value of the config as an int
// Returns an error if the key is not found
func (m *Workflow) GetConfigTokenInt(key string, model *TemplateData, required bool) (int, error) {
	if val, ok := model.CurrentAction.Config[key]; ok {
		value, err := m.GetTokenInt(val, model)
		if err != nil {
			return value, err
		}
		if m.Verbose == LOG_VERBOSE {
			lib.LogVerbose(fmt.Sprintf("GetConfigTokenInt Param(%s) Value(%v) Result(%v)\n", key, val, value))
		}
		return value, nil
	}
	if required {
		return 0, fmt.Errorf("required parameter %s missing", key)
	}
	return 0, nil
}

// GetTokenInt will parse the token string and convert to an int
// value - the value to parse
// model - the template data
// Returns the parsed int
func (m *Workflow) GetTokenInt(value interface{}, model *TemplateData) (int, error) {
	//Need to check the type is a string
	switch val := value.(type) {
	case string:
		new_val, err := m.ParseToken(model, val)
		if err != nil {
			return 0, err
		}
		int_val, err := strconv.Atoi(new_val)
		if err != nil {
			return 0, err
		}
		lib.LogVerbose(fmt.Sprintf("GetTokenInt Value(%v) Result(%v)\n", val, int_val))
		return int_val, nil
	case int:
		return val, nil
	case float64:
		return int(val), nil
	default:
		// User defined types work as well
		return 0, nil
	}
}
func (m *Workflow) GetTemplateFuncMap() template.FuncMap {
	//*********************
	//Create a function map
	//*********************
	funcMap := template.FuncMap{
		"base64enc":   lib.Base64EncString,
		"base64dec":   lib.Base64DecString,
		"gzip_base64": lib.GzipBase64String,
		"lc":          strings.ToLower,
		"uc":          strings.ToUpper,
		"domain":      lib.GetDomainOrIP,
		"port_string": lib.GetPortString,
		"port_int":    lib.GetPortInt,
		"clean":       lib.Clean,
		"concat":      lib.Concat,
		"replace":     strings.ReplaceAll,
		"contains":    lib.CommaListContainsString,
		"not":         lib.NOT,
		"or":          lib.OR,
		"and":         lib.AND,
		"plus":        lib.Plus,
		"minus":       lib.Minus,
		"multiply":    lib.Multiply,
		"divide":      lib.Divide,
		"get_stk_val": m.GetStackVariable,       //get stack variable
		"get_param":   m.GetParamValue,          //Gets the parameter value
		"get_store":   m.GetValueFromDataBucket, //Gets the value from the data bucket
	}
	return funcMap
}

// ParseToken will parse the token string and replace any tokens with the values from the model
// data - the template data
// value - the value to parse
// Returns the parsed string
// Returns an error if the key is not found
func (m *Workflow) ParseToken(data *TemplateData, value string) (string, error) {

	//*****************
	//Pase the template
	//*****************
	if m.templateFuncMap == nil {
		m.templateFuncMap = m.GetTemplateFuncMap()
	}

	tmpl, err := template.New("CodeRun").Funcs(m.templateFuncMap).Parse(value)
	if err != nil {
		return value, err
	}

	//**************************************
	//Run the template to verify the output.
	//**************************************
	var tpl bytes.Buffer
	err = tmpl.Execute(&tpl, data)
	if err != nil {
		return value, err
	}

	return tpl.String(), nil
}
