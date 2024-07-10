package kubedtn

import (
	"context"
	"fmt"
	"net"
	"os"

	common "dslab.sjtu/kube-dtn-sidecar/common"
	"github.com/digitalocean/go-openvswitch/ovs"
	dhcp "github.com/krolaw/dhcp4"
	dhcpconn "github.com/krolaw/dhcp4/conn"
)

type KubeDTN struct {
	ctx          context.Context
	ovsClient    *ovs.Client
	dhcpServer   *common.DHCPServer
	name         string
	vmInterfaces map[string]common.VMInterface
	stopChan     chan struct{}
}

func (m *KubeDTN) DHCPServe() error {
	if err := common.MakeVeth(m.dhcpServer.DhcpInterface, m.dhcpServer.DhcpBridgeInterface); err != nil {
		return err
	}

	if err := common.AddInterfaceToDefaultBridge(m.ovsClient, m.dhcpServer.DhcpBridgeInterface.IntfName); err != nil {
		return err
	}

	if err := common.AddDHCPFlowFromPortToDefaultBridge(m.ovsClient, m.dhcpServer.DhcpBridgeInterface.IntfName); err != nil {
		return err
	}

	conn, err := dhcpconn.NewUDP4BoundListener(m.dhcpServer.DhcpInterface.IntfName, ":67")
	if err != nil {
		return err
	}
	return dhcp.Serve(conn, m.dhcpServer.DhcpHandler)
}

func New() (*KubeDTN, error) {
	hostname := os.Getenv("HOSTNAME")
	if hostname == "" {
		return nil, fmt.Errorf("failed to get hostname")
	}

	ctx := context.Background()

	ovsClient := ovs.New()
	if err := common.CreateDefaultOVSBridge(ovsClient); err != nil {
		return nil, fmt.Errorf("failed to create default OVS bridge: %s", err)
	}

	serverIP := net.IP{172, 30, 0, 1}
	dhcpHandler := &common.DHCPHandler{
		Ip:      serverIP,
		Leases:  make([]common.Lease, 10),
		Options: dhcp.Options{
			// dhcp.OptionSubnetMask:       []byte{255, 255, 240, 0},
			// dhcp.OptionRouter:           []byte(serverIP), // Presuming Server is also your router
			// dhcp.OptionDomainNameServer: []byte(serverIP), // Presuming Server is also your DNS server
		},
	}

	dhcpServer := &common.DHCPServer{
		DhcpHandler: dhcpHandler,
		DhcpInterface: common.NetworkInterface{
			IntfName: common.DhcpInterfaceDefaultName,
			Mac:      common.RandomMac(),
		},
		DhcpBridgeInterface: common.NetworkInterface{
			IntfName: common.DhcpBrInterfaceDefaultName,
			Mac:      common.RandomMac(),
		},
	}

	m := &KubeDTN{
		ctx:          ctx,
		ovsClient:    ovsClient,
		dhcpServer:   dhcpServer,
		name:         hostname,
		vmInterfaces: make(map[string]common.VMInterface),
		stopChan:     make(chan struct{}),
	}

	return m, nil
}

func (m *KubeDTN) Destroy() error {
	m.ovsClient.VSwitch.DeleteBridge(common.DefaultBridgeName)

	// Clean tap interfaces
	for _, vmInterface := range m.vmInterfaces {
		common.DeleteNetworkInterface(vmInterface.TapInterface.IntfName)
	}

	common.DeleteNetworkInterface(m.dhcpServer.DhcpInterface.IntfName)
	common.DeleteNetworkInterface(m.dhcpServer.DhcpBridgeInterface.IntfName)

	return nil
}

func (m *KubeDTN) GracefulStop() {
	close(m.stopChan)
	m.Destroy()
}
