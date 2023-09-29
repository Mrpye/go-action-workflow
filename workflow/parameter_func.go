package workflow

import (
	"github.com/Mrpye/golib/convert"
)

// GetValue returns the value of the parameter
// returns the value of the parameter
func (m *Parameter) GetValue() interface{} {
	if m.answer != nil {
		return m.answer
	}
	return m.Value
}

// ValueString returns the value of the parameter as a string
// returns the value of the parameter as a string
func (m *Parameter) ValueString() string {
	if m.answer != nil {
		return convert.ToString(m.answer)
	}
	return convert.ToString(m.Value)
}

// ValueInt returns the value of the parameter as an int
// returns the value of the parameter as an int
func (m *Parameter) ValueInt() int {
	if m.answer != nil {
		return convert.ToInt(m.answer)
	}
	return convert.ToInt(m.Value)
}

// ValueBool returns the value of the parameter as a bool
// returns the value of the parameter as a bool
func (m *Parameter) ValueBool() bool {
	if m.answer != nil {
		return convert.ToBool(m.answer)
	}
	return convert.ToBool(m.Value)
}

// ValueFloat returns the value of the parameter as a float64
// returns the value of the parameter as a float64
func (m *Parameter) ValueFloat() float64 {
	if m.answer != nil {
		return convert.ToFloat64(m.answer)
	}
	return convert.ToFloat64(m.Value)
}

// SetAnswer sets the answer of the parameter
// - answer is the answer to set
func (m *Parameter) SetAnswer(answer interface{}) {
	m.answer = answer
}
