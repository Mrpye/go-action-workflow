package action_govc

import (
	"log"
	"os"
	"strconv"
	"strings"

	_ "github.com/vmware/govmomi/govc/about"
	"github.com/vmware/govmomi/govc/cli"
	_ "github.com/vmware/govmomi/govc/cluster"
	_ "github.com/vmware/govmomi/govc/cluster/group"
	_ "github.com/vmware/govmomi/govc/cluster/module"
	_ "github.com/vmware/govmomi/govc/cluster/override"
	_ "github.com/vmware/govmomi/govc/cluster/rule"
	_ "github.com/vmware/govmomi/govc/datacenter"
	_ "github.com/vmware/govmomi/govc/datastore"
	_ "github.com/vmware/govmomi/govc/datastore/cluster"
	_ "github.com/vmware/govmomi/govc/datastore/disk"
	_ "github.com/vmware/govmomi/govc/datastore/maintenance"
	_ "github.com/vmware/govmomi/govc/datastore/vsan"
	_ "github.com/vmware/govmomi/govc/device"
	_ "github.com/vmware/govmomi/govc/device/cdrom"
	_ "github.com/vmware/govmomi/govc/device/clock"
	_ "github.com/vmware/govmomi/govc/device/floppy"
	_ "github.com/vmware/govmomi/govc/device/pci"
	_ "github.com/vmware/govmomi/govc/device/scsi"
	_ "github.com/vmware/govmomi/govc/device/serial"
	_ "github.com/vmware/govmomi/govc/device/usb"
	_ "github.com/vmware/govmomi/govc/disk"
	_ "github.com/vmware/govmomi/govc/disk/snapshot"
	_ "github.com/vmware/govmomi/govc/dvs"
	_ "github.com/vmware/govmomi/govc/dvs/portgroup"
	_ "github.com/vmware/govmomi/govc/env"
	_ "github.com/vmware/govmomi/govc/events"
	_ "github.com/vmware/govmomi/govc/export"
	_ "github.com/vmware/govmomi/govc/extension"
	_ "github.com/vmware/govmomi/govc/fields"
	_ "github.com/vmware/govmomi/govc/folder"
	_ "github.com/vmware/govmomi/govc/host"
	_ "github.com/vmware/govmomi/govc/host/account"
	_ "github.com/vmware/govmomi/govc/host/autostart"
	_ "github.com/vmware/govmomi/govc/host/cert"
	_ "github.com/vmware/govmomi/govc/host/date"
	_ "github.com/vmware/govmomi/govc/host/esxcli"
	_ "github.com/vmware/govmomi/govc/host/firewall"
	_ "github.com/vmware/govmomi/govc/host/maintenance"
	_ "github.com/vmware/govmomi/govc/host/option"
	_ "github.com/vmware/govmomi/govc/host/portgroup"
	_ "github.com/vmware/govmomi/govc/host/service"
	_ "github.com/vmware/govmomi/govc/host/storage"
	_ "github.com/vmware/govmomi/govc/host/vnic"
	_ "github.com/vmware/govmomi/govc/host/vswitch"
	_ "github.com/vmware/govmomi/govc/importx"
	_ "github.com/vmware/govmomi/govc/library"
	_ "github.com/vmware/govmomi/govc/library/policy"
	_ "github.com/vmware/govmomi/govc/library/session"
	_ "github.com/vmware/govmomi/govc/library/subscriber"
	_ "github.com/vmware/govmomi/govc/library/trust"
	_ "github.com/vmware/govmomi/govc/license"
	_ "github.com/vmware/govmomi/govc/logs"
	_ "github.com/vmware/govmomi/govc/ls"
	_ "github.com/vmware/govmomi/govc/metric"
	_ "github.com/vmware/govmomi/govc/metric/interval"
	_ "github.com/vmware/govmomi/govc/namespace/cluster"
	_ "github.com/vmware/govmomi/govc/namespace/service"
	_ "github.com/vmware/govmomi/govc/option"
	_ "github.com/vmware/govmomi/govc/permissions"
	_ "github.com/vmware/govmomi/govc/pool"
	_ "github.com/vmware/govmomi/govc/role"
	_ "github.com/vmware/govmomi/govc/session"
	_ "github.com/vmware/govmomi/govc/sso/group"
	_ "github.com/vmware/govmomi/govc/sso/idp"
	_ "github.com/vmware/govmomi/govc/sso/lpp"
	_ "github.com/vmware/govmomi/govc/sso/service"
	_ "github.com/vmware/govmomi/govc/sso/user"
	_ "github.com/vmware/govmomi/govc/storage/policy"
	_ "github.com/vmware/govmomi/govc/tags"
	_ "github.com/vmware/govmomi/govc/tags/association"
	_ "github.com/vmware/govmomi/govc/tags/category"
	_ "github.com/vmware/govmomi/govc/task"
	_ "github.com/vmware/govmomi/govc/vapp"
	_ "github.com/vmware/govmomi/govc/vcsa/access/consolecli"
	_ "github.com/vmware/govmomi/govc/vcsa/access/dcui"
	_ "github.com/vmware/govmomi/govc/vcsa/access/shell"
	_ "github.com/vmware/govmomi/govc/vcsa/access/ssh"
	_ "github.com/vmware/govmomi/govc/vcsa/log"
	_ "github.com/vmware/govmomi/govc/vcsa/proxy"
	_ "github.com/vmware/govmomi/govc/vcsa/shutdown"
	_ "github.com/vmware/govmomi/govc/version"
	_ "github.com/vmware/govmomi/govc/vm"
	_ "github.com/vmware/govmomi/govc/vm/disk"
	_ "github.com/vmware/govmomi/govc/vm/guest"
	_ "github.com/vmware/govmomi/govc/vm/network"
	_ "github.com/vmware/govmomi/govc/vm/option"
	_ "github.com/vmware/govmomi/govc/vm/rdm"
	_ "github.com/vmware/govmomi/govc/vm/snapshot"
	_ "github.com/vmware/govmomi/govc/volume"
	_ "github.com/vmware/govmomi/govc/vsan"
)

// Test the connection
func (m *VCenter) Test() (bool, string, error) {
	err := m.Run("about", true)
	if err != nil {
		return false, "", nil
	}
	return true, "", nil
}

// Run GOVC commands
func (m *VCenter) Run(command string, silent bool) error {
	quoted := false
	args := strings.FieldsFunc(command, func(r rune) bool {
		if r == '"' {
			quoted = !quoted
		}
		return !quoted && r == ' '
	})

	for x := range args {
		args[x] = strings.ReplaceAll(args[x], "\"", "")
	}
	if !silent {
		for _, x := range args {
			log.Println(x)
		}
	}
	os.Setenv("GOVC_URL", m.Host)
	os.Setenv("GOVC_USERNAME", m.User)
	os.Setenv("GOVC_PASSWORD", m.Password)
	os.Setenv("GOVC_INSECURE", strconv.FormatBool(m.IgnoreSSL))
	os.Setenv("GOVC_DATASTORE", m.DataStore)
	os.Setenv("GOVC_DATACENTER", m.DataCenter)
	os.Setenv("GOVC_NETWORK", m.Network)
	os.Setenv("GOVC_RESOURCE_POOL", m.ResourcePool)

	cli.Run(args)

	return nil
}

/*
// Create a new vcenter client
func (m *VCenter) NewClient(ctx context.Context) (*vim25.Client, error) {
	conn := ""
	//***************************
	//Setup the connection string
	//***************************
	m.User = strings.ReplaceAll(m.User, "@", "%40")
	tpass := strings.ReplaceAll(m.AuthDecrypt()+"", "@", "%40")
	host := m.Host
	host = strings.ReplaceAll(host, "https://", "")
	host = strings.ReplaceAll(host, "http://", "")
	if strings.HasPrefix(m.Host, "https") {
		conn = fmt.Sprintf("https://%s:%s@%s/sdk", m.User, tpass, host)
	} else {
		conn = fmt.Sprintf("http://%s:%s@%s/sdk", m.User, tpass, host)
	}
	u, err := url.Parse(conn)
	if err != nil {
		return nil, err
	}
	//***************************
	// Share govc's session cache
	//***************************
	s := &cache.Session{
		URL:      u,
		Insecure: m.IgnoreSSL,
	}

	//***************************
	// Create a new client
	//***************************
	c := new(vim25.Client)
	err = s.Login(ctx, c, nil)
	if err != nil {
		return nil, err
	}

	return c, nil
}

// Get a list of all the VCenter data centers
func (m *VCenter) GetDataCenters() (bool, []string, error) {
	var results []string
	//********************************
	//Create the connection to vcenter
	//********************************
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client, err := m.NewClient(ctx)
	if err != nil {
		return false, results, err
	}

	//***********
	// Data center
	//***********
	finder := find.NewFinder(client)
	dcs, err := finder.DatacenterList(ctx, "*")
	if err != nil {
		return false, results, err
	}
	for _, dc := range dcs {
		results = append(results, strings.Replace(dc.InventoryPath, "/", "", 1))
	}

	return true, results, nil
}

// Get the VCenter data center
func (m *VCenter) GetDataCenter(data_center string) (bool, *object.Datacenter, error) {

	//********************************
	//Create the connection to vcenter
	//********************************
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client, err := m.NewClient(ctx)
	if err != nil {
		return false, nil, err
	}

	//***********
	// Data center
	//***********
	finder := find.NewFinder(client)
	dcs, err := finder.DatacenterList(ctx, data_center)
	if err != nil {
		return false, nil, err
	}
	if len(dcs) > 0 {
		return true, dcs[0], nil
	}

	return false, nil, nil
}

// Get the VCenter networks
func (m *VCenter) GetNetworks(data_center string) (bool, []string, error) {
	var results []string
	//********************************
	//Create the connection to vcenter
	//********************************
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client, err := m.NewClient(ctx)
	if err != nil {
		return false, results, err
	}

	//***********
	// Data center
	//***********
	finder := find.NewFinder(client)
	//Set the datastore
	if data_center != "" {
		success, dc, _ := m.GetDataCenter(data_center)
		if success {
			finder.SetDatacenter(dc)
		}
	}
	nws, err := finder.NetworkList(ctx, "*")
	if err != nil {
		return false, results, err
	}
	for _, nw := range nws {
		results = append(results, nw.GetInventoryPath())
	}

	return true, results, nil
}

// Get the VCenter Clusters
func (m *VCenter) GetClusters(data_center string) (bool, []string, error) {
	var results []string
	//********************************
	//Create the connection to vcenter
	//********************************
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client, err := m.NewClient(ctx)
	if err != nil {
		return false, results, err
	}

	//**************
	// Data clusters
	//**************
	finder := find.NewFinder(client)
	//Set the datastore
	if data_center != "" {
		success, dc, _ := m.GetDataCenter(data_center)
		if success {
			finder.SetDatacenter(dc)
		}
	}
	clusters, err := finder.ClusterComputeResourceList(ctx, "*")
	if err != nil {
		return false, results, err
	}
	for _, cluster := range clusters {
		results = append(results, cluster.Name())
	}

	return true, results, nil
}

// Get the VCenter data stores
func (m *VCenter) GetDataStores(data_center string) (bool, []string, error) {
	var results []string
	//********************************
	//Create the connection to vcenter
	//********************************
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client, err := m.NewClient(ctx)
	if err != nil {
		return false, results, err
	}

	//***********
	// Data center
	//***********
	finder := find.NewFinder(client)
	//Set the datastore
	if data_center != "" {
		success, dc, _ := m.GetDataCenter(data_center)
		if success {
			finder.SetDatacenter(dc)
		}
	}
	dss, err := finder.DatastoreList(ctx, "*")
	if err != nil {
		return false, results, err
	}
	for _, ds := range dss {
		results = append(results, ds.Name())
	}

	return true, results, nil
}

// Get the VCenter resource pools
func (m *VCenter) GetResourcePools(data_center string) (bool, []string, error) {
	var results []string
	//********************************
	//Create the connection to vcenter
	//********************************
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client, err := m.NewClient(ctx)
	if err != nil {
		return false, results, err
	}

	//***********
	// Data center
	//***********
	finder := find.NewFinder(client)
	//Set the datastore
	if data_center != "" {
		success, dc, _ := m.GetDataCenter(data_center)
		if success {
			finder.SetDatacenter(dc)
		}
	}
	dss, err := finder.ResourcePoolList(ctx, "*")
	if err != nil {
		return false, results, err
	}
	for _, ds := range dss {
		results = append(results, ds.InventoryPath)
	}

	return true, results, nil
}

// get the vcenter folders
func (m *VCenter) GetFolders(data_center string) (bool, []string, error) {
	var results []string
	//********************************
	//Create the connection to vcenter
	//********************************
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client, err := m.NewClient(ctx)
	if err != nil {
		return false, results, err
	}

	//***********
	// Data center
	//***********
	finder := find.NewFinder(client)
	//Set the datastore
	if data_center != "" {
		success, dc, _ := m.GetDataCenter(data_center)
		if success {
			finder.SetDatacenter(dc)
		}
	}
	dss, err := finder.FolderList(ctx, "*")
	if err != nil {
		return false, results, err
	}
	for _, ds := range dss {
		results = append(results, ds.InventoryPath)
	}

	return true, results, nil
}
*/
