package action_govc

import (
	"fmt"
)

// Vcenter Target struct
type VCenter struct {
	Host         string `json:"host" yaml:"host" flag:"host h" desc:"Url to vcenter"`
	User         string `json:"user" yaml:"user" flag:"user u" desc:"The user name for the vcenter"`
	Password     string `json:"password" yaml:"password" flag:"password p" desc:"The password for the vcenter"`
	IgnoreSSL    bool   `json:"ignore_ssl" yaml:"ignore_ssl" flag:"ignore_ssl i" desc:"Ignore SSL"`
	DataStore    string `json:"data_store" yaml:"data_store" flag:"data_store s" desc:"Data Store"`
	DataCenter   string `json:"data_center" yaml:"data_center" flag:"data_center c" desc:"Data Center"`
	Network      string `json:"network" yaml:"network" flag:"network n" desc:"Network"`
	ResourcePool string `json:"resource_pool" yaml:"resource_pool"  flag:"resource_pool r" desc:"Resource Pool"`
}

// Load the config
type VCenterOption func(*VCenter)

func OptionVCenterHost(host string) VCenterOption {
	return func(h *VCenter) {
		h.Host = host
	}
}
func OptionVCenterUser(user string) VCenterOption {
	return func(h *VCenter) {
		h.User = user
	}
}
func OptionVCenterPassword(password string) VCenterOption {
	return func(h *VCenter) {
		h.Password = password
	}
}
func OptionVCenterIgnoreSSL(ignore_ssl bool) VCenterOption {
	return func(h *VCenter) {
		h.IgnoreSSL = ignore_ssl
	}
}

func OptionVCenterDataCenter(data_center string) VCenterOption {
	return func(h *VCenter) {
		h.DataCenter = data_center
	}
}
func OptionVCenterDataStore(data_store string) VCenterOption {
	return func(h *VCenter) {
		h.DataStore = data_store
	}
}
func OptionVCenterNetwork(network string) VCenterOption {
	return func(h *VCenter) {
		h.Network = network
	}
}
func OptionVCenterResourcePool(resource_pool string) VCenterOption {
	return func(h *VCenter) {
		h.ResourcePool = resource_pool
	}
}

func (m *VCenter) Update(opts ...VCenterOption) {
	// Loop through each option
	for _, opt := range opts {
		opt(m)
	}
}

func (m *VCenter) String() string {
	return fmt.Sprintf("%s,ignore_ssl=%v", m.Host, m.IgnoreSSL)
}

func CreateVCenterOptions(env string, opts ...VCenterOption) *VCenter {
	obj := &VCenter{}
	obj.Update(opts...)
	return obj
}

func CreateVCenter(env string, host string, user string, password string, ignore_ssl bool, data_store string, data_center string, network string, resource_pool string) *VCenter {
	return &VCenter{
		Host:         host,
		User:         user,
		Password:     password,
		IgnoreSSL:    ignore_ssl,
		DataStore:    data_store,
		DataCenter:   data_center,
		Network:      network,
		ResourcePool: resource_pool,
	}
}
