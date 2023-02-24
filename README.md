# go-workflow is a basic automation workflow library for GO

## Description
go-workflow allows you to embed a basic workflow engine into your GO projects supporting basic flow operations such as loops, goto, pause and end. you can create custom actions to perform tasks, and using go's template engine you can manipulate values in your workflow manifest. with the ability to run **js** you can manipulate data returned from your actions and store data so that it can be used in later actions. 

---

## When to use go-workflow
- go-workflow can be used in your project when you need automate tasks
- To create automation tooling

---

## Requirements
* go 1.8 [https://go.dev/doc/install](https://go.dev/doc/install) to run and install helm-api

---

## Project folders
Below is a description helm-api project folders and what they contain
|   Folder        | Description  | 
|-----------|---|
| workflow    | go-workflow lib  |
| examples    | example go-workflow use cases  |
| actions    | some prebuilt actions that can be included in your workflow  |

---

## Installation
You can install go-workflow using the following command
<details>
<summary>1. Install go-workflow</summary>

```go
go get github.com/Mrpye/go-workflow
```
</details>

<details>
<summary>2. Add go-workflow to your project</summary>

```go
    import "github.com/Mrpye/go-workflow/workflow"
```
</details>

---

## go-workflow manifest format
The go-workflow manifest is a YAML file used to define a workflow and actions to be performed, you can then load this into go-workflow and run a job. see the [`basic-example`](#how-to-use-go-workflow) under **How to use go-workflow** for more information on loading and running workflow manifest file. 

<details>
<summary>example manifest</summary>

```yaml

meta_data:
    name: simple-example
    description: This is a simple example workflow that demonstrates the basic features of the workflow engine.
    version: 1.0.0
    author: Andrew Pye
    contact: 
    create_date: "2022-11-13 11:39:44"
    update_date: "2022-11-13 11:39:44"
    vars: 
      example_value: "This is an example value"
jobs:
    - key: simple-example
      title:  Simple example
      description:  This job how to use loops and variables
      actions:
        - action: "print; {{ .Meta.Vars.example_value }}"
        - action: "for;i;0;{{ get_param `times_to_loop`}}"
        - action: "print; Hello World {{ get_stk_val `i` }}"
        - action: "next"
parameters: 
    - key: times_to_loop
      value: 100
    
```


</details>

<details>
<summary>1. meta_data</summary>
Meta data is used to describe the workflow it contains information about the package. the vars is used to allow the creator the option to sets values that can be used throughout the package in programming terms think of them as constant values. 

|   Field        | Description  | 
|-----------------|--------------|
| name    | Name of the workflow package  |
| description    | A description of the work flow manifest  |
| version    | The version of the workflow manifest  |
| author    | who wrote the manifest  |
| contact    | contact details of who wrote the manifest   |
| create_date    | The date the manifest was created  |
| update_date    |  The date the manifest was updated |
| vars    | variables that can be used in the manifest  |


</details>

<details>
<summary>2. jobs</summary>
The Job section is used to store jobs and its associates actions that will be performed.

|   Field        | Description  | 
|-----------------|--------------|
| *key    | A unique key for this job this is what is referenced to run the job  |
| title    | A title for the job  |
| description    | A description of the job  |
| *Actions    | A list of actions to perform  |


</details>

<details>
<summary>3. actions</summary>

The action section is used to configure the actions that are associated with the job

|   Field        | Description  | 
|-----------------|--------------|
| *key    | A unique key for the parameter that is used when injecting values into the payload  |
| action    | The name of the action to run  |
| description    | A description of the job  |
| fail    | what to do if the action fails (default end) |
| continue_on_error    | if action fails the continue to run (if true ignores fail field)  |
| config    | key pair value of configuration options for the action  |
| disabled    | disables the action (default false)  |
</details>

<details>
<summary>4. parameters</summary>

The parameters section is used to create parameters that can be specified by the end user and used to inject values into workflow

|   Field        | Description  | 
|-----------------|--------------|
| *key    | A unique key for the parameter that is used when injecting values into the payload  |
| title    | A title for the job  |
| description    | A description of the job  |
| InputType    | The data type of the parameter (this is currently not used but future plane to add validation) |
| *Value    | Default value to be used for this parameter  |

</details>

---

## Template Engine
To facilitate the ability to inject values into the manifest go-workflow uses golang's template engine. template tokens that are wrapped with **{{ }}** are used to inject values.You can get more information on the go template [here](https://pkg.go.dev/text/template). go-workflow parses each parameter and if present replaces the token with the required value.
To access the data in the manifest a model is passed to the template engine.

<details>
<summary>Template Data Model</summary>

``` Go
type TemplateData struct {
	Meta          *MetaData
	Manifest       *Manifest
	CurrentAction *Action
}
```

|   Field        | Description  | 
|-----------------|--------------|
| Meta    | Gives quick access to the Meta data in the manifest |
| Manifest    | Gives access to the manifest  |
| CurrentAction    | Gives access to the currently running action  |

Notice that the way to access the data is slightly different than how the manifest fields are written. the only exception is when you are specifying the var key (see example below **example_value**).

``` Yaml
    #Get value from the Vars in the Meta section of the manifest
    - action: "print; {{ .Meta.Vars.example_value }}"
    #Get the name of the manifest
    - action: "print; {{ .Manifest.Name }}"

```

### Reference to the manifest object model and valid data values.
Below shows the field names and the allowed data type for each field you can use to access the data in the manifest, the field names are case sensitive.

``` yaml
Meta:
    Name: (string)
    Description: (string)
    Version: (string)
    Author: (string)
    Contact: (string)
    CreatedDate: (date string)
    UpdateDate: (date string)
    Vars: {}  (key pair value)
Jobs:
    - Key: (string)
      Title: (string)
      Description: (string)
      Actions:
        - Key: (string)
          Action: (string)
          Description: (string)
          Fail:  (string/token)
          ContinueOnError: (bool/token)
          Config: (key pair value/token can be used as value)
          Disabled: (bool/token)
        
Parameters: 
    - Key: (string)
      Value: (any)
```

</details>

<details>
<summary>Template functions</summary>

you can use template functions to manipulate or access data. below is a table of the currently implements function, but it is possible to define your own see [here example 4](#how-to-use-go-workflow)

```yaml
#get a parameter value
- action: "for;i;0;{{ get_param `times_to_loop`}}"
#example of nesting functions
- action: "for;i;0;{{ get_param (lc `times_to_loop`)}}"
```

|   function        |params| Description  | 
|-----------------|-------|-------|
| read_file  |[string]|Reads a text file|
| base64enc    | [string]|base64 encode a string |
| base64dec    |[string] |base64 decode a string |
| gzip_base64    |[string] |zip a string and base 64 encode|
| lc    |[string] |make string lowercase  |
| uc    | [string]|make string uppercase  |
| domain    | [url string] |get the domain or ip from a url  |
| port_string    |[url string] |get the port of from a url as a string |
| port_int    |[url string] |get the port of from a url and returns an int |
| clean    | [string] [replace] |clean a string of spaces and special charts  |
| concat    |[string] ...[string]  |concatenate two or more strings together  |
| replace    |  [string] [find] [replace]|replace a value in a string  |
| contains    |  [string] [find] |searches a comma separated list for an instance returns a bool  |
| and    | [bool] [bool]|and two bool values |
| or    | [bool] [bool]|or two bool values |
| not    | [bool] [bool]|not two bool values |
| plus    | [int] [int]|add two int values together |
| minus    | [int] [int]|subtracts two int values together |
| multiply    | [int] [int]|multiply two int values together |
| get_stk_val    | [loop variable name string]|used to access the value in a loop variable |
| get_param    | [string]|get a value from a parameter specify the key value |
| get_store    | [key string] [name string]|get a value from the bucket store|
| get_data    | [key string]|Allows you to access the custom data from the manifest|
</details>

---
## Custom data

You can add arbitrary custom data to the manifest and then use this data either to inject using the template engine or in code when writing your actions.

There is a more detailed example in the examples folder:

**Example** example/workflow-custom-data

<details>
<summary>Examples</summary>

```yaml
...
jobs:
  - key: custom-example
    actions:
      - action: "print; **From the custom action**"
      - action:  custom
      - action: "print; **From the template**"
      - action: "print; {{(get_data `items`).GetArrayCount }}"
      - action: "for;i;0;{{minus (get_data `items`).GetArrayCount 1 }}"
      - action: "print;-{{(index .Manifest.Data.items (get_stk_val `i`)).msg}}"
      - action: "next"
data:
  items:
    - msg: "item 1"
    - msg: "item 2"
    - msg: "item 3"
    - msg: "item 4"
    - msg: "item 5"
  data_test: "y"
```

```go
val := m.Manifest.DataModel().GetMapItem("data_test").ToBool()
	fmt.Println(val)
	//*******************
	//Get the custom data
	//*******************
	item := m.Manifest.DataModel().GetMapItem("items").GetArray()

```

</details>

---

## Custom data code reference

<details>
<summary>reference</summary>

```go

// The main entry point to the data
w.DataModel() *Data

// GetArrayItem returns an item from the array
// - index: the index of the item to get
GetArrayItem(index int) *Data

// GetArrayCount returns the number of items in the array
GetArrayCount() int

// GetArray returns an array of Data objects
GetArray() []Data

// GetMap returns a map of Data objects
GetMap() map[string]Data

// GetMapItem a Data object
// - Key: the key to get
GetMapItem(key string) *Data

// GetInterface returns the data as an interface{}
GetInterface() interface{}

// GetType returns the type of the data as a string
GetType() string

// ToString returns the data as a string
ToString() string

// ToBool returns the data as a bool
ToBool() bool 

// ToFloat32 returns the data as a float32
ToFloat32() float32

// ToFloat64 returns the data as a float64
ToFloat64() float64

// ToInt returns the data as an int
ToInt() int 

// ToInt8 returns the data as an int8
ToInt8() int8

// ToInt32 returns the data as an int32
ToInt32() int32

// ToInt64 returns the data as an int64
ToInt64() int64


```
</details>

---
## inbuilt actions
go-workflow comes with some basic actions mainly around handling the flow. each parameter is separated with **;**. You can also use end or goto in the Fail field of the action

|   action        |params| Description  |
|-----------------|-------|-------|
| end    |N/A|Ends the work flow |
| print    |[any] |print a value |
| goto    |[action key] |jump to an action with the matching key |
| wait-seconds or wait    |[int] |wait for x seconds |
| wait-minutes    |[int] |wait for x minutes |
| for    |[variable];[start_value];[end_value]|for loop |
| next   |N/A|used withe the for loop to denote the end of the loop |
| error   |[message]]|causes the job to fail |

---

<details>
<summary>Examples</summary>

``` yaml

jobs:
    - key: demo-of-inbuilt-actions
      title: demo of actions
      description:  Demo of inbuilt actions
      actions:
        #Print examples
        - action: "print;hello world"
          fail: "goto; end_action"
        - action: "print;{{get_store `target` `port`}}"
        #Nested loops
        - action: "for;i;0;10"
          fail: "end"
        - action: "for;j;0;2"
        - action: "print;{{ get_stk_val `i` }} {{ get_stk_val `j` }}"
        - action: "next"
        - action: "next"
        #infante loop
        - action: "for"
        - action: "goto; end_action"
        - action: "next"
        #end action
        - action: "end"
            key: "end_action"
        
```
</details>

---

## Actions available in the package
To help get you started go-workflow has some Actions that if you choose can be added to the workflow engine.

|   action        |package| Description  |
|-----------------|-------|-------|
| CallApi    |actions\api|action for rest api calls |
| Action_Copy    |actions\file|Copies a file |
| Action_Rename    |actions\file|Renames a file |
| Action_Delete    |actions\file|Deletes a file |
| Action_RunJS    |actions\js|Run java script |
| Action_Store    |actions\store|Gives the ability to store values is the bucket store |
| Action_Store    |actions\store|Gives the ability to store values is the bucket store |
| Action_Condition    |actions\condition|Gives the ability to evaluate a condition and jump to a task |
| Action_SubWorkflow    |actions\sub_workflow|Gives the ability to run sub-workflow jobs |
| Action_Parallel    |actions\parallel_workflow|Allows you to run multiple actions in parallel |

## Actions Examples and Parameters

<details>
<summary>1. CallApi Action</summary>

This action enables you to make Rest API Calls

**Example:**  examples/workflow-call-api-action

|   field     |  Options | Description|
|-----------------|-------|-------|
| url    |   |target url |
| method    | POST,GET,PATCH,PUT,DELETE  | |
| body_type    | none,form-data,raw | |
| body    |   |body payload data |
|body_from_file||Allows you to read the body from file|
| header_    |   |set the headers |
| result_action    |none,print,js   |what to do with the returned results |
| result_format    |none,json,yaml,toml,xml,plain   |how to format the result|
| result_js    |   |js to process the result can be a file or inline |

### POST Example
```yaml
- action: api
    description: "This is an example of calling an API POST request."
    config:
      method: POST
      url: https://gorest.co.in/public/v2/users
      body_type: raw
      body: |
        {"name":"Agent Smith", "gender":"male", "email":"agent.smith@15ce.com", "status":"active"}
      header_Content-Type: application/json
      header_Authorization: "Bearer {{get_param `token`}}"
      result_action: "js"
      result_js: |
        function ActionResults(model,result){
          var obj=JSON.parse(result);
          store_value("api_result","user_id",obj.id);
          console(result);
          return true;
        }
```

### GET Example
```yaml
- action: api
    description: "This is an example of calling an API GET request."
    config:
      method: GET
      url: https://gorest.co.in/public/v2/users
      header_Content-Type: application/json
      header_Authorization: "Bearer {{get_param `token`}}"
      result_action: "js"
      result_js: |
        function ActionResults(model,result){
          console(result);
          return true;
        }
```
</details>

<details>
<summary>2. Copy Action</summary>

This action copies a file

**Example:**  examples/workflow-copy-file-action

|   field     |  Options | Description|
|-----------------|-------|-------|
| source_file    |   |Source file to copy |
| dest_file    |  | where to copy to (Use fullt path and filename) |



### Copy Example
```yaml
- action: copy
    config:
      source_file: "./source_file.txt"
      dest_file: "./destination_file.txt"
```

</details>

<details>
<summary>3. Delete Action</summary>

This action deletes a file

**Example:**  examples/workflow-copy-file-action

|   field     |  Options | Description|
|-----------------|-------|-------|
| source_file    |   |The file to delete |

### Copy Example
```yaml
- action: delete
    config:
      source_file: "./my_file.txt"
```

</details>

<details>
<summary>4. Rename Action</summary>

This action renames a file also acts as a move if your path is different

**Example:**  examples/workflow-copy-file-action

|   field     |  Options | Description|
|-----------------|-------|-------|
| source_file    |   |The file to rename |
| dest_file    |   |What to name it to (Full path)  |

### Copy Example
```yaml
- action: rename
    config:
      source_file: "./destination_file.txt"
      dest_file: "./my_file.txt"
```

</details>

<details>
<summary>5. JS Action</summary>

This action runs javascript

**Example:**  examples/workflow-js-action

|   field     |  Options | Description|
|-----------------|-------|-------|
| js    |   |javascript code to run|
| js_file    |   |read js in from a file |

```yaml
 actions:
  - action: js
    config:
      js_file: "./code2.js;./code3.js"
  - action: js
    config:
      js: |
        console(model.Meta.Name );
  - action: js
    config:
      js_file: "./code1.js"
```

</details>


<details>
<summary>6. Store Action</summary>

This action stores values in the data bucket

**Example:**  examples/workflow-store-action

|   field     |  Options | Description|
|-----------------|-------|-------|
| bucket    |   |the data bucket to use|
| key    |   |the key to store the value against |
| value    |   |value to save |

```yaml
 actions:
    - action: store
      config:
        bucket: "my_data"
        key: "name"
        value: "Andrew"
    - action: "print;{{get_store `my_data` `name`}}"
```

</details>

<details>
<summary>7. Condition Action</summary>

This action check for a condition and runs a task based on pass or fail

**Example:**  examples/workflow-store-action

|   field     |  Options | Description|
|-----------------|-------|-------|
| condition    |   |The condition that need to be evaluate must resolve to a bool|
| pass    |end [action key]   |the action to run or end |
| fail    |end [action key]   |the action to run or end |

```yaml
 - action: condition
    config:
      condition: "{{ get_param `times_to_loop`}} > 1 && {{ get_param `times_to_loop`}} < 5"
      pass: A
      fail: B
```

</details>


<details>
<summary>8. Sub-Workflow Action</summary>

This action runs a sub-workflow job the job it runs must have the 
- is_sub_workflow: true
- and inputs defines 

```YAML
is_sub_workflow: true
    inputs:
      value1:
      value2: test
      value3: test
```

**Example:**  examples/workflow-store-action

|   field     |  Options | Description|
|-----------------|-------|-------|
| job    |   |The sub job to run|
| inputs    |   |key value pair |


```yaml
jobs:
 - key: main-workflow-example
    title: Main Workflow
    description: This is the main workflow
    actions:
      - action: "sub-workflow"
        config:
          job: sub-workflow-example
          inputs:
            value1: "This is the second value"
  - key: sub-workflow-example
    title: Sub Workflow
    description: This is the Sub workflow
    is_sub_workflow: true
    inputs:
      value1:
    actions:
      - action: "print {{get_input `value1`}}"
 
```

</details>

<details>
<summary>9. Parallel Action</summary>

This action allows you to run multiple actions in parallel

```YAML
- key: parallel-example
    actions:
      - action: "parallel"
        config:
          actions:
            - action: "wait-seconds;2"
            - action: "wait-seconds;10"
            - action: print;Hello World"
      - action: "This is run after the parallel action"
```

**Example:**  examples/workflow-store-action

|   field     |  Options | Description|
|-----------------|-------|-------|
| actions    |   |List of actions|
| - action:    |   |the action and config the same as it would be under the job:  actions: |


```yaml
jobs:
 - key: main-workflow-example
    title: Main Workflow
    description: This is the main workflow
    actions:
      - action: "sub-workflow"
        config:
          job: sub-workflow-example
          inputs:
            value1: "This is the second value"
  - key: sub-workflow-example
    title: Sub Workflow
    description: This is the Sub workflow
    is_sub_workflow: true
    inputs:
      value1:
    actions:
      - action: "print {{get_input `value1`}}"
 
```

</details>
---

## How to use go-workflow
The quickest and easiest way to get started is by creating a workflow manifest. you can create a workflow programmatically but it is far easier to write the manifest in a YML file and load it into using the library. 

The [`git repo`](https://github.com/Mrpye/go-workflow) comes with some example that will cover the basics to more advanced features.
You can find the example in the examples folder


<details>
<summary>1. simple-example</summary>

This example creates a workflow that loops x number of time based on the value in times_to_loop parameter.

you can locate the example under: examples/simple-example

### main.go

To be able to run the workflow you need to create an 
- instance of go-workflow using **workflow.CreateWorkflow()** 
- load the workflow **wf.LoadManifest("./workflow.yaml")** 
- run the workflow **wf.RunJob("simple-example")**

also this example has set the logging to quiet so you will see the print action and any error.
try changing the logging levels
- workflow.LOG_QUIET   = 0
- workflow.LOG_INFO    = 1
- workflow.LOG_VERBOSE = 2



```go
package main

import "github.com/Mrpye/go-workflow/workflow"

func main() {
	//*****************
	//create a workflow
	//*****************
	wf := workflow.CreateWorkflow()

    //**********************************
	//Only show errors and print actions
	//**********************************
	wf.Verbose = workflow.LOG_QUIET

	//*************************
	//load the workflow manifest
	//*************************
	err := wf.LoadManifest("./workflow.yaml")
	if err != nil {
		println(err.Error())
	}

	//********************
	//Run the workflow job
	//********************
	err = wf.RunJob("simple-example")
	if err != nil {
		println(err.Error())
	}
}

```

### workflow.yaml


```yaml

meta_data:
    name: simple-example
    description: This is a simple example workflow that demonstrates the basic features of the workflow engine.
    version: 1.0.0
    author: Andrew Pye
    contact: 
    create_date: "2022-11-13 11:39:44"
    update_date: "2022-11-13 11:39:44"
    vars: {}
jobs:
    - key: simple-example
      title:  Simple example
      description:  This job how to use loops and variables
      actions:
        - action: "for;i;0;{{ get_param `times_to_loop`}}"
        - action: "print; Hello World {{ get_stk_val `i` }}"
        - action: "next"
parameters: 
    - key: times_to_loop
      value: 10
    
```
## Result

```bash
This is an example value
 Hello World 0 
 Hello World 1 
 Hello World 2 
 Hello World 3 
 Hello World 4 
 Hello World 5 
 Hello World 6 
 Hello World 7 
 Hello World 8 
 Hello World 9 
 Hello World 10
```

</details>

<details>
<summary>2. add-custom-actions-example</summary>

This example creates a workflow that shows how you can create your own custom actions

you can locate the example under: examples/add-custom-actions-example

### main.go

In this example we have created a custom action called **MultiPrint** and we added it to the workflow engine
- The function must use the following definition **func [FunctionName](w *workflow.Workflow, m *workflow.TemplateData) error**
- To add the function to the workflow engine use **wf.ActionList["FunctionName"] = FunctionName**
- To get values from the config use one of the following based on the data type you with to receive 
    - w.GetConfigTokenString
    - w.GetConfigTokenInt
    - w.GetConfigTokenBool
    - w.GetConfigTokenMap  
    - w.GetConfigTokenInterface  
- You can make the w.GetConfigToken optional by setting the 3rd parameter to false w.GetConfigTokenString("string_value", m, false)
- w.GetConfigTokenString([name of the config key]], [model to pass], [required])


```go
package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/Mrpye/go-workflow/workflow"
)

func main() {
	//*****************
	//create a workflow
	//*****************
	wf := workflow.CreateWorkflow()

	//**********************************
	//Only show errors and print actions
	//**********************************
	wf.Verbose = workflow.LOG_QUIET

	//*******************
	//Add a custom action
	//*******************
	wf.ActionList["MultiPrint"] = MultiPrint

	//*************************
	//load the workflow manifest
	//*************************
	err := wf.LoadManifest("./workflow.yaml")
	if err != nil {
		println(err.Error())
	}

	//********************
	//Run the workflow job
	//********************
	err = wf.RunJob("add-custom-actions-example")
	if err != nil {
		println(err.Error())
	}

}

//**************************
//print will print a message
//**************************
func MultiPrint(w *workflow.Workflow, m *workflow.TemplateData) error {
	//**********************************
	//Get a string value from the config
	//**********************************
	string_value, err := w.GetConfigTokenString("string_value", m, true)
	if err != nil {
		return err
	}
	//*******************************
	//Get a int value from the config
	//*******************************
	int_value, err := w.GetConfigTokenInt("int_value", m, true)
	if err != nil {
		return err
	}
	//********************************
	//Get a bool value from the config
	//********************************
	bool_value, err := w.GetConfigTokenBool("bool_value", m, true)
	if err != nil {
		return err
	}
	//*******************************
	//Get a map value from the config
	//*******************************
	map_value, err := w.GetConfigTokenMap("map_value", m, true)
	if err != nil {
		return err
	}

	//****************
	//Print the values
	//****************
	println(string_value)
	println(int_value)
	println(bool_value)
	b, err := json.Marshal(map_value)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(b))

	return nil
}


```

### workflow.yaml


```yaml

meta_data:
    name: add-custom-actions-example
    description: This example show how you can add your own custom actions to the workflow engine
    version: 1.0.0
    author: Andrew Pye
    contact: 
    create_date: "2022-11-13 11:39:44"
    update_date: "2022-11-13 11:39:44"
    vars: 
      example_value: "This is an example value"
jobs:
    - key: add-custom-actions-example
      title: Run the custom action
      description: This runs the custom action
      actions:
        - action: MultiPrint
          config:
            string_value: "{{.Meta.Vars.example_value }}"
            int_value: 100
            bool_value: true
            map_value:
              map_value1: "{{get_param `value1`}}"
              map_value2: "{{get_param `value2`}}"
              map_value3: "{{get_param `value3`}}"
parameters: 
    - key: value1
      value: "Hello world"
    - key: value2
      value: 55
    - key: value3
      value: false


```
## Result

```bash
This is an example value
100
true
{"map_value1":"Hello world","map_value2":"55","map_value3":"false"}
```

</details>

<details>
<summary>3. handling-result-data-storing-results</summary>

This example shows you how to process results and store the values for later use

you can locate the example under: examples/handling-result-data-storing-results

### main.go

In this example we use the **w.ActionProcessResults** process the results using js to extract data from the payload and store the data for use in a later action, we also use the model parameter in the **js** to get data from the manifest and store this.

- to use the **w.ActionProcessResults** you must supply parameters in the config section of your action
    - **result_action** (print/js) what to do
        - **print** just print the result
        - **js** allows you to run js to process the result
        - **default** just print the result
    - **result_format** (none/(json default)/yaml/toml/xml/plain) pre process the data into a particular format
    - **result_js** the js you wish to run you can use **result_js: |** for multi line 
- When you run **js** you need to wrap your code in the function **function ActionResults(model,result){}**
- also you need to return true for pass or false for fail the action
- you will need to pass the dat you wish to process as a parameter of **w.ActionProcessResults(data)** and make sure to return the error so the process fails on an error
- to store values in **js** you use the **store_value** function  **store_value([the data bucket],[the key],[value]);**
- you can retrieve the value using **get_store** function  **get_store([the data bucket],[the key]);**
- When using template function you can use **{{ get_store "the data bucket" "the key" }}** 

```go
package main

import (
	"encoding/json"

	"github.com/Mrpye/go-workflow/workflow"
)

func main() {
	//*****************
	//create a workflow
	//*****************
	wf := workflow.CreateWorkflow()

	//**********************************
	//Only show errors and print actions
	//**********************************
	wf.Verbose = workflow.LOG_QUIET

	//*******************
	//Add a custom action
	//*******************
	wf.ActionList["MultiPrint"] = MultiPrint

	//*************************
	//load the workflow manifest
	//*************************
	err := wf.LoadManifest("./workflow.yaml")
	if err != nil {
		println(err.Error())
	}

	//********************
	//Run the workflow job
	//********************
	err = wf.RunJob("handling-result-data-storing-results")
	if err != nil {
		println(err.Error())
	}

}

//**************************
//print will print a message
//**************************
func MultiPrint(w *workflow.Workflow, m *workflow.TemplateData) error {
	//*******************************
	//Get a map value from the config
	//*******************************
	map_value, err := w.GetConfigTokenMap("map_value", m, true)
	if err != nil {
		return err
	}

	//***************************************
	//Convert to json and process the results
	//***************************************
	b, err := json.Marshal(map_value)
	if err != nil {
		return err
	}

	//***********************************
	//This function processes the results
	//***********************************
	err = w.ActionProcessResults(string(b))
	if err != nil {
		return err
	}
	return nil
}



```

### workflow.yaml


```yaml
meta_data:
    name: handling-result-data-storing-results
    description: This example show how you can handle result data and store it in the store for later use.
    version: 1.0.0
    author: Andrew Pye
    contact: 
    create_date: "2022-11-13 11:39:44"
    update_date: "2022-11-13 11:39:44"
jobs:
    - key: handling-result-data-storing-results
      title: Run the custom action
      description: This runs the custom action
      actions:
        - action: MultiPrint
          config:
            map_value:
              map_value1: "{{get_param `value1`}}"
              map_value2: "{{get_param `value2`}}"
              map_value3: "{{get_param `value3`}}"
              map_value4: "This is a value from the config"
            result_action: "js"
            result_js: |
              function ActionResults(model,result){
                //parse the result
                var obj=JSON.parse(result);
                //print the value
                console(obj.map_value1);
                //Store the value1
                store_value("my_bucket","my_key",obj.map_value1);
                //print the value from the config
                console(model.CurrentAction.Config.map_value.map_value4)
                //store the value from the config
                store_value("my_bucket","stored_config_value",model.CurrentAction.Config.map_value.map_value4);
                // return true for success
                return true;
              }
        - action: "print;{{get_store `my_bucket` `my_key`}}"
        - action: "print;{{get_store `my_bucket` `stored_config_value`}}"
parameters: 
    - key: value1
      value: "Hello world"
    - key: value2
      value: 55
    - key: value3
      value: false



```
## Result

```bash
Hello world
This is a value from the config
Hello world
This is a value from the config
```

</details>


<details>
<summary>4. adding-start-and-cleanup-handlers</summary>

This example show how you can add start and cleanup handlers to your workflow. Also add new template functions to the workflow.

you can locate the example under: examples/adding-start-and-cleanup-handlers

### main.go

In this example we use the **wf.InitFunc** and **wf.CleanFunc** to handle startup and cleanup states. You function must follow the function patter **func FunctionName(w *workflow.Workflow, m *workflow.TemplateData) error**

- **InitFunc** is run at the start of a work flow job and **CleanFunc** is run at the end
- you can use InitFunc to read a config file and store the values so that the workflow can use them
- we used the **w.GetTemplateFuncMap()** to get the function map then added a new function in this case **AddHttp** with the key **add_http**
- to use function map use **w.SetTemplateFuncMap(template_functions)** and then you can use it {{add_http "localhost"}}
- if you use nested function calls you will need to wrap the child function call in () like so **{{add_http (get_store "target" "host")}}**

```go
package main

import (
	"fmt"

	"github.com/Mrpye/go-workflow/workflow"
)

func main() {
	//*****************
	//create a workflow
	//*****************
	wf := workflow.CreateWorkflow()

	//**********************************
	//Only show errors and print actions
	//**********************************
	wf.Verbose = workflow.LOG_QUIET

	//*************************************
	//Add the startup and cleanup functions
	//*************************************
	wf.InitFunc = Startup
	wf.CleanFunc = Clean

	//*************************
	//load the workflow manifest
	//*************************
	err := wf.LoadManifest("./workflow.yaml")
	if err != nil {
		println(err.Error())
	}

	//********************
	//Run the workflow job
	//********************
	err = wf.RunJob("adding-start-and-cleanup-handlers")
	if err != nil {
		println(err.Error())
	}

}

func AddHttp(val string) string {
	return fmt.Sprintf("http://%s", val)
}

func Startup(w *workflow.Workflow, m *workflow.TemplateData) error {
	//****************************************************
	//Save some values to the data bucket
	//You could read values from a config file or database
	//and save them to the data bucket
	//****************************************************
	w.SetValueToDataBucket("target", "host", "localhost")
	w.SetValueToDataBucket("target", "port", 8080)

	//**************************************************************
	// Add a custom function to the template function map
	//You don't necessarily have to do this in the startup function
	//You could do it in the main function
	//**************************************************************
	template_functions := w.GetTemplateFuncMap()
	template_functions["add_http"] = AddHttp

	//***********************************
	//Save it back to the workflow engine
	//***********************************
	w.SetTemplateFuncMap(template_functions)

	return nil
}

func Clean(w *workflow.Workflow, m *workflow.TemplateData) error {
	fmt.Println("Cleaning up")
	return nil
}


```

### workflow.yaml


```yaml
meta_data:
    name: adding-start-and-cleanup-handlers
    description: This example show how you can add start and cleanup handlers to your workflow. Also add new template functions to the workflow.
    version: 1.0.0
    author: Andrew Pye
    contact: 
    create_date: "2022-11-13 11:39:44"
    update_date: "2022-11-13 11:39:44"
jobs:
    - key: adding-start-and-cleanup-handlers
      title: Print the values
      description:  Get the values from the store that we set in the start handler
      actions:
        - action: "print;{{add_http (get_store `target` `host`)}}"
        - action: "print;{{get_store `target` `port`}}"
parameters: 

```
## Result

```bash
http://localhost
8080
Cleaning up
```

</details>


<details>
<summary>5. full-test-example</summary>

This example we are using workflow engine to test its features and validate the value are correct.
This example makes use of actions/tests package where there are some helper action that we can use to test new features.
ActionTest contains the tests and if an error occurs this will be passed by the RunJob function

you can locate the example under: examples/full-test-example

### main.go

```go
package main

import (
	"github.com/Mrpye/go-workflow/actions/store"
	"github.com/Mrpye/go-workflow/actions/tests"
	"github.com/Mrpye/go-workflow/workflow"
)

func main() {
	//*****************
	//create a workflow
	//*****************
	wf := workflow.CreateWorkflow()

	//**********************************
	//Only show errors and print actions
	//**********************************
	wf.Verbose = workflow.LOG_INFO

	//*******************
	//Add a custom action
	//*******************
	wf.ActionList["ActionStore"] = store.ActionStore
	wf.ActionList["ActionTest"] = tests.ActionTest
	wf.ActionList["ActionFailTest"] = tests.ActionFailTest
	wf.ActionList["ActionJSAndMap"] = tests.ActionJSAndMap

	//*************************
	//load the workflow manifest
	//*************************
	err := wf.LoadManifest("./workflow.yaml")
	if err != nil {
		println(err.Error())
		return
	}

	//********************
	//Run the workflow job
	//********************
	err = wf.RunJob("test-example")
	if err != nil {
		println(err.Error())
		return
	}

	println("Test Passed")

}


```

### workflow.yaml


```yaml
meta_data:
  name: test-example
  description: This is used for testing
  version: 1.0.0
  author: Andrew Pye
  contact: test@test.com
  create_date: "2022-11-13 11:39:44"
  update_date: "2022-11-13 11:39:44"
  vars:
    example_value: "This is an example value"
jobs:
  - key: test-example
    title: Simple example
    description: This job will test features of the workflow engine
    actions:
      - action: ActionJSAndMap
        config:
          map_value:
            map_value1: "{{get_param `times_to_loop`}}"
            map_value2: "{{get_param `test_string`}}"
            map_value3: "{{get_param `test_bool`}}"
            map_value4: "This is a value from the config"
          result_action: "js"
          result_js: |
              function ActionResults(model,result){
                //parse the result
                var obj=JSON.parse(result);
                //Store the value1
                store_value("test","js_map_value1",obj.map_value1);
                store_value("test","js_map_value2",obj.map_value2.toUpperCase());
                store_value("test","js_map_value3",obj.map_value3);
                store_value("test","js_map_value4",obj.map_value4.toUpperCase());
                return true;
              }
      # This action will store the value of the meta var example_value in the store
      - action: ActionStore
        config:
          bucket: "test"
          key: "meta_var"
          value: "{{ .Meta.Vars.example_value }}"
      - action: ActionStore
        config:
          bucket: "test"
          key: "meta_name"
          value: "{{ .Meta.Name }}"
      - action: ActionStore
        config:
          bucket: "test"
          key: "meta_description"
          value: "{{ .Meta.Description }}"
      - action: ActionStore
        config:
          bucket: "test"
          key: "meta_version"
          value: "{{ .Meta.Version }}"
      - action: ActionStore
        config:
          bucket: "test"
          key: "meta_author"
          value: "{{ .Meta.Author }}"
      - action: ActionStore
        config:
          bucket: "test"
          key: "meta_contact"
          value: "{{ .Meta.Contact }}"
      - action: ActionStore
        config:
          bucket: "test"
          key: "meta_create_date"
          value: "{{ .Meta.CreatedDate }}"
      - action: ActionStore
        config:
          bucket: "test"
          key: "meta_update_date"
          value: "{{ .Meta.UpdateDate }}"
      - action: ActionStore
        config:
          bucket: "test"
          key: "param_test_int"
          value: "{{ get_param `times_to_loop` }}"
      - action: ActionStore
        config:
          bucket: "test"
          key: "param_test_string"
          value: "{{ get_param `test_string` }}"
      - action: ActionStore
        config:
          bucket: "test"
          key: "param_test_bool"
          value: "{{ get_param `test_bool` }}"

      - action: "for;i;0;{{ get_param `times_to_loop`}}"
      - action: ActionStore
        config:
          bucket: "test"
          key: "loop_increment{{ get_stk_val `i`}}"
          value: "{{ get_stk_val `i`}}"
      - action: "next"

      - action: "for;i;{{ get_param `times_to_loop`}};0"
      - action: ActionStore
        config:
          bucket: "test"
          key: "loop_decrement{{ get_stk_val `i`}}"
          value: "{{ get_stk_val `i`}}"
      - action: "next"

      - action: "for;i;0;{{ get_param `times_to_loop`}}"
      - action: "for;j;0;{{ get_param `times_to_loop`}}"
      - action: ActionStore
        config:
          bucket: "test"
          key: "nested_loop_{{ get_stk_val `i`}}-{{ get_stk_val `j`}}"
          value: "{{ get_stk_val `i`}}-{{ get_stk_val `j`}}"
      - action: "next"
      - action: "next"
      - action: "goto;ActionTest"
      - action: ActionFailTest
        key: ActionFailTest
      - action: ActionTest
        key: ActionTest
parameters:
  - key: times_to_loop
    value: 3
  - key: test_string
    value: "this is a test string"
  - key: test_bool
    value: true

```
## Result

```bash
http://localhost
8080
Cleaning up
```

</details>

<details>
<summary>6. condition-example</summary>

This example shows how we can check for conditions and then jump to a action. if you set the times_to_loop between 1 and 5 it will use the pass condition
else it will fail. This also shows how you can manipulate the next task to be run using a Custom action.

you can locate the example under: examples/condition-example

### main.go

```go
package main

import (
	"github.com/Mrpye/go-workflow/actions/condition"
	"github.com/Mrpye/go-workflow/workflow"
)

func main() {
	//*****************
	//create a workflow
	//*****************
	wf := workflow.CreateWorkflow()

	//**********************************
	//Only show errors and print actions
	//**********************************
	wf.Verbose = workflow.LOG_INFO
	wf.ActionList["condition"] = condition.Action_Condition

	//*************************
	//load the workflow manifest
	//*************************
	err := wf.LoadManifest("./workflow.yaml")
	if err != nil {
		println(err.Error())
	}

	//********************
	//Run the workflow job
	//********************
	err = wf.RunJob("condition-example")
	if err != nil {
		println(err.Error())
	}
}

```

### workflow.yaml


```yaml
meta_data:
    name: condition-example
    description: This is a example we use an action to evaluate a condition then use the result to decide what to do next
    version: 1.0.0
    author: Andrew Pye
    contact: 
    create_date: "2022-11-13 11:39:44"
    update_date: "2022-11-13 11:39:44"
jobs:
    - key: condition-example
      title:  Simple example
      description:  This job how to use loops and variables
      actions:
        - action: condition
          config:
            condition: "{{ get_param `times_to_loop`}} > 1 && {{ get_param `times_to_loop`}} < 5"
            pass: A
            fail: B
        - action: end
        - action: "print; This is A"
          key: A
        - action: end
        - action: "print; This is B"
          key: B
        
parameters: 
    - key: times_to_loop
      value: 1

```
## Result

```bash
2023/02/23 10:16:04 *************************
2023/02/23 10:16:04 ** Action: ->condition **
2023/02/23 10:16:04 *************************
2023/02/23 10:16:04 condition failed going to action B: Fail
2023/02/23 10:16:04 ---------------------------------------
2023/02/23 10:16:04 ** Action Completed: ->condition: OK **
2023/02/23 10:16:04 ---------------------------------------


2023/02/23 10:16:04 *********************************
2023/02/23 10:16:04 ** Action: B->print; This is B **
2023/02/23 10:16:04 *********************************
 This is B
2023/02/23 10:16:04 -----------------------------------------------
2023/02/23 10:16:04 ** Action Completed: B->print; This is B: OK **
2023/02/23 10:16:04 -----------------------------------------------
```

</details>

<details>
<summary>7. Custom data and Call Custom Action from withing a Custom Action</summary>

This action allows you to run multiple actions in parallel

**Example:**  examples/workflow-custom-data

</details>

---
## Change Log
### v0.1.0
  First build 

### v0.2.0
  - Added lots of examples
  - updated the readme-
  - Added some actions that can be used
    - API Call
    - Copy File
    - Rename File
    - Delete File
    - Store data
    - Run JS

### v0.3.0
  - Added Condition Action
  - Added Sub Workflow Action
  - Added ability for actions to change the next action to run
  - Added Sub-Workflow option and inputs to job
  - Added a CreateSubWorkflow function so actions can run a workflow inside the workflow
  - Added Parallel Action so you can run multiple actions at the same time
  - Added Custom Data to the manifest and added the ability to add via code to template engine
  - Added More examples

---

## To Do
- Unit Test for actions
- Documentation the package


--- 


## Some notable 3rd party Libraries
- [https://github.com/dop251/goja](https://github.com/dop251/goja) JS engine


---

## license
go-workflow is Apache 2.0 licensed.