package main

import (
	"bytes"
	"reflect"
	"strings"
	"text/template"

	"github.com/Mrpye/go-action-workflow/actions/action_docker"
	"github.com/Mrpye/go-action-workflow/actions/action_k8"
	"github.com/Mrpye/go-action-workflow/workflow"
	"github.com/Mrpye/golib/file"
)

type DocGen struct {
	TargetScheme   map[string]*workflow.TargetSchema   //map of target schemas
	ActionScheme   map[string]*workflow.ActionSchema   //map of action schemas
	FunctionScheme map[string]*workflow.FunctionSchema //map of action schemas
	Template       string
	Output         string
}

func CreateWorkflow() *DocGen {
	doc := &DocGen{}
	doc.TargetScheme = make(map[string]*workflow.TargetSchema)
	doc.ActionScheme = make(map[string]*workflow.ActionSchema)
	doc.FunctionScheme = make(map[string]*workflow.FunctionSchema)
	return doc
}

func (m *DocGen) GetFunctionsByLibrary(lib string) map[string]*workflow.FunctionSchema {
	grouped := make(map[string]*workflow.FunctionSchema)
	for key, f := range m.FunctionScheme {
		parts := strings.Split(f.Target, ".")
		group := parts[0]
		if group == "" {
			group = "General"
		}
		if strings.EqualFold(group, lib) {
			grouped[key] = f
		}
	}
	return grouped
}

func (m *DocGen) GetGroupedActions() map[string]map[string]*workflow.ActionSchema {
	grouped := make(map[string]map[string]*workflow.ActionSchema)
	for key, action := range m.ActionScheme {
		parts := strings.Split(action.Target, ".")
		group := parts[0]
		if group == "" {
			group = "General"
		}
		if _, ok := grouped[group]; !ok {
			grouped[group] = make(map[string]*workflow.ActionSchema)
		}
		grouped[group][key] = action
	}
	return grouped
}

func GetLib(TargetName string) string {
	parts := strings.Split(TargetName, ".")
	return parts[0]
}

func getStructTag(f reflect.StructField, tagName string) string {
	return string(f.Tag.Get(tagName))
}

func GetTag(t interface{}, field_name string, tag_name string) string {
	field, ok := reflect.TypeOf(t).Elem().FieldByName(field_name)
	if !ok {
		panic("Field not found")
	}
	return getStructTag(field, tag_name)
}

func GetObjectName(t interface{}) string {
	type_name := strings.ReplaceAll(strings.ToLower(reflect.TypeOf(t).String()), "*", "")
	return type_name
}

func (m *DocGen) AddActionSchema(sch workflow.SchemaEndpoint) {

	//***********************************
	//Add the target schema to the client
	//***********************************
	scheme := sch.GetTargetSchema()
	for key := range scheme {
		target_schema := scheme[key]
		m.TargetScheme[key] = &target_schema
	}

	//***********************************
	//Add the function map to the client
	//***********************************
	function_schemes := sch.GetFunctionMap()
	for key := range function_schemes {
		function_scheme := function_schemes[key]
		m.FunctionScheme[key] = &function_scheme
	}

	//********************************
	//Add the actions to the workflow
	//********************************
	actions := sch.GetActionSchema()
	for key := range actions {
		action_schema := actions[key]
		m.ActionScheme[key] = &action_schema
	}

}

func (m *DocGen) funcMap() template.FuncMap {
	return template.FuncMap{
		"lc":       strings.ToLower, //Lowercase a string
		"uc":       strings.ToUpper, //Uppercase a string
		"tag":      GetTag,
		"lib":      GetLib,
		"replace":  strings.ReplaceAll,
		"obj_name": GetObjectName,
	}
}
func (m *DocGen) BuildDoc() error {
	template_file, err := file.ReadFileToString(m.Template)
	if err != nil {
		return err
	}
	//********************************
	//Create a new template and parse
	//********************************
	tmpl, err := template.New("CodeRun").Funcs(m.funcMap()).Parse(template_file)
	if err != nil {
		return err
	}

	//**************************************
	//Run the template to verify the output.
	//**************************************
	var tpl bytes.Buffer
	err = tmpl.Execute(&tpl, m)
	if err != nil {
		return err
	}

	//*************
	//Save the file
	//*************
	err = file.SaveStringToFile(m.Output, tpl.String())
	if err != nil {
		return err
	}

	//******************
	//Return the result
	//******************
	return nil
}

func main() {
	// Generate the documentation for the package
	// and write it to the file "doc.go"
	//doc.GenDoc("doc.go", "github.com/yourname/yourpackage")
	doc := CreateWorkflow()
	doc.Template = "action_doc_template.md"
	doc.Output = "action_doc.md"
	doc.AddActionSchema(action_k8.GetSchema())
	doc.AddActionSchema(action_docker.GetSchema())
	err := doc.BuildDoc()
	if err != nil {
		panic(err)
	}

}
