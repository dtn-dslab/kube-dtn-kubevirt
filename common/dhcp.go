package common

import (
	"fmt"
	"net"
	"sync"
	"time"

	dhcp "github.com/krolaw/dhcp4"
)

const (
	DhcpInterfaceDefaultName   = "dhcp"
	DhcpBrInterfaceDefaultName = "dhcpbr"
	DefaultLeaseDuration       = time.Minute
)

type Lease struct {
	HwAddr string
	IPAddr net.IP
}

type DHCPHandler struct {
	Ip         net.IP
	Options    dhcp.Options
	Leases     []Lease
	LeaseMutex sync.RWMutex
}

type DHCPServer struct {
	DhcpHandler         *DHCPHandler
	DhcpInterface       NetworkInterface
	DhcpBridgeInterface NetworkInterface
}

func (h *DHCPServer) AddLease(link Link) error {
	netIP := net.ParseIP(link.LocalIP)
	if netIP != nil {
		return fmt.Errorf("failed to parse IP from lease")
	}

	lease := Lease{
		HwAddr: link.LocalMAC,
		IPAddr: netIP,
	}

	h.DhcpHandler.LeaseMutex.Lock()
	h.DhcpHandler.Leases = append(h.DhcpHandler.Leases, lease)
	h.DhcpHandler.LeaseMutex.Unlock()

	return nil
}

func (h *DHCPHandler) ServeDHCP(p dhcp.Packet, msgType dhcp.MessageType, options dhcp.Options) (d dhcp.Packet) {
	switch msgType {

	case dhcp.Discover:
		hwAddr := p.CHAddr().String()
		var properLease Lease
		h.LeaseMutex.RLock()
		for _, lease := range h.Leases {
			if lease.HwAddr == hwAddr {
				properLease = lease
				goto reply
			}
		}
		h.LeaseMutex.Unlock()
		return
	reply:
		return dhcp.ReplyPacket(p, dhcp.Offer, h.Ip, properLease.IPAddr, DefaultLeaseDuration,
			h.Options.SelectOrderOrAll(options[dhcp.OptionParameterRequestList]))

	case dhcp.Request:
		if server, ok := options[dhcp.OptionServerIdentifier]; ok && !net.IP(server).Equal(h.Ip) {
			return nil // Message not for this dhcp server
		}
		reqIP := net.IP(options[dhcp.OptionRequestedIPAddress])
		if reqIP == nil {
			reqIP = net.IP(p.CIAddr())
		}

		if len(reqIP) == 4 && !reqIP.Equal(net.IPv4zero) {
			return dhcp.ReplyPacket(p, dhcp.ACK, h.Ip, reqIP, DefaultLeaseDuration,
				h.Options.SelectOrderOrAll(options[dhcp.OptionParameterRequestList]))
		}

		return dhcp.ReplyPacket(p, dhcp.NAK, h.Ip, nil, 0, nil)

	case dhcp.Release, dhcp.Decline:

	}
	return nil
}
