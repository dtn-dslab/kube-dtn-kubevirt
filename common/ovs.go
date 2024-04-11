package common

import "github.com/digitalocean/go-openvswitch/ovs"

const (
	DefaultBridgeName = "dtnbr0"
)

func CreateDefaultOVSBridge(c *ovs.Client) error {
	if err := c.VSwitch.AddBridge(DefaultBridgeName); err != nil {
		return err
	}
	return nil
}
