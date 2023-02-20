package workflow

import (
	"fmt"
	"strings"

	"github.com/Mrpye/golib/lib"
	"github.com/dop251/goja"
)

/*func (m *Workflow) action_RunJS() error {
	//*********************
	//Get the config values
	//*********************
	js_file, err := m.GetConfigTokenString("js_file", m.Model, true)
	if err != nil {
		return err
	}
	//*****************************
	//See if we have multiple files
	//*****************************
	files := strings.Split(js_file, ";")
	code := ""
	for _, o := range files {
		file_data, err := lib.ReadFileToString(o)
		if err != nil {
			return err
		}
		code = file_data + "\n"
	}
	//************
	//Run our code
	//************
	vm := m.createJSEngine()
	vm.Set("model", m.Model)
	_, err = vm.RunString(code)
	if err != nil {
		return err
	}
	return nil
}*/

//Set up the js engine
func (m *Workflow) createJSEngine() *goja.Runtime {

	//************************************************
	//Setup the js engine with all the passed function
	//************************************************
	vm := goja.New()
	//vm.Set("target_config", m.GetTargetConfig)
	vm.Set("get_param", m.GetParamValue)
	vm.Set("console", fmt.Println)
	vm.Set("print_ok", lib.PrintOK)
	vm.Set("print_fail", lib.PrintFail)
	vm.Set("action_log", lib.ActionLog)
	vm.Set("action_log_ok", lib.ActionLogOK)
	//vm.Set("action_log_fail", lib.ActionLogFail)
	//vm.Set("file_path", lib.FilePath)
	vm.Set("clear_bucket", m.ClearDataBucket)
	vm.Set("store_value", m.SetValueToDataBucket)
	vm.Set("get_string", m.GetValueFromDataBucketAsStrings)
	vm.Set("get_int", m.GetValueFromDataBucketAsInt)
	vm.Set("get_bool", m.GetValueFromDataBucketAsBool)
	vm.Set("get_store", m.GetValueFromDataBucket)
	vm.Set("domain", lib.GetDomainOrIP)
	vm.Set("port_string", lib.GetPortString)
	vm.Set("port_int", lib.GetPortInt)
	vm.Set("clean", lib.Clean)
	vm.Set("clean_space", strings.ReplaceAll)
	return vm
}
