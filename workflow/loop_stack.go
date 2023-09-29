package workflow

import (
	"errors"
	"fmt"
)

// Loop is used to track the current loop
type loop struct {
	VariableName string
	Index        int
	From         int
	To           int
	CurrentValue int
	Dec          bool
	Exceeded     bool
}

// loopStack is a stack of loops
type loopStack struct {
	stackList []loop
}

// CreateLoopStack creates a new loop stack
// returns the loop stack
func CreateLoopStack() loopStack {
	workflow := loopStack{}
	return workflow
}

// Increment increments the current variable of the loop
// returns true if the loop has exceeded the to value
// returns an error if the stack is empty
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

// Push pushes a new loop onto the stack
// - variable_name is the name of the variable
// - index is the index of action in the job
// - from is the from value of the loop
// - to is the to value of the loop
// returns the loop
// returns an error if the variable name is empty
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
		m.stackList = append(m.stackList, loop{VariableName: variable_name, Index: index, From: from, To: to, CurrentValue: from, Dec: false})
	} else {
		m.stackList = append(m.stackList, loop{VariableName: variable_name, Index: index, From: from, To: to, CurrentValue: from, Dec: true})
	}
	return &m.stackList[len(m.stackList)-1], nil
}

// Pop pops the top loop off the stack
// returns the loop
// returns an error if the stack is empty
func (m *loopStack) Pop() (*loop, error) {
	if len(m.stackList) > 0 {
		temp_loop := m.stackList[len(m.stackList)-1]
		m.stackList = m.stackList[:len(m.stackList)-1]
		return &temp_loop, nil
	}
	return nil, errors.New("the stack is empty")
}

// Peek peeks at the top loop on the stack
// returns the loop
// returns an error if the stack is empty
func (m *loopStack) Peek() (*loop, error) {
	if len(m.stackList) > 0 {
		return &m.stackList[len(m.stackList)-1], nil
	}
	return nil, errors.New("the stack is empty")
}

// PeekVariable peeks at the loop on the stack with the specified variable name
// - variable_name is the name of the variable
// returns the loop
// returns an error if the variable name does not exist
func (m *loopStack) PeekVariable(variable_name string) (*loop, error) {
	for _, l := range m.stackList {
		if l.VariableName == variable_name {
			return &l, nil
		}
	}
	return nil, fmt.Errorf("variable %s does not exist exists", variable_name)
}
