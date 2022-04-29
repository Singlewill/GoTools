package main

import (
	"fmt"
	"syscall"

	"edr_server/drv_comm"
	"edr_server/netlink"
)

type NetlinkSocket struct {
	fd       int
	protocol int
	groups   uint32
}

func main() {
	nl_socket, err := netlink.NewNetlinkSocket(syscall.NETLINK_GENERIC, 0)
	if err != nil {
		fmt.Println("NewNetlinkSocket Failed")
		return
	}
	defer nl_socket.Close()

	family, err := nl_socket.GetFamily("EdrFamily")
	if err != nil {
		fmt.Println("GetFamily EDRFamily Failed : ", err)
		return
	}

	fmt.Println(family.ID)
	fmt.Println(family.Version)
	fmt.Println(family.Name)

	go drv_comm.HeartBeat(nl_socket, family.ID)
	for {
		drv_comm.DrvRecv(nl_socket)
	}

}
