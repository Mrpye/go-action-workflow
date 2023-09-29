// This is the go workflow action plugin package for managing k8 cluster
package action_k8

import (
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/Mrpye/go-action-workflow/workflow"
	"github.com/Mrpye/go_k8_helm"

	"github.com/Mrpye/golib/convert"
	"github.com/Mrpye/golib/log"
	"github.com/Mrpye/golib/str"
	"gopkg.in/yaml.v2"
)

type K8Schema struct {
}

// GetTargetSchema returns the target schema for this action
func (s K8Schema) GetTargetSchema() map[string]workflow.TargetSchema {
	//Build a test schema
	return map[string]workflow.TargetSchema{
		"action_k8.k8": workflow.BuildTargetConfig("k8", "k8", &go_k8_helm.K8{}),
	}
}

// GetActions returns the actions for this schema
func (s K8Schema) GetActionSchema() map[string]workflow.ActionSchema {

	//no short d f h
	return map[string]workflow.ActionSchema{
		"k8_yaml": {
			Action:         Action_K8ApplyDeleteYaml,
			Short:          "Apply or delete a yaml file to a k8 cluster",
			Long:           "Apply or delete a yaml file to a k8 cluster",
			Target:         "action_k8.k8",
			ProcessResults: false,
			ConfigSchema: map[string]*workflow.Schema{
				"target_name": {
					Type:        workflow.TypeString,
					Partial:     false,
					Required:    false,
					Description: "The target name to use if not default target type",
					ConfigKey:   "target_name",
					Short:       "t",
					Default:     "",
				},
				"namespace": {
					Type:        workflow.TypeString,
					Partial:     false,
					Required:    false,
					Description: "Namespace to use",
					Short:       "s",
					Default:     "default",
				},
				"process_tokens": {
					Type:        workflow.TypeBool,
					Partial:     false,
					Required:    false,
					Description: "Should tokens be processed in the k8 manifest",
					Short:       "p",
					Default:     true,
				},
				"delete": {
					Type:        workflow.TypeBool,
					Partial:     false,
					Required:    false,
					Description: "If true, delete the deployment",
					Short:       "d",
					Default:     false,
				},
				"manifest": {
					Type:        workflow.TypeString,
					Partial:     false,
					Required:    true,
					Description: "k8 manifest to apply or delete as a file path or object",
					Short:       "m",
					Default:     "",
				},
			},
		},
		"k8_create_ns": {
			Action:         Action_K8CreateNS,
			Short:          "Create a namespace",
			Long:           "Create a namespace",
			Target:         "action_k8.k8",
			ProcessResults: false,
			ConfigSchema: map[string]*workflow.Schema{
				"target_name": {
					Type:        workflow.TypeString,
					Partial:     false,
					Required:    false,
					Description: "The target name to use if not default target type",
					Short:       "t",
					Default:     "",
				},
				"namespace": {
					Type:        workflow.TypeString,
					Partial:     false,
					Required:    false,
					Description: "Namespace to use",
					Short:       "s",
					Default:     "default",
				},
			},
		},
		"k8_delete_ns": {
			Action:         Action_K8DeleteNS,
			Short:          "Delete a namespace",
			Long:           "Delete a namespace",
			Target:         "action_k8.k8",
			ProcessResults: false,
			ConfigSchema: map[string]*workflow.Schema{
				"target_name": {
					Type:        workflow.TypeString,
					Partial:     false,
					Required:    false,
					Description: "The target name to use if not default target type",
					Short:       "t",
					Default:     "",
				},
				"namespace": {
					Type:        workflow.TypeString,
					Partial:     false,
					Required:    false,
					Description: "Namespace to use",
					Short:       "s",
					Default:     "default",
				},
			},
		},
		"k8_copy": {
			Action:         Action_K8Copy,
			Short:          "Copy a file from a pod to the local machine or vice versa",
			Long:           "Copy a file from a pod to the local machine or vice versa",
			ProcessResults: true,
			Target:         "action_k8.k8",
			ConfigSchema: map[string]*workflow.Schema{
				"target_name": {
					Type:        workflow.TypeString,
					Partial:     false,
					Required:    false,
					Description: "The target name to use if not default target type",
					Short:       "t",
					Default:     "",
				},
				"namespace": {
					Type:        workflow.TypeString,
					Partial:     false,
					Required:    false,
					Description: "Namespace to use",
					Short:       "s",
					Default:     "default",
				},
				"src": {
					Type:        workflow.TypeString,
					Partial:     false,
					Required:    true,
					Description: "source file path",
					Short:       "f",
					Default:     "",
				},
				"dest": {
					Type:        workflow.TypeString,
					Partial:     false,
					Required:    true,
					Description: "destination file path",
					Short:       "d",
					Default:     "",
				},
				"container_name": {
					Type:        workflow.TypeString,
					Partial:     false,
					Required:    true,
					Description: "container name",
					Short:       "c",
					Default:     "",
				},
			},
		},
		"k8_pod_exec": {
			Action:         Action_K8PodExec,
			Short:          "Execute a command in a pod",
			Long:           "Execute a command in a pod",
			ProcessResults: true,
			Target:         "action_k8.k8",
			ConfigSchema: map[string]*workflow.Schema{
				"target_name": {
					Type:        workflow.TypeString,
					Partial:     false,
					Required:    false,
					Description: "The target name to use if not default target type",
					Short:       "t",
					Default:     "",
				},
				"namespace": {
					Type:        workflow.TypeString,
					Partial:     false,
					Required:    false,
					Description: "Namespace to use",
					Short:       "s",
					Default:     "default",
				},
				"pod_name": {
					Type:        workflow.TypeString,
					Partial:     false,
					Required:    true,
					Description: "pod name",
					Short:       "n",
					Default:     "",
				},
				"command": {
					Type:        workflow.TypeString,
					Partial:     false,
					Required:    true,
					Description: "command to execute",
					Short:       "c",
					Default:     "",
				},
			},
		},
		"k8_get_service_ip": {
			Action:         Action_K8GetServiceIP,
			Short:          "Get the IP of a service",
			Long:           "Get the IP of a service",
			ProcessResults: true,
			Target:         "action_k8.k8",
			ConfigSchema: map[string]*workflow.Schema{
				"target_name": {
					Type:        workflow.TypeString,
					Partial:     false,
					Required:    false,
					Description: "The target name to use if not default target type",
					Short:       "t",
					Default:     "",
				},
				"namespace": {
					Type:        workflow.TypeString,
					Partial:     false,
					Required:    false,
					Description: "Namespace to use",
					Short:       "s",
					Default:     "default",
				},
				"name": {
					Type:        workflow.TypeString,
					Partial:     false,
					Required:    true,
					Description: "name of service can use regex",
					Short:       "n",
					Default:     "",
				},
			},
		},
		"k8_get_ws_items": {
			Action:         Action_K8GetWorkspace,
			Short:          "Get items in a workspace",
			Long:           "Get items in a workspace",
			ProcessResults: true,
			Target:         "action_k8.k8",
			ConfigSchema: map[string]*workflow.Schema{
				"target_name": {
					Type:        workflow.TypeString,
					Partial:     false,
					Required:    false,
					Description: "The target name to use if not default target type",
					Short:       "t",
					Default:     "",
				},
				"namespace": {
					Type:        workflow.TypeString,
					Partial:     false,
					Required:    false,
					Description: "Namespace to use",
					Short:       "S",
					Default:     "default",
				},
				"workspace": {
					Type:        workflow.TypeString,
					Partial:     false,
					Required:    true,
					Description: "workspace to use",
					Short:       "w",
					Default:     "",
				},
			},
		},
		"k8_wait": {
			Action:         Action_K8WaitCompleteStatus,
			Short:          "Wait for a k8 resource to be in a complete state",
			Long:           "Wait for a k8 resource to be in a complete state",
			ProcessResults: false,
			Target:         "action_k8.k8",
			ConfigSchema: map[string]*workflow.Schema{
				"target_name": {
					Type:        workflow.TypeString,
					Partial:     false,
					Required:    false,
					Description: "The target name to use if not default target type",
					Short:       "t",
					Default:     "",
				},
				"namespace": {
					Type:        workflow.TypeString,
					Partial:     false,
					Required:    false,
					Description: "Namespace to use",
					Short:       "s",
					Default:     "default",
				},
				"retry": {
					Type:        workflow.TypeInt,
					Partial:     false,
					Required:    false,
					Description: "retry count",
					Short:       "r",
					Default:     10,
				},
				"not_running": {
					Type:        workflow.TypeBool,
					Partial:     false,
					Required:    false,
					Description: "All checks not running",
					Short:       "x",
					Default:     false,
				},
				"checks": {
					Type:        workflow.TypeList,
					Partial:     false,
					Required:    false,
					Description: "Checks to run",
					Short:       "c",
					Default: []string{"replica:nginx2(.*)",
						"stateful:nginx3(.*)",
						"demon:nginx4(.*)",
						"service:nginx(.*)"},
				},
			},
		},
		"k8_helm_deploy_upgrade": {
			Action:         Action_DeployUpgradeHelmChart,
			Short:          "Deploy or upgrade a helm chart",
			Long:           "Deploy or upgrade a helm chart",
			ProcessResults: false,
			Target:         "action_k8.k8",
			ConfigSchema: map[string]*workflow.Schema{
				"target_name": {
					Type:        workflow.TypeString,
					Partial:     false,
					Required:    false,
					Description: "The target name to use if not default target type",
					Short:       "t",
					Default:     "",
				},
				"namespace": {
					Type:        workflow.TypeString,
					Partial:     false,
					Required:    false,
					Description: "Namespace to use",
					Short:       "s",
					Default:     "default",
				},
				"chart_Path": {
					Type:        workflow.TypeString,
					Partial:     false,
					Required:    false,
					Description: "chart path",
					Short:       "c",
					Default:     "",
				},
				"release_name": {
					Type:        workflow.TypeString,
					Partial:     false,
					Required:    false,
					Description: "release name",
					Short:       "n",
					Default:     "",
				},
				/*"helm_config": {
					Type:        workflow.TypeMap,
					Partial:     false,
					Required:    false,
					Description: "helm config",
					Short:       "h",
					Default: map[string]interface{}{
						"name": "test",
					},
				},*/
				"upgrade": {
					Type:        workflow.TypeBool,
					Partial:     false,
					Required:    false,
					Description: "chart path",
					Short:       "u",
					Default:     false,
				},
			},
		},
		"k8_helm_delete": {
			Action:         Action_DeleteHelmChart,
			Short:          "Delete a helm chart",
			Long:           "Delete a helm chart",
			ProcessResults: false,
			Target:         "action_k8.k8",
			ConfigSchema: map[string]*workflow.Schema{
				"target_name": {
					Type:        workflow.TypeString,
					Partial:     false,
					Required:    false,
					Description: "The target name to use if not default target type",
					Short:       "t",
					Default:     "",
				},
				"namespace": {
					Type:        workflow.TypeString,
					Partial:     false,
					Required:    false,
					Description: "Namespace to use",
					Short:       "s",
					Default:     "default",
				},
				"release_name": {
					Type:        workflow.TypeString,
					Partial:     false,
					Required:    true,
					Description: "release name to use",
					Short:       "n",
					Default:     "",
				},
			},
		},
		"k8_helm_add_repo": {
			Action:         Action_AddHelmRepo,
			Short:          "Add a helm repo",
			Long:           "Add a helm repo",
			ProcessResults: false,
			Target:         "action_k8.k8",
			ConfigSchema: map[string]*workflow.Schema{
				"target_name": {
					Type:        workflow.TypeString,
					Partial:     false,
					Required:    false,
					Description: "The target name to use if not default target type",
					Short:       "t",
					Default:     "",
				},
				"namespace": {
					Type:        workflow.TypeString,
					Partial:     false,
					Required:    false,
					Description: "Namespace to use",
					Short:       "s",
					Default:     "default",
				},
				"name": {
					Type:        workflow.TypeString,
					Partial:     false,
					Required:    false,
					Description: "repo name",
					Short:       "a",
					Default:     "",
				},
				"url": {
					Type:        workflow.TypeString,
					Partial:     false,
					Required:    false,
					Description: "Url of repo",
					Short:       "",
					Default:     "h",
				},
				"username": {
					Type:        workflow.TypeString,
					Partial:     false,
					Required:    false,
					Description: "Username",
					Short:       "u",
					Default:     "",
				},
				"password": {
					Type:        workflow.TypeString,
					Partial:     false,
					Required:    false,
					Description: "Password",
					Short:       "p",
					Default:     "",
				},
				"use_config": {
					Type:        workflow.TypeBool,
					Partial:     false,
					Required:    false,
					Description: "Use the target  config",
					Short:       "c",
					Default:     false,
				},
			},
		},
		"k8_delete_service": {
			Action:         Action_K8DeleteService,
			Short:          "Delete Service",
			Long:           "Delete Service",
			ProcessResults: false,
			Target:         "action_k8.k8",
			ConfigSchema: map[string]*workflow.Schema{
				"target_name": {
					Type:        workflow.TypeString,
					Partial:     false,
					Required:    false,
					Description: "The target name to use if not default target type",
					Short:       "t",
					Default:     "",
				},
				"namespace": {
					Type:        workflow.TypeString,
					Partial:     false,
					Required:    false,
					Description: "Namespace to use",
					Short:       "s",
					Default:     "default",
				},
				"name": {
					Type:        workflow.TypeString,
					Partial:     false,
					Required:    true,
					Description: "Name of the service",
					Short:       "n",
					Default:     "",
				},
			},
		},
		"k8_delete_deployment": {
			Action:         Action_K8DeleteDeployment,
			Short:          "Delete Deployment",
			Long:           "Delete Deployment",
			ProcessResults: false,
			Target:         "action_k8.k8",
			ConfigSchema: map[string]*workflow.Schema{
				"target_name": {
					Type:        workflow.TypeString,
					Partial:     false,
					Required:    false,
					Description: "The target name to use if not default target type",
					Short:       "t",
					Default:     "",
				},
				"namespace": {
					Type:        workflow.TypeString,
					Partial:     false,
					Required:    false,
					Description: "Namespace to use",
					Short:       "s",
					Default:     "default",
				},
				"name": {
					Type:        workflow.TypeString,
					Partial:     false,
					Required:    true,
					Description: "Name of the Deployment",
					Short:       "n",
					Default:     "",
				},
			},
		},
		"k8_delete_pod": {
			Action:         Action_K8DeletePod,
			Short:          "Delete Pod",
			Long:           "Delete Pod",
			ProcessResults: false,
			Target:         "action_k8.k8",
			ConfigSchema: map[string]*workflow.Schema{
				"target_name": {
					Type:        workflow.TypeString,
					Partial:     false,
					Required:    false,
					Description: "The target name to use if not default target type",
					Short:       "t",
					Default:     "",
				},
				"namespace": {
					Type:        workflow.TypeString,
					Partial:     false,
					Required:    false,
					Description: "Namespace to use",
					Short:       "s",
					Default:     "default",
				},
				"name": {
					Type:        workflow.TypeString,
					Partial:     false,
					Required:    true,
					Description: "Name of the Pod",
					Short:       "n",
					Default:     "",
				},
			},
		},
		"k8_delete_secret": {
			Action:         Action_K8DeleteSecrets,
			Short:          "Delete a secret",
			Long:           "Delete a secret",
			ProcessResults: false,
			Target:         "action_k8.k8",
			ConfigSchema: map[string]*workflow.Schema{
				"target_name": {
					Type:        workflow.TypeString,
					Partial:     false,
					Required:    false,
					Description: "The target name to use if not default target type",
					Short:       "t",
					Default:     "",
				},
				"namespace": {
					Type:        workflow.TypeString,
					Partial:     false,
					Required:    false,
					Description: "Namespace to use",
					Short:       "s",
					Default:     "default",
				},
				"name": {
					Type:        workflow.TypeString,
					Partial:     false,
					Required:    true,
					Description: "Name of the secret",
					Short:       "n",
					Default:     "",
				},
			},
		},
		"k8_delete_demon_set": {
			Action:         Action_K8DeleteDemonSet,
			Short:          "Delete a demon set",
			Long:           "Delete a demon set",
			ProcessResults: false,
			Target:         "action_k8.k8",
			ConfigSchema: map[string]*workflow.Schema{
				"target_name": {
					Type:        workflow.TypeString,
					Partial:     false,
					Required:    false,
					Description: "The target name to use if not default target type",
					Short:       "t",
					Default:     "",
				},
				"namespace": {
					Type:        workflow.TypeString,
					Partial:     false,
					Required:    false,
					Description: "Namespace to use",
					Short:       "s",
					Default:     "default",
				},
				"name": {
					Type:        workflow.TypeString,
					Partial:     false,
					Required:    true,
					Description: "Name of the demon set",
					Short:       "n",
					Default:     "",
				},
			},
		},
		"k8_delete_stateful_set": {
			Action:         Action_K8DeleteStatefulSet,
			Short:          "Delete a stateful set",
			Long:           "Delete a stateful set",
			ProcessResults: false,
			Target:         "action_k8.k8",
			ConfigSchema: map[string]*workflow.Schema{
				"target_name": {
					Type:        workflow.TypeString,
					Partial:     false,
					Required:    false,
					Description: "The target name to use if not default target type",
					Short:       "t",
					Default:     "",
				},
				"namespace": {
					Type:        workflow.TypeString,
					Partial:     false,
					Required:    false,
					Description: "Namespace to use",
					Short:       "s",
					Default:     "default",
				},
				"name": {
					Type:        workflow.TypeString,
					Partial:     false,
					Required:    true,
					Description: "Name of the stateful set",
					Short:       "n",
					Default:     "",
				},
			},
		},
		"k8_delete_pvc": {
			Action:         Action_K8DeletePVC,
			Short:          "Delete a PVC",
			Long:           "Delete a PVC",
			ProcessResults: false,
			Target:         "action_k8.k8",
			ConfigSchema: map[string]*workflow.Schema{
				"target_name": {
					Type:        workflow.TypeString,
					Partial:     false,
					Required:    false,
					Description: "The target name to use if not default target type",
					Short:       "t",
					Default:     "",
				},
				"namespace": {
					Type:        workflow.TypeString,
					Partial:     false,
					Required:    false,
					Description: "Namespace to use",
					Short:       "s",
					Default:     "default",
				},
				"name": {
					Type:        workflow.TypeString,
					Partial:     false,
					Required:    true,
					Description: "Name of the PVC",
					Short:       "n",
					Default:     "",
				},
			},
		},
		"k8_delete_pv": {
			Action:         Action_K8DeletePV,
			Short:          "Delete a PV",
			Long:           "Delete a PV",
			ProcessResults: false,
			Target:         "action_k8.k8",
			ConfigSchema: map[string]*workflow.Schema{
				"target_name": {
					Type:        workflow.TypeString,
					Partial:     false,
					Required:    false,
					Description: "The target name to use if not default target type",
					Short:       "t",
					Default:     "",
				},
				"namespace": {
					Type:        workflow.TypeString,
					Partial:     false,
					Required:    false,
					Description: "Namespace to use",
					Short:       "s",
					Default:     "default",
				},
				"name": {
					Type:        workflow.TypeString,
					Partial:     false,
					Required:    true,
					Description: "Name of the PV",
					Short:       "n",
					Default:     "",
				},
			},
		},
	}
}

// GetFunctionMap returns the function map for this schema
func (s K8Schema) GetFunctionMap() map[string]workflow.FunctionSchema {
	return map[string]workflow.FunctionSchema{}
}

// GetAction returns the action for the target
func GetSchema() workflow.SchemaEndpoint {
	return K8Schema{}
}

// Action_K8ApplyDeleteYaml - Apply a yaml file or delete k8 manifest file
// - w - the workflow
// - m - the template data
// - returns - error
func Action_K8ApplyDeleteYaml(w *workflow.Workflow, m *workflow.TemplateData) error {

	//********************
	//Create the k8 client
	//********************
	k8, err := go_k8_helm.CreateK8()
	if err != nil {
		return err
	}
	git_obj, err := w.MapTargetConfigValue(m, k8)
	if err != nil {
		return err
	}
	k8_client := git_obj.(*go_k8_helm.K8) //cast it as a git type
	k8_client.CreateConfigAndContext()

	file_data := ""

	//See if to process
	process, err := w.GetConfigTokenBool("process_tokens", m, false)
	if err != nil {
		return err
	}

	namespace, err := w.GetConfigTokenString("namespace", m, false)
	if err != nil {
		return err
	}

	delete, err := w.GetConfigTokenBool("delete", m, false)
	if err != nil {
		return err
	}

	if w.LogLevel == workflow.LOG_VERBOSE {
		log.LogVerbose(fmt.Sprintf("namespace: %s\n", namespace))
		log.LogVerbose(fmt.Sprintf("process_tokens: %v\n", process))
		log.LogVerbose(fmt.Sprintf("delete: %v\n", delete))
	}

	//************************
	//Get the param for the k8
	//************************
	manifest, err := w.GetConfigTokenInterface("manifest", m, true)
	if err != nil {
		return err
	}

	if convert.IsMapInterface(manifest) {
		//check if the manifest is a map
		//and process the tokens
		//if so then convert to yaml
		//manifest, err = w.GetConfigTokenMap("manifest", m, true)
		//if err != nil {
		//	return err
		//}

		//convert to yaml
		file, err := yaml.Marshal(manifest)
		if err != nil {
			return err
		}
		file_data = string(file)
	} else if convert.IsString(manifest) {
		if w.LogLevel == workflow.LOG_VERBOSE {
			log.LogVerbose(fmt.Sprintf("manifest: %s\n", manifest))
		}

		//Check if the manifest is a string
		//read the file and process the tokens if needed
		manifest_path := w.BuildPath(manifest.(string))
		b, err := ioutil.ReadFile(manifest_path) // just pass the file name
		if err != nil {
			return err
		}

		//*********************
		//Now process the token
		//*********************
		file_data = string(b)
	}

	//*********************
	//Now process the token
	//*********************
	if process {
		file_data, err = w.ParseToken(m, file_data)
		if err != nil {
			return err
		}
	}

	//*************
	//Get the items
	//*************
	err = k8_client.ProcessK8File([]byte(file_data), namespace, !delete)
	if err != nil {
		return err
	}

	//**************
	//Return the err
	//**************
	return err
}

// Action_K8DeleteNS - Delete a namespace
// - w - the workflow
// - m - the template data
// - returns - error
func Action_K8DeleteNS(w *workflow.Workflow, m *workflow.TemplateData) error {

	//********************
	//Create the k8 client
	//********************
	k8, err := go_k8_helm.CreateK8()
	if err != nil {
		return err
	}
	git_obj, err := w.MapTargetConfigValue(m, k8)
	if err != nil {
		return err
	}
	k8_client := git_obj.(*go_k8_helm.K8) //cast it as a git type
	k8_client.CreateConfigAndContext()

	namespace, _ := w.GetConfigTokenString("namespace", m, true)

	err = str.CheckNotBlank(namespace)
	if err != nil {
		return err
	}

	if w.LogLevel == workflow.LOG_VERBOSE {
		log.LogVerbose(fmt.Sprintf("namespace: %s\n", namespace))
	}
	//*************
	//Get the items
	//*************
	err = k8_client.DeleteNS(namespace)

	//**************
	//Return the err
	//**************
	return err
}

// Action_K8CreateNS - Create a namespace
// - w - the workflow
// - m - the template data
// - returns - error
func Action_K8CreateNS(w *workflow.Workflow, m *workflow.TemplateData) error {

	//********************
	//Create the k8 client
	//********************
	k8, err := go_k8_helm.CreateK8()
	if err != nil {
		return err
	}
	git_obj, err := w.MapTargetConfigValue(m, k8)
	if err != nil {
		return err
	}
	k8_client := git_obj.(*go_k8_helm.K8) //cast it as a git type
	k8_client.CreateConfigAndContext()

	namespace, _ := w.GetConfigTokenString("namespace", m, true)

	err = str.CheckNotBlank(namespace)
	if err != nil {
		return err
	}

	if w.LogLevel == workflow.LOG_VERBOSE {
		log.LogVerbose(fmt.Sprintf("namespace: %s\n", namespace))
	}
	//*************
	//Get the items
	//*************
	err = k8_client.CreateNS(namespace)

	//**************
	//Return the err
	//**************
	return err
}

// Action_K8Copy - Copy a file to and from a k8 pod
// - w - the workflow
// - m - the template data
// - returns - error
func Action_K8Copy(w *workflow.Workflow, m *workflow.TemplateData) error {

	//********************
	//Create the k8 client
	//********************
	k8, err := go_k8_helm.CreateK8()
	if err != nil {
		return err
	}
	git_obj, err := w.MapTargetConfigValue(m, k8)
	if err != nil {
		return err
	}
	k8_client := git_obj.(*go_k8_helm.K8) //cast it as a git type
	k8_client.CreateConfigAndContext()

	namespace, _ := w.GetConfigTokenString("namespace", m, false)

	src, err := w.GetConfigTokenString("src", m, true)
	if err != nil {
		return err
	}
	dest, err := w.GetConfigTokenString("dest", m, true)
	if err != nil {
		return err
	}
	container_name, err := w.GetConfigTokenString("container_name", m, true)
	if err != nil {
		return err
	}
	err = str.CheckNotBlank(src, dest, container_name)
	if err != nil {
		return err
	}

	if w.LogLevel == workflow.LOG_VERBOSE {
		log.LogVerbose(fmt.Sprintf("src: %s\n", src))
		log.LogVerbose(fmt.Sprintf("dest: %s\n", dest))
		log.LogVerbose(fmt.Sprintf("container_name: %s\n", container_name))
	}
	//*************
	//Get the items
	//*************
	result, err := k8_client.PodCopy(namespace, src, dest, container_name)
	if err != nil {
		return err
	}

	//****************
	//Process the pods
	//****************
	err = w.ActionProcessResults(m, result)

	//**************
	//Return the err
	//**************
	return err
}

// Action_K8Exec - Execute a command on a k8 pod
// - w - the workflow
// - m - the template data
// - returns - error
func Action_K8PodExec(w *workflow.Workflow, m *workflow.TemplateData) error {

	//********************
	//Create the k8 client
	//********************
	k8, err := go_k8_helm.CreateK8()
	if err != nil {
		return err
	}
	git_obj, err := w.MapTargetConfigValue(m, k8)
	if err != nil {
		return err
	}
	k8_client := git_obj.(*go_k8_helm.K8) //cast it as a git type
	k8_client.CreateConfigAndContext()

	namespace, _ := w.GetConfigTokenString("namespace", m, false)
	pod_name, err := w.GetConfigTokenString("pod_name", m, true)
	if err != nil {
		return err
	}
	command, err := w.GetConfigTokenString("command", m, true)
	if err != nil {
		return err
	}
	err = str.CheckNotBlank(pod_name, command)
	if err != nil {
		return err
	}

	if w.LogLevel == workflow.LOG_VERBOSE {
		log.LogVerbose(fmt.Sprintf("namespace: %s\n", namespace))
		log.LogVerbose(fmt.Sprintf("pod_name: %s\n", pod_name))
		log.LogVerbose(fmt.Sprintf("command: %s\n", command))
	}
	//*************
	//Get the items
	//*************
	result, err := k8_client.PodExec(namespace, pod_name, command)
	if err != nil {
		return err
	}

	//****************
	//Process the pods
	//****************
	err = w.ActionProcessResults(m, result)

	//**************
	//Return the err
	//**************
	return err
}

// Action_K8GetServiceIP - Get the IP address of a service
// - w - the workflow
// - m - the template data
// - returns - error
func Action_K8GetServiceIP(w *workflow.Workflow, m *workflow.TemplateData) error {

	//********************
	//Create the k8 client
	//********************
	k8, err := go_k8_helm.CreateK8()
	if err != nil {
		return err
	}
	git_obj, err := w.MapTargetConfigValue(m, k8)
	if err != nil {
		return err
	}
	k8_client := git_obj.(*go_k8_helm.K8) //cast it as a git type
	k8_client.CreateConfigAndContext()

	namespace, _ := w.GetConfigTokenString("namespace", m, false)
	service_name, err := w.GetConfigTokenString("name", m, true)
	if err != nil {
		return err
	}
	if w.LogLevel == workflow.LOG_VERBOSE {
		log.LogVerbose(fmt.Sprintf("namespace: %s\n", namespace))
		log.LogVerbose(fmt.Sprintf("service_name: %s\n", service_name))
	}
	//*************
	//Get the items
	//*************

	items, err := k8_client.GetServiceIP(namespace, service_name)
	if err != nil {
		return err
	}

	//****************
	//Process the pods
	//****************
	err = w.ActionProcessResults(m, items)

	//**************
	//Return the err
	//**************
	return err
}

// Action_K8GetWorkspace - Get data from a k8 workspace e.g pods, deployments, services
// - w - the workflow
// - m - the template data
// - returns - error
func Action_K8GetWorkspace(w *workflow.Workflow, m *workflow.TemplateData) error {

	//********************
	//Create the k8 client
	//********************
	k8, err := go_k8_helm.CreateK8()
	if err != nil {
		return err
	}
	git_obj, err := w.MapTargetConfigValue(m, k8)
	if err != nil {
		return err
	}
	k8_client := git_obj.(*go_k8_helm.K8) //cast it as a git type
	k8_client.CreateConfigAndContext()

	namespace, _ := w.GetConfigTokenString("namespace", m, false)
	workspace, err := w.GetConfigTokenString("workspace", m, true)
	if err != nil {
		return err
	}
	if w.LogLevel == workflow.LOG_VERBOSE {
		log.LogVerbose(fmt.Sprintf("namespace: %s\n", namespace))
		log.LogVerbose(fmt.Sprintf("workspace: %s\n", workspace))
	}
	//*************
	//Get the items
	//*************
	var items interface{}
	switch strings.ToLower(workspace) {
	case "pods":
		items, err = k8_client.GetPods(namespace)
		if err != nil {
			return err
		}
	case "deployments", "dep":
		items, err = k8_client.GetDeployments(namespace)
		if err != nil {
			return err
		}
	case "daemonsets", "ds":
		items, err = k8_client.GetDemonSet(namespace)
		if err != nil {
			return err
		}
	case "statefulsets", "sts", "ss":
		items, err = k8_client.GetStatefulSets(namespace)
		if err != nil {
			return err
		}
	case "secret", "sec":
		items, err = k8_client.GetSecrets(namespace)
		if err != nil {
			return err
		}
	case "services", "service", "svc":
		items, err = k8_client.GetServices(namespace)
		if err != nil {
			return err
		}
	}

	//****************
	//Process the pods
	//****************
	err = w.ActionProcessResults(m, items)

	//**************
	//Return the err
	//**************
	return err
}

// Action_K8WaitCompleteStatus - Wait for a k8 deployment to complete
// - w - the workflow
// - m - the template data
// - returns - error
func Action_K8WaitCompleteStatus(w *workflow.Workflow, m *workflow.TemplateData) error {

	/*
		Example:
			checks := []interface{}{
				"deployment:nginx(.*)",
				"replica:nginx2(.*)",
				"stateful:nginx3(.*)",
				"demon:nginx4(.*)",
				"service:nginx(.*)",
			}
	*/

	//********************
	//Create the k8 client
	//********************
	k8, err := go_k8_helm.CreateK8()
	if err != nil {
		return err
	}
	git_obj, err := w.MapTargetConfigValue(m, k8)
	if err != nil {
		return err
	}
	k8_client := git_obj.(*go_k8_helm.K8) //cast it as a git type
	k8_client.CreateConfigAndContext()

	namespace, _ := w.GetConfigTokenString("namespace", m, false)
	checks, err := w.GetConfigTokenInterface("checks", m, true)
	if err != nil {
		return err
	}

	not_running, err := w.GetConfigTokenBool("not_running", m, false)
	if err != nil {
		return err
	}

	retry, err := w.GetConfigTokenInt("retry", m, false)
	if err != nil {
		return err
	}

	if w.LogLevel == workflow.LOG_VERBOSE {
		log.LogVerbose(fmt.Sprintf("namespace: %s\n", namespace))
		log.LogVerbose(fmt.Sprintf("retry: %v\n", retry))
	}
	//*************************
	//Check the status of items
	//*************************
	array_checks := convert.ToArrayInterface(checks)

	if retry == 0 {
		retry = 1
	}
	if retry > 1000 {
		retry = 1000
	}

	//***************
	//Run the command
	//***************
	sum := 0
	for {
		sum++ // repeated forever

		//****************
		//Check the status
		//****************
		completed, results, err := k8_client.CheckStatusOf(namespace, array_checks, not_running)
		if err != nil {
			return err
		}
		//*****************
		//Print the results
		//*****************
		fmt.Println(results)

		//*******************************
		//See if the deployment completed
		//*******************************

		if completed {
			return nil
		}
		if sum >= retry {
			return fmt.Errorf("retry limit reached")
		}
		time.Sleep(10 * time.Second)
	}
}

// Action_DeployUpgradeHelmChart - Deploy or upgrade a helm chart
// - w - the workflow
// - m - the template data
// - returns - error
func Action_DeployUpgradeHelmChart(w *workflow.Workflow, m *workflow.TemplateData) error {
	//********************
	//Create the k8 client
	//********************
	k8, err := go_k8_helm.CreateK8()
	if err != nil {
		return err
	}
	git_obj, err := w.MapTargetConfigValue(m, k8)
	if err != nil {
		return err
	}
	k8_client := git_obj.(*go_k8_helm.K8) //cast it as a git type
	k8_client.CreateConfigAndContext()

	namespace, _ := w.GetConfigTokenString("namespace", m, false)

	// Getting the deployment name from the config file.
	chart_Path, err := w.GetConfigTokenString("chart_Path", m, true)
	chart_Path = w.BuildPath(chart_Path)
	if err != nil {
		return err
	}
	release_name, err := w.GetConfigTokenString("release_name", m, true)
	if err != nil {
		return err
	}
	helm_config, err := w.GetConfigTokenMap("helm_config", m, true)
	if err != nil {
		return err
	}

	upgrade, err := w.GetConfigTokenBool("upgrade", m, true)
	if err != nil {
		return err
	}

	if w.LogLevel == workflow.LOG_VERBOSE {
		log.LogVerbose(fmt.Sprintf("chart_Path: %s\n", chart_Path))
		log.LogVerbose(fmt.Sprintf("namespace: %s\n", namespace))
		log.LogVerbose(fmt.Sprintf("release_name: %s\n", release_name))
		log.LogVerbose(fmt.Sprintf("upgrade: %v\n", upgrade))
	}

	if upgrade {
		//**********************
		//Upgrade the helm chart
		//**********************
		return k8_client.UpgradeHelmChart(chart_Path, release_name, namespace, helm_config)
	}
	//************************
	//Deploying the helm chart
	//************************
	return k8_client.DeployHelmChart(chart_Path, release_name, namespace, helm_config)
}

// Action_DeleteHelmChart - Delete a helm chart
// - w - the workflow
// - m - the template data
// - returns - error
func Action_DeleteHelmChart(w *workflow.Workflow, m *workflow.TemplateData) error {
	//********************
	//Create the k8 client
	//********************
	k8, err := go_k8_helm.CreateK8()
	if err != nil {
		return err
	}
	git_obj, err := w.MapTargetConfigValue(m, k8)
	if err != nil {
		return err
	}
	k8_client := git_obj.(*go_k8_helm.K8) //cast it as a git type
	k8_client.CreateConfigAndContext()

	namespace, _ := w.GetConfigTokenString("namespace", m, false)
	release_name, err := w.GetConfigTokenString("release_name", m, true)
	if err != nil {
		return err
	}

	if w.LogLevel == workflow.LOG_VERBOSE {
		log.LogVerbose(fmt.Sprintf("namespace: %s\n", namespace))
		log.LogVerbose(fmt.Sprintf("release_name: %s\n", release_name))
	}

	//************************
	//Deploying the helm chart
	//************************
	return k8_client.UninstallHelmChart(release_name, namespace)
}

// Action_AddHelmRepo - Add a helm repo
// - w - the workflow
// - m - the template data
// - returns - error
func Action_AddHelmRepo(w *workflow.Workflow, m *workflow.TemplateData) error {
	//********************
	//Create the k8 client
	//********************
	k8, err := go_k8_helm.CreateK8()
	if err != nil {
		return err
	}
	git_obj, err := w.MapTargetConfigValue(m, k8)
	if err != nil {
		return err
	}
	k8_client := git_obj.(*go_k8_helm.K8) //cast it as a git type
	k8_client.CreateConfigAndContext()

	name, err := w.GetConfigTokenString("name", m, true)
	if err != nil {
		return err
	}

	//See if to use target credentials
	use_config, err := w.GetConfigTokenBool("use_config", m, false)
	if err != nil {
		return err
	}
	url := ""
	username := ""
	password := ""
	if use_config {
		//*****************************
		//Get the env from runtime_vars
		//*****************************
		env, err := w.GetRuntimeVar("env")
		if err != nil {
			return err
		}
		env = env.(string)

		//*******************
		//Get the target_name
		//*******************
		target_name, err := w.GetConfigTokenString("target_name", m, false)
		if err != nil {
			return err
		}
		if target_name == "" {
			target_name = "helm-repo"
		} else {
			target_name = "helm-repo" + target_name
		}

		//***********************
		//Get the url from config
		//***********************
		config_url, err := w.GetConfigValue("", target_name+".url", "string", env.(string))
		if err != nil {
			return err
		}
		url = config_url.(string)
		//*****************************
		//Get the username from config
		//*****************************
		config_username, err := w.GetConfigValue("", target_name+".username", "string", env.(string))
		if err != nil {
			return err
		}
		username = config_username.(string)
		//*****************************
		//Get the password from config
		//*****************************
		config_password, err := w.GetConfigValue("", target_name+".password", "string", env.(string))
		if err != nil {
			return err
		}
		password = config_password.(string)

	} else {
		url, err = w.GetConfigTokenString("url", m, true)
		if err != nil {
			return err
		}
		username, err = w.GetConfigTokenString("username", m, false)
		if err != nil {
			return err
		}
		password, err = w.GetConfigTokenString("password", m, false)
		if err != nil {
			return err
		}
	}
	if w.LogLevel == workflow.LOG_VERBOSE {
		log.LogVerbose(fmt.Sprintf("name: %s\n", name))
		log.LogVerbose(fmt.Sprintf("url: %s\n", url))
	}

	//************************
	//Deploying the helm chart
	//************************
	err = k8_client.RepoAdd(name, url, username, password)
	if err != nil {
		return err
	}
	k8_client.RepoUpdate()
	return nil
}

// Action_K8DeleteService - Delete a k8 service
// - w - the workflow
// - m - the template data
// - returns - error
func Action_K8DeleteService(w *workflow.Workflow, m *workflow.TemplateData) error {

	//********************
	//Create the k8 client
	//********************
	k8, err := go_k8_helm.CreateK8()
	if err != nil {
		return err
	}
	git_obj, err := w.MapTargetConfigValue(m, k8)
	if err != nil {
		return err
	}
	k8_client := git_obj.(*go_k8_helm.K8) //cast it as a git type
	k8_client.CreateConfigAndContext()

	namespace, _ := w.GetConfigTokenString("namespace", m, false)
	service_name, err := w.GetConfigTokenString("name", m, true)
	if err != nil {
		return err
	}
	if w.LogLevel == workflow.LOG_VERBOSE {
		log.LogVerbose(fmt.Sprintf("namespace: %s\n", namespace))
		log.LogVerbose(fmt.Sprintf("service_name: %s\n", service_name))
	}
	//*************
	//Get the items
	//*************

	err = k8_client.DeleteService(namespace, service_name)
	if err != nil {
		return err
	}

	//**************
	//Return the err
	//**************
	return err
}

// Action_K8DeleteDeployment - Delete a k8 deployment
// - w - the workflow
// - m - the template data
// - returns - error
func Action_K8DeleteDeployment(w *workflow.Workflow, m *workflow.TemplateData) error {

	//********************
	//Create the k8 client
	//********************
	k8, err := go_k8_helm.CreateK8()
	if err != nil {
		return err
	}
	git_obj, err := w.MapTargetConfigValue(m, k8)
	if err != nil {
		return err
	}
	k8_client := git_obj.(*go_k8_helm.K8) //cast it as a git type
	k8_client.CreateConfigAndContext()

	namespace, _ := w.GetConfigTokenString("namespace", m, false)
	service_name, err := w.GetConfigTokenString("name", m, true)
	if err != nil {
		return err
	}
	if w.LogLevel == workflow.LOG_VERBOSE {
		log.LogVerbose(fmt.Sprintf("namespace: %s\n", namespace))
		log.LogVerbose(fmt.Sprintf("service_name: %s\n", service_name))
	}
	//*************
	//Get the items
	//*************

	err = k8_client.DeleteDeployment(namespace, service_name)
	if err != nil {
		return err
	}

	//**************
	//Return the err
	//**************
	return err
}

// Action_K8DeletePod - Delete a k8 pod
// - w - the workflow
// - m - the template data
// - returns - error
func Action_K8DeletePod(w *workflow.Workflow, m *workflow.TemplateData) error {

	//********************
	//Create the k8 client
	//********************
	k8, err := go_k8_helm.CreateK8()
	if err != nil {
		return err
	}
	git_obj, err := w.MapTargetConfigValue(m, k8)
	if err != nil {
		return err
	}
	k8_client := git_obj.(*go_k8_helm.K8) //cast it as a git type
	k8_client.CreateConfigAndContext()

	namespace, _ := w.GetConfigTokenString("namespace", m, false)
	service_name, err := w.GetConfigTokenString("name", m, true)
	if err != nil {
		return err
	}
	if w.LogLevel == workflow.LOG_VERBOSE {
		log.LogVerbose(fmt.Sprintf("namespace: %s\n", namespace))
		log.LogVerbose(fmt.Sprintf("service_name: %s\n", service_name))
	}
	//*************
	//Get the items
	//*************

	err = k8_client.DeletePod(namespace, service_name)
	if err != nil {
		return err
	}

	//**************
	//Return the err
	//**************
	return err
}

// Action_K8DeleteSecrets - Delete a k8 secret
// - w - the workflow
// - m - the template data
// - returns - error
func Action_K8DeleteSecrets(w *workflow.Workflow, m *workflow.TemplateData) error {

	//********************
	//Create the k8 client
	//********************
	k8, err := go_k8_helm.CreateK8()
	if err != nil {
		return err
	}
	git_obj, err := w.MapTargetConfigValue(m, k8)
	if err != nil {
		return err
	}
	k8_client := git_obj.(*go_k8_helm.K8) //cast it as a git type
	k8_client.CreateConfigAndContext()

	namespace, _ := w.GetConfigTokenString("namespace", m, false)
	secret_name, err := w.GetConfigTokenString("name", m, true)
	if err != nil {
		return err
	}
	if w.LogLevel == workflow.LOG_VERBOSE {
		log.LogVerbose(fmt.Sprintf("namespace: %s\n", namespace))
		log.LogVerbose(fmt.Sprintf("secret_name: %s\n", secret_name))
	}
	//*************
	//Get the items
	//*************

	err = k8_client.DeleteSecrets(namespace, secret_name)
	if err != nil {
		return err
	}

	//**************
	//Return the err
	//**************
	return err
}

// Action_K8DeleteConfigMap - Delete a k8 config map
// - w - the workflow
// - m - the template data
// - returns - error
func Action_K8DeleteStatefulSet(w *workflow.Workflow, m *workflow.TemplateData) error {

	//********************
	//Create the k8 client
	//********************
	k8, err := go_k8_helm.CreateK8()
	if err != nil {
		return err
	}
	git_obj, err := w.MapTargetConfigValue(m, k8)
	if err != nil {
		return err
	}
	k8_client := git_obj.(*go_k8_helm.K8) //cast it as a git type
	k8_client.CreateConfigAndContext()

	namespace, _ := w.GetConfigTokenString("namespace", m, false)
	stateful_set_name, err := w.GetConfigTokenString("name", m, true)
	if err != nil {
		return err
	}
	if w.LogLevel == workflow.LOG_VERBOSE {
		log.LogVerbose(fmt.Sprintf("namespace: %s\n", namespace))
		log.LogVerbose(fmt.Sprintf("stateful_set_name: %s\n", stateful_set_name))
	}
	//*************
	//Get the items
	//*************

	err = k8_client.DeleteStatefulSets(namespace, stateful_set_name)
	if err != nil {
		return err
	}

	//**************
	//Return the err
	//**************
	return err
}

// Action_K8DeleteConfigMap - Delete a k8 config map
// - w - the workflow
// - m - the template data
// - returns - error
func Action_K8DeleteDemonSet(w *workflow.Workflow, m *workflow.TemplateData) error {

	//********************
	//Create the k8 client
	//********************
	k8, err := go_k8_helm.CreateK8()
	if err != nil {
		return err
	}
	git_obj, err := w.MapTargetConfigValue(m, k8)
	if err != nil {
		return err
	}
	k8_client := git_obj.(*go_k8_helm.K8) //cast it as a git type
	k8_client.CreateConfigAndContext()

	namespace, _ := w.GetConfigTokenString("namespace", m, false)
	demon_set_name, err := w.GetConfigTokenString("name", m, true)
	if err != nil {
		return err
	}
	if w.LogLevel == workflow.LOG_VERBOSE {
		log.LogVerbose(fmt.Sprintf("namespace: %s\n", namespace))
		log.LogVerbose(fmt.Sprintf("demon_set_name: %s\n", demon_set_name))
	}
	//*************
	//Get the items
	//*************

	err = k8_client.DeleteDemonSet(namespace, demon_set_name)
	if err != nil {
		return err
	}

	//**************
	//Return the err
	//**************
	return err
}

// Action_K8DeletePVC - Delete a k8 pvc
// - w - the workflow
// - m - the template data
// - returns - error
func Action_K8DeletePVC(w *workflow.Workflow, m *workflow.TemplateData) error {

	//********************
	//Create the k8 client
	//********************
	k8, err := go_k8_helm.CreateK8()
	if err != nil {
		return err
	}
	git_obj, err := w.MapTargetConfigValue(m, k8)
	if err != nil {
		return err
	}
	k8_client := git_obj.(*go_k8_helm.K8) //cast it as a git type
	k8_client.CreateConfigAndContext()

	namespace, _ := w.GetConfigTokenString("namespace", m, false)
	pvc_name, err := w.GetConfigTokenString("name", m, true)
	if err != nil {
		return err
	}
	if w.LogLevel == workflow.LOG_VERBOSE {
		log.LogVerbose(fmt.Sprintf("namespace: %s\n", namespace))
		log.LogVerbose(fmt.Sprintf("pvc_name: %s\n", pvc_name))
	}
	//*************
	//Get the items
	//*************

	err = k8_client.DeletePVC(namespace, pvc_name)
	if err != nil {
		return err
	}

	//**************
	//Return the err
	//**************
	return err
}

// Action_K8DeletePV - Delete a k8 pv
// - w - the workflow
// - m - the template data
// - returns - error
func Action_K8DeletePV(w *workflow.Workflow, m *workflow.TemplateData) error {

	//********************
	//Create the k8 client
	//********************
	k8, err := go_k8_helm.CreateK8()
	if err != nil {
		return err
	}
	git_obj, err := w.MapTargetConfigValue(m, k8)
	if err != nil {
		return err
	}
	k8_client := git_obj.(*go_k8_helm.K8) //cast it as a git type
	k8_client.CreateConfigAndContext()

	namespace, _ := w.GetConfigTokenString("namespace", m, false)
	pv_name, err := w.GetConfigTokenString("name", m, true)
	if err != nil {
		return err
	}
	if w.LogLevel == workflow.LOG_VERBOSE {
		log.LogVerbose(fmt.Sprintf("namespace: %s\n", namespace))
		log.LogVerbose(fmt.Sprintf("pv_name: %s\n", pv_name))
	}
	//*************
	//Get the items
	//*************

	err = k8_client.DeletePV(namespace, pv_name)
	if err != nil {
		return err
	}

	//**************
	//Return the err
	//**************
	return err
}
