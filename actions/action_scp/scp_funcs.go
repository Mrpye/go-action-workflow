package action_scp

import (
	"context"
	"fmt"
	"os"

	scp "github.com/bramvdbogaerde/go-scp"
	"github.com/bramvdbogaerde/go-scp/auth"
	"golang.org/x/crypto/ssh"
)

func (m *SCP) Test() (bool, string, error) {
	// Use SSH key authentication from the auth package
	// we ignore the host key in this example, please change this if you use this library
	clientConfig, _ := auth.PasswordKey(m.User, m.Password, ssh.InsecureIgnoreHostKey())

	// For other authentication methods see ssh.ClientConfig and ssh.AuthMethod

	// Create a new SCP client
	client := scp.NewClient(m.Host, &clientConfig)

	// Connect to the remote server
	err := client.Connect()
	if err != nil {
		return false, "", fmt.Errorf("couldn't establish a connection to the remote server %s ", err)
	}

	// Close client connection after the file has been copied
	defer client.Close()

	return true, "", nil
}

func (m *SCP) SCPFileUpload(source string, target string) error {
	// Use SSH key authentication from the auth package
	// we ignore the host key in this example, please change this if you use this library
	clientConfig, _ := auth.PasswordKey(m.User, m.Password, ssh.InsecureIgnoreHostKey())

	// For other authentication methods see ssh.ClientConfig and ssh.AuthMethod

	// Create a new SCP client
	client := scp.NewClient(m.Host, &clientConfig)

	// Connect to the remote server
	err := client.Connect()
	if err != nil {
		return fmt.Errorf("couldn't establish a connection to the remote server %s ", err)
	}

	// Open a file
	f, err := os.Open(source)
	if err != nil {
		return fmt.Errorf("could not open file %s  %s", source, err)
	}
	// Close client connection after the file has been copied
	defer client.Close()

	// Close the file after it has been copied
	defer f.Close()

	ctx := context.Background()
	err = client.CopyFile(ctx, f, target, "0655")

	if err != nil {
		return fmt.Errorf("error while copying file %s", err)
	}
	return nil
}
func (m *SCP) SCPFileDownload(source string, target string) error {
	// Use SSH key authentication from the auth package
	// we ignore the host key in this example, please change this if you use this library
	clientConfig, _ := auth.PasswordKey(m.User, m.Password, ssh.InsecureIgnoreHostKey())

	// For other authentication methods see ssh.ClientConfig and ssh.AuthMethod

	// Create a new SCP client
	client := scp.NewClient(m.Host, &clientConfig)

	// Connect to the remote server
	err := client.Connect()
	if err != nil {
		return fmt.Errorf("couldn't establish a connection to the remote server %s ", err)
	}

	// Open a file
	f, err := os.Open(target)
	if err != nil {
		return fmt.Errorf("could not open file %s  %s", target, err)
	}
	// Close client connection after the file has been copied
	defer client.Close()

	// Close the file after it has been copied
	defer f.Close()
	ctx := context.Background()
	err = client.CopyFromRemote(ctx, f, source)

	if err != nil {
		return fmt.Errorf("error while copying file %s", err)
	}
	return nil
}
