package workflow

import (
	"reflect"
)

type ValueType int
type ValueOperand int

const (
	TypeInvalid ValueType = iota
	TypeBool
	TypeInt
	TypeFloat
	TypeString
	TypeList
	TypeMap
	TypeSet
	TypeObject
	TypeWorkflow
)

func (v ValueType) String() string {
	return [...]string{"Invalid", "Bool", "Int", "Float", "String", "List", "Map", "Set", "Object", "Workflow"}[v]
}

// Schema is the schema definition for the action plugin
type Schema struct {
	Type        ValueType
	Partial     bool
	Required    bool
	Description string
	ConfigKey   string
	Short       string
	Default     interface{}
	Value       string
}

// WriteConfigFunc is the function definition that needs to be implemented
// to be able to write config values to the config store of choice
// - key: the key to write
// - value: the value to write
// - custom: the custom data that is used to to pass data to the config function
type WriteConfigFunc func(key string, value interface{}, custom ...string) error

// DeleteConfigFunc is the function definition that needs to be implemented
// to be able to delete config values to the config store of choice
// - key: the key to delete
// - custom: the custom data that is used to to pass data to the config function
type DeleteConfigFunc func(key string, custom ...string) error

// ReadConfigFunc is the function definition that needs to be implemented for when you
// want to get config values from the config store of choice
// - key: the key to get
// - data_type: the data type to convert the value to
// - custom: the custom data that is used to to pass data to the config function
// returns the value of the config or nil if the config does not exist
// returns an error
type ReadConfigFunc func(key string, data_type string, custom ...string) (interface{}, error)

// InlineFormatter is used with the action schema to format the inline params
// inline params as passed with the action name e.g action "print;Hello World"
// Normally params are passed as a map[string]interface{} in the action config
// but if the action schema has InlineParams set to true then the params are passed as a string
// - cfg: the config map
// returns the formatted string
type InlineFormatter func(cfg map[string]interface{}) string

// ActionFunc is the function definition that needs to be implemented
// to be able to execute an action, this is called by the workflow
// - w: the workflow
// - m: the template data
// returns an error if the action fails
type ActionFunc func(w *Workflow, m *TemplateData) error

// EventFunc is the function definition is used as part of the workflow event system
// Start and a cleanup function can be added to the workflow and you implement code to handle
// the event
// - w: the workflow
// returns an error if the event fails
type EventFunc func(w *Workflow) error

// TargetMapFunc is the function that is called to map config values to a target type
// - w: the workflow
// - m: the config values
// - target: the target type
type TargetMapFunc func(w *Workflow, m interface{}, target interface{}) (interface{}, error)

// TargetSchema is the schema definition for the target plugin
type TargetSchema struct {
	Action string      //The action that the target is used for e.g action_git
	Short  string      //The short name of the target
	Long   string      //The long name of the target
	Target interface{} //The target type
}

func (m *TargetSchema) GetTargetMap() map[string]interface{} {
	out := make(map[string]interface{})

	v := reflect.ValueOf(m.Target)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct { // Non-structural return error
		return nil
	}

	t := v.Type()
	// Traversing structure fields
	// Specify the tagName value as the key in the map; the field value as the value in the map
	for i := 0; i < v.NumField(); i++ {
		fi := t.Field(i)
		if fi.IsExported() {
			out[fi.Name] = v.Field(i).Interface()
		}
	}
	return out
}

// ActionSchema is the schema definition for the action plugin
type ActionSchema struct {
	Short  string
	Long   string
	Target string //The module and target that the action is for e.g action_git_target
	//Config         map[string]interface{}
	ConfigSchema   map[string]*Schema
	ProcessResults bool            //Uses the process results function to process the results so auto add the results to the template data
	InlineParams   bool            //If true then the params are passed as a string in the action name e.g action: "print;Hello World"
	InlineFormat   InlineFormatter //If InlineParams is true then this is used to format the inline params into a string this will need to be implemented by user
	Action         ActionFunc      //The action function that is called to execute the action
}

type FunctionSchema struct {
	Cmd             string
	Description     string
	Target          string //The module that the function is for e.g action_git
	ParameterSchema map[string]*Schema
	Function        any //The function that is called to execute the function
}

// SchemaEndpoint is the interface that needs to be implemented by the plugin
type SchemaEndpoint interface {
	GetTargetSchema() map[string]TargetSchema  //Get the target schema
	GetActionSchema() map[string]ActionSchema  //Get the action schema
	GetFunctionMap() map[string]FunctionSchema //Get the custom function map for the template
	//GetFunctionMapDocs() map[string]ActionSchema //Get the custom function map for the template
}

// BuildTargetConfig is a helper function to build a target schema
// - short: the short name of the target
// - long: the long name of the target
// - target: the target type
// returns the target schema
func BuildTargetConfig(short string, long string, target interface{}) TargetSchema {
	//Build a test schema
	var targetSchema = TargetSchema{
		//Action: action_target,
		Short:  short,
		Long:   long,
		Target: target,
	}

	return targetSchema
}
