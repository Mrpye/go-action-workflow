// This is the workflow action for calling an api
package action_docker

import (
	"fmt"
	"io/ioutil"
	slog "log"
	"path"
	"strings"

	"github.com/Mrpye/go-action-workflow/workflow"
	"github.com/Mrpye/golib/log"
	dockerparser "github.com/novln/docker-parser"
)

type DockerRegSchema struct {
}

func (s DockerRegSchema) GetFunctionMap() map[string]workflow.FunctionSchema {
	//Params w *workflow.Workflow, image string, target_name string, no_tag bool, use_original bool, original string
	return map[string]workflow.FunctionSchema{
		"image_name": {
			Cmd:         "image_name",
			Description: "gets the name of the image from the image path",
			Target:      "action_docker",
			Function:    ImageName,
			ParameterSchema: map[string]*workflow.Schema{
				"image": {
					Type:        workflow.TypeString,
					Required:    true,
					Description: "docker.io/circleci/[slim-base]:latest",
				},
			},
		},
		"image_name_tag": {
			Cmd:         "image_name_tag",
			Description: "gets the name and tag from the image path",
			Target:      "action_docker",
			Function:    ImageNameTag,
			ParameterSchema: map[string]*workflow.Schema{
				"image": {
					Type:        workflow.TypeString,
					Required:    true,
					Description: "docker.io/circleci/[slim-base:latest]",
				},
			},
		},
		"image_account": {
			Cmd:         "image_account",
			Description: "gets the account name from the image path",
			Target:      "action_docker",
			Function:    ImageAccount,
			ParameterSchema: map[string]*workflow.Schema{
				"image": {
					Type:        workflow.TypeString,
					Required:    true,
					Description: "docker.io/[circleci]/slim-base:latest",
				},
			},
		},
		"image_shortname": {
			Cmd:         "image_shortname",
			Description: "gets the account name from the image path",
			Target:      "action_docker",
			Function:    ImageShortName,
			ParameterSchema: map[string]*workflow.Schema{
				"image": {
					Type:        workflow.TypeString,
					Required:    true,
					Description: "docker.io/[circleci/slim-base]:latest",
				},
			},
		},
		"image_registry": {
			Cmd:         "image_registry",
			Description: "gets the registry name from the image path",
			Target:      "action_docker",
			Function:    ImageRegistry,
			ParameterSchema: map[string]*workflow.Schema{
				"image": {
					Type:        workflow.TypeString,
					Required:    true,
					Description: "docker.io/circleci/slim-base]:latest",
				},
			},
		},
		"image_tag": {
			Cmd:         "image_tag",
			Description: "gets the tag  from the image path",
			Target:      "action_docker",
			Function:    ImageTag,
			ParameterSchema: map[string]*workflow.Schema{
				"image": {
					Type:        workflow.TypeString,
					Required:    true,
					Description: "docker.io/circleci/slim-base:[latest]",
				},
			},
		},
		"remap_image": {
			Cmd:         "remap_image",
			Description: "remaps the docker image",
			Target:      "action_docker",
			Function:    RemapImage,
			ParameterSchema: map[string]*workflow.Schema{
				"*workflow.Workflow": {
					Type:        workflow.TypeWorkflow,
					Required:    true,
					Description: "work flow engine use (get_wf) to get the workflow engine",
				},
				"image": {
					Type:        workflow.TypeString,
					Required:    true,
					Description: "name of the image",
				},
				"target_name": {
					Type:        workflow.TypeString,
					Required:    false,
					Description: "name of the target or just use ``",
				},
				"no_tag": {
					Type:        workflow.TypeBool,
					Required:    true,
					Description: "true not include the tag in the returned value",
				},
				"use_original": {
					Type:        workflow.TypeBool,
					Required:    true,
					Description: "use the original image path",
				},
				"original": {
					Type:        workflow.TypeString,
					Required:    false,
					Description: "the original image path to use if use_original=true else build path based on target",
				},
			},
		},
		"remap_image2": {
			Cmd:         "remap_image2",
			Description: "remaps the docker image no option for using default image path",
			Target:      "action_docker",
			Function:    RemapImage2,
			ParameterSchema: map[string]*workflow.Schema{
				"*workflow.Workflow": {
					Type:        workflow.TypeWorkflow,
					Required:    true,
					Description: "work flow engine use (get_wf) to get the workflow engine",
				},
				"image": {
					Type:        workflow.TypeString,
					Required:    true,
					Description: "name of the image",
				},
				"target_name": {
					Type:        workflow.TypeString,
					Required:    false,
					Description: "name of the target or just use ``",
				},
				"no_tag": {
					Type:        workflow.TypeBool,
					Required:    true,
					Description: "true not include the tag in the returned value",
				},
			},
		},
	}
}

// GetTargetSchema returns the target schema for this action
func (s DockerRegSchema) GetTargetSchema() map[string]workflow.TargetSchema {
	//Build a test schema
	return map[string]workflow.TargetSchema{
		"action_docker.registry": workflow.BuildTargetConfig("Registry", "Registry", &Registry{}),
	}
}

// GetActions returns the actions for this schema
func (s DockerRegSchema) GetActionSchema() map[string]workflow.ActionSchema {

	//no short d f h
	return map[string]workflow.ActionSchema{
		"docker_reg_download": {
			Action:         Action_DockerRegDownloadImages,
			Short:          "Download a docker image from a registry",
			Long:           "Download a docker image from a registry",
			Target:         "action_docker.registry",
			ProcessResults: false,
			ConfigSchema: map[string]*workflow.Schema{
				"target_name": {
					Type:        workflow.TypeString,
					Partial:     false,
					Required:    false,
					Description: "The target to use",
					Short:       "t",
					Default:     "",
				},
				"images": {
					Type:        workflow.TypeList,
					Partial:     false,
					Required:    true,
					Description: "List of images to download",
					Short:       "i",
					Default:     []string{},
				},
				"folder": {
					Type:        workflow.TypeString,
					Partial:     false,
					Required:    true,
					Description: "Folder to save the images to",
					Short:       "f",
					Default:     "",
				},
			},
		},
		"docker_reg_upload": {
			Action:         Action_DockerRegUploadImages,
			Short:          "Upload a docker image to a registry",
			Long:           "Upload a docker image to a registry",
			ProcessResults: false,
			Target:         "action_docker.registry",
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
				"images": {
					Type:        workflow.TypeList,
					Partial:     false,
					Required:    true,
					Description: "List of images to upload",
					Short:       "i",
					Default:     []string{},
				},
				"folder": {
					Type:        workflow.TypeString,
					Partial:     false,
					Required:    true,
					Description: "Folder to load the images from",
					Short:       "f",
					Default:     "",
				},
				"import_all": {
					Type:        workflow.TypeBool,
					Partial:     false,
					Required:    false,
					Description: "Import all images in the folder",
					Short:       "a",
					Default:     false,
				},
			},
		},
	}
}

// GetAction returns the action for the target
func GetSchema() workflow.SchemaEndpoint {
	return DockerRegSchema{}
}

func Action_DockerRegDownloadImages(w *workflow.Workflow, m *workflow.TemplateData) error {

	docker_reg_obj, err := w.MapTargetConfigValue(m, &Registry{})
	if err != nil {
		return err
	}
	docker_reg := docker_reg_obj.(*Registry)

	//**************
	//Get the images
	//**************
	images, err := w.GetConfigTokenInterfaceArray("images", m, true)
	if err != nil {
		return err
	}

	//***********************************
	//Get the folder where to save images
	//***********************************
	folder, err := w.GetConfigTokenString("folder", m, true)
	if err != nil {
		return err
	}
	//******************************
	//if folder is blank, use images
	//this is the default
	//******************************
	if folder == "" {
		folder = "images"
	}

	//***************
	//Download images
	//***************
	for _, image := range images {

		if w.LogLevel >= workflow.LOG_INFO {

			slog.Printf("Downloading image: %s \n", image)

			//log.LogVerbose(fmt.Sprintf("image: %s\n", image))
		}

		//**************************
		//Create the image file name
		//**************************
		parse, err := dockerparser.Parse(image.(string))
		if err != nil {
			return err
		}
		parts := strings.Split(parse.ShortName(), "/")
		//get the last part
		image_name := parts[len(parts)-1]
		image_file_name := fmt.Sprintf("%s_%s.tar", image_name, parse.Tag())
		image_file_name = path.Join(folder, image_file_name)

		//**************
		//Download image
		//**************
		err = docker_reg.Download(image.(string), w.BuildPath(image_file_name))
		if err != nil {
			return err
		}

		if w.LogLevel >= workflow.LOG_INFO {
			log.PrintlnOK(fmt.Sprintf("Downloaded image: %s", image))
		}

	}

	//**************
	//Return the err
	//**************
	return nil
}

func Action_DockerRegUploadImages(w *workflow.Workflow, m *workflow.TemplateData) error {

	docker_reg_obj, err := w.MapTargetConfigValue(m, &Registry{})
	if err != nil {
		return err
	}
	docker_reg := docker_reg_obj.(*Registry)

	//***********************************
	//Get the folder where to save images
	//***********************************
	import_all, err := w.GetConfigTokenBool("import_all", m, false)
	if err != nil {
		return err
	}

	//**************
	//Get the images
	//**************
	images, err := w.GetConfigTokenInterfaceArray("images", m, !import_all)
	if err != nil {
		return err
	}

	//***********************************
	//Get the folder where to save images
	//***********************************
	folder, err := w.GetConfigTokenString("folder", m, true)
	if err != nil {
		return err
	}

	//******************************
	//if folder is blank, use images
	//this is the default
	//******************************
	if folder == "" {
		folder = "images"
	}

	if import_all {
		//***********************************
		//Jut import all images in the folder
		//***********************************
		files, err := ioutil.ReadDir(w.BuildPath(folder))
		if err != nil {
			return err
		}
		for _, f := range files {
			if strings.HasSuffix(f.Name(), ".tar") {
				//if w.LogLevel == workflow.LOG_INFO {
				//	slog.Printf("Uploading image: %s\n", f.Name())
				//}

				registry_target, err := docker_reg.Upload(w.BuildPath(path.Join(folder, f.Name())), w.LogLevel)
				if err != nil {
					return err
				}
				if w.LogLevel == workflow.LOG_INFO {
					log.PrintlnOK(fmt.Sprintf("Uploaded image: %s,%s", f.Name(), registry_target))
				}
			}
		}
	} else {
		//*************
		//Upload images
		//*************
		for _, image := range images {

			image_n := strings.TrimSuffix(image.(string), ".tar")
			//**************************
			//Create the image file name
			//**************************
			parse, err := dockerparser.Parse(image_n)
			if err != nil {
				return err
			}
			parts := strings.Split(parse.ShortName(), "/")
			//get the last part
			image_name := parts[len(parts)-1]
			image_file_name := fmt.Sprintf("%s_%s.tar", image_name, parse.Tag())
			image_file_name = path.Join(folder, image_file_name)

			//if w.LogLevel == workflow.LOG_INFO {
			//	slog.Printf("Uploading image tar: %s\n", image_file_name)
			//}

			//**************
			//Upload image
			//**************
			registry_target, err := docker_reg.Upload(w.BuildPath(image_file_name), w.LogLevel)
			if err != nil {
				return err
			}
			if w.LogLevel == workflow.LOG_INFO {
				log.PrintlnOK(fmt.Sprintf("Uploaded image: %s,%s", image_file_name, registry_target))
			}
		}
	}

	//**************
	//Return the err
	//**************
	return nil
}
