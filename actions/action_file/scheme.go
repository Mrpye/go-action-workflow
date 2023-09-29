package action_file

import (
	"github.com/Mrpye/go-action-workflow/workflow"
)

type FileSchema struct {
}

// GetTargetSchema returns the target schema for this action
func (s FileSchema) GetTargetSchema() map[string]workflow.TargetSchema {
	//Build a test schema
	return map[string]workflow.TargetSchema{}
}

// GetFunctionMap returns the function map for this schema
func (s FileSchema) GetFunctionMap() map[string]workflow.FunctionSchema {
	return map[string]workflow.FunctionSchema{}
}

// GetActions returns the actions for this schema
func (s FileSchema) GetActionSchema() map[string]workflow.ActionSchema {

	//no short d f h
	return map[string]workflow.ActionSchema{
		"file_template": {
			Action:         Action_Template,
			Short:          "Create a file using a template",
			Long:           "Create a file using a template",
			ProcessResults: true,
			ConfigSchema: map[string]*workflow.Schema{
				"template": {
					Type:        workflow.TypeString,
					Partial:     false,
					Required:    true,
					Description: "Template file",
					Short:       "t",
					Default:     "",
				},
				"file": {
					Type:        workflow.TypeString,
					Partial:     false,
					Required:    true,
					Description: "file to create",
					Short:       "f",
					Default:     "",
				},
				"data": {
					Type:        workflow.TypeMap,
					Partial:     false,
					Required:    true,
					Description: "Map Data to use in the template",
					Short:       "d",
					Default:     []string{},
				},
			},
		},
		"file_create": {
			Action:         Action_Create,
			Short:          "Create a file",
			Long:           "Create a file",
			ProcessResults: true,
			ConfigSchema: map[string]*workflow.Schema{
				"source_file": {
					Type:        workflow.TypeString,
					Partial:     false,
					Required:    true,
					Description: "Source file to create",
					Short:       "s",
					Default:     "",
				},
				"content": {
					Type:        workflow.TypeString,
					Partial:     false,
					Required:    true,
					Description: "Content to add to the file",
					Short:       "c",
					Default:     "",
				},
			},
		},
		"file_append": {
			Action:         Action_Append,
			Short:          "Append a file",
			Long:           "Append a file",
			ProcessResults: true,
			ConfigSchema: map[string]*workflow.Schema{
				"source_file": {
					Type:        workflow.TypeString,
					Partial:     false,
					Required:    true,
					Description: "Source file to append to",
					Short:       "s",
					Default:     "",
				},
				"content": {
					Type:        workflow.TypeString,
					Partial:     false,
					Required:    true,
					Description: "Content to append to the file",
					Short:       "c",
					Default:     "",
				},
			},
		},
		"file_copy": {
			Action:         Action_Copy,
			Short:          "Copy a file",
			Long:           "Copy a file",
			ProcessResults: true,
			ConfigSchema: map[string]*workflow.Schema{
				"source_file": {
					Type:        workflow.TypeString,
					Partial:     false,
					Required:    true,
					Description: "Source file to copy",
					Short:       "s",
					Default:     "",
				},
				"dest_file": {
					Type:        workflow.TypeString,
					Partial:     false,
					Required:    true,
					Description: "Destination file to copy to",
					Short:       "d",
					Default:     "",
				},
			},
		},
		"file_rename": {
			Action:         Action_Rename,
			Short:          "Rename a file",
			Long:           "Rename a file",
			ProcessResults: true,
			ConfigSchema: map[string]*workflow.Schema{
				"source_file": {
					Type:        workflow.TypeString,
					Partial:     false,
					Required:    true,
					Description: "Source file to rename",
					Short:       "s",
					Default:     "",
				},
				"dest_file": {
					Type:        workflow.TypeString,
					Partial:     false,
					Required:    true,
					Description: "Destination file to rename to",
					Short:       "d",
					Default:     "",
				},
			},
		},
		"file_delete": {
			Action:         Action_Delete,
			Short:          "Delete a file",
			Long:           "Delete a file",
			ProcessResults: true,
			ConfigSchema: map[string]*workflow.Schema{
				"source_file": {
					Type:        workflow.TypeString,
					Partial:     false,
					Required:    true,
					Description: "Source file to rename",
					Short:       "s",
					Default:     "",
				},
			},
		},
	}
}

// GetAction returns the action for the target
func GetSchema() workflow.SchemaEndpoint {
	return FileSchema{}
}
