package common

import (
	"encoding/xml"
	"fmt"
	"net/url"
	"time"

	libvirt "github.com/digitalocean/go-libvirt"
)

type LibvirtMAC struct {
	XMLName xml.Name `xml:"mac"`
	Address string   `xml:"address,attr"`
}

type LibvirtTarget struct {
	XMLName xml.Name `xml:"target"`
	Dev     string   `xml:"dev,attr"`
	Managed string   `xml:"managed,attr"`
}

type LibvirtModel struct {
	XMLName xml.Name `xml:"model"`
	Type    string   `xml:"type,attr"`
}

type LibvirtMTU struct {
	XMLName xml.Name `xml:"mtu"`
	Size    string   `xml:"size,attr"`
}

type LibvirtROM struct {
	XMLName xml.Name `xml:"rom"`
	Enabled string   `xml:"enabled,attr"`
}

type LibvirtAlias struct {
	XMLName xml.Name `xml:"alias"`
	Name    string   `xml:"name,attr"`
}

type LibvirtInterface struct {
	XMLName xml.Name `xml:"interface"`
	Type    string   `xml:"type,attr"`
	MAC     LibvirtMAC
	Target  LibvirtTarget
	Model   LibvirtModel
	Alias   LibvirtAlias
	MTU     LibvirtMTU
	ROM     LibvirtROM
}

const (
	maxRetries = 30
)

func ConnectLibvirt() (*libvirt.Libvirt, error) {
	uri, _ := url.Parse("qemu+unix:///session?socket=/var/run/libvirt/virtqemud-sock")
	l, err := libvirt.ConnectToURI(uri)
	if err != nil {
		return nil, err
	}
	return l, nil
}

func ConnectLibvirtBlock() (*libvirt.Libvirt, error) {
	for i := 0; i < maxRetries; i++ {
		l, err := ConnectLibvirt()
		if err == nil {
			return l, nil
		}
		time.Sleep(1 * time.Second)
	}

	return nil, fmt.Errorf("failed to connect to libvirt after %d retries", maxRetries)
}

func GenerateLibvirtXML(vmInterface VMInterface) string {
	libvirtInterface := LibvirtInterface{
		Type: "ethernet",
		MAC: LibvirtMAC{
			Address: vmInterface.VirtInterface.Mac,
		},
		Target: LibvirtTarget{
			Dev:     vmInterface.TapInterface.IntfName,
			Managed: "no",
		},
		Model: LibvirtModel{
			Type: "virtio-non-transitional",
		},
		Alias: LibvirtAlias{
			Name: "ua-" + vmInterface.VirtInterface.IntfName,
		},
		MTU: LibvirtMTU{
			Size: "1500",
		},
		ROM: LibvirtROM{
			Enabled: "no",
		},
	}

	xmlData, err := xml.MarshalIndent(libvirtInterface, "", "\t")
	if err != nil {
		fmt.Println("Error marshaling XML:", err)
		return ""
	}
	return string(xmlData)
}

func AttachDeviceByLink(libvirtClient *libvirt.Libvirt, vmInterface VMInterface) error {
	// Get all domains
	domains, _, err := libvirtClient.ConnectListAllDomains(1, libvirt.ConnectListDomainsActive|libvirt.ConnectListDomainsPersistent)
	if err != nil {
		return err
	}

	if len(domains) == 0 {
		return fmt.Errorf("no KubeVirt domains found")
	}

	xml := GenerateLibvirtXML(vmInterface)

	err = libvirtClient.DomainAttachDevice(domains[0], xml)
	if err != nil {
		return err
	}

	return nil

}

func AttachDeviceByLinkBlock(libvirtClient *libvirt.Libvirt, vmInterface VMInterface) error {
	var err error

	for i := 0; i < maxRetries; i++ {
		err = AttachDeviceByLink(libvirtClient, vmInterface)
		if err == nil {
			return nil
		}
		time.Sleep(1 * time.Second)
	}

	return err
}
