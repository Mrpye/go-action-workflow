package file

import (
	"fmt"
	"log"
	"os"

	"github.com/Mrpye/go-workflow/workflow"
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

	dest_file, err := w.GetConfigTokenString("dest_file", m, true)
	if err != nil {
		return err
	}

	e := os.Rename(source_file, dest_file)
	if e != nil {
		log.Fatal(e)
	}

	if w.Verbose > workflow.LOG_QUIET {
		fmt.Printf("file %s renamed to %s\n", source_file, dest_file)
	}

	return nil
}
