package action_docker

import "fmt"

type Registry struct {
	Host      string `json:"host" yaml:"host"  flag:"host h" desc:"the host url"`
	UserName  string `json:"user" yaml:"user" flag:"user u" desc:"Username for the registry"`
	Password  string `json:"password" yaml:"password" flag:"password p" desc:"Password for the registry"`
	IgnoreSSL bool   `json:"ignore_ssl" yaml:"ignore_ssl"  flag:"ignore_ssl i" desc:"Ignore SSL"`
	Library   string `json:"library" yaml:"library" flag:"library l" desc:"Library to use"`
}

type Manifest []struct {
	Config   string   `json:"Config"`
	RepoTags []string `json:"RepoTags"`
	Layers   []string `json:"Layers"`
}

// K8Option is the option for the k8 connection
type DockerRegistryOption func(*Registry)

//OptionDockerRegHost is the option for the host
func OptionDockerRegHost(host string) DockerRegistryOption {
	return func(h *Registry) {
		h.Host = host
	}
}

// OptionDockerRegUser is the option for the user
func OptionDockerRegUser(user string) DockerRegistryOption {
	return func(h *Registry) {
		h.UserName = user
	}
}

// OptionDockerRegPassword is the option for the password
func OptionDockerRegPassword(password string) DockerRegistryOption {
	return func(h *Registry) {
		h.Password = password
	}
}

// OptionDockerRegIgnoreSSL is the option for the ignore ssl
func OptionDockerRegIgnoreSSL(ignore_ssl bool) DockerRegistryOption {
	return func(h *Registry) {
		h.IgnoreSSL = ignore_ssl
	}
}

// Update the k8 Type with the options
func (m *Registry) Update(opts ...DockerRegistryOption) {
	// Loop through each option
	for _, opt := range opts {
		// Call the option giving the instantiated
		opt(m)
	}
}

// String returns the string representation of the k8 type
func (m *Registry) String() string {
	return fmt.Sprintf("%s,%s,%v", m.Host, m.UserName, m.IgnoreSSL)
}

// Function to create the Registry struct
func CreateDockerRegistry(host string, UserName string, Password string, library string, IgnoreSSL bool) (*Registry, error) {
	return &Registry{
		Host:      host,
		UserName:  UserName,
		Password:  Password,
		IgnoreSSL: IgnoreSSL,
		Library:   library,
	}, nil
}
