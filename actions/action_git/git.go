// This is the workflow action for calling an api
package action_git

import (
	"fmt"
	"io/ioutil"
	slog "log"
	"strings"

	"github.com/Mrpye/go-action-workflow/workflow"
	"github.com/Mrpye/golib/dir"
	"github.com/Mrpye/golib/log"
)

type GitSchema struct {
}

// GetTargetSchema returns the target schema for this action
func (s GitSchema) GetTargetSchema() map[string]workflow.TargetSchema {
	//Build a test schema
	return map[string]workflow.TargetSchema{
		"action_git.git": workflow.BuildTargetConfig("git", "git", &Git{}),
	}
}

// GetActions returns the actions for this schema
func (s GitSchema) GetActionSchema() map[string]workflow.ActionSchema {

	//no short d f h
	return map[string]workflow.ActionSchema{
		"git_download": {
			Action:         Action_GitDownload,
			Short:          "Apply or delete a yaml file to a k8 cluster",
			Long:           "Apply or delete a yaml file to a k8 cluster",
			ProcessResults: false,
			Target:         "action_git.git",
			ConfigSchema: map[string]*workflow.Schema{
				"target_name": {
					Type:        workflow.TypeString,
					Partial:     false,
					Required:    false,
					Description: "The target to use",
					ConfigKey:   "target_name",
					Short:       "t",
					Default:     "",
				},
				"files": {
					Type:        workflow.TypeList,
					Partial:     false,
					Required:    false,
					Description: "List of files to download",
					Short:       "f",
					Default:     []string{},
				},
			},
		},
	}
}

// GetFunctionMap returns the function map for this schema
func (s GitSchema) GetFunctionMap() map[string]workflow.FunctionSchema {
	return map[string]workflow.FunctionSchema{}
}

// GetAction returns the action for the target
func GetSchema() workflow.SchemaEndpoint {
	return GitSchema{}
}

func Action_GitDownload(w *workflow.Workflow, m *workflow.TemplateData) error {

	//*************************
	//Create the git client
	//And map the config values
	//*************************
	git_obj, err := w.MapTargetConfigValue(m, &Git{})
	if err != nil {
		return err
	}
	git := git_obj.(*Git) //cast it as a git type

	//*******************
	//Get the files list
	//*******************
	files, _ := w.GetConfigTokenMapArray("files", m, true)

	for _, file := range files {

		cfg_file := ""
		cfg_project := ""
		cfg_service := ""
		cfg_dest := ""
		cfg_branch := ""
		cfg_short := ""
		//*******************
		//Get the file config
		//*******************
		if file["short"] != nil {
			//********************************
			//short file config
			//service;project;file;dest;branch
			//********************************
			cfg_short = file["short"].(string)
			parts := strings.Split(cfg_short, ";")
			if len(parts) != 5 {
				return fmt.Errorf("short file config must have 5 parts: service;project;file;dest;branch")
			}
			cfg_service = parts[0]
			cfg_project = parts[1]
			cfg_file = parts[2]
			cfg_dest = parts[3]
			cfg_branch = parts[4]

		} else {
			if file["service"] != nil {
				cfg_service = file["service"].(string)
			}
			if file["project"] != nil {
				cfg_project = file["project"].(string)
			}
			if file["file"] != nil {
				cfg_file = file["file"].(string)
			}
			if file["dest"] != nil {
				cfg_dest = file["dest"].(string)
			}
			if file["branch"] != nil {
				cfg_branch = file["branch"].(string)
			}
		}
		if w.LogLevel >= workflow.LOG_VERBOSE {
			log.LogVerbose(fmt.Sprintf("file: %s\n", cfg_file))
			log.LogVerbose(fmt.Sprintf("project: %s\n", cfg_project))
			log.LogVerbose(fmt.Sprintf("service: %s\n", cfg_service))
			log.LogVerbose(fmt.Sprintf("dest: %s\n", cfg_dest))
			log.LogVerbose(fmt.Sprintf("branch: %s\n", cfg_branch))
		}

		//download the file
		if w.LogLevel >= workflow.LOG_INFO {
			slog.Printf("Downloading file: %s from %s branch %s\n", cfg_file, cfg_project, cfg_branch)
		}

		data, err := git.DownloadGitFile(cfg_service, cfg_project, cfg_file, cfg_branch)
		if err != nil {
			return err
		}
		//************************
		//Make sure the dir exists
		//************************
		err = dir.MakeDirAll(w.BuildPath(cfg_dest))
		if err != nil {
			return err
		}

		//**************
		//Write the file
		//**************
		err = ioutil.WriteFile(w.BuildPath(cfg_dest), []byte(data), 0644)
		if err != nil {
			return err
		}
		if w.LogLevel >= workflow.LOG_INFO {
			log.PrintlnOK(fmt.Sprintf("Downloaded file: %s from %s branch %s", cfg_file, cfg_project, cfg_branch))
		}

	}

	//**************
	//Return the err
	//**************
	return nil
}
