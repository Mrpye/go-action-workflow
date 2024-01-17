package action_file

import (
	"fmt"
	"io/ioutil"

	"github.com/Mrpye/go-action-workflow/workflow"
	"github.com/Mrpye/golib/dir"
	"github.com/Mrpye/golib/file"
)

// Action_Copy is a custom action that copies a file
/*
- action: ActionStore
	config:
		source_file: the file to copy
		dest_file: 	the destination file
*/
func Action_Template(w *workflow.Workflow, m *workflow.TemplateData) error {

	//*************
	//Template file
	//*************
	template_file, err := w.GetConfigTokenString("template", m, true)
	if err != nil {
		return err
	}
	template_file = w.BuildPath(template_file)

	//**************
	//Where to write
	//**************
	dest_file, err := w.GetConfigTokenString("file", m, true)
	if err != nil {
		return err
	}
	dest_file = w.BuildPath(dest_file)

	//************************************
	//Get the data that has been passed in
	//************************************
	data, err := w.GetConfigTokenMap("data", m, false)
	if err != nil {
		return err
	}
	m.Data = data

	if w.LogLevel > workflow.LOG_QUIET {
		fmt.Printf("Template: %s", template_file)
		fmt.Printf("Destination: %s", template_file)
		fmt.Printf("Data: %s", template_file)
	}

	//**********************
	//Read the template file
	//**********************
	template, err := file.ReadFileToString(template_file)
	if err != nil {
		return err
	}

	//*********************
	//Render the template
	//*********************
	result, err := w.ParseToken(m, template)
	if err != nil {
		return err
	}

	//**************
	//Write the file
	//**************
	dir.MakeDirAll(dest_file)
	err = ioutil.WriteFile(dest_file, []byte(result), 0644)
	if err != nil {
		return err
	}

	return nil
}
