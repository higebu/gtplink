package main

import (
	"flag"
	"net"
	"os"
	"os/signal"

	"github.com/vishvananda/netlink"
)

var (
	v0addr = flag.String("v0addr", "0.0.0.0:3386", "v0 addr")
	v1addr = flag.String("v1addr", "0.0.0.0:2152", "v1 addr")
	name   = flag.String("name", "gtp0", "gtp dev name")
	mtu    = flag.Int("mtu", 1280, "mtu for gtp dev")
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
	conn1, err := net.ListenUDP("udp", a0)
	if err != nil {
		panic(err)
	}
	defer conn1.Close()
	conn2, err := net.ListenUDP("udp", a1)
	if err != nil {
		panic(err)
	}
	defer conn2.Close()
	fd1, _ := conn1.File()
	fd2, _ := conn2.File()
	gtp := &netlink.GTP{
		LinkAttrs: netlink.LinkAttrs{
			Name: *name,
		},
		FD0: int(fd1.Fd()),
		FD1: int(fd2.Fd()),
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
