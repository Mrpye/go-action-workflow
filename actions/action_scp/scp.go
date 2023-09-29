package action_scp

import (
	"fmt"

	"github.com/Mrpye/go-action-workflow/workflow"
)

type ScpSchema struct {
}

// GetFunctionMap returns the function map for this schema
func (s ScpSchema) GetFunctionMap() map[string]workflow.FunctionSchema {
	return map[string]workflow.FunctionSchema{}
}

// GetTargetSchema returns the target schema for this action
func (s ScpSchema) GetTargetSchema() map[string]workflow.TargetSchema {
	//Build a test schema
	return map[string]workflow.TargetSchema{
		"action_scp.scp": workflow.BuildTargetConfig("scp", "scp", &SCP{}),
	}
}

// GetActions returns the actions for this schema
func (s ScpSchema) GetActionSchema() map[string]workflow.ActionSchema {

	//no short d f h
	return map[string]workflow.ActionSchema{
		"scp_upload": {
			Action:         Action_ScpUpload,
			Short:          "Upload file to scp server",
			Long:           "Upload file to scp server",
			Target:         "action_scp.scp",
			ProcessResults: true,
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
					Description: "Source file to upload",
					Short:       "f",
					Default:     "",
				},
				"dest": {
					Type:        workflow.TypeString,
					Partial:     false,
					Required:    true,
					Description: "Destination file to upload to on the server",
					Short:       "d",
					Default:     "",
				},
			},
		},
		"scp_download": {
			Action:         Action_ScpDownload,
			Short:          "Download file to scp server",
			Long:           "Download file to scp server",
			ProcessResults: true,
			Target:         "action_scp.scp",
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
					Description: "Source file to download from the server",
					Short:       "f",
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
	return ScpSchema{}
}

func Action_ScpUpload(w *workflow.Workflow, m *workflow.TemplateData) error {

	scp_obj, err := w.MapTargetConfigValue(m, &SCP{})
	if err != nil {
		return err
	}
	scp := scp_obj.(*SCP)
	fmt.Print(scp.Host)

	//*********************
	//Get the config values
	//*********************
	source_file, err := w.GetConfigTokenString("file", m, true)
	if err != nil {
		return err
	}
	source_file = w.BuildPath(source_file)

	destination, err := w.GetConfigTokenString("dest", m, true)
	if err != nil {
		return err
	}
	destination = w.BuildPath(destination)

	err = scp.SCPFileUpload(source_file, destination)
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

func Action_ScpDownload(w *workflow.Workflow, m *workflow.TemplateData) error {

	scp_obj, err := w.MapTargetConfigValue(m, &SCP{})
	if err != nil {
		return err
	}
	scp := scp_obj.(*SCP)
	fmt.Print(scp.Host)

	//*********************
	//Get the config values
	//*********************
	source_file, err := w.GetConfigTokenString("file", m, true)
	if err != nil {
		return err
	}
	source_file = w.BuildPath(source_file)

	destination, err := w.GetConfigTokenString("dest", m, true)
	if err != nil {
		return err
	}
	destination = w.BuildPath(destination)

	err = scp.SCPFileDownload(source_file, destination)
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
