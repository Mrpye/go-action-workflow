# Plugin Library Documentation
Below are the instructions for using the plugins library in your workflow. Below are the setting and parameters for each of the actions, target and function included in the library.

---

{{range $action_group_key, $action_group := .GetGroupedActions}}
# **Plugin Library**:  {{$action_group_key}}

# Actions:

{{range $action_key, $action := $action_group}}
## **Action**: {{$action_key}}
{{$action.Long}}

<details>
<summary>Info</summary>

|Property|Description|
|----|----|
|Action|{{$action_key}}|
|Description|{{$action.Short}}|
|Target|[{{$action.Target}}](#target-{{replace $action.Target `.` ``}})|
|InlineParams|{{$action.InlineParams}}
|ProcessResults|{{$action.ProcessResults}}|


### Config Parameters

|parameter|Type|Description|Required|Short Flag|Default|
|----|----|----|----|----|----|
{{range $i, $e := $action.ConfigSchema}}|{{$i}}|{{$e.Type.String}}|{{$e.Description}}|{{$e.Required}}|{{$e.Short}}||{{$e.Default}}|
{{end}}

</details>

---
{{end}}

# Functions:

## {{$action_group_key}} Functions:

{{range $function_key, $function := $.GetFunctionsByLibrary $action_group_key}}
## **Function**: {{$function_key}}
{{$function.Description}}

<details>
<summary>Info</summary>

|Property|Description|
|----|----|
|Function|{{$function_key}}|
|Description|{{$function.Description}}|

### Parameters
|parameter|Type|Description|Required|
|----|----|----|----|
{{range $i, $e := $function.ParameterSchema}}|{{$i}}|{{$e.Type.String}}|{{$e.Description}}|{{$e.Required}}|
{{end}}

</details>

---
{{end}}
---
{{end}}

# Targets:
{{range $target_key, $target := .TargetScheme}}
## **Target** {{$target_key}}

<details>
<summary>Info</summary>

|Field|json|cli flag|Desc|
|----|----|----|----|
{{range $i, $e := $target.GetTargetMap}}|{{$i}}|{{tag $target.Target $i `yaml`}}|{{tag $target.Target $i `flag`}}|{{tag $target.Target $i `desc`}}|
{{end}}

</details>

---
{{end}}


