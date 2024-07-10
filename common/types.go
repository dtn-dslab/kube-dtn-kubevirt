package common

const (
	DefaultPort          string = "51112"
	DefaultMainInterface string = "eth0"
)

type NetworkInterface struct {
	IntfName string
	Mac      string
}

type VMInterface struct {
	CNIInterface  NetworkInterface
	TapInterface  NetworkInterface
	VirtInterface NetworkInterface
}
