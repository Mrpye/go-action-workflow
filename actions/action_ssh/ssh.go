package action_ssh

import (
	"fmt"

	"github.com/Mrpye/go-action-workflow/workflow"
)

type SSHSchema struct {
}

// GetFunctionMap returns the function map for this schema
func (s SSHSchema) GetFunctionMap() map[string]workflow.FunctionSchema {
	return map[string]workflow.FunctionSchema{}
}

// GetTargetSchema returns the target schema for this action
func (s SSHSchema) GetTargetSchema() map[string]workflow.TargetSchema {
	//Build a test schema
	return map[string]workflow.TargetSchema{
		"action_ssh.ssh": workflow.BuildTargetConfig("ssh", "ssh", &SSH{}),
	}
}

// GetActions returns the actions for this schema
func (s SSHSchema) GetActionSchema() map[string]workflow.ActionSchema {
	//no short d f h
	return map[string]workflow.ActionSchema{
		"ssh_run_cmd": {
			Action:         Action_SSHRunCMD,
			Short:          "Run a command on a SSH server",
			Long:           "Run a command on a SSH server",
			ProcessResults: true,
			Target:         "action_ssh.ssh",
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
					Description: "command to run on the server",
					Short:       "c",
					Default:     "",
				},
			},
		},
		"ssh_run_script_file": {
			Action:         Action_SSHRunScriptFile,
			Short:          "Run a script file on a SSH server",
			Long:           "Run a script file on a SSH server",
			ProcessResults: true,
			Target:         "action_ssh.ssh",
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
				"file": {
					Type:        workflow.TypeString,
					Partial:     false,
					Required:    true,
					Description: "Script file to run on the server",
					Short:       "f",
					Default:     "",
				},
			},
		},
		"ssh_run_script": {
			Action:         Action_SSHRunScript,
			Short:          "Run a script on a SSH server",
			Long:           "Run a script on a SSH server",
			ProcessResults: true,
			Target:         "action_ssh.ssh",
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
				"script": {
					Type:        workflow.TypeString,
					Partial:     false,
					Required:    true,
					Description: "Script to run on the server",
					Short:       "s",
					Default:     "",
				},
			},
		},
		"ssh_upload": {
			Action:         Action_SSHUpload,
			Short:          "Upload a file to a SSH server",
			Long:           "Upload a file to a SSH server",
			ProcessResults: true,
			Target:         "action_ssh.ssh",
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
				"source": {
					Type:        workflow.TypeString,
					Partial:     false,
					Required:    true,
					Description: "Source file to upload",
					Short:       "s",
					Default:     "",
				},
				"dest": {
					Type:        workflow.TypeString,
					Partial:     false,
					Required:    true,
					Description: "Destination file to download to",
					Short:       "d",
					Default:     "",
				},
			},
		},
		"ssh_download": {
			Action:         Action_SSHDownload,
			Short:          "Download a file to a SSH server",
			Long:           "Download a file to a SSH server",
			ProcessResults: true,
			Target:         "action_ssh.ssh",
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
				"source": {
					Type:        workflow.TypeString,
					Partial:     false,
					Required:    true,
					Description: "Source file to download",
					Short:       "s",
					Default:     "",
				},
				"dest": {
					Type:        workflow.TypeString,
					Partial:     false,
					Required:    true,
					Description: "Destination file to download to",
					Short:       "d",
					Default:     "",
				},
			},
		},
	}
}

// GetAction returns the action for the target
func GetSchema() workflow.SchemaEndpoint {
	return SSHSchema{}
}
func Action_SSHRunCMD(w *workflow.Workflow, m *workflow.TemplateData) error {

	ssh_obj, err := w.MapTargetConfigValue(m, &SSH{})
	if err != nil {
		return err
	}
	ssh := ssh_obj.(*SSH)

	//*********************
	//Get the config values
	//*********************
	command, err := w.GetConfigTokenString("command", m, true)

	if err != nil {
		return err
	}

	out, err := ssh.SSHRunCMD(command)
	if err != nil {
		return err
	}

	if w.LogLevel > workflow.LOG_QUIET {
		fmt.Printf("cmd %s result %s ", command, out)
	}

	//**************
	//Return the err
	//**************
	return w.ActionProcessResults(m, out)
}

func Action_SSHRunScriptFile(w *workflow.Workflow, m *workflow.TemplateData) error {

	ssh_obj, err := w.MapTargetConfigValue(m, &SSH{})
	if err != nil {
		return err
	}
	ssh := ssh_obj.(*SSH)

	//*********************
	//Get the config values
	//*********************
	file, err := w.GetConfigTokenString("file", m, true)
	if err != nil {
		return err
	}
	file = w.BuildPath(file)
	out, err := ssh.SSHRunScriptFile(file)
	if err != nil {
		return err
	}

	if w.LogLevel > workflow.LOG_QUIET {
		fmt.Printf("script %s result %s ", file, out)
	}

	//**************
	//Return the err
	//**************
	return w.ActionProcessResults(m, out)
}

func Action_SSHRunScript(w *workflow.Workflow, m *workflow.TemplateData) error {

	ssh_obj, err := w.MapTargetConfigValue(m, &SSH{})
	if err != nil {
		return err
	}
	ssh := ssh_obj.(*SSH)

	//*********************
	//Get the config values
	//*********************
	script, err := w.GetConfigTokenString("script", m, true)
	if err != nil {
		return err
	}

	out, err := ssh.SSHRunScript(script)
	if err != nil {
		return err
	}

	if w.LogLevel > workflow.LOG_QUIET {
		fmt.Printf("script %s result %s ", script, out)
	}

	//**************
	//Return the err
	//**************
	return w.ActionProcessResults(m, out)
}

func Action_SSHUpload(w *workflow.Workflow, m *workflow.TemplateData) error {
	ssh_obj, err := w.MapTargetConfigValue(m, &SSH{})
	if err != nil {
		return err
	}
	ssh := ssh_obj.(*SSH)
	fmt.Print(ssh.Host)

	//*********************
	//Get the config values
	//*********************
	source_file, err := w.GetConfigTokenString("source", m, true)
	if err != nil {
		return err
	}
	source_file = w.BuildPath(source_file)

	destination, err := w.GetConfigTokenString("dest", m, true)
	if err != nil {
		return err
	}
	destination = w.BuildPath(destination)

	err = ssh.SSHUploadFile(source_file, destination)
	if err != nil {
		return err
	}

	if w.LogLevel > workflow.LOG_QUIET {
		fmt.Printf("file %s copied to %s ", source_file, destination)
	}

	//**************
	//Return the err
	//**************
	return nil
}

func Action_SSHDownload(w *workflow.Workflow, m *workflow.TemplateData) error {

	ssh_obj, err := w.MapTargetConfigValue(m, &SSH{})
	if err != nil {
		return err
	}
	ssh := ssh_obj.(*SSH)
	fmt.Print(ssh.Host)

	//*********************
	//Get the config values
	//*********************
	source_file, err := w.GetConfigTokenString("source", m, true)
	if err != nil {
		return err
	}

	destination, err := w.GetConfigTokenString("dest", m, true)
	if err != nil {
		return err
	}
	destination = w.BuildPath(destination)

	err = ssh.SSHDownloadFile(source_file, destination)
	if err != nil {
		return err
	}

	if w.LogLevel > workflow.LOG_QUIET {
		fmt.Printf("file %s copied to %s ", source_file, destination)
	}

	//**************
	//Return the err
	//**************
	return nil
}
