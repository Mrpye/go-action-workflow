package workflow

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Mrpye/golib/lib"
	"github.com/dop251/goja"
	"github.com/drewstinnett/gout/v2"
	"github.com/drewstinnett/gout/v2/formats/json"
	"github.com/drewstinnett/gout/v2/formats/plain"
	"github.com/drewstinnett/gout/v2/formats/toml"
	"github.com/drewstinnett/gout/v2/formats/xml"
	"github.com/drewstinnett/gout/v2/formats/yaml"
)

//Function to Process results from an action
func (m *Workflow) ActionProcessResults(data interface{}) error {
	//***********************
	//Get how to store result
	//***********************

	result_action, err := m.GetConfigTokenString("result_action", m.Model, false)
	if err != nil {
		return err
	}

	result_format, err := m.GetConfigTokenString("result_format", m.Model, false)
	if err != nil {
		return err
	}

	//*************
	//Nothing to do
	//*************
	if result_action == "" || result_action == "none" {
		return nil
	}
	//*************************
	//Create the gout formatter
	//*************************
	w, err := gout.New()
	if err != nil {
		panic(err)
	}

	//***************
	//Format the data
	//***************
	if result_format != "" && result_format != "none" {
		switch strings.ToLower(result_format) {
		case "json":
			w.SetFormatter(json.Formatter{})
		case "yaml":
			w.SetFormatter(yaml.Formatter{})
		case "toml":
			w.SetFormatter(toml.Formatter{})
		case "xml":
			w.SetFormatter(xml.Formatter{})
		case "plain":
			w.SetFormatter(plain.Formatter{})
		default:
			w.SetFormatter(json.Formatter{})
		}

		b := new(strings.Builder)
		w.SetWriter(b)
		err = w.Print(data)
		if err != nil {
			return err
		}
		data = b.String()
	}

	//*****************
	//Print the results
	//*****************
	if result_action == "print" || result_action == "default" || result_action == "" {
		fmt.Println(data)
		return nil
	}

	//***********************
	//What other action to do
	//***********************
	if result_action == "js" {
		//**************
		//run js results
		//**************
		result_js, err := m.GetConfigTokenString("result_js", m.Model, true)
		if err != nil {
			return err
		}
		result_js_path := result_js
		//If code will use else it will get overwritten
		js_code := result_js
		//*******************
		//Test if file exists
		//*******************
		if lib.FileExists(result_js_path) {
			//*************
			//Read the code
			//*************
			js_code, err = lib.ReadFileToString(result_js_path)
			if err != nil {
				return err
			}
		}

		//**********
		//Run the js
		//**********
		vm := m.createJSEngine()
		_, err = vm.RunString(js_code)
		if err != nil {
			return err
		}
		action_results, ok := goja.AssertFunction(vm.Get("ActionResults"))
		if !ok {
			return errors.New("no function found 'ActionResults'(model,result)")
		}
		res, err := action_results(goja.Undefined(), vm.ToValue(m.Model), vm.ToValue(data))

		if err != nil {
			return err
		}
		if !res.ToBoolean() {
			return errors.New("ActionResults returned false")
		}
	} else {
		return fmt.Errorf("unknown result_action %s", result_action)
	}
	return nil
}
