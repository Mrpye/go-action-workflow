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
func Action_Create(w *workflow.Workflow, m *workflow.TemplateData) error {

	//*********************
	//Get the config values
	//*********************
	source_file, err := w.GetConfigTokenString("source_file", m, true)
	if err != nil {
		return err
	}
	source_file = w.BuildPath(source_file)

	content, err := w.GetConfigTokenString("content", m, true)
	if err != nil {
		return err
	}

	f, err := os.OpenFile(source_file,
		os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := f.WriteString(content); err != nil {
		return err
	}

	if w.LogLevel > workflow.LOG_QUIET {
		fmt.Printf("content written to file %s\n", source_file)
	}

	return nil
}
