package workflow

import (
	"fmt"
	"strings"

	"github.com/Mrpye/golib/lib"
	"github.com/dop251/goja"
)

//Set up the js engine
func (m *Workflow) CreateJSEngine() *goja.Runtime {

	//************************************************
	//Setup the js engine with all the passed function
	//************************************************
	vm := goja.New()
	//vm.Set("target_config", m.GetTargetConfig)
	//vm.Set("read_file", lib.ReadFileToString)
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
