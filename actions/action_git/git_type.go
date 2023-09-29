package action_git

import (
	"fmt"
)

// K8 is the struct for the k8 connection
// If you want to use the kube config file, you need to set the DefaultContext and ConfigPath
// DefaultContext is the default context to use
// ConfigPath is the path to the kube config file
// Host is the host to connect to
// You will need to set the Host and Authorization if you want to use the token connection
// Authorization is the authorization token
// UseTokenConnection if true, use the token connection, otherwise use the kube config file
// Ignore_ssl if true, ignore the ssl connection
type Git struct {
	Host          string `json:"host" yaml:"host" flag:"host h" desc:"the host url"`
	User          string `json:"user" yaml:"user" flag:"user u" desc:"The user name"`
	Authorization string `json:"authorization" yaml:"authorization" flag:"auth a" desc:"The auth token"`
	Ssh           string `json:"ssh" yaml:"ssh" flag:"ssh s" desc:"The SSH key"`
	Email         string `json:"email" yaml:"email" flag:"email e" desc:"The Email address"`
}

// K8Option is the option for the k8 connection
type GitOption func(*Git)

// OptionK8DefaultContext is the option for the default context
func OptionGitHost(host string) GitOption {
	return func(h *Git) {
		h.Host = host
	}
}

// OptionK8ConfigPath is the option for the config path
func OptionGitUser(user string) GitOption {
	return func(h *Git) {
		h.User = user
	}
}

// OptionK8Host is the option for the host
func OptionGitPassword(password string) GitOption {
	return func(h *Git) {
		h.Authorization = password
	}
}

// OptionK8Auth is the option for the authorization
func OptionGitSSH(ssh string) GitOption {
	return func(h *Git) {
		h.Ssh = ssh
	}
}

// OptionK8IgnoreSSL is the option for the ignore ssl
func OptionGitEmail(email string) GitOption {
	return func(h *Git) {
		h.Email = email
	}
}

// Update the k8 Type with the options
func (m *Git) Update(opts ...GitOption) {
	// Loop through each option
	for _, opt := range opts {
		// Call the option giving the instantiated
		opt(m)
	}
}

// String returns the string representation of the k8 type
func (m *Git) String() string {
	return fmt.Sprintf("%s,%s,%s", m.Host, m.User, m.Email)
}

// Create a instance of the k8 type
// default_context is the default context to use
// config_path is the path to the kube config file
func CreateGit(host string, authorization string) (*Git, error) {
	git := &Git{
		Host:          host,
		Authorization: authorization,
	}
	return git, nil

}

// CreateK8Options creates a instance of the k8 type
// opts are the options for the k8 type
func CreateGitOptions(opts ...GitOption) (*Git, error) {
	git := &Git{}
	git.Update(opts...)
	return git, nil
}
