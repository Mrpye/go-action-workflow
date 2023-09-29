package action_file

import (
	"fmt"

	"github.com/Mrpye/go-action-workflow/workflow"
	"github.com/Mrpye/golib/file"
)

// Action_Copy is a custom action that copies a file
/*
- action: ActionStore
	config:
		source_file: the file to copy
		dest_file: 	the destination file
*/
func Action_Copy(w *workflow.Workflow, m *workflow.TemplateData) error {

	//*********************
	//Get the config values
	//*********************
	source_file, err := w.GetConfigTokenString("source_file", m, true)
	if err != nil {
		return err
	}
	source_file = w.BuildPath(source_file)

	dest_file, err := w.GetConfigTokenString("dest_file", m, true)
	if err != nil {
		return err
	}
	dest_file = w.BuildPath(dest_file)

	data_copy, err := file.CopyFile(source_file, dest_file)
	if err != nil {
		return err
	}

	if w.LogLevel > workflow.LOG_QUIET {
		fmt.Printf("file %s copied to %s bytes %v\n", source_file, dest_file, data_copy)
	}

	return nil
}
