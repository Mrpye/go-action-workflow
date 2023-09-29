// This is the workflow action for calling an api
package action_api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"strings"

	"github.com/Mrpye/go-action-workflow/workflow"
	"github.com/Mrpye/golib/file"
	"github.com/Mrpye/golib/net"
)

// Constants for the body data type
const (
	BODY_DATA_TYPE_FORM_DATA             = "form-data"
	BODY_DATA_TYPE_NONE                  = "none"
	BODY_DATA_TYPE_RAW                   = "raw"
	BODY_DATA_TYPE_BINARY                = "binary"
	BODY_DATA_TYPE_X_WWW_FORM_URLENCODED = "x-www-form-urlencoded"
)

type ApiSchema struct {
}

// GetTargetSchema returns the target schema for this action
func (s ApiSchema) GetTargetSchema() map[string]workflow.TargetSchema {
	//Build a test schema
	return map[string]workflow.TargetSchema{}
}

// GetActions returns the actions for this schema
func (s ApiSchema) GetActionSchema() map[string]workflow.ActionSchema {

	//no short d f h
	return map[string]workflow.ActionSchema{
		"api": {
			Action:         Action_CallApi,
			Short:          "Rest Api Call",
			Long:           "Rest Api Call",
			ProcessResults: true,
			ConfigSchema: map[string]*workflow.Schema{
				"url": {
					Type:        workflow.TypeString,
					Partial:     false,
					Required:    true,
					Description: "Url to call",
					Short:       "u",
					Default:     "",
				},
				"method": {
					Type:        workflow.TypeString,
					Partial:     false,
					Required:    true,
					Description: "Method to call POST, GET, PUT, DELETE",
					Short:       "m",
					Default:     "",
				},
				"ignore_ssl": {
					Type:        workflow.TypeBool,
					Partial:     false,
					Required:    false,
					Description: "Ignore SSL errors",
					Short:       "i",
					Default:     false,
				},
				"body_type": {
					Type:        workflow.TypeString,
					Partial:     false,
					Required:    false,
					Description: "Type of body to send",
					Short:       "t",
					Default:     false,
				},
				"body": {
					Type:        workflow.TypeString,
					Partial:     false,
					Required:    false,
					Description: "Body to send",
					Short:       "b",
					Default:     false,
				},
				"body_from_file": {
					Type:        workflow.TypeString,
					Partial:     false,
					Required:    false,
					Description: "Body to send from file",
					Short:       "r",
					Default:     false,
				},
			},
		},
	}
}

// GetFunctionMap returns the function map for this schema
func (s ApiSchema) GetFunctionMap() map[string]workflow.FunctionSchema {
	return map[string]workflow.FunctionSchema{}
}

// GetAction returns the action for the target
func GetSchema() workflow.SchemaEndpoint {
	return ApiSchema{}
}

func create_payload(w *workflow.Workflow, m *workflow.TemplateData) (interface{}, string, error) {
	//TODO need to add raw type
	//*********************************************
	//Get the body data to see what type of payload
	//*********************************************
	body_data_type, err := w.GetConfigTokenString("body_type", m, false)
	if err != nil {
		return nil, "", err
	}

	//*******************************************
	//if Body data does not exist default to none
	//*******************************************
	if body_data_type == "" {
		body_data_type = BODY_DATA_TYPE_NONE
	}

	//********************
	//Handle the body data
	//********************
	if body_data_type == BODY_DATA_TYPE_FORM_DATA {
		//****************
		//Convert the data
		//****************
		i_body := m.CurrentAction.GetConfig("body")
		//TODO:Need to check of accepted type see if yaml
		data, err := json.Marshal(i_body)
		if err != nil {
			return nil, "", err
		}
		read_body_js := false
		switch val := i_body.(type) {
		case string:
			if val == "" {
				read_body_js = true
			}
		}
		//*********************
		//Read the data from js
		//*********************
		if read_body_js {
			body_js, err := w.GetConfigTokenString("body_from_file", m, false)
			body_js = w.BuildPath(body_js)
			if err != nil {
				return nil, "", err
			}
			if body_js == "" && !file.FileExists(body_js) {
				return nil, "", fmt.Errorf("no file found for %s", body_js)
			}
			file_data, err := file.ReadFileToString(body_js)
			if err != nil {
				return nil, "", err
			}
			data = []byte(file_data)
		}

		//****************************
		//Read form data into key pair
		//****************************
		var form_data []net.Header
		err = json.Unmarshal(data, &form_data)
		if err != nil {
			return nil, "", err
		}

		//***********************
		//Lets replace any tokens
		//***********************
		for i := range form_data {
			parsed_str, err := w.ParseToken(m, string(form_data[i].Value))
			if err != nil {
				return nil, "", err
			}
			form_data[i].Value = parsed_str
		}

		//****************************
		//Convert the data to a reader
		//****************************
		payload := &bytes.Buffer{}
		writer := multipart.NewWriter(payload)
		for _, o := range form_data {
			err = writer.WriteField(o.Key, o.Value)
			if err != nil {
				return nil, "", err
			}
		}
		err = writer.Close()
		if err != nil {
			return nil, "", err
		}
		return payload, writer.FormDataContentType(), nil
	} else if body_data_type == BODY_DATA_TYPE_RAW {
		i_body := m.CurrentAction.GetConfig("body")
		read_body_js := false
		switch val := i_body.(type) {
		case string:
			if val == "" {
				read_body_js = true
			}
			val, err := w.ParseToken(m, string(val))
			if err != nil {
				return nil, "", err
			}
			return strings.NewReader(val), "", nil
		case map[string]interface{}:
			val = w.ParseInterfaceMap(m, val)
			data, err := json.Marshal(val)
			if err != nil {
				return nil, "", err
			}
			return bytes.NewReader(data), "", nil
		}
		if read_body_js {
			body_js, err := w.GetConfigTokenString("body_from_file", m, false)
			if err != nil {
				return nil, "", err
			}
			body_js = w.BuildPath(body_js)
			if body_js == "" && !file.FileExists(body_js) {
				return nil, "", fmt.Errorf("no file found for %s", body_js)
			}
			file_data, err := file.ReadFileToString(body_js)
			if err != nil {
				return nil, "", err
			}
			file_data, err = w.ParseToken(m, file_data)
			if err != nil {
				return nil, "", err
			}
			return strings.NewReader(file_data), "", nil
		}
	} else if body_data_type == BODY_DATA_TYPE_NONE {
		//Body Data None
		return nil, "", nil
	}

	return nil, "", fmt.Errorf("no body_type of type %s found", body_data_type)
}

// CallApi is the main function for the action
/*
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
		function ActionResults(m,result){
			var obj=JSON.parse(result);
			store_value("api_result","user_id",obj.id);
			console(result);
			return true;
		}
*/
func Action_CallApi(w *workflow.Workflow, m *workflow.TemplateData) error {

	//*********************
	//Get the config values
	//*********************
	method, err := w.GetConfigTokenString("method", m, true)
	if err != nil {
		return err
	}

	url, err := w.GetConfigTokenString("url", m, true)
	if err != nil {
		return err
	}

	ignore_ssl, err := w.GetConfigTokenBool("ignore_ssl", m, false)
	if err != nil {
		return err
	}

	//*********************
	//Deal with the headers
	//*********************
	var headers []net.Header
	for k, e := range m.CurrentAction.Config {
		if strings.Contains(k, "header_") {
			key := strings.ReplaceAll(k, "header_", "")
			//******************
			//Replace the tokens
			//******************
			value, err := w.ParseToken(m, string(e.(string)))
			if err != nil {
				return err
			}
			header := net.Header{
				Key:   key,
				Value: value,
			}
			headers = append(headers, header)
		}
	}
	//******************
	//Create the payload
	//******************
	payload, content_type, err := create_payload(w, m)
	if err != nil {
		return err
	}

	//**************************************************************
	//If a content type is returned then set the content type header
	//**************************************************************
	if content_type != "" {
		header := net.Header{
			Key:   "Content-Type",
			Value: content_type,
		}
		headers = append(headers, header)
	}

	//********************************
	//See if we need to pass a payload
	//********************************
	var passed_payload io.Reader
	if payload != nil {
		switch val := payload.(type) {
		case *bytes.Buffer:
			passed_payload = val
		case io.Reader:
			passed_payload = val
		}

	}

	//**********************
	//Make a call to the api
	//**********************
	data, _, err := net.CallApi(url, method, headers, passed_payload, ignore_ssl)
	if err != nil {
		return err
	}

	//*******************
	//Process the results
	//*******************
	err = w.ActionProcessResults(m, string(data))
	if err != nil {
		return err
	}

	if w.LogLevel > workflow.LOG_INFO {
		fmt.Printf("file %s\n", data)
	}

	return nil
}
