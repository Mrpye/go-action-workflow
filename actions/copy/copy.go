package copy

import (
	"fmt"

	"github.com/Mrpye/go-workflow/workflow"
	"github.com/Mrpye/golib/lib"
)

func Action_Copy(w *workflow.Workflow) error {

	//*********************
	//Get the config values
	//*********************
	source_file, err := w.GetConfigTokenString("source_file", w.Model, true)
	if err != nil {
		return err
	}

	dest_file, err := w.GetConfigTokenString("dest_file", w.Model, true)
	if err != nil {
		return err
	}

	data_copy, err := lib.CopyFile(source_file, dest_file)
	if err != nil {
		return err
	}

	if w.Verbose > workflow.LOG_QUIET {
		fmt.Printf("file %s copied to %s bytes %v", source_file, dest_file, data_copy)
	}

	return nil
}
