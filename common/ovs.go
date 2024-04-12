package common

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/digitalocean/go-openvswitch/ovs"
)

const (
	DefaultBridgeName = "dtnbr0"
)

func GetPortID(bridge, port string) (int, error) {
	// sudo ovs-vsctl get Interface port_name ofport
	cmd := exec.Command("ovs-vsctl", "get", "Interface", port, "ofport")
	output, err := cmd.Output()
	if err != nil {
		return -1, fmt.Errorf("failed to get port %s id on OVS bridge %s: %v", port, bridge, err)
	}
	resultStr := strings.TrimSpace(string(output))
	resultInt, err := strconv.Atoi(resultStr)
	if err != nil {
		return -1, fmt.Errorf("error converting port %s id %s to int: %v", port, resultStr, err)
	}
	return resultInt, nil
}

func CreateDefaultOVSBridge(c *ovs.Client) error {
	if err := c.VSwitch.AddBridge(DefaultBridgeName); err != nil {
		return err
	}
	if err := c.OpenFlow.DelFlows(DefaultBridgeName, nil); err != nil {
		return err
	}
	return nil
}

func AddInterfaceToBridge(c *ovs.Client, bridgeName, ifaceName string) error {
	if err := c.VSwitch.AddPort(bridgeName, ifaceName); err != nil {
		return err
	}
	return nil
}

func AddInterfaceToDefaultBridge(c *ovs.Client, ifaceName string) error {
	if err := AddInterfaceToBridge(c, DefaultBridgeName, ifaceName); err != nil {
		return err
	}
	return nil
}

func DeleteInterfaceFromBridge(c *ovs.Client, bridgeName, ifaceName string) error {
	if err := c.VSwitch.DeletePort(bridgeName, ifaceName); err != nil {
		return err
	}
	return nil
}

func DeleteInterfaceFromDefaultBridge(c *ovs.Client, ifaceName string) error {
	if err := DeleteInterfaceFromBridge(c, DefaultBridgeName, ifaceName); err != nil {
		return err
	}
	return nil
}

func AddFlowToBridge(c *ovs.Client, bridgeName string, vmInterface VMInterface) error {
	cniPortID, err := GetPortID(bridgeName, vmInterface.CNIInterface.LocalIntf)
	if err != nil {
		return err
	}
	tapPortID, err := GetPortID(bridgeName, vmInterface.TapInterface.IntfName)
	if err != nil {
		return err
	}

	flowToTap := ovs.Flow{
		InPort: cniPortID,
		Actions: []ovs.Action{
			ovs.Output(tapPortID),
		},
	}

	flowToCNI := ovs.Flow{
		InPort: tapPortID,
		Actions: []ovs.Action{
			ovs.SetField(vmInterface.CNIInterface.LocalMAC, "eth_src"),
			ovs.Output(cniPortID),
		},
	}

	if err := c.OpenFlow.AddFlow(bridgeName, &flowToTap); err != nil {
		return err
	}

	if err := c.OpenFlow.AddFlow(bridgeName, &flowToCNI); err != nil {
		return err
	}

	return nil
}

func AddFlowToDefaultBridge(c *ovs.Client, vmInterface VMInterface) error {
	if err := AddFlowToBridge(c, DefaultBridgeName, vmInterface); err != nil {
		return err
	}
	return nil
}
