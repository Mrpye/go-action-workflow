package workflow

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Mrpye/golib/lib"
	"github.com/gookit/color"
)

// *****************************************
//RunJob will run the job with the given key
// key - the key of the job to run
// returns an error if there is one
// *****************************************
func (m *Workflow) RunJob(key string) error {

	// *******************************
	// Run the job with the given key.
	// *******************************
	err := m.executeJob(key)

	//**************************************
	//Initialize the the actions and targets
	//**************************************
	if m.CleanFunc != nil {
		err := m.CleanFunc(m)
		if err != nil {
			return err
		}
	}
	return err
}

// **********************************************
// executeJob will run the job with the given key
// key - the key of the job to run
// returns an error if there is one
// **********************************************
func (m *Workflow) executeJob(key string) error {

	//**********************************************
	//Lets get the Job for the selected app profile
	//**********************************************
	Job := m.Manifest.GetJob(key)
	if Job == nil {
		return fmt.Errorf("cannot find job %s", key)
	}

	//**************************************
	//Initialize the the actions and targets
	//**************************************
	if m.InitFunc != nil {
		err := m.InitFunc(m)
		if err != nil {
			return err
		}
	}

	//********************************
	//Create our loop stack
	//*NOTE* This is a global variable
	//********************************
	m.stack = loopStack{}

	//**************
	//create a Model
	//**************
	m.Model = m.CreateTemplateData(nil)

	//********************
	//Setup some variables
	//********************
	Job_action_count := len(Job.Actions)
	current_item := 0

	//*********************
	//Loop over the actions
	//*********************
	for current_item = 0; current_item < Job_action_count; current_item++ {
		current_action := &Job.Actions[current_item]
		//********************************************
		//Set the current action to our template Model
		//********************************************
		m.Model.SetAction(current_action)

		//********************************
		//See if this action is disabled
		//********************************
		is_disabled, err := m.GetTokenBool(current_action.Disabled, m.Model)
		if err != nil {
			return err
		}
		if is_disabled {
			if m.Verbose > LOG_QUIET {
				color := color.FgRed.Render
				lib.ActionLog(fmt.Sprintf("Action Ignored: %s->%s : %s", current_action.Key, current_action.Action, color("Disabled")), '*')
			}
			continue
		}
		if m.Verbose > LOG_QUIET {
			lib.ActionLog(fmt.Sprintf("Action: %s->%s", current_action.Key, current_action.Action), '*')
		}
		//********************************
		//Parse any variable in the action
		//********************************
		lowercase_action, err := m.ParseToken(m.Model, current_action.Action)
		action_parts := strings.Split(lowercase_action, ";")
		action_parts[0] = strings.ToLower(action_parts[0])
		if err != nil {
			return err
		}

		//**************************
		//Find the action to perform
		//**************************
		switch action_parts[0] {
		case "end":
			//End the job
			if m.Verbose > LOG_QUIET {
				lib.ActionLog("End", '*')
			}
			return nil
		case "print":
			fmt.Println(action_parts[1])
			if m.Verbose > LOG_QUIET {
				lib.ActionLogOK(fmt.Sprintf("Action Completed: %s->%s", current_action.Key, current_action.Action), '-')
			}
		case "goto":
			if len(action_parts) < 2 {
				return errors.New("not enough args for goto")
			}
			if m.Verbose > LOG_QUIET {
				lib.ActionLog("goto: "+action_parts[1], '>')
			}
			index := Job.GetKeyIndex(action_parts[1])
			if index == -1 {
				return fmt.Errorf("cannot find label %s", action_parts[1])
			}
			current_item = index - 1
		case "wait-seconds", "wait":
			if len(action_parts) < 2 {
				return errors.New("not enough args for wait-seconds wait ")
			}
			if m.Verbose > LOG_QUIET {
				lib.ActionLog("wait-seconds/wait: "+action_parts[1], '*')
			}
			count, err := strconv.Atoi(action_parts[1])
			if err != nil {
				return fmt.Errorf("value should be an int for wait-seconds wait %s", action_parts[1])
			}
			time.Sleep(time.Duration(count) * time.Second)
			if m.Verbose > LOG_QUIET {
				lib.ActionLogOK(fmt.Sprintf("Action Completed: %s->%s", current_action.Key, current_action.Action), '-')
			}
		case "wait-minutes":
			if len(action_parts) < 2 {
				return errors.New("not enough args for wait-minutes")
			}
			if m.Verbose > LOG_QUIET {
				lib.ActionLog("wait-minutes: "+action_parts[1], '*')
			}
			count, err := strconv.Atoi(action_parts[1])
			if err != nil {
				return fmt.Errorf("value should be an int for wait-minutes %s", action_parts[1])
			}
			time.Sleep(time.Duration(count) * time.Minute)
			if m.Verbose > LOG_QUIET {
				lib.ActionLogOK(fmt.Sprintf("Action Completed: %s->%s", current_action.Key, current_action.Action), '-')
			}
		case "for":
			//this will loop a section x times
			if len(action_parts) < 2 {
				return errors.New("not enough args for loop")
			}
			from, err := strconv.Atoi(action_parts[2])
			if err != nil {
				return fmt.Errorf("value should be an int for loop %s", action_parts[2])
			}
			to, err := strconv.Atoi(action_parts[3])
			if err != nil {
				return fmt.Errorf("value should be an int for loop %s", action_parts[3])
			}
			temp_loop, err := m.stack.Push(action_parts[1], current_item, from, to)
			if err != nil {
				return err
			}
			if m.Verbose > LOG_QUIET {
				lib.ActionLog("loop: "+action_parts[1]+"["+strconv.FormatInt(int64(temp_loop.CurrentValue), 10)+"]", '>')
			}
		case "next":
			//Get the item off the stack
			temp_loop, err := m.stack.Peek()
			if err != nil {
				return err
			}
			if m.Verbose > LOG_QUIET {
				lib.ActionLog("loop-end: "+temp_loop.VariableName+"["+strconv.FormatInt(int64(temp_loop.CurrentValue), 10)+"]", '<')
			}
			//Increment the loop counter
			inc_result, err := m.stack.Increment()
			if err != nil {
				return err
			}

			//See if this is the end of the loop
			if inc_result {
				m.stack.Pop()
			} else {
				current_item = temp_loop.Index - 1
			}
		default:

			continue_on_error := lib.ConvertToBool(current_action.ContinueOnError)
			var err error
			if val, ok := m.ActionList[current_action.Action]; ok {
				err = val(m)
			}
			if err != nil {
				if continue_on_error {
					lib.PrintlnFail("!! There was an error but Continue On Error was set to true !!")
					continue
				}
				//************************************
				//lets see if to end or goto an action
				//************************************
				lowercase_fail, err := m.ParseToken(m.Model, current_action.Fail)
				fail_parts := strings.Split(lowercase_fail, ";")
				fail_parts[0] = strings.ToLower(fail_parts[0])
				if fail_parts[0] == "end" || fail_parts[0] == "" {
					return err
				} else if fail_parts[0] == "goto" {
					lib.PrintlnFail("!! There but goto is set for fail !!")

					if m.Verbose > LOG_QUIET {
						lib.LogVerbose(fmt.Sprintf("The follow error occurred for action: %s->%s:%s", current_action.Key, current_action.Action, err.Error()))
					}

					if len(fail_parts) < 2 {
						return errors.New("not enough args for goto")
					}
					if m.Verbose > LOG_QUIET {
						lib.ActionLog("goto: "+fail_parts[1], '>')
					}
					index := Job.GetKeyIndex(fail_parts[1])
					if index == -1 {
						return fmt.Errorf("cannot find label %s", fail_parts[1])
					}
					current_item = index - 1
				}

			}
			if m.Verbose > LOG_QUIET {
				lib.ActionLogOK(fmt.Sprintf("Action Completed: %s->%s", current_action.Key, current_action.Action), '-')
			}
		}
	}

	return nil
}
