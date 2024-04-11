package main

import (
	"log"
	"os"

	"github.com/containernetworking/cni/pkg/skel"
	"github.com/containernetworking/cni/pkg/version"
)

var logger = log.New(os.Stderr, "", 0)

func cmdAdd(args *skel.CmdArgs) error {
	logger.Println("kubedtn-hack: cmdAdd")
	return nil
}

func cmdCheck(args *skel.CmdArgs) error {
	logger.Println("kubedtn-hack: cmdCheck")
	return nil
}

func cmdDel(args *skel.CmdArgs) error {
	logger.Println("kubedtn-hack: cmdDel")
	return nil
}

func main() {
	skel.PluginMain(cmdAdd, cmdCheck, cmdDel, version.All, "0.1.0")
}
