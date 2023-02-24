//********************************************************************************************
//This module contains common functions that can be used in workflows for testing all features
//********************************************************************************************
package tests

import (
	"encoding/json"
	"fmt"

	"github.com/Mrpye/go-workflow/workflow"
)

func ActionJSAndMap(w *workflow.Workflow, m *workflow.TemplateData) error {

	//*******************************
	//Get a map value from the config
	//*******************************
	map_value, err := w.GetConfigTokenMap("map_value", m, true)
	if err != nil {
		return err
	}

	//***************************************
	//Convert to json and process the results
	//***************************************
	b, err := json.Marshal(map_value)
	if err != nil {
		return err
	}

	//***********************************
	//This function processes the results
	//***********************************
	err = w.ActionProcessResults(m, string(b))
	if err != nil {
		return err
	}
	return nil
}

func ActionTest(w *workflow.Workflow, m *workflow.TemplateData) error {
	//**********************************************
	//Get the model if m is passed then its parallel
	//**********************************************
	//****************************************************************
	//This will test the data and return an error if it is not correct
	//****************************************************************
	test_data := w.GetDataBucketContent("test")

	//*******************
	//test the map and js
	//*******************
	err := MapCheck(test_data, "js_map_value1", "3")
	if err != nil {
		return err
	}

	//*******************
	//Action test
	//*******************
	err = MapCheck(test_data, "sub_workflow_value1", "sub-workflow first value")
	if err != nil {
		return err
	}
	err = MapCheck(test_data, "sub_workflow_value2", "sub-workflow second value")
	if err != nil {
		return err
	}
	err = MapCheck(test_data, "sub_workflow_value3", "sub-workflow third value")
	if err != nil {
		return err
	}

	err = MapCheck(test_data, "store_test_key", "This is a value from store test")
	if err != nil {
		return err
	}

	err = MapCheck(test_data, "condition_test", "pass")
	if err != nil {
		return err
	}

	err = MapCheck(test_data, "js_map_value2", "THIS IS A TEST STRING")
	if err != nil {
		return err
	}
	err = MapCheck(test_data, "js_map_value3", "true")
	if err != nil {
		return err
	}
	err = MapCheck(test_data, "js_map_value4", "THIS IS A VALUE FROM THE CONFIG")
	if err != nil {
		return err
	}
	//******************
	//test the meta data
	//******************
	err = MapCheck(test_data, "meta_var", "This is an example value")
	if err != nil {
		return err
	}

	err = MapCheck(test_data, "meta_name", "test-example")
	if err != nil {
		return err
	}

	err = MapCheck(test_data, "meta_description", "This is used for testing")
	if err != nil {
		return err
	}

	err = MapCheck(test_data, "meta_version", "1.0.0")
	if err != nil {
		return err
	}

	err = MapCheck(test_data, "meta_author", "Andrew Pye")
	if err != nil {
		return err
	}

	err = MapCheck(test_data, "meta_contact", "test@test.com")
	if err != nil {
		return err
	}
	err = MapCheck(test_data, "meta_create_date", "2022-11-13 11:39:44")
	if err != nil {
		return err
	}
	err = MapCheck(test_data, "meta_update_date", "2022-11-13 11:39:44")
	if err != nil {
		return err
	}

	//***************
	//test the params
	//***************
	err = MapCheck(test_data, "param_test_int", "3")
	if err != nil {
		return err
	}

	err = MapCheck(test_data, "param_test_string", "this is a test string")
	if err != nil {
		return err
	}

	err = MapCheck(test_data, "param_test_bool", "true")
	if err != nil {
		return err
	}

	for i := 0; i <= 3; i++ {
		err = MapCheck(test_data, fmt.Sprintf("loop_increment%d", i), fmt.Sprintf("%d", i))
		if err != nil {
			return err
		}
	}

	for i := 3; i >= 0; i-- {
		err = MapCheck(test_data, fmt.Sprintf("loop_decrement%d", i), fmt.Sprintf("%d", i))
		if err != nil {
			return err
		}
	}

	for i := 0; i <= 3; i++ {
		for j := 0; j <= 3; j++ {
			err = MapCheck(test_data, fmt.Sprintf("nested_loop_%d-%d", i, j), fmt.Sprintf("%d-%d", i, j))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func MapCheck(map_data map[string]interface{}, key string, value any) error {
	if map_data[key] == nil {
		return fmt.Errorf("%s is not set", key)
	}
	if map_data[key] != value {
		return fmt.Errorf("%s is not set to = %v: is %v", key, value, map_data[key])
	}
	return nil
}
func ActionFailTest(w *workflow.Workflow, m *workflow.TemplateData) error {
	//**********************************************
	//Get the model if m is passed then its parallel
	//**********************************************
	return fmt.Errorf("this action should not run")
}
