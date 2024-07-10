package kubedtn

import (
	"context"
	"log"
	"strings"
	"syscall"

	common "dslab.sjtu/kube-dtn-sidecar/common"
	"github.com/vishvananda/netlink"
	health "google.golang.org/grpc/health/grpc_health_v1"
)

func (m *KubeDTN) Check(context.Context, *health.HealthCheckRequest) (*health.HealthCheckResponse, error) {
	return &health.HealthCheckResponse{Status: health.HealthCheckResponse_SERVING}, nil
}

func (m *KubeDTN) Watch(*health.HealthCheckRequest, health.Health_WatchServer) error {
	return nil
}

func IsIntfManaged(intfName string) bool {
	// skip loopback
	if intfName == "lo" {
		return false
	}
	// skip default intf
	if intfName == "eth0" ||
		intfName == "eth0-nic" ||
		intfName == "tap0" ||
		intfName == "k6t-eth0" {
		return false
	}
	// skip ovs
	if intfName == "ovs-system" ||
		intfName == common.DefaultBridgeName {
		return false
	}

	// skip related
	if strings.HasPrefix(intfName, "tap") {
		return false
	}

	return true
}

func (m *KubeDTN) InitNetworkConfigure() error {
	intfs, err := netlink.LinkList()
	if err != nil {
		return err
	}

	for _, intf := range intfs {
		if IsIntfManaged(intf.Attrs().Name) {
			m.AddVMIntf(intf)
		}
	}

	return nil
}

func (m *KubeDTN) ServeNetworkConfigure() error {

	if err := m.InitNetworkConfigure(); err != nil {
		log.Fatalf("failed to init network configure: %v", err)
	}

	done := make(chan struct{})
	updates := make(chan netlink.LinkUpdate)
	err := netlink.LinkSubscribe(updates, done)
	if err != nil {
		log.Fatalf("failed to subscribe to intf updates: %v", err)
	}

	log.Println("Listening for intf updates")

	for {
		select {
		case update := <-updates:
			switch update.Header.Type {
			case syscall.RTM_NEWLINK:
				log.Printf("New intf: %v\n", update.Link)
				if IsIntfManaged(update.Link.Attrs().Name) {
					m.AddVMIntf(update.Link)
				}
			case syscall.RTM_DELLINK:
				log.Printf("Deleted intf: %v\n", update.Link)
			case syscall.RTM_SETLINK:
				log.Printf("Updated intf: %v\n", update.Link)
			}
		case <-m.stopChan:
			log.Println("Stopping intf updates")
			close(done)
			return nil
		}
	}

}
