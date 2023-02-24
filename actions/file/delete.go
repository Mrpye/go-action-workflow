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
func Action_Delete(w *workflow.Workflow, m *workflow.TemplateData) error {

	//*********************
	//Get the config values
	//*********************
	source_file, err := w.GetConfigTokenString("source_file", m, true)
	if err != nil {
		return err
	}

	e := os.Remove(source_file)
	if e != nil {
		log.Fatal(e)
	}

	if w.Verbose > workflow.LOG_QUIET {
		fmt.Printf("file %s deleted\n", source_file)
	}

	return nil
}
