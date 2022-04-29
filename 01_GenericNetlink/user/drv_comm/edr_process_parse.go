package drv_comm

import (
	"edr_server/netlink"
	"encoding/binary"
	"fmt"
)

type ProcessInfo struct {
	TimeStamp uint64
	TaskName  string
	curPID    uint64
	childPID  uint64
}

func handle_clone_exit_info(data []byte) {
	ad, err := netlink.NewAttributeDecoder(data, binary.LittleEndian)
	if err != nil {
		return
	}
	info := ProcessInfo{}
	for ad.Next() {
		switch ad.Type() {
		case EDR_ATTR_TIMESTAMP:
			info.TimeStamp = ad.ByteOrder.Uint64(ad.AttrData())
		case EDR_ATTR_TASKNAME:
			info.TaskName = string(ad.AttrData())
		case EDR_ATTR_CUR_PID:
			info.curPID = ad.ByteOrder.Uint64(ad.AttrData())
		case EDR_ATTR_CHILD_PID:
			info.childPID = ad.ByteOrder.Uint64(ad.AttrData())
		default:
			fmt.Println("clone default")
			//其他的属性都不需要
		}
	}
	fmt.Println(info)

}
