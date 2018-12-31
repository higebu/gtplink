package main

import (
	"flag"
	"net"
	"os"
	"os/signal"

	"github.com/vishvananda/netlink"
)

var (
	v0addr      = flag.String("v0-addr", "0.0.0.0:3386", "v0 addr")
	v1addr      = flag.String("v1-addr", "0.0.0.0:2152", "v1 addr")
	name        = flag.String("name", "gtp0", "gtp dev name")
	mtu         = flag.Int("mtu", 1280, "mtu for gtp dev")
	sgsn        = flag.Bool("sgsn-mode", false, "sgsn mode")
	pdpHashSize = flag.Int("pdp-hash-size", 1024, "pdp hash size")
)

func main() {
	flag.Parse()

	a0, err := net.ResolveUDPAddr("udp", *v0addr)
	if err != nil {
		panic(err)
	}
	a1, err := net.ResolveUDPAddr("udp", *v1addr)
	if err != nil {
		panic(err)
	}
	conn0, err := net.ListenUDP("udp", a0)
	if err != nil {
		panic(err)
	}
	defer conn0.Close()
	conn1, err := net.ListenUDP("udp", a1)
	if err != nil {
		panic(err)
	}
	defer conn1.Close()
	fd0, _ := conn0.File()
	fd1, _ := conn1.File()
	gtp := &netlink.GTP{
		LinkAttrs: netlink.LinkAttrs{
			Name: *name,
		},
		FD0:         int(fd0.Fd()),
		FD1:         int(fd1.Fd()),
		PDPHashsize: *pdpHashSize,
	}
	if *sgsn {
		gtp.Role = 1
	}
	if err := netlink.LinkAdd(gtp); err != nil {
		panic(err)
	}
	if err := netlink.LinkSetUp(gtp); err != nil {
		panic(err)
	}
	if err := netlink.LinkSetMTU(gtp, *mtu); err != nil {
		panic(err)
	}
	defer netlink.LinkDel(gtp)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
}
