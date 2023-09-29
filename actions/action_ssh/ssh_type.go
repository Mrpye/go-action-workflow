package action_ssh

import (
	"fmt"
)

type SSH struct {
	Host           string `json:"host" yaml:"host" flag:"host h" desc:"SSH host"`
	User           string `json:"user" yaml:"user" flag:"user u" desc:"User for SSH"`
	Password       string `json:"password" yaml:"password" flag:"password p" desc:"Password for SSH"`
	PrivateKeyFile string `json:"private_key_file" yaml:"private_key_file" flag:"private_key_file s" desc:"private key file for SSH"`
}

type SSHOptions func(*SSH)

func OptionSSHHost(host string) SSHOptions {
	return func(h *SSH) {
		h.Host = host
	}
}
func OptionSSHPrivateKeyFile(privatekeyfile string) SSHOptions {
	return func(h *SSH) {
		h.PrivateKeyFile = privatekeyfile
	}
}
func OptionSSHUser(user string) SSHOptions {
	return func(h *SSH) {
		h.User = user
	}
}
func OptionSSHPassword(password string) SSHOptions {
	return func(h *SSH) {
		h.Password = password
	}
}

func (m *SSH) Update(opts ...SSHOptions) {
	// Loop through each option
	for _, opt := range opts {
		// Call the option giving the instantiated
		opt(m)
	}
}

func (m *SSH) String() string {
	return fmt.Sprintf("%s,%s", m.Host, m.User)
}

// Create a instance of the docker registry
func CreateSSH(env string, host string, user string, password string) *SSH {
	return &SSH{
		Host:     host,
		User:     user,
		Password: password,
	}
}

func CreateSSHOptions(env string, opts ...SSHOptions) *SSH {
	obj := &SSH{}
	obj.Update(opts...)
	return obj
}
