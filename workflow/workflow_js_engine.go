package workflow

import (
	"fmt"
	"log"
	"strings"

	go_log "github.com/Mrpye/golib/log"
	"github.com/Mrpye/golib/str"
	"github.com/dop251/goja"
)

// CreateJSEngine creates a new js engine
// this is used to run the js code
// returns a new js engine
func (m *Workflow) CreateJSEngine() *goja.Runtime {

	//************************************************
	//Setup the js engine with all the passed function
	//************************************************
	vm := goja.New()
	//vm.Set("target_config", m.GetTargetConfig)
	//vm.Set("read_file", lib.ReadFileToString)
	vm.Set("get_param", m.GetParamValue)
	vm.Set("get_input", m.GetInputValue)
	vm.Set("console", fmt.Println)
	vm.Set("log", log.Println)
	vm.Set("print_ok", go_log.PrintOK)
	vm.Set("print_fail", go_log.PrintFail)
	vm.Set("action_log", go_log.ActionLog)
	vm.Set("action_log_ok", go_log.ActionLogOK)
	//vm.Set("action_log_fail", lib.ActionLogFail)
	//vm.Set("file_path", lib.FilePath)
	vm.Set("clear_bucket", m.ClearDataBucket)
	vm.Set("store_value", m.SetValueToDataBucket)
	vm.Set("get_string", m.GetValueFromDataBucketAsStrings)
	vm.Set("get_int", m.GetValueFromDataBucketAsInt)
	vm.Set("get_bool", m.GetValueFromDataBucketAsBool)
	vm.Set("get_store", m.GetValueFromDataBucket)
	vm.Set("domain", str.GetDomainOrIP)
	vm.Set("port_string", str.GetPortString)
	vm.Set("port_int", str.GetPortInt)
	vm.Set("clean", str.Clean)
	vm.Set("clean_space", strings.ReplaceAll)
	vm.Set("log_title_green", go_log.ActionLogGreen)
	vm.Set("log_title_red", go_log.ActionLogRed)
	return vm
}
