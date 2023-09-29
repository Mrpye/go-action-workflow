package action_scp

import (
	"fmt"
)

type SCP struct {
	Host     string `json:"host" yaml:"host" flag:"host h" desc:"SCP host"`
	User     string `json:"user" yaml:"user" flag:"user u" desc:"User for SCP"`
	Password string `json:"password" yaml:"password" flag:"password p" desc:"Password for SCP"`
}

type SCPOptions func(*SCP)

func OptionSCPHost(host string) SCPOptions {
	return func(h *SCP) {
		h.Host = host
	}
}
func OptionSCPUser(user string) SCPOptions {
	return func(h *SCP) {
		h.User = user
	}
}
func OptionSCPPassword(password string) SCPOptions {
	return func(h *SCP) {
		h.Password = password
	}
}

func (m *SCP) Update(opts ...SCPOptions) {
	// Loop through each option
	for _, opt := range opts {
		// Call the option giving the instantiated
		opt(m)
	}
}

func (m *SCP) String() string {
	return fmt.Sprintf("%s,%s", m.Host, m.User)
}

// Create a instance of the docker registry
func CreateSCP(env string, host string, user string, password string) *SCP {
	return &SCP{
		Host:     host,
		User:     user,
		Password: password,
	}
}

func CreateSCPOptions(env string, opts ...SCPOptions) *SCP {
	obj := &SCP{}
	obj.Update(opts...)
	return obj
}
