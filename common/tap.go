package common

import (
	"crypto/rand"
	"fmt"
	"os/exec"
)

const (
	local     = 0b10
	multicast = 0b1
)

func RandomMac() string {
	buf := make([]byte, 6)
	_, err := rand.Read(buf)
	if err != nil {
		fmt.Println("error:", err)
		return ""
	}
	// clear multicast bit (&^), ensure local bit (|)
	buf[0] = buf[0]&^multicast | local
	return fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x", buf[0], buf[1], buf[2], buf[3], buf[4], buf[5])
}

func AddTapInterface(linkName, mac string) error {
	cmd := exec.Command("sudo", "ip", "tuntap", "add", "mode", "tap", linkName)
	if err := cmd.Run(); err != nil {
		// return fmt.Errorf("failed to add a new tap interface: %s", err)
	}

	cmd = exec.Command("sudo", "ip", "link", "set", "dev", linkName, "address", mac)
	if err := cmd.Run(); err != nil {
		// return fmt.Errorf("failed to set mac on tap interface: %s", err)
	}

	cmd = exec.Command("sudo", "ip", "link", "set", linkName, "up")
	if err := cmd.Run(); err != nil {
		// return fmt.Errorf("failed to set tap interface up: %s", err)
	}

	return nil
}
