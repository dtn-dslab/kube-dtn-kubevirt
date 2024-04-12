package common

import (
	"net"

	"github.com/vishvananda/netlink"
)

func MakeVeth(intf NetworkInterface, peer NetworkInterface) error {
	veth1HardwareAddr, err := net.ParseMAC(intf.Mac)
	if err != nil {
		return err
	}

	veth2HardwareAddr, err := net.ParseMAC(peer.Mac)
	if err != nil {
		return err
	}

	veth1 := &netlink.Veth{
		LinkAttrs: netlink.LinkAttrs{
			Name:         intf.IntfName,
			HardwareAddr: veth1HardwareAddr,
		},
		PeerName:         peer.IntfName,
		PeerHardwareAddr: veth2HardwareAddr,
	}

	if err := netlink.LinkAdd(veth1); err != nil {
		return err
	}

	return nil
}
