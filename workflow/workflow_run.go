package workflow

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/Mrpye/golib/convert"
	"github.com/Mrpye/golib/log"
	"github.com/gookit/color"
)

// RunJob will run the job with the given key
// key - the key of the job to run
// returns - error
func (m *Workflow) RunJob(key string, run_time_vars map[string]interface{}, params map[string]interface{}) error {

	//************************
	//Log the start of the job
	//************************
	log.ActionLogDT(fmt.Sprintf("Running Job %s", key), '*')

	//************************
	//Set the run time variables
	//************************
	m.runtime_vars = run_time_vars

	//***************
	//Set the answers
	//***************
	if params != nil {
		m.SetAnswers(params)
	}

	// *******************************
	// Run the job with the given key.
	// *******************************
	err := m.executeJob(key, false, nil)

	//**************************************
	//Initialize the the actions and targets
	//**************************************
	if m.CleanFunc != nil {
		err := m.CleanFunc(m)
		if err != nil {
			return err
		}
	}
	//************************
	//Clear Internal Variables
	//************************
	m.current_job = nil
	m.current_index = 0
	m.model = nil
	m.stack = CreateLoopStack()
	m.runtime_vars = nil

	if err != nil {
		log.ActionLogRed(fmt.Sprintf("Job %s failed with error %s", key, err.Error()), '*')
	} else {
		log.ActionLogGreen(fmt.Sprintf("Job %s completed successfully", key), '*')
	}

	fmt.Println("")
	fmt.Println("")
	fmt.Println("")

	return err
}

// RunSubWorkflow will run a SubWorkflow job with the given key
// key - the key of the job to run
// inputs - the inputs to the job
// returns - error
func (m *Workflow) RunSubWorkflow(key string, inputs map[string]interface{}) error {

	// *******************************
	// Run the job with the given key.
	// *******************************
	err := m.executeJob(key, true, inputs)

	//************************
	//Clear Internal Variables
	//************************
	m.current_job = nil
	m.current_index = 0
	m.model = nil
	m.stack = CreateLoopStack()

	return err
}

// executeJob will run the job with the given key
// key - the key of the job to run
// returns - error
func (m *Workflow) executeJob(key string, sub_process bool, input_values map[string]interface{}) error {

	//**********************************************
	//Lets get the m.current_job for the selected app profile
	//**********************************************
	m.current_job = m.Manifest.GetJob(key)
	if m.current_job == nil {
		return fmt.Errorf("cannot find job %s", key)
	}

	//**********
	//Map values
	//**********
	if sub_process {
		//******************************************************
		// See if this is a sub process is run by a sub-workflow
		//******************************************************
		if !(sub_process && m.current_job.IsSubWorkflow) {
			return fmt.Errorf("this job %s cannot be run directly you need to run via the sub-workflow action", key)
		}
		err := m.MapValuesToInput(input_values)
		if err != nil {
			return err
		}
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
	m.stack = CreateLoopStack()

	//**************
	//create a Model
	//**************
	m.model = m.CreateTemplateData(nil)

	//********************
	//Setup some variables
	//********************
	Job_action_count := len(m.current_job.Actions)
	m.current_index = 0

	//*********************
	//Loop over the actions
	//*********************
	for m.current_index = 0; m.current_index < Job_action_count; m.current_index++ {
		current_action := &m.current_job.Actions[m.current_index]
		//********************************************
		//Set the current action to our template Model
		//********************************************
		m.model.SetAction(current_action)

		//******************************
		//See if this action is disabled
		//******************************
		is_disabled, err := m.GetTokenBool(current_action.Disabled, m.model)
		if err != nil {
			return err
		}
		if is_disabled {
			if m.LogLevel > LOG_QUIET {
				color := color.FgRed.Render
				log.ActionLogDT(fmt.Sprintf("Action Ignored: %s->%s : %s", current_action.Key, current_action.Action, color("Disabled")), '*')
			}
			continue
		}
		if m.LogLevel > LOG_QUIET {
			log.ActionLogDT(fmt.Sprintf("Action: %s->%s", current_action.Key, current_action.Action), '*')
		}
		//********************************
		//Parse any variable in the action
		//********************************
		action_parts, err := m.splitActionParams(current_action.Action, true)
		if err != nil {
			return err
		}

		//**************************
		//Find the action to perform
		//**************************
		switch action_parts[0] {
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
			temp_loop, err := m.stack.Push(action_parts[1], m.current_index, from, to)
			if err != nil {
				return err
			}
			if m.LogLevel > LOG_QUIET {
				log.ActionLogDT("loop: "+action_parts[1]+"["+strconv.FormatInt(int64(temp_loop.CurrentValue), 10)+"]", '>')
			}
		case "next":
			//Get the item off the stack
			temp_loop, err := m.stack.Peek()
			if err != nil {
				return err
			}
			if m.LogLevel > LOG_QUIET {
				log.ActionLogDT("loop-end: "+temp_loop.VariableName+"["+strconv.FormatInt(int64(temp_loop.CurrentValue), 10)+"]", '<')
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
				m.current_index = temp_loop.Index - 1
			}
		default:
			//***************
			//Run the action
			//***************
			err = m.RunAction(current_action)
			if errors.Is(err, ErrEndWorkflow) {
				if m.LogLevel > LOG_QUIET {
					log.ActionLogDateOK(fmt.Sprintf("Action Completed: %s->%s", current_action.Key, current_action.Action), '-')
				}
				return nil
			}
			if err != nil {
				return err
			}
			if m.LogLevel > LOG_QUIET {
				log.ActionLogDateOK(fmt.Sprintf("Action Completed: %s->%s", current_action.Key, current_action.Action), '-')
			}
		}
	}

	return nil
}

// RunAction will run the action
// - current_action - the action to run
// - returns - error
func (m *Workflow) RunAction(current_action *Action) error {

	//********************************
	//Parse any variable in the action
	//********************************
	action_parts, err := m.splitActionParams(current_action.Action, false)
	if err != nil {
		return err
	}
	switch action_parts[0] {
	case "end":
		//End the job
		if m.LogLevel > LOG_QUIET {
			log.ActionLogDT("End", '*')
		}
		return ErrEndWorkflow
	case "fail":
		//End the job
		if m.LogLevel > LOG_QUIET {
			log.ActionLogDateFail("fail", '*')
		}
		return fmt.Errorf("%s", action_parts[1])
	case "print":
		fmt.Println(action_parts[1])
		/*if m.LogLevel > LOG_QUIET {
			log.ActionLogDateOK(fmt.Sprintf("Action Completed: %s->%s", current_action.Key, current_action.Action), '-')
		}*/
	case "goto":
		if len(action_parts) < 2 {
			return errors.New("not enough args for goto")
		}
		if m.LogLevel > LOG_QUIET {
			log.ActionLogDT("goto: "+action_parts[1], '>')
		}
		index := m.current_job.GetKeyIndex(action_parts[1])
		if index == -1 {
			return fmt.Errorf("cannot find label %s", action_parts[1])
		}
		m.current_index = index - 1
	case "wait-seconds", "wait":
		if len(action_parts) < 2 {
			return errors.New("not enough args for wait-seconds wait ")
		}
		/*if m.LogLevel > LOG_QUIET {
			log.ActionLogDT("wait-seconds/wait: "+action_parts[1], '*')
		}*/
		count, err := strconv.Atoi(action_parts[1])
		if err != nil {
			return fmt.Errorf("value should be an int for wait-seconds wait %s", action_parts[1])
		}
		time.Sleep(time.Duration(count) * time.Second)
	//	if m.LogLevel > LOG_QUIET {
	//	log.ActionLogDateOK(fmt.Sprintf("Action Completed: %s->%s", current_action.Key, current_action.Action), '-')
	//	}
	case "wait-minutes":
		if len(action_parts) < 2 {
			return errors.New("not enough args for wait-minutes")
		}
		/*if m.LogLevel > LOG_QUIET {
			log.ActionLogDT("wait-minutes: "+action_parts[1], '*')
		}*/
		count, err := strconv.Atoi(action_parts[1])
		if err != nil {
			return fmt.Errorf("value should be an int for wait-minutes %s", action_parts[1])
		}
		time.Sleep(time.Duration(count) * time.Minute)
		//if m.LogLevel > LOG_QUIET {
	//		log.ActionLogDateOK(fmt.Sprintf("Action Completed: %s->%s", current_action.Key, current_action.Action), '-')
	//}
	case "wait-hours":
		if len(action_parts) < 2 {
			return errors.New("not enough args for wait-minutes")
		}
		/*if m.LogLevel > LOG_QUIET {
			log.ActionLogDT("wait-minutes: "+action_parts[1], '*')
		}*/
		count, err := strconv.Atoi(action_parts[1])
		if err != nil {
			return fmt.Errorf("value should be an int for wait-minutes %s", action_parts[1])
		}
		time.Sleep(time.Duration(count) * time.Hour)
		//if m.LogLevel > LOG_QUIET {
		//	log.ActionLogDateOK(fmt.Sprintf("Action Completed: %s->%s", current_action.Key, current_action.Action), '-')
		//}
	case "action":
		if len(action_parts) < 2 {
			return errors.New("not enough args for action")
		}
		current_action = m.Manifest.GetGlobalAction(action_parts[1])
		if current_action == nil {
			return fmt.Errorf("global action %s not found", action_parts[1])
		}
		return m.RunAction(current_action)
	default:
		continue_on_error := convert.ToBool(current_action.ContinueOnError)
		var err error

		if val, ok := m.ActionList[current_action.Action]; ok {
			m.model = m.CreateTemplateData(current_action)
			err = val(m, m.model)
		} else {
			return fmt.Errorf("action %s not found", current_action.Action)
		}

		if err != nil {
			//**************************************
			//See if this is an end workflow error
			// Just stop the workflow and return nil
			//**************************************
			if errors.Is(err, ErrEndWorkflow) {
				return ErrEndWorkflow
			}
			//****************************
			// See if to continue on error
			//****************************
			if continue_on_error {
				log.PrintlnFail("!! There was an error but Continue On Error was set to true !!")
				return nil
			}
			//************************************
			//lets see if to end or goto an action
			//************************************
			fail_parts, Slit_err := m.splitActionParams(current_action.Fail, true)
			if Slit_err != nil {
				return Slit_err
			}
			if fail_parts[0] == "end" || fail_parts[0] == "" {
				return err
			} else if fail_parts[0] == "goto" {
				log.PrintlnFail("!! There but goto is set for fail !!")

				if m.LogLevel > LOG_QUIET {
					log.LogVerbose(fmt.Sprintf("The follow error occurred for action: %s->%s:%s", current_action.Key, current_action.Action, err.Error()))
				}

				if len(fail_parts) < 2 {
					return errors.New("not enough args for goto")
				}
				if m.LogLevel > LOG_QUIET {
					log.ActionLogDT("goto: "+fail_parts[1], '>')
				}
				index := m.current_job.GetKeyIndex(fail_parts[1])
				if index == -1 {
					return fmt.Errorf("cannot find label %s", fail_parts[1])
				}
				m.current_index = index - 1
			}
		}
	}
	return nil
}
