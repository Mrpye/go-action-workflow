package count

import (
	"fmt"

	"github.com/Mrpye/go-workflow/workflow"
)

func Action_Count(w *workflow.Workflow) error {

	//*********************
	//Get the config values
	//*********************
	error_on_count, err := w.GetConfigTokenInt("error_on_count", w.Model, true)
	if err != nil {
		return err
	}

	result_store, err := w.GetConfigTokenString("result_store", w.Model, true)
	if err != nil {
		return err
	}

	count_var, err := w.GetConfigTokenString("count_var", w.Model, true)
	if err != nil {
		return err
	}
	count := w.GetValueFromDataBucketAsInt(result_store, count_var)
	count++
	w.SetValueToDataBucket(result_store, count_var, count)
	if count > error_on_count {
		return fmt.Errorf("error on count [%s] %v > %v %v", result_store, count, error_on_count, count > error_on_count)
	}

	if w.Verbose > workflow.LOG_QUIET {
		fmt.Printf("error on count [%s] %v > %v %v", result_store, count, error_on_count, count > error_on_count)
	}

	return nil
}
