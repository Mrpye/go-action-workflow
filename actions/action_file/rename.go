package action_file

import (
	"fmt"
	"os"

	"github.com/Mrpye/go-action-workflow/workflow"
)

// Action_Copy is a custom action that copies a file
/*
- action: ActionStore
	config:
		source_file: the file to copy
		dest_file: 	the destination file
*/
func Action_Rename(w *workflow.Workflow, m *workflow.TemplateData) error {

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

	e := os.Rename(source_file, dest_file)
	if e != nil {
		return err
	}

	if w.LogLevel > workflow.LOG_QUIET {
		fmt.Printf("file %s renamed to %s\n", source_file, dest_file)
	}

	return nil
}
