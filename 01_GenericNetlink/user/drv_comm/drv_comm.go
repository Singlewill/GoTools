package drv_comm

import (
	"edr_server/netlink"
	"fmt"
	"syscall"
	"time"
)

func HeartBeat(ns *netlink.NetlinkSocket, id uint16) {
	genl_ctrl_attr := &netlink.Attribute{
		NlAttr: syscall.NlAttr{
			//这里需要额外+1，表示'\0'
			Len:  uint16(netlink.NlaHeaderLen + len(SYNC_STR) + 1),
			Type: EDR_ATTR_SYNC,
		},
		Data: []byte(SYNC_STR),
	}
	genl_req := netlink.NetlinkGenericMessage{
		Header: netlink.GenlMsgHdr{
			Command: EDR_CMD_SYNC,
		},
	}

	genl_req.AddAttr(genl_ctrl_attr)

	//netlink_req := ns.NewNetlinkMessage(id, syscall.NLM_F_REQUEST)
	netlink_req := ns.NewNetlinkMessage(id, syscall.NLMSG_DONE)
	netlink_req.AddData(genl_req.Serial())

	for {
		ns.Send(netlink_req)
		time.Sleep(5 * time.Second)
	}

}

func DrvRecv(ns *netlink.NetlinkSocket) error {
	genls, err := ns.GenlRecv()
	if err != nil {
		return err

	}
	for _, genl := range genls {
		switch genl.Header.Command {
		case EDR_CMD_OPENAT_UPLOAD:
			handle_openat_info(genl.Data)
		case EDR_CMD_CLONE_UPLOAD:
			handle_clone_exit_info(genl.Data)
		case EDR_CMD_EXIT_GROUP_UPLOAD:
			handle_clone_exit_info(genl.Data)
		default:
			fmt.Println("Got default cmd")
		}
	}
	return nil
}
