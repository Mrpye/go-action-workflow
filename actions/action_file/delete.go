package action_file

import (
	"fmt"
	"os"

	"github.com/Mrpye/go-action-workflow/workflow"
	"github.com/Mrpye/golib/log"
)

// Action_Copy is a custom action that copies a file
/*
- action: ActionStore
	config:
		source_file: the file to copy
		dest_file: 	the destination file
*/
func Action_Delete(w *workflow.Workflow, m *workflow.TemplateData) error {

	//*********************
	//Get the config values
	//*********************
	source_file, err := w.GetConfigTokenString("source_file", m, true)
	if err != nil {
		return err
	}
	source_file = w.BuildPath(source_file)
	//check if source file is a directory
	fileInfo, err := os.Stat(source_file)
	if err != nil {
		return err
	}

	if fileInfo.IsDir() {
		// is a directory
		err := os.RemoveAll(source_file)
		if err != nil {
			return err
		}
		if w.LogLevel > workflow.LOG_QUIET {
			log.PrintlnOK(fmt.Sprintf("Folder %s deleted", source_file))
		}
	} else {
		// is not a directory
		e := os.Remove(source_file)
		if e != nil {
			return err
		}
		if w.LogLevel > workflow.LOG_QUIET {
			log.PrintlnOK(fmt.Sprintf("File %s deleted", source_file))
		}
	}

	return nil
}
