package store

import "github.com/Mrpye/go-workflow/workflow"

func Action_Store(w *workflow.Workflow) error {
	//**********************************
	//Get a string value from the config
	//**********************************
	bucket, err := w.GetConfigTokenString("bucket", w.Model, true)
	if err != nil {
		return err
	}

	key, err := w.GetConfigTokenString("key", w.Model, true)
	if err != nil {
		return err
	}

	string_value, err := w.GetConfigTokenInterface("value", w.Model, true)
	if err != nil {
		return err
	}

	w.SetValueToDataBucket(bucket, key, string_value)

	return nil
}
