package kubedtn

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"

	"github.com/digitalocean/go-openvswitch/ovs"
	common "github.com/dtn-dslab/kube-dtn-sidecar/common"
	pb "github.com/dtn-dslab/kube-dtn-sidecar/proto/v1"
	"github.com/go-redis/redis/v8"
	dhcp "github.com/krolaw/dhcp4"
	dhcpconn "github.com/krolaw/dhcp4/conn"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"kubevirt.io/client-go/log"
)

type KubeDTN struct {
	pb.UnimplementedVMSidecarServer
	s            *grpc.Server
	lis          net.Listener
	redisClient  *redis.Client
	ctx          context.Context
	ovsClient    *ovs.Client
	dhcpServer   *common.DHCPServer
	name         string
	vmInterfaces map[string]common.VMInterface
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

type Config struct {
	Port     int
	GRPCOpts []grpc.ServerOption
}

func New(cfg Config) (*KubeDTN, error) {
	hostname := os.Getenv("HOSTNAME")
	if hostname == "" {
		return nil, fmt.Errorf("failed to get hostname")
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	redisClient := common.GenerateRedisClient()

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
		s:            grpc.NewServer(cfg.GRPCOpts...),
		lis:          lis,
		ctx:          ctx,
		redisClient:  redisClient,
		ovsClient:    ovsClient,
		dhcpServer:   dhcpServer,
		name:         hostname,
		vmInterfaces: make(map[string]common.VMInterface),
	}

	pb.RegisterVMSidecarServer(m.s, m)
	reflection.Register(m.s)
	return m, nil
}

func (m *KubeDTN) AddLink(link common.Link) error {
	newVmInterface := common.VMInterface{
		CNIInterface: link,
		TapInterface: common.NetworkInterface{
			IntfName: "tap" + link.LocalIntf,
			Mac:      common.RandomMac(),
		},
		VirtInterface: common.NetworkInterface{
			IntfName: link.LocalIntf,
			Mac:      link.LocalMAC,
		},
	}

	newVmInterface.CNIInterface.LocalMAC = common.RandomMac()
	common.SetInterfaceMac(link.LocalIntf, newVmInterface.CNIInterface.LocalMAC)

	if err := common.AddInterfaceToDefaultBridge(m.ovsClient, link.LocalIntf); err != nil {
		return err
	}

	log.Log.Infof("Added CNI interface to bridge: %s", link.LocalIntf)

	// Create Tap and add it to the bridge
	if err := common.AddTapInterface(newVmInterface.TapInterface.IntfName, newVmInterface.TapInterface.Mac); err != nil {
		common.DeleteInterfaceFromDefaultBridge(m.ovsClient, link.LocalIntf)
		return err
	}

	log.Log.Infof("Added tap interface: %s", newVmInterface.TapInterface.IntfName)

	if err := common.AddInterfaceToDefaultBridge(m.ovsClient, newVmInterface.TapInterface.IntfName); err != nil {
		common.DeleteInterfaceFromDefaultBridge(m.ovsClient, link.LocalIntf)
		return err
	}

	log.Log.Infof("Added tap interface to bridge: %s", newVmInterface.TapInterface.IntfName)

	libvirtClient, err := common.ConnectLibvirtBlock()
	if err != nil {
		return err
	}

	if err := common.AttachDeviceByLinkBlock(libvirtClient, newVmInterface); err != nil {
		return err
	}

	log.Log.Infof("Attached device to libvirt: %s", newVmInterface.VirtInterface.IntfName)

	if err := common.AddFlowToDefaultBridge(m.ovsClient, newVmInterface); err != nil {
		return err
	}

	log.Log.Infof("Added flow rules: %s", newVmInterface.VirtInterface.IntfName)

	// Remove IP address from the CNI interface
	common.RemoveInterfaceAddress(link.LocalIntf)

	m.vmInterfaces[link.LocalIntf] = newVmInterface
	// if err := m.dhcpServer.AddLease(link); err != nil {
	// 	log.Log.Warningf("Failed to add lease")
	// }

	return nil
}

func (m *KubeDTN) InitStatus() error {
	topoSpec, err := common.GetTopoSpecFromRedis(m.ctx, m.redisClient, m.name)
	if err != nil {
		return err
	}
	for _, link := range topoSpec.Links {
		if !common.IsInterfaceExist(link.LocalIntf) {
			continue
		}

		m.AddLink(link)
	}

	return nil
}

func (m *KubeDTN) SetupStatus() error {
	mainInterfaceAddr, err := common.GetMainInterfaceAddress()
	if err != nil {
		return fmt.Errorf("failed to get main interface address: %s", err)
	}

	topoStatus, err := common.GetTopoStatusFromRedis(m.ctx, m.redisClient, m.name)
	if err != nil {
		return fmt.Errorf("failed to get topology status from redis: %s", err)
	}

	topoStatus.PodIP = mainInterfaceAddr

	statusJSON, err := json.Marshal(topoStatus)
	if err != nil {
		return fmt.Errorf("failed to marshal topology status: %s", err)
	}

	if err := m.redisClient.Set(m.ctx, "cni_"+m.name+"_status", statusJSON, 0).Err(); err != nil {
		return fmt.Errorf("failed to set topology status to redis: %s", err)
	}
	return nil
}

func (m *KubeDTN) Destroy() error {
	m.ovsClient.VSwitch.DeleteBridge(common.DefaultBridgeName)
	for _, vmInterface := range m.vmInterfaces {
		common.DeleteNetworkInterface(vmInterface.TapInterface.IntfName)
	}

	common.DeleteNetworkInterface(m.dhcpServer.DhcpInterface.IntfName)
	common.DeleteNetworkInterface(m.dhcpServer.DhcpBridgeInterface.IntfName)

	return nil
}

func (m *KubeDTN) Serve() error {
	if err := m.InitStatus(); err != nil {
		return err
	}

	if err := m.SetupStatus(); err != nil {
		return err
	}

	return m.s.Serve(m.lis)
}

func (m *KubeDTN) GracefulStop() {
	m.s.GracefulStop()
	m.Destroy()
}
