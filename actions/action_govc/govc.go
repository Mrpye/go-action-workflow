package action_govc

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/Mrpye/go-action-workflow/workflow"
)

type GovcSchema struct {
}

// GetFunctionMap returns the function map for this schema
func (s GovcSchema) GetFunctionMap() map[string]workflow.FunctionSchema {
	return map[string]workflow.FunctionSchema{}
}

// GetTargetSchema returns the target schema for this action
func (s GovcSchema) GetTargetSchema() map[string]workflow.TargetSchema {
	//Build a test schema
	return map[string]workflow.TargetSchema{
		"action_govc.vcenter": workflow.BuildTargetConfig("vcenter", "vcenter", &VCenter{}),
	}
}

// GetActions returns the actions for this schema
func (s GovcSchema) GetActionSchema() map[string]workflow.ActionSchema {

	//no short d f h
	return map[string]workflow.ActionSchema{
		"govc": {
			Action:         Action_GOVC,
			Short:          "Run GOVC commands",
			Long:           "Run GOVC commands",
			ProcessResults: true,
			Target:         "action_govc.vcenter",
			ConfigSchema: map[string]*workflow.Schema{
				"target_name": {
					Type:        workflow.TypeString,
					Partial:     false,
					Required:    false,
					Description: "The target to use",
					ConfigKey:   "target_name",
					Short:       "t",
					Default:     "",
				},
				"command": {
					Type:        workflow.TypeString,
					Partial:     false,
					Required:    true,
					Description: "Command to run",
					Short:       "c",
					Default:     "",
				},
			},
		},
	}
}

// GetAction returns the action for the target
func GetSchema() workflow.SchemaEndpoint {
	return GovcSchema{}
}
func PrintSomethingOut(out io.Writer) {
	fmt.Fprintln(out, "print something to io.Writer")
}

/*func capture() func() (string, error) {
	r, w, err := os.Pipe()
	if err != nil {
		panic(err)
	}

	done := make(chan error, 1)

	save := os.Stdout
	os.Stdout = w

	var buf strings.Builder

	go func() {
		_, err := io.Copy(&buf, r)
		r.Close()
		done <- err
	}()

	return func() (string, error) {
		os.Stdout = save
		w.Close()
		err := <-done
		return buf.String(), err
	}
}*/

func Action_GOVC(w *workflow.Workflow, m *workflow.TemplateData) error {

	govc_obj, err := w.MapTargetConfigValue(m, &VCenter{})
	if err != nil {
		return err
	}
	vcenter_target := govc_obj.(*VCenter)

	//**********************
	//Get the command string
	//**********************
	command, err := w.GetConfigTokenString("command", m, true)
	if err != nil {
		return err
	}
	//remove the carriage character
	command = strings.ReplaceAll(command, "\n", " ")

	//***********************
	//Run the vcenter command
	//***********************
	rescueStdout := os.Stdout
	rescueStderr := os.Stderr

	//Error Pipe
	er, ew, _ := os.Pipe()
	os.Stderr = ew

	//******************************************************************************
	// Read the output in a separate goroutine so printing can't block indefinitely.
	//******************************************************************************
	var buf strings.Builder
	scanner := bufio.NewScanner(io.Reader(rescueStdout))
	scanner.Split(bufio.ScanWords)

	for scanner.Scan() {
		m := scanner.Text()
		buf.WriteString(string(m))
		fmt.Println(string(m))
	}

	//******************
	//Run the command
	//******************
	err = vcenter_target.Run(command, false)
	if err != nil {
		return err
	}

	//******************
	//Close the pipes
	//******************
	ew.Close()
	err_out, _ := ioutil.ReadAll(er)
	os.Stderr = rescueStderr
	os.Stdout = rescueStdout

	//******************
	//Display the result
	//******************
	if w.LogLevel > workflow.LOG_QUIET {
		fmt.Println(string(buf.String()))
		fmt.Println(string(err_out))
	}

	//******************
	//Check for errors
	//******************
	if string(err_out) != "" {
		return errors.New(string(err_out))
	}
	//*******************
	//Process the results
	//*******************
	err = w.ActionProcessResults(m, buf.String())
	if err != nil {
		return err
	}

	//**************
	//Return the err
	//**************
	return nil
}
