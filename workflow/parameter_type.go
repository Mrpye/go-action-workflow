package workflow

// Parameter data types
const (
	TYPE_STRING = "string"
	TYPE_INT    = "int"
	TYPE_FLOAT  = "float"
	TYPE_BOOL   = "bool"
)

// Parameter is a struct that holds the information for a parameter
type Parameter struct {
	Key         string      `json:"key" yaml:"key" flag:"key k" desc:"key for the parameter, must be unique within the workflow, used as the parameter name"`
	Title       string      `json:"title" yaml:"title" flag:"title t" desc:"title for the parameter, used as the parameter title"`
	Description string      `json:"description" yaml:"description" flag:"desc d" desc:"description for the parameter, used as the parameter description"`
	InputType   string      `json:"type" yaml:"type" flag:"type y" desc:"type of the parameter, must be one of string, int, float, bool"`
	Value       interface{} `json:"value" yaml:"value" flag:"value v" desc:"value of the parameter, must be a string, int, float, or bool"`
	answer      interface{}
}

// ParameterOption is a function that sets a parameter option
type ParameterOption func(*Parameter)

// OptionParameterKey sets the key of the parameter
// - v is the key of the parameter
func OptionParameterKey(v string) ParameterOption {
	return func(h *Parameter) {
		h.Key = v
	}
}

// OptionParameterTitle sets the title of the parameter
// - v is the title of the parameter
func OptionParameterTitle(v string) ParameterOption {
	return func(h *Parameter) {
		h.Title = v
	}
}

// OptionParameterDescription sets the description of the parameter
// - v is the description of the parameter
func OptionParameterDescription(v string) ParameterOption {
	return func(h *Parameter) {
		h.Description = v
	}
}

// OptionParameterType sets the type of the parameter
// - v is the type of the parameter
func OptionParameterType(v string) ParameterOption {
	return func(h *Parameter) {
		h.InputType = v
	}
}

// OptionParameterValue sets the value of the parameter
// - v is the value of the parameter
func OptionParameterValue(v string) ParameterOption {
	return func(h *Parameter) {
		h.Value = v
	}
}

func CreateParameter(opts ...ParameterOption) *Parameter {
	parameter := &Parameter{}
	for _, opt := range opts {
		opt(parameter)
	}
	return parameter
}
