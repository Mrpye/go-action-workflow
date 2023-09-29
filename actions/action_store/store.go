package action_store

import (
	"github.com/Mrpye/go-action-workflow/workflow"
)

type StoreSchema struct {
}

// GetFunctionMap returns the function map for this schema
func (s StoreSchema) GetFunctionMap() map[string]workflow.FunctionSchema {
	return map[string]workflow.FunctionSchema{}
}

// GetTargetSchema returns the target schema for this action
func (s StoreSchema) GetTargetSchema() map[string]workflow.TargetSchema {
	//Build a test schema
	return map[string]workflow.TargetSchema{}
}

// GetActions returns the actions for this schema
func (s StoreSchema) GetActionSchema() map[string]workflow.ActionSchema {

	//no short d f h
	return map[string]workflow.ActionSchema{
		"store": {
			Action:         Action_Store,
			Short:          "Store a value in the data bucket",
			Long:           "Store a value in the data bucket",
			ProcessResults: false,
			ConfigSchema: map[string]*workflow.Schema{
				"bucket": {
					Type:        workflow.TypeString,
					Partial:     false,
					Required:    true,
					Description: "The bucket to store the value in",
					Short:       "b",
					Default:     "",
				},
				"key": {
					Type:        workflow.TypeString,
					Partial:     false,
					Required:    true,
					Description: "The key to store the value with",
					Short:       "i",
					Default:     "",
				},
				"value": {
					Type:        workflow.TypeString,
					Partial:     false,
					Required:    true,
					Description: "The value to store",
					Short:       "v",
					Default:     "",
				},
			},
		},
	}
}

// GetAction returns the action for the target
func GetSchema() workflow.SchemaEndpoint {
	return StoreSchema{}
}

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
