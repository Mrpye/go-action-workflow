package action_git

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Mrpye/golib/net"
)

const (
	TARGET_TYPE_GITLAB = "gitlab"
	TARGET_TYPE_GITHUB = "github"
)

func (m *Git) DownloadFromGitLab(project string, File string, branch string) (string, error) {
	url := fmt.Sprintf("%s/api/v4/projects/%s/repository/files/%s/raw?ref=%s", m.Host, project, File, branch)
	headers := []net.Header{
		{Key: "Content-Type", Value: "application/json"},
		{Key: "Accept", Value: "application/json"},
		{Key: "PRIVATE-TOKEN", Value: m.Authorization},
	}
	res, success, err := net.CallApi(url, "GET", headers, nil, false)
	if !success {
		return "", errors.New(string(res))
	}
	//if string(res) == "{\"message\":\"401 Unauthorized\"}" {
	//	return string(res), errors.New("invalid auth token")
	//}
	return string(res), err
}

func (m *Git) DownloadFromGitHub(project string, File string, branch string) (string, error) {
	url := fmt.Sprintf("%s/repos/%s/contents/%s?ref=%s", m.Host, project, File, branch)
	headers := []net.Header{
		{Key: "Content-Type", Value: "application/json"},
		{Key: "Accept", Value: "application/vnd.github+json"},
		//{Key: "Content-Type", Value: "application/json"},
		//{Key: "Accept", Value: "*/*"},
		{Key: "Authorization", Value: "Bearer " + m.Authorization},
	}
	res, success, err := net.CallApi(url, "GET", headers, nil, false)
	if !success {
		return "", errors.New(string(res))
	}
	data := make(map[string]interface{})
	json.Unmarshal(res, &data)
	if data["content"] != nil {
		content := data["content"].(string)
		res, err = base64.StdEncoding.DecodeString(content)
		if err != nil {
			return "", err
		}
	}

	return string(res), err
}

func (m *Git) DownloadGitFile(target_type string, project string, File string, branch string) (string, error) {
	if target_type == TARGET_TYPE_GITLAB {
		//gitlab
		data, err := m.DownloadFromGitLab(project, File, branch)
		if err != nil {
			return "", err
		}
		return data, nil
	} else if target_type == TARGET_TYPE_GITHUB {
		//github
		data, err := m.DownloadFromGitHub(project, File, branch)
		if err != nil {
			return "", err
		}
		return data, nil
	} else {
		return "", errors.New("invalid url")
	}
}
