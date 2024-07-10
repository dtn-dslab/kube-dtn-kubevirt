package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"os"

	"github.com/spf13/pflag"

	vmSchema "kubevirt.io/api/core/v1"

	"kubevirt.io/kubevirt/pkg/virt-launcher/virtwrap/api"
)

var logger *log.Logger

func onDefineDomain(vmiJSON, domainXML []byte) (string, error) {

	vmiSpec := vmSchema.VirtualMachineInstance{}
	if err := json.Unmarshal(vmiJSON, &vmiSpec); err != nil {
		return "", fmt.Errorf("failed to unmarshal given VMI spec: %s %s", err, string(vmiJSON))
	}

	domainSpec := api.DomainSpec{}
	if err := xml.Unmarshal(domainXML, &domainSpec); err != nil {
		return "", fmt.Errorf("failed to unmarshal given Domain spec: %s %s", err, string(domainXML))
	}

	newDomainXML, err := xml.Marshal(domainSpec)
	if err != nil {
		return "", fmt.Errorf("failed to marshal new Domain spec: %s %+v", err, domainSpec)
	}

	return string(newDomainXML), nil
}

func main() {
	var vmiJSON, domainXML string
	pflag.StringVar(&vmiJSON, "vmi", "", "VMI to change in JSON format")
	pflag.StringVar(&domainXML, "domain", "", "Domain spec in XML format")
	pflag.Parse()

	// make log file
	logFile, err := os.OpenFile("./kube-dtn-sidecar.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	logger = log.New(logFile, "kubedtn", log.Ldate)

	logger.Printf("Starting kubedtn sidecar...")

	if vmiJSON == "" || domainXML == "" {
		logger.Printf("Bad input vmi=%d, domain=%d", len(vmiJSON), len(domainXML))
		os.Exit(1)
	}

	domainXML, err = onDefineDomain([]byte(vmiJSON), []byte(domainXML))
	if err != nil {
		logger.Fatalf("onDefineDomain failed: %s", err)
		panic(err)
	}
	fmt.Println(domainXML)
}
