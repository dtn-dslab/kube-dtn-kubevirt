package kubedtn

import (
	"context"

	common "dslab.sjtu/kube-dtn-sidecar/common"
	"github.com/vishvananda/netlink"
	"kubevirt.io/client-go/log"
)

func (m *KubeDTN) AddVMIntf(intf netlink.Link) error {
	intfName := intf.Attrs().Name
	intfMac := intf.Attrs().HardwareAddr.String()

	newVmInterface := common.VMInterface{
		CNIInterface: common.NetworkInterface{
			IntfName: intfName,
			Mac:      intfMac,
		},
		TapInterface: common.NetworkInterface{
			IntfName: "tap" + intfName,
			Mac:      common.RandomMac(),
		},
		VirtInterface: common.NetworkInterface{
			IntfName: intfName,
			Mac:      intfMac,
		},
	}

	// Modify intf created by cni to avoid MAC conflicts
	newVmInterface.CNIInterface.Mac = common.RandomMac()
	common.SetInterfaceMac(intfName, newVmInterface.CNIInterface.Mac)

	// Connect CNI intf to the bridge
	if err := common.AddInterfaceToDefaultBridge(m.ovsClient, intfName); err != nil {
		return err
	}
	log.Log.Infof("Added CNI interface to bridge: %s", intfName)

	// Create Tap and add it to the bridge
	if err := common.AddTapInterface(newVmInterface.TapInterface.IntfName, newVmInterface.TapInterface.Mac); err != nil {
		common.DeleteInterfaceFromDefaultBridge(m.ovsClient, intfName)
		return err
	}

	log.Log.Infof("Added tap interface: %s", newVmInterface.TapInterface.IntfName)

	if err := common.AddInterfaceToDefaultBridge(m.ovsClient, newVmInterface.TapInterface.IntfName); err != nil {
		common.DeleteInterfaceFromDefaultBridge(m.ovsClient, intfName)
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
	common.RemoveInterfaceAddress(intfName)

	m.vmInterfaces[intfName] = newVmInterface

	return nil
}

func (m *KubeDTN) DeleteVMIntf(ctx context.Context, intf netlink.Link) error {
	return nil
}
