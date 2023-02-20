package workflow

import (
	"errors"
	"fmt"
)

type loop struct {
	VariableName string
	Index        int
	From         int
	To           int
	CurrentValue int
	Dec          bool
	Exceeded     bool
}

type loopStack struct {
	stackList []loop
}

func (m *loopStack) Increment() (bool, error) {
	if len(m.stackList) > 0 {
		index := len(m.stackList) - 1
		if m.stackList[index].Dec {
			m.stackList[index].CurrentValue = m.stackList[index].CurrentValue - 1
			if m.stackList[index].CurrentValue < m.stackList[index].To {
				m.stackList[index].Exceeded = true
				return true, nil
			}
		} else {
			m.stackList[index].CurrentValue = m.stackList[index].CurrentValue + 1
			if m.stackList[index].CurrentValue > m.stackList[index].To {
				m.stackList[index].Exceeded = true
				return true, nil
			}
		}

		return false, nil

	}
	return false, errors.New("the stack is empty")
}

func (m *loopStack) Push(variable_name string, index int, from int, to int) (*loop, error) {
	if variable_name == "" {
		return nil, errors.New("variable name cannot be empty")
	}
	for _, l := range m.stackList {
		if l.VariableName == variable_name {
			return &m.stackList[len(m.stackList)-1], nil
		}
	}
	if from < to {
		m.stackList = append(m.stackList, loop{VariableName: variable_name, Index: index, From: from, To: to, CurrentValue: 0, Dec: false})
	} else {
		m.stackList = append(m.stackList, loop{VariableName: variable_name, Index: index, From: from, To: to, CurrentValue: 0, Dec: false})
	}
	return &m.stackList[len(m.stackList)-1], nil
}

func (m *loopStack) Pop() (*loop, error) {
	if len(m.stackList) > 0 {
		temp_loop := m.stackList[len(m.stackList)-1]
		m.stackList = m.stackList[:len(m.stackList)-1]
		return &temp_loop, nil
	}
	return nil, errors.New("the stack is empty")
}
func (m *loopStack) Peek() (*loop, error) {
	if len(m.stackList) > 0 {
		return &m.stackList[len(m.stackList)-1], nil
	}
	return nil, errors.New("the stack is empty")
}
func (m *loopStack) PeekVariable(variable_name string) (*loop, error) {
	for _, l := range m.stackList {
		if l.VariableName == variable_name {
			return &l, nil
		}
	}
	return nil, fmt.Errorf("variable %s does not exist exists", variable_name)
}
