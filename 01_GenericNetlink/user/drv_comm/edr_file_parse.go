package drv_comm

import (
	"edr_server/netlink"
	"encoding/binary"
	"fmt"
)

type OpenatInfo struct {
	Flag     uint64
	Filename string
}

func handle_openat_info(data []byte) {
	ad, err := netlink.NewAttributeDecoder(data, binary.LittleEndian)
	if err != nil {
		return
	}
	info := OpenatInfo{}
	for ad.Next() {
		switch ad.Type() {
		case EDR_ATTR_FLAG:
			info.Flag = ad.ByteOrder.Uint64(ad.AttrData())
		case EDR_ATTR_FILENAME:
			info.Filename = string(ad.AttrData())
		default:
			fmt.Println("openat default")
			//其他的属性都不需要
		}
	}
	fmt.Println(info)

}
