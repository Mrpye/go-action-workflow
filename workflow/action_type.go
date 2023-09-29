package workflow

// Action is the struct that defines an action
type Action struct {
	Key         string `json:"key,omitempty" yaml:"key,omitempty" flag:"key k" desc:"key for the job or global action"`
	Description string `json:"description,omitempty" yaml:"description,omitempty" flag:"desc" desc:"description for the job"`
	Action      string `json:"action" yaml:"action" flag:"-"`
	//Label           string                 `json:"label" yaml:"label"`
	Fail            string                 `json:"fail,omitempty" yaml:"fail,omitempty" flag:"fail" desc:"The action to run if this action fails"`
	ContinueOnError interface{}            `json:"continue_on_error,omitempty" yaml:"continue_on_error,omitempty" flag:"continue_on_error" desc:"Ignore errors and continue to the next action (true/false) or template token to evaluate to true/false"`
	Config          map[string]interface{} `json:"config,omitempty" yaml:"config,omitempty"`
	Disabled        interface{}            `json:"disabled,omitempty" yaml:"disabled,omitempty" flag:"disable" desc:"Disable this action (true/false) or template token to evaluate to true/false"`
}

// ActionOption is a function that sets an option on the action
type ActionOption func(*Action)

// OptionActionDisabled sets the disabled flag on the action
func OptionActionDisabled(v bool) ActionOption {
	return func(h *Action) {
		h.Disabled = v
	}
}

// OptionActionFail sets the fail flag on the action
func OptionActionFail(v string) ActionOption {
	return func(h *Action) {
		h.Fail = v
	}
}

// OptionActionKey sets the key on the action
func OptionActionKey(v string) ActionOption {
	return func(h *Action) {
		h.Key = v
	}
}

//	OptionActionDescription sets the description on the action
func OptionActionDescription(v string) ActionOption {
	return func(h *Action) {
		h.Description = v
	}
}

// OptionActionAction sets the action on the action
func OptionActionAction(v string) ActionOption {
	return func(h *Action) {
		h.Action = v
	}
}

// OptionActionContinueOnError sets the continue on error flag on the action
func OptionActionContinueOnError(v interface{}) ActionOption {
	return func(h *Action) {
		h.ContinueOnError = v
	}
}

// OptionActionConfig sets the config on the action
func OptionActionConfig(v map[string]interface{}) ActionOption {
	return func(h *Action) {
		h.Config = v
	}
}

// CreateAction creates a new action
func CreateAction(opts ...ActionOption) *Action {
	action := &Action{}
	action.Config = make(map[string]interface{})
	for _, opt := range opts {
		opt(action)
	}
	return action
}
