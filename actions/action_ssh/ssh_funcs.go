package action_ssh

import (
	"fmt"

	ssh_client "github.com/helloyi/go-sshclient"
)

func (m *SSH) CreateSSHClient() (*ssh_client.Client, error) {
	if m.PrivateKeyFile == "" && m.Password != "" {
		client, err := ssh_client.DialWithPasswd(m.Host, m.User, m.Password)
		if err != nil {
			return nil, err
		}
		return client, nil
	} else if m.PrivateKeyFile != "" && m.Password != "" {
		client, err := ssh_client.DialWithKeyWithPassphrase(m.Host, m.User, m.PrivateKeyFile, m.Password)
		if err != nil {
			return nil, err
		}
		return client, nil
	} else if m.PrivateKeyFile != "" {
		client, err := ssh_client.DialWithKey(m.Host, m.User, m.PrivateKeyFile)
		if err != nil {
			return nil, err
		}
		return client, nil
	}
	return nil, fmt.Errorf("no valid SSH credentials found")
}

func (m *SSH) SSHRunCMD(cmd string) (string, error) {

	client, err := m.CreateSSHClient()
	if err != nil {
		return "", err
	}
	defer client.Close()

	// run one command
	out, err := client.Cmd(cmd).Output()
	if err != nil {
		return string(out), err
	}

	return string(out), nil
}

func (m *SSH) SSHRunScript(script string) (string, error) {

	client, err := m.CreateSSHClient()
	if err != nil {
		return "", err
	}
	defer client.Close()

	// run one command
	out, err := client.Script(script).Output()
	if err != nil {
		return string(out), err
	}

	return string(out), nil
}

func (m *SSH) SSHRunScriptFile(script_file string) (string, error) {

	client, err := m.CreateSSHClient()
	if err != nil {
		return "", err
	}
	defer client.Close()

	// run one command
	out, err := client.ScriptFile(script_file).Output()
	if err != nil {
		return string(out), err
	}

	return string(out), nil
}

func (m *SSH) SSHUploadFile(src string, dest string) error {

	client, err := m.CreateSSHClient()
	if err != nil {
		return err
	}
	defer client.Close()

	// run one command
	sftp := client.Sftp()
	// upload
	if err := sftp.Upload(src, dest); err != nil {
		return err
	}
	if err := sftp.Close(); err != nil {
		return err
	}

	return nil
}
func (m *SSH) SSHDownloadFile(src string, dest string) error {

	client, err := m.CreateSSHClient()
	if err != nil {
		return err
	}
	defer client.Close()

	// run one command
	sftp := client.Sftp()
	// upload

	if err := sftp.Download(src, dest); err != nil {
		return err
	}
	if err := sftp.Close(); err != nil {
		return err
	}

	return nil
}
