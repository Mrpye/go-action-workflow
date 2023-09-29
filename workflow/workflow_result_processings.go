package workflow

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Mrpye/golib/file"
	"github.com/Mrpye/golib/log"
	"github.com/dop251/goja"
	"github.com/drewstinnett/gout/v2"
	"github.com/drewstinnett/gout/v2/formats/json"
	"github.com/drewstinnett/gout/v2/formats/plain"
	"github.com/drewstinnett/gout/v2/formats/toml"
	"github.com/drewstinnett/gout/v2/formats/xml"
	"github.com/drewstinnett/gout/v2/formats/yaml"
)

// ActionProcessResults processes the results of the action
// - m is the template data model
// - data is the data to process
// - returns an error if there is an error
func (w *Workflow) ActionProcessResults(m *TemplateData, data interface{}) error {

	//************************
	// what to do with results
	//************************
	result_action, err := w.GetConfigTokenString("result_action", m, false)
	if err != nil {
		return err
	}

	//**************************
	// how to format the results
	//**************************
	result_format, err := w.GetConfigTokenString("result_format", m, false)
	if err != nil {
		return err
	}
	if w.LogLevel == LOG_VERBOSE {
		log.LogVerbose(fmt.Sprintf("result_action: %s\n", result_action))
		log.LogVerbose(fmt.Sprintf("result_format: %s\n", result_format))
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
	g, err := gout.New()
	if err != nil {
		return err
	}

	//*******************************************************
	//Format the data in the format specified default is json
	//*******************************************************
	if result_format != "" && result_format != "none" {
		switch strings.ToLower(result_format) {
		case "json":
			g.SetFormatter(json.Formatter{})
		case "yaml":
			g.SetFormatter(yaml.Formatter{})
		case "toml":
			g.SetFormatter(toml.Formatter{})
		case "xml":
			g.SetFormatter(xml.Formatter{})
		case "plain":
			g.SetFormatter(plain.Formatter{})
		default:
			g.SetFormatter(json.Formatter{})
		}

		//******************************
		//Format the results to a string
		//******************************
		b := new(strings.Builder)
		g.SetWriter(b)
		err = g.Print(data)
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

	//************************
	// Run the js results code
	//************************
	if result_action == "js" {

		//***********************
		//get the js code or file
		//***********************
		js_code, err := w.GetConfigTokenString("result_js", m, true)
		if err != nil {
			return err
		}

		//****************************
		//Test if it is a file or code
		//****************************
		if file.FileExists(w.BuildPath(js_code)) {
			//*************
			//Read the code
			//*************
			js_code, err = file.ReadFileToString(js_code)
			if err != nil {
				return err
			}
		}

		//*****************
		// Create JS engine
		//*****************
		vm := w.CreateJSEngine()
		_, err = vm.RunString(js_code)
		if err != nil {
			return err
		}

		//************************************************
		// Get the ActionResults function from the js code
		//************************************************
		action_results, ok := goja.AssertFunction(vm.Get("ActionResults"))
		if !ok {
			return errors.New("no function found 'ActionResults'(model,result)")
		}

		//********************************************************
		// Call the ActionResults function with the model and data
		//********************************************************
		res, err := action_results(goja.Undefined(), vm.ToValue(m), vm.ToValue(data))
		if err != nil {
			return err
		}
		//********************************
		// Check if the js returned false
		//********************************
		if !res.ToBoolean() {
			return errors.New("ActionResults JS returned false")
		}
	} else {
		return fmt.Errorf("unknown result_action %s", result_action)
	}
	return nil
}
