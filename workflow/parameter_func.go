package workflow

import "fmt"

func (m *Parameter) GetValue() interface{} {
	if m.answer != nil {
		return m.answer
	}
	return m.Value
}

func (m *Parameter) ValueString() string {
	if m.InputType == INPUT_TYPE_TEXT {
		if m.answer != nil {
			return fmt.Sprint(m.answer)
		}
		return fmt.Sprint(m.Value)
	}
	return fmt.Sprint(m.Value)
}

func (m *Parameter) ValueInt() int {
	if m.InputType == INPUT_TYPE_INT {
		if m.answer != nil {
			return m.answer.(int)
		}
		return m.Value.(int)
	}
	return -1
}

func (m *Parameter) ValueBool() bool {
	if m.InputType == INPUT_TYPE_BOOL {
		if m.answer != nil {
			return m.answer.(bool)
		}
		return m.Value.(bool)
	}
	return false
}
func (m *Parameter) ValueFloat() float64 {
	if m.InputType == INPUT_TYPE_FLOAT {
		if m.answer != nil {
			return m.answer.(float64)
		}
		return m.Value.(float64)
	}
	return -1
}

func (m *Parameter) SetAnswer(answer interface{}) {
	m.answer = answer
}
