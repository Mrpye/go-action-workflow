# Plugin Library Documentation
Below are the instructions for using the plugins library in your workflow. Below are the setting and parameters for each of the actions, target and function included in the library.

---


# **Plugin Library**:  action_docker

# Actions:


## **Action**: docker_reg_download
Download a docker image from a registry

<details>
<summary>Info</summary>

|Property|Description|
|----|----|
|Action|docker_reg_download|
|Description|Download a docker image from a registry|
|Target|[action_docker.registry](#target-action_dockerregistry)|
|InlineParams|false
|ProcessResults|false|


### Config Parameters

|parameter|Type|Description|Required|Short Flag|Default|
|----|----|----|----|----|----|
|folder|String|Folder to save the images to|true|f|||
|images|List|List of images to download|true|i||[]|
|target_name|String|The target to use|false|t|||


</details>

---

## **Action**: docker_reg_upload
Upload a docker image to a registry

<details>
<summary>Info</summary>

|Property|Description|
|----|----|
|Action|docker_reg_upload|
|Description|Upload a docker image to a registry|
|Target|[action_docker.registry](#target-action_dockerregistry)|
|InlineParams|false
|ProcessResults|false|


### Config Parameters

|parameter|Type|Description|Required|Short Flag|Default|
|----|----|----|----|----|----|
|folder|String|Folder to load the images from|true|f|||
|images|List|List of images to upload|true|i||[]|
|import_all|Bool|Import all images in the folder|false|a||false|
|target_name|String|The target to use|false|t|||


</details>

---


# Functions:

## action_docker Functions:


## **Function**: image_account
gets the account name from the image path

<details>
<summary>Info</summary>

|Property|Description|
|----|----|
|Function|image_account|
|Description|gets the account name from the image path|

### Parameters
|parameter|Type|Description|Required|
|----|----|----|----|
|image|String|docker.io/[circleci]/slim-base:latest|true|


</details>

---

## **Function**: image_name
gets the name of the image from the image path

<details>
<summary>Info</summary>

|Property|Description|
|----|----|
|Function|image_name|
|Description|gets the name of the image from the image path|

### Parameters
|parameter|Type|Description|Required|
|----|----|----|----|
|image|String|docker.io/circleci/[slim-base]:latest|true|


</details>

---

## **Function**: image_name_tag
gets the name and tag from the image path

<details>
<summary>Info</summary>

|Property|Description|
|----|----|
|Function|image_name_tag|
|Description|gets the name and tag from the image path|

### Parameters
|parameter|Type|Description|Required|
|----|----|----|----|
|image|String|docker.io/circleci/[slim-base:latest]|true|


</details>

---

## **Function**: image_registry
gets the registry name from the image path

<details>
<summary>Info</summary>

|Property|Description|
|----|----|
|Function|image_registry|
|Description|gets the registry name from the image path|

### Parameters
|parameter|Type|Description|Required|
|----|----|----|----|
|image|String|docker.io/circleci/slim-base]:latest|true|


</details>

---

## **Function**: image_shortname
gets the account name from the image path

<details>
<summary>Info</summary>

|Property|Description|
|----|----|
|Function|image_shortname|
|Description|gets the account name from the image path|

### Parameters
|parameter|Type|Description|Required|
|----|----|----|----|
|image|String|docker.io/[circleci/slim-base]:latest|true|


</details>

---

## **Function**: image_tag
gets the tag  from the image path

<details>
<summary>Info</summary>

|Property|Description|
|----|----|
|Function|image_tag|
|Description|gets the tag  from the image path|

### Parameters
|parameter|Type|Description|Required|
|----|----|----|----|
|image|String|docker.io/circleci/slim-base:[latest]|true|


</details>

---

## **Function**: remap_image
remaps the docker image

<details>
<summary>Info</summary>

|Property|Description|
|----|----|
|Function|remap_image|
|Description|remaps the docker image|

### Parameters
|parameter|Type|Description|Required|
|----|----|----|----|
|*workflow.Workflow|Workflow|work flow engine use (get_wf) to get the workflow engine|true|
|image|String|name of the image|true|
|no_tag|Bool|true not include the tag in the returned value|true|
|original|String|the original image path to use if use_original=true else build path based on target|false|
|target_name|String|name of the target or just use ``|false|
|use_original|Bool|use the original image path|true|


</details>

---

---

# **Plugin Library**:  action_k8

# Actions:


## **Action**: k8_copy
Copy a file from a pod to the local machine or vice versa

<details>
<summary>Info</summary>

|Property|Description|
|----|----|
|Action|k8_copy|
|Description|Copy a file from a pod to the local machine or vice versa|
|Target|[action_k8.k8](#target-action_k8k8)|
|InlineParams|false
|ProcessResults|true|


### Config Parameters

|parameter|Type|Description|Required|Short Flag|Default|
|----|----|----|----|----|----|
|container_name|String|container name|true|c|||
|dest|String|destination file path|true|d|||
|namespace|String|Namespace to use|false|s||default|
|src|String|source file path|true|f|||
|target_name|String|The target name to use if not default target type|false|t|||


</details>

---

## **Action**: k8_create_ns
Create a namespace

<details>
<summary>Info</summary>

|Property|Description|
|----|----|
|Action|k8_create_ns|
|Description|Create a namespace|
|Target|[action_k8.k8](#target-action_k8k8)|
|InlineParams|false
|ProcessResults|false|


### Config Parameters

|parameter|Type|Description|Required|Short Flag|Default|
|----|----|----|----|----|----|
|namespace|String|Namespace to use|false|s||default|
|target_name|String|The target name to use if not default target type|false|t|||


</details>

---

## **Action**: k8_delete_demon_set
Delete a demon set

<details>
<summary>Info</summary>

|Property|Description|
|----|----|
|Action|k8_delete_demon_set|
|Description|Delete a demon set|
|Target|[action_k8.k8](#target-action_k8k8)|
|InlineParams|false
|ProcessResults|false|


### Config Parameters

|parameter|Type|Description|Required|Short Flag|Default|
|----|----|----|----|----|----|
|name|String|Name of the demon set|true|n|||
|namespace|String|Namespace to use|false|s||default|
|target_name|String|The target name to use if not default target type|false|t|||


</details>

---

## **Action**: k8_delete_deployment
Delete Deployment

<details>
<summary>Info</summary>

|Property|Description|
|----|----|
|Action|k8_delete_deployment|
|Description|Delete Deployment|
|Target|[action_k8.k8](#target-action_k8k8)|
|InlineParams|false
|ProcessResults|false|


### Config Parameters

|parameter|Type|Description|Required|Short Flag|Default|
|----|----|----|----|----|----|
|name|String|Name of the Deployment|true|n|||
|namespace|String|Namespace to use|false|s||default|
|target_name|String|The target name to use if not default target type|false|t|||


</details>

---

## **Action**: k8_delete_ns
Delete a namespace

<details>
<summary>Info</summary>

|Property|Description|
|----|----|
|Action|k8_delete_ns|
|Description|Delete a namespace|
|Target|[action_k8.k8](#target-action_k8k8)|
|InlineParams|false
|ProcessResults|false|


### Config Parameters

|parameter|Type|Description|Required|Short Flag|Default|
|----|----|----|----|----|----|
|namespace|String|Namespace to use|false|s||default|
|target_name|String|The target name to use if not default target type|false|t|||


</details>

---

## **Action**: k8_delete_pod
Delete Pod

<details>
<summary>Info</summary>

|Property|Description|
|----|----|
|Action|k8_delete_pod|
|Description|Delete Pod|
|Target|[action_k8.k8](#target-action_k8k8)|
|InlineParams|false
|ProcessResults|false|


### Config Parameters

|parameter|Type|Description|Required|Short Flag|Default|
|----|----|----|----|----|----|
|name|String|Name of the Pod|true|n|||
|namespace|String|Namespace to use|false|s||default|
|target_name|String|The target name to use if not default target type|false|t|||


</details>

---

## **Action**: k8_delete_pv
Delete a PV

<details>
<summary>Info</summary>

|Property|Description|
|----|----|
|Action|k8_delete_pv|
|Description|Delete a PV|
|Target|[action_k8.k8](#target-action_k8k8)|
|InlineParams|false
|ProcessResults|false|


### Config Parameters

|parameter|Type|Description|Required|Short Flag|Default|
|----|----|----|----|----|----|
|name|String|Name of the PV|true|n|||
|namespace|String|Namespace to use|false|s||default|
|target_name|String|The target name to use if not default target type|false|t|||


</details>

---

## **Action**: k8_delete_pvc
Delete a PVC

<details>
<summary>Info</summary>

|Property|Description|
|----|----|
|Action|k8_delete_pvc|
|Description|Delete a PVC|
|Target|[action_k8.k8](#target-action_k8k8)|
|InlineParams|false
|ProcessResults|false|


### Config Parameters

|parameter|Type|Description|Required|Short Flag|Default|
|----|----|----|----|----|----|
|name|String|Name of the PVC|true|n|||
|namespace|String|Namespace to use|false|s||default|
|target_name|String|The target name to use if not default target type|false|t|||


</details>

---

## **Action**: k8_delete_secret
Delete a secret

<details>
<summary>Info</summary>

|Property|Description|
|----|----|
|Action|k8_delete_secret|
|Description|Delete a secret|
|Target|[action_k8.k8](#target-action_k8k8)|
|InlineParams|false
|ProcessResults|false|


### Config Parameters

|parameter|Type|Description|Required|Short Flag|Default|
|----|----|----|----|----|----|
|name|String|Name of the secret|true|n|||
|namespace|String|Namespace to use|false|s||default|
|target_name|String|The target name to use if not default target type|false|t|||


</details>

---

## **Action**: k8_delete_service
Delete Service

<details>
<summary>Info</summary>

|Property|Description|
|----|----|
|Action|k8_delete_service|
|Description|Delete Service|
|Target|[action_k8.k8](#target-action_k8k8)|
|InlineParams|false
|ProcessResults|false|


### Config Parameters

|parameter|Type|Description|Required|Short Flag|Default|
|----|----|----|----|----|----|
|name|String|Name of the service|true|n|||
|namespace|String|Namespace to use|false|s||default|
|target_name|String|The target name to use if not default target type|false|t|||


</details>

---

## **Action**: k8_delete_stateful_set
Delete a stateful set

<details>
<summary>Info</summary>

|Property|Description|
|----|----|
|Action|k8_delete_stateful_set|
|Description|Delete a stateful set|
|Target|[action_k8.k8](#target-action_k8k8)|
|InlineParams|false
|ProcessResults|false|


### Config Parameters

|parameter|Type|Description|Required|Short Flag|Default|
|----|----|----|----|----|----|
|name|String|Name of the stateful set|true|n|||
|namespace|String|Namespace to use|false|s||default|
|target_name|String|The target name to use if not default target type|false|t|||


</details>

---

## **Action**: k8_get_service_ip
Get the IP of a service

<details>
<summary>Info</summary>

|Property|Description|
|----|----|
|Action|k8_get_service_ip|
|Description|Get the IP of a service|
|Target|[action_k8.k8](#target-action_k8k8)|
|InlineParams|false
|ProcessResults|true|


### Config Parameters

|parameter|Type|Description|Required|Short Flag|Default|
|----|----|----|----|----|----|
|name|String|name of service can use regex|true|n|||
|namespace|String|Namespace to use|false|s||default|
|target_name|String|The target name to use if not default target type|false|t|||


</details>

---

## **Action**: k8_get_ws_items
Get items in a workspace

<details>
<summary>Info</summary>

|Property|Description|
|----|----|
|Action|k8_get_ws_items|
|Description|Get items in a workspace|
|Target|[action_k8.k8](#target-action_k8k8)|
|InlineParams|false
|ProcessResults|true|


### Config Parameters

|parameter|Type|Description|Required|Short Flag|Default|
|----|----|----|----|----|----|
|namespace|String|Namespace to use|false|S||default|
|target_name|String|The target name to use if not default target type|false|t|||
|workspace|String|workspace to use|true|w|||


</details>

---

## **Action**: k8_helm_add_repo
Add a helm repo

<details>
<summary>Info</summary>

|Property|Description|
|----|----|
|Action|k8_helm_add_repo|
|Description|Add a helm repo|
|Target|[action_k8.k8](#target-action_k8k8)|
|InlineParams|false
|ProcessResults|false|


### Config Parameters

|parameter|Type|Description|Required|Short Flag|Default|
|----|----|----|----|----|----|
|name|String|repo name|false|a|||
|namespace|String|Namespace to use|false|s||default|
|password|String|Password|false|p|||
|target_name|String|The target name to use if not default target type|false|t|||
|url|String|Url of repo|false|||h|
|use_config|Bool|Use the target  config|false|c||false|
|username|String|Username|false|u|||


</details>

---

## **Action**: k8_helm_delete
Delete a helm chart

<details>
<summary>Info</summary>

|Property|Description|
|----|----|
|Action|k8_helm_delete|
|Description|Delete a helm chart|
|Target|[action_k8.k8](#target-action_k8k8)|
|InlineParams|false
|ProcessResults|false|


### Config Parameters

|parameter|Type|Description|Required|Short Flag|Default|
|----|----|----|----|----|----|
|namespace|String|Namespace to use|false|s||default|
|release_name|String|release name to use|true|n|||
|target_name|String|The target name to use if not default target type|false|t|||


</details>

---

## **Action**: k8_helm_deploy_upgrade
Deploy or upgrade a helm chart

<details>
<summary>Info</summary>

|Property|Description|
|----|----|
|Action|k8_helm_deploy_upgrade|
|Description|Deploy or upgrade a helm chart|
|Target|[action_k8.k8](#target-action_k8k8)|
|InlineParams|false
|ProcessResults|false|


### Config Parameters

|parameter|Type|Description|Required|Short Flag|Default|
|----|----|----|----|----|----|
|chart_Path|String|chart path|false|c|||
|namespace|String|Namespace to use|false|s||default|
|release_name|String|release name|false|n|||
|target_name|String|The target name to use if not default target type|false|t|||
|upgrade|Bool|chart path|false|u||false|


</details>

---

## **Action**: k8_pod_exec
Execute a command in a pod

<details>
<summary>Info</summary>

|Property|Description|
|----|----|
|Action|k8_pod_exec|
|Description|Execute a command in a pod|
|Target|[action_k8.k8](#target-action_k8k8)|
|InlineParams|false
|ProcessResults|true|


### Config Parameters

|parameter|Type|Description|Required|Short Flag|Default|
|----|----|----|----|----|----|
|command|String|command to execute|true|c|||
|namespace|String|Namespace to use|false|s||default|
|pod_name|String|pod name|true|n|||
|target_name|String|The target name to use if not default target type|false|t|||


</details>

---

## **Action**: k8_wait
Wait for a k8 resource to be in a complete state

<details>
<summary>Info</summary>

|Property|Description|
|----|----|
|Action|k8_wait|
|Description|Wait for a k8 resource to be in a complete state|
|Target|[action_k8.k8](#target-action_k8k8)|
|InlineParams|false
|ProcessResults|false|


### Config Parameters

|parameter|Type|Description|Required|Short Flag|Default|
|----|----|----|----|----|----|
|checks|List|Checks to run|false|c||[replica:nginx2(.*) stateful:nginx3(.*) demon:nginx4(.*) service:nginx(.*)]|
|namespace|String|Namespace to use|false|s||default|
|not_running|Bool|All checks not running|false|x||false|
|retry|Int|retry count|false|r||10|
|target_name|String|The target name to use if not default target type|false|t|||


</details>

---

## **Action**: k8_yaml
Apply or delete a yaml file to a k8 cluster

<details>
<summary>Info</summary>

|Property|Description|
|----|----|
|Action|k8_yaml|
|Description|Apply or delete a yaml file to a k8 cluster|
|Target|[action_k8.k8](#target-action_k8k8)|
|InlineParams|false
|ProcessResults|false|


### Config Parameters

|parameter|Type|Description|Required|Short Flag|Default|
|----|----|----|----|----|----|
|delete|Bool|If true, delete the deployment|false|d||false|
|manifest|String|k8 manifest to apply or delete as a file path or object|true|m|||
|namespace|String|Namespace to use|false|s||default|
|process_tokens|Bool|Should tokens be processed in the k8 manifest|false|p||true|
|target_name|String|The target name to use if not default target type|false|t|||


</details>

---


# Functions:

## action_k8 Functions:


---


# Targets:

## **Target** action_docker.registry

<details>
<summary>Info</summary>

|Field|json|cli flag|Desc|
|----|----|----|----|
|Host|host|host h|the host url|
|IgnoreSSL|ignore_ssl|ignore_ssl i|Ignore SSL|
|Library|library|library l|Library to use|
|Password|password|password p|Password for the registry|
|UserName|user|user u|Username for the registry|


</details>

---

## **Target** action_k8.k8

<details>
<summary>Info</summary>

|Field|json|cli flag|Desc|
|----|----|----|----|
|Authorization|authorization|auth a|The authorization token|
|ConfigPath|config_path|config_path p|The path to the kube config file|
|DefaultContext|default_context|context c|The default context to use|
|Host|host|host h|The host to connect to|
|Ignore_ssl|ignore_ssl|ignore_ssl i|If true, ignore the ssl connection|
|UseTokenConnection|use_token_connection|conn-type u|Connection type if true, use the token connection, otherwise use the kube config file|


</details>

---



