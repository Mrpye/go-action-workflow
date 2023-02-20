package workflow

const (
	INPUT_TYPE_TEXT  = "text"
	INPUT_TYPE_INT   = "int"
	INPUT_TYPE_FLOAT = "float"
	INPUT_TYPE_BOOL  = "bool"
)

// Parameter is a struct that holds the information for a parameter
type Parameter struct {
	Key         string      `json:"key" yaml:"key"`
	Title       string      `json:"title" yaml:"title"`
	Description string      `json:"description" yaml:"description"`
	InputType   string      `json:"type" yaml:"type"`
	Value       interface{} `json:"value" yaml:"value"`
	answer      interface{}
}

type ParameterOption func(*Parameter)

func OptionParameterKey(v string) ParameterOption {
	return func(h *Parameter) {
		h.Key = v
	}
}

func OptionParameterTitle(v string) ParameterOption {
	return func(h *Parameter) {
		h.Title = v
	}
}

func OptionParameterDescription(v string) ParameterOption {
	return func(h *Parameter) {
		h.Description = v
	}
}

func OptionParameterType(v string) ParameterOption {
	return func(h *Parameter) {
		h.InputType = v
	}
}

func OptionParameterValue(v string) ParameterOption {
	return func(h *Parameter) {
		h.Value = v
	}
}
