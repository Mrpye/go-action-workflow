package store

import "github.com/Mrpye/go-workflow/workflow"

func Action_Store(w *workflow.Workflow, m *workflow.TemplateData) error {

	//**********************************
	//Get a string value from the config
	//**********************************
	bucket, err := w.GetConfigTokenString("bucket", m, true)
	if err != nil {
		return err
	}

	key, err := w.GetConfigTokenString("key", m, true)
	if err != nil {
		return err
	}

	string_value, err := w.GetConfigTokenInterface("value", m, true)
	if err != nil {
		return err
	}

	w.SetValueToDataBucket(bucket, key, string_value)

	return nil
}
