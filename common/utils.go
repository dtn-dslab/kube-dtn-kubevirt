package common

import (
	"fmt"
	"net"
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
