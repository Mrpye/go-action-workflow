package workflow

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"text/template"

	"github.com/Mrpye/golib/convert"
	"github.com/Mrpye/golib/encrypt"
	"github.com/Mrpye/golib/file"
	"github.com/Mrpye/golib/log"
	"github.com/Mrpye/golib/map_utils"
	"github.com/Mrpye/golib/math"
	"github.com/Mrpye/golib/str"
	"golang.org/x/crypto/bcrypt"
)

// TemplateData is the data that is passed to the template engine
type TemplateData struct {
	Meta          *MetaData
	Manifest      *Manifest
	CurrentAction *Action
	Data          map[string]interface{}
}

// SetAction sets the current action
// current_action - the current action
func (m *TemplateData) SetAction(current_action *Action) {
	m.CurrentAction = current_action
}

// CreateTemplateData creates the template data
// current_action - the current action
// returns the template data
func (m *Workflow) CreateTemplateData(current_action *Action) *TemplateData {
	return &TemplateData{
		Meta:          &m.Manifest.Meta,
		Manifest:      &m.Manifest,
		CurrentAction: current_action,
		Data:          make(map[string]interface{}),
	}
}

// GetStackVariable gets the value of a variable from the stack
// variable_name - the name of the variable
// returns the value of the variable as an int or -1 if the variable is not found
func (m *Workflow) GetStackVariable(variable_name string) int {
	value, err := m.stack.PeekVariable(variable_name)
	if err != nil {
		return -1
	}
	return value.CurrentValue
}

// ParseInterfaceMap will parse the interface map and replace its tokens
// model - the template data
// val - the map to parse
// returns the value of the map with the tokens replaced
// returns an error if the key is not found
func (m *Workflow) ParseInterfaceMap(model *TemplateData, val map[string]interface{}) map[string]interface{} {
	for k, v := range val {
		_ = k
		switch data_val := v.(type) {
		case string:
			parsed_str, err := m.ParseToken(model, string(data_val))
			if err != nil {
				log.PrintFail(fmt.Sprintf("Error ParseToken Value(%v) Error(%v)\n", data_val, err.Error()))
				//log.LogError(fmt.Sprintf("Error ParseToken Value(%v) Error(%v)\n", data_val, err.Error()))
				//log.LogVerbose(fmt.Sprintf("Error ParseToken Value(%v) Error(%v)\n", data_val, err.Error()))
			}
			if m.LogLevel == LOG_VERBOSE {
				log.LogVerbose(fmt.Sprintf("ParseToken Value(%v) Result(%v)\n", data_val, parsed_str))
			}
			val[k] = parsed_str
		case map[string]interface{}:
			m.ParseInterfaceMap(model, data_val)
		}
	}
	return val
}

// GetConfigTokenMap will get the config token map from the action and replaces its tokens
// key - the key of the config
// model - the template data
// required - if the key is required
// returns the value of the map with the tokens replaced
// returns an error if the key is not found
func (m *Workflow) GetConfigTokenMap(key string, model *TemplateData, required bool) (map[string]interface{}, error) {
	if val, ok := model.CurrentAction.Config[key]; ok {
		switch val := val.(type) {
		case string:
			//************************************************
			//See if we are getting data from the data of meta
			//************************************************
			if strings.HasPrefix(val, "$") {
				//The we read the data
				result, err := m.GetDataFromString(val)
				if err != nil {
					return nil, err
				}
				//**************
				//Convert to map
				//**************
				switch result_data := result.(type) {
				case map[string]interface{}:
					return result_data, nil

				case []interface{}:
					data := make(map[string]interface{})
					for i, v := range result_data {
						data[strconv.Itoa(i)] = v
					}
					return data, nil
				}
			}
			// User defined types work as well
			return nil, nil
		case map[string]interface{}:
			new_val := map_utils.CopyMap(val)
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

// GetConfigTokenMap will get the config token map from the action and replaces its tokens
// key - the key of the config
// model - the template data
// required - if the key is required
// returns the value of the map with the tokens replaced
// returns an error if the key is not found
func (m *Workflow) GetConfigTokenMapArray(key string, model *TemplateData, required bool) ([]map[string]interface{}, error) {
	if val, ok := model.CurrentAction.Config[key]; ok {
		switch val := val.(type) {

		case []interface{}:
			data := make([]map[string]interface{}, len(val))
			for i, v := range val {
				new_val := map_utils.CopyMap(v.(map[string]interface{}))
				new_val = m.ParseInterfaceMap(model, new_val)
				data[i] = new_val
			}
			return data, nil
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

// GetConfigTokenMap will get the config token map from the action and replaces its tokens
// key - the key of the config
// model - the template data
// required - if the key is required
// returns the value of the map with the tokens replaced
// returns an error if the key is not found
func (m *Workflow) GetConfigTokenInterfaceArray(key string, model *TemplateData, required bool) ([]interface{}, error) {
	if val, ok := model.CurrentAction.Config[key]; ok {
		switch val := val.(type) {
		case string:
			//************************************************
			//See if we are getting data from the data of meta
			//************************************************
			if strings.HasPrefix(val, "$") {
				//The we read the data
				result, err := m.GetDataFromString(val)
				if err != nil {
					return nil, err
				}
				//****************
				//Convert to array
				//****************
				switch result_data := result.(type) {
				case map[string]interface{}:
					data := make([]interface{}, len(result_data))
					count := 0
					for _, v := range result_data {
						data[count] = v
						count++
					}
					return data, nil
				case []interface{}:
					return result_data, nil
				}
			}
		case []interface{}:
			data := make([]interface{}, len(val))
			for i, v := range val {
				switch v := v.(type) {
				case string:
					new_val, err := m.ParseToken(model, v)
					if m.LogLevel == LOG_VERBOSE {
						log.LogVerbose(fmt.Sprintf("GetTokenString Value(%v) Result(%v)\n", v, new_val))
					}
					if err != nil {
						return nil, err
					}
					data[i] = new_val
				default:
					data[i] = v
				}
			}
			return data, nil
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

// GetConfigTokenMap will get the config token map from the action and replaces its tokens
// key - the key of the config
// model - the template data
// required - if the key is required
// returns the value of the map with the tokens replaced
// returns an error if the key is not found
func (m *Workflow) GetConfigTokenStringArray(key string, model *TemplateData, required bool) ([]string, error) {
	if val, ok := model.CurrentAction.Config[key]; ok {
		switch val := val.(type) {
		case string:
			//************************************************
			//See if we are getting data from the data of meta
			//************************************************
			if strings.HasPrefix(val, "$") {
				//The we read the data
				result, err := m.GetDataFromString(val)
				if err != nil {
					return nil, err
				}
				//****************
				//Convert to array
				//****************
				switch result_data := result.(type) {
				case map[string]interface{}:
					data := make([]string, len(result_data))
					count := 0
					for _, v := range result_data {
						new_val, err := m.ParseToken(model, fmt.Sprintf("%v", v))
						if m.LogLevel == LOG_VERBOSE {
							log.LogVerbose(fmt.Sprintf("GetTokenString Value(%v) Result(%v)\n", v, new_val))
						}
						if err != nil {
							return nil, err
						}
						data[count] = new_val
						count++
					}
					return data, nil
				case []interface{}:
					data := make([]string, len(result_data))
					count := 0
					for _, v := range result_data {
						new_val, err := m.ParseToken(model, fmt.Sprintf("%v", v))
						if m.LogLevel == LOG_VERBOSE {
							log.LogVerbose(fmt.Sprintf("GetTokenString Value(%v) Result(%v)\n", v, new_val))
						}
						if err != nil {
							return nil, err
						}
						data[count] = new_val
						count++
					}
					return data, nil
				}
			}
		case []interface{}:
			data := make([]string, len(val))
			for i, v := range val {
				switch v := v.(type) {
				case string:
					new_val, err := m.ParseToken(model, v)
					if m.LogLevel == LOG_VERBOSE {
						log.LogVerbose(fmt.Sprintf("GetTokenString Value(%v) Result(%v)\n", v, new_val))
					}
					if err != nil {
						return nil, err
					}
					data[i] = new_val
				default:
					new_val, err := m.ParseToken(model, fmt.Sprintf("%v", v))
					if m.LogLevel == LOG_VERBOSE {
						log.LogVerbose(fmt.Sprintf("GetTokenString Value(%v) Result(%v)\n", v, new_val))
					}
					if err != nil {
						return nil, err
					}
					data[i] = new_val
				}
			}
			return data, nil
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
// Returns the value of the config as an interface
// Returns an error if the key is not found
func (m *Workflow) GetConfigToken(key string, model *TemplateData) (interface{}, error) {
	if val, ok := model.CurrentAction.Config[key]; ok {
		//Need to check the type is a string
		switch val := val.(type) {
		case string:
			new_val, err := m.ParseToken(model, val)
			if m.LogLevel == LOG_VERBOSE {
				log.LogVerbose(fmt.Sprintf("GetConfigToken Value(%v) Result(%v)\n", val, new_val))
			}
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
// required - if the key should exist and a value should be returned
// Returns the value of the config as an interface
// Returns an error if the key is not found
func (m *Workflow) GetConfigTokenInterface(key string, model *TemplateData, required bool) (interface{}, error) {
	if value, ok := model.CurrentAction.Config[key]; ok {
		switch val := value.(type) {
		case string:
			new_val, err := m.ParseToken(model, val)
			if m.LogLevel == LOG_VERBOSE {
				log.LogVerbose(fmt.Sprintf("GetConfigTokenInterface Value(%v) Result(%v)\n", val, new_val))
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
					if m.LogLevel == LOG_VERBOSE {
						log.LogVerbose(fmt.Sprintf("GetTokenString Value(%v) Result(%v)\n", v, new_val))
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
// required - if the key should exist and a value should be returned
// Returns the value of the config as a string
// Returns an error if the key is not found
func (m *Workflow) GetConfigTokenString(key string, model *TemplateData, required bool) (string, error) {
	if model == nil {
		return "", fmt.Errorf("model is nil")
	}
	if val, ok := model.CurrentAction.Config[key]; ok {
		value, err := m.GetTokenString(val, model)
		if err != nil {
			return value, fmt.Errorf("error with key %s, err:%v", key, err)
		}
		if m.LogLevel == LOG_VERBOSE {
			log.LogVerbose(fmt.Sprintf("GetConfigTokenString Param(%s) Value(%s) Result(%s)\n", key, val, value))
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
// Returns an error if the value is not a string
func (m *Workflow) GetTokenString(value interface{}, model *TemplateData) (string, error) {

	//Need to check the type is a string
	switch val := value.(type) {
	case string:
		new_val, err := m.ParseToken(model, val)
		if m.LogLevel == LOG_VERBOSE {
			log.LogVerbose(fmt.Sprintf("GetTokenString Value(%v) Result(%v)\n", val, new_val))
		}
		if err != nil {
			return val, err
		}
		return fmt.Sprintf("%v", new_val), nil
	default:
		// User defined types work as well
		return val.(string), nil
	}
}

// GetConfigTokenBool will get the config token from the action and replaces its tokens and converts to a bool
// key - the key of the config
// model - the template data
// required - if the key should exist and a value should be returned
// Returns the value of the config as a bool
// Returns an error if the key is not found
func (m *Workflow) GetConfigTokenBool(key string, model *TemplateData, required bool) (bool, error) {
	if model == nil {
		return false, fmt.Errorf("model is nil")
	}
	if val, ok := model.CurrentAction.Config[key]; ok {
		//Need to check the type is a string
		value, err := m.GetTokenBool(val, model)
		if err != nil {
			return value, err
		}
		if m.LogLevel == LOG_VERBOSE {
			log.LogVerbose(fmt.Sprintf("GetConfigTokenBool Param(%s) Value(%v) Result(%v)\n", key, val, value))
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
// Returns an error if the value is not a bool
func (m *Workflow) GetTokenBool(value interface{}, model *TemplateData) (bool, error) {
	//Need to check the type is a string
	switch val := value.(type) {
	case string:
		new_val, err := m.ParseToken(model, val)
		if m.LogLevel == LOG_VERBOSE {
			log.LogVerbose(fmt.Sprintf("GetTokenBool Value(%v) Result(%v)\n", val, new_val))
		}
		if err != nil {
			return false, err
		}
		return convert.ToBool(new_val), nil
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
// required - if the key should exist and a value should be returned
// Returns the value of the config as an int
// Returns an error if the key is not found
func (m *Workflow) GetConfigTokenInt(key string, model *TemplateData, required bool) (int, error) {
	if model == nil {
		return -1, fmt.Errorf("model is nil")
	}
	if val, ok := model.CurrentAction.Config[key]; ok {
		value, err := m.GetTokenInt(val, model)
		if err != nil {
			return value, err
		}
		if m.LogLevel == LOG_VERBOSE {
			log.LogVerbose(fmt.Sprintf("GetConfigTokenInt Param(%s) Value(%v) Result(%v)\n", key, val, value))
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
// Returns an error if the value is not an int
func (m *Workflow) GetTokenInt(value interface{}, model *TemplateData) (int, error) {
	//Need to check the type is a string
	switch val := value.(type) {
	case string:
		new_val, err := m.ParseToken(model, val)
		if err != nil {
			return 0, err
		}
		int_val, err := strconv.Atoi(fmt.Sprintf("%v", new_val))
		if err != nil {
			return 0, err
		}
		if m.LogLevel == LOG_VERBOSE {
			log.LogVerbose(fmt.Sprintf("GetTokenInt Value(%v) Result(%v)\n", val, int_val))
		}
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

// GetConfigTokenInt will get the config token from the action and replaces its tokens and converts to an int
// key - the key of the config
// model - the template data
// required - if the key should exist and a value should be returned
// Returns the value of the config as an int
// Returns an error if the key is not found
func (m *Workflow) GetConfigTokenFloat(key string, model *TemplateData, required bool) (float64, error) {
	if model == nil {
		return -1, fmt.Errorf("model is nil")
	}
	if val, ok := model.CurrentAction.Config[key]; ok {
		value, err := m.GetTokenFloat64(val, model)
		if err != nil {
			return value, err
		}
		if m.LogLevel == LOG_VERBOSE {
			log.LogVerbose(fmt.Sprintf("GetConfigTokenInt Param(%s) Value(%v) Result(%v)\n", key, val, value))
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
// Returns an error if the value is not an int
func (m *Workflow) GetTokenFloat64(value interface{}, model *TemplateData) (float64, error) {
	//Need to check the type is a string
	switch val := value.(type) {
	case string:
		new_val, err := m.ParseToken(model, val)
		if err != nil {
			return 0, err
		}
		int_val, err := strconv.ParseFloat(fmt.Sprintf("%v", new_val), 64)
		if err != nil {
			return 0, err
		}
		if m.LogLevel == LOG_VERBOSE {
			log.LogVerbose(fmt.Sprintf("GetTokenInt Value(%v) Result(%v)\n", val, int_val))
		}
		return int_val, nil
	case int:
		return float64(val), nil
	case float64:
		return val, nil
	default:
		// User defined types work as well
		return 0, nil
	}
}

// GetWorkflow will return the workflow
func (m *Workflow) GetWorkflow() *Workflow {
	return m
}

// GenerateToken will generate a token from a password
// password - the password to generate the token from
// Returns the token
// Returns an error if the token could not be generated
func (m *Workflow) GenerateToken(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(hash), nil
}

// FileTemplate will read a file and parse the tokens
// file - the file to read
// custom - the custom data to use
// Returns the parsed string
// Returns an error if the file could not be read
func (m *Workflow) FileTemplate(file string, custom ...interface{}) (string, error) {
	//read the file
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}
	return m.Template(string(data), custom...)
}

// TemplateMap will parse the token string and replace the tokens with the custom data
// val - the value to parse
// custom - the custom data to use
// Returns the parsed string
// Returns an error if the value could not be parsed
func (m *Workflow) TemplateMap(val string, custom map[string]interface{}) (string, error) {
	m.model.Data = custom

	new_val, err := m.ParseToken(m.model, val)
	if err != nil {
		return "", err
	}
	return new_val, nil

}

// Template will parse the token string and replace the tokens with the custom data
// val - the value to parse
// custom - the custom data to use as a string `key=value` or a map
// Returns the parsed string
// Returns an error if the value could not be parsed
func (m *Workflow) Template(val string, custom ...interface{}) (string, error) {
	model_data := make(map[string]interface{})
	for _, value := range custom {
		switch v := value.(type) {
		case string:
			parts := strings.Split(v, "=")
			if len(parts) > 1 {
				model_data[parts[0]] = parts[1]
			}
		case map[string]interface{}:
			return m.TemplateMap(val, v)
		}
	}
	return m.TemplateMap(val, model_data)
}

// BuildStoreMap will build a map of values to store in the data bucket or from key/value pairs as string `key=value`
// key - the key of the data bucket
// custom - the custom values to add to the map
// Returns the map of values
func (m *Workflow) BuildStoreMap(key string, custom ...string) map[string]interface{} {
	model_data := make(map[string]interface{})

	for _, v := range custom {
		if strings.Contains(v, "=") {
			parts := strings.Split(v, "=")
			if len(parts) > 1 {
				model_data[parts[0]] = parts[1]
			}
		} else {
			if key != "" {
				value := m.GetValueFromDataBucket(key, v)
				model_data[v] = value
			}
		}
	}
	return model_data
}

// KeyPair will return a key pair when string is passed as `key=value`
// key - the key
// value - the value
// Returns the key pair
func (m *Workflow) KeyPair(key string, value interface{}) string {
	return fmt.Sprintf("%s=%v", key, value)
}

// KeyPairMap will return a map of the key pairs when string is passed as `key=value`
// custom - the key pairs
func (m *Workflow) KeyPairMap(custom ...string) map[string]interface{} {
	model_data := make(map[string]interface{})
	for _, v := range custom {
		if strings.Contains(v, "=") {
			parts := strings.Split(v, "=")
			if len(parts) > 1 {
				model_data[parts[0]] = parts[1]
			}
		}
	}
	return model_data
}

// GetTemplateFuncMap will return the template function map
func (m *Workflow) GetInbuiltTemplateFuncMap() template.FuncMap {
	//*********************
	//Create a function map
	//*********************
	funcMap := template.FuncMap{
		"kp":          m.KeyPair,                   //Create a key pair string
		"kps":         m.KeyPairMap,                //Create a key pair map
		"kps_store":   m.BuildStoreMap,             //Create a key pair map from the data bucket
		"gen_token":   m.GenerateToken,             //Generate a token
		"tpl":         m.Template,                  //Render a template
		"tpl_file":    m.FileTemplate,              //Read file and parse template a string
		"read_file":   file.ReadFileToString,       //Read a file to a string
		"file":        file.ReadFile,               //Read a file to a  []byte
		"base64enc":   encrypt.Base64EncString,     //Base64 encode a string
		"base64dec":   encrypt.Base64DecString,     //Base64 decode a string
		"gzip_base64": str.GzipBase64,              //Gzip and Base64 encode a string
		"lc":          strings.ToLower,             //Lowercase a string
		"uc":          strings.ToUpper,             //Uppercase a string
		"domain":      str.GetDomainOrIP,           //Get the domain or IP from a string
		"port_string": str.GetPortString,           //Get the port from a string
		"port_int":    str.GetPortInt,              // Get the port from a string as an int
		"clean":       str.Clean,                   //Clean a string removing spaces and special characters
		"concat":      str.Concat,                  //Concatenate a two or more strings
		"replace":     strings.ReplaceAll,          //Replace a value in a string
		"contains":    str.CommaListContainsString, //Check if a string is in a comma separated list
		"not":         math.NOT,                    //Not a boolean
		"or":          math.OR,                     //Or two booleans
		"and":         math.AND,                    //And two booleans
		"plus":        math.Plus,                   //Add two integers
		"minus":       math.Minus,                  //Subtract two integers
		"multiply":    math.Multiply,               //Multiply two integers
		"divide":      math.Divide,                 //Divide two integers
		"get_stk_val": m.GetStackVariable,          //get stack variable
		"get_param":   m.GetParamValue,             //Gets the parameter value
		"get_input":   m.GetInputValue,             //Gets the input value
		"get_store":   m.GetValueFromDataBucket,    //Gets the value from the data bucket
		"get_data":    m.GetDataItem,               //Get the custom item
		"get_config":  m.GetConfigValue,            //Gets the config value
		"get_wf":      m.GetWorkflow,               //Gets the workflow
		"build_path":  m.BuildPath,                 //Build a path from a list of strings
		"count_array": m.CountArray,                //Build a path from a list of strings
	}
	return funcMap
}

// ParseToken will parse the token string and replace any tokens with the values from the model
// data - the template data
// value - the value to parse
// Returns the parsed string
// Returns an error if the key is not found
func (m *Workflow) ParseToken(data *TemplateData, value string) (string, error) {

	//****************************************
	//Create the func map if it does not exist
	//****************************************
	if m.templateFuncMap == nil {
		m.templateFuncMap = m.GetInbuiltTemplateFuncMap()
	}

	//********************************
	//Create a new template and parse
	//********************************
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

	//******************
	//Return the result
	//******************
	return tpl.String(), nil
}
