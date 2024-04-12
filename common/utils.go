package common

import (
	"fmt"
	"net"

	"github.com/vishvananda/netlink"
)

func Map[T, U any](ts []T, f func(T) U) []U {
	us := make([]U, len(ts))
	for i := range ts {
		us[i] = f(ts[i])
	}
	return us
}

func GetMainInterfaceAddress() (string, error) {
	mainInterface, err := net.InterfaceByName(DefaultMainInterface)
	if err != nil {
		return "", fmt.Errorf("failed to get main interface: %s", err)
	}
	mainInterfaceAddrs, err := mainInterface.Addrs()
	if err != nil {
		return "", fmt.Errorf("failed to get main interface addresses: %s", err)
	}
	var mainInterfaceAddr net.IP
	for _, addr := range mainInterfaceAddrs {
		if mainInterfaceAddr = addr.(*net.IPNet).IP.To4(); mainInterfaceAddr != nil {
			break
		}
	}

	if mainInterfaceAddr == nil {
		return "", fmt.Errorf("main interface address does not have a IPv4 address")
	}

	return mainInterfaceAddr.String(), nil
}

func IsInterfaceExist(name string) bool {
	_, err := net.InterfaceByName(name)
	return err == nil
}

func SetInterfaceMac(name, mac string) error {
	link, err := netlink.LinkByName(name)
	if err != nil {
		return err
	}

	hwAddr, err := net.ParseMAC(mac)
	if err != nil {
		return err
	}

	return netlink.LinkSetHardwareAddr(link, hwAddr)
}

func RemoveInterfaceAddress(name string) error {
	link, err := netlink.LinkByName(name)
	if err != nil {
		return err
	}
	addrList, err := netlink.AddrList(link, netlink.FAMILY_V4)
	if err != nil {
		return err
	}

	for _, addr := range addrList {
		if err := netlink.AddrDel(link, &addr); err != nil {
			// return err
		}
	}

	return nil
}

func DeleteNetworkInterface(name string) error {
	link, err := netlink.LinkByName(name)
	if err != nil {
		return nil
	}

	netlink.LinkDel(link)

	return nil
}
