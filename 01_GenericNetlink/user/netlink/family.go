package netlink

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"syscall"
)

var errInvalidFamilyVersion = errors.New("invalid family version attribute")

// A Family is a generic netlink family.
type Family struct {
	ID      uint16
	Version uint8
	Name    string
}

func (ns *NetlinkSocket) GetFamily(name string) (Family, error) {
	genl_ctrl_attr := &Attribute{
		NlAttr: syscall.NlAttr{
			//这里需要额外+1，表示'\0'
			Len:  uint16(NlaHeaderLen + len(name) + 1),
			Type: CTRL_ATTR_FAMILY_NAME,
		},
		Data: []byte(name),
	}
	genl_req := NetlinkGenericMessage{
		Header: GenlMsgHdr{
			Command: CTRL_CMD_GETFAMILY,
			Version: CTRL_CMD_VERSION,
		},
	}

	genl_req.AddAttr(genl_ctrl_attr)

	netlink_req := ns.NewNetlinkMessage(GENL_ID_CTRL, syscall.NLM_F_REQUEST)
	netlink_req.AddData(genl_req.Serial())
	err := ns.Send(netlink_req)
	if err != nil {
		return Family{}, err
	}
	genls, err := ns.GenlRecv()
	if err != nil {
		return Family{}, err

	}
	families := make([]Family, 0, len(genls))
	for _, genl := range genls {
		ad, err := NewAttributeDecoder(genl.Data, binary.LittleEndian)
		if err != nil {
			return Family{}, err
		}
		//var f Family
		f := Family{}
		for ad.Next() {
			switch ad.Type() {
			case CTRL_ATTR_FAMILY_ID:
				f.ID = ad.ByteOrder.Uint16(ad.a.Data)
			case CTRL_ATTR_FAMILY_NAME:
				f.Name = string(ad.a.Data)
			case CTRL_ATTR_VERSION:
				v := ad.ByteOrder.Uint32(ad.a.Data)
				if v > math.MaxUint8 {
					return Family{}, errInvalidFamilyVersion
				}
				f.Version = uint8(v)
			default:
				//其他的属性都不需要
			}
		}
		if f.ID > 0 && f.Name != "" && f.Version > 0 {
			families = append(families, f)
		}

	}
	if len(families) == 0 {
		return Family{}, fmt.Errorf("response has no Family Info")
	}
	if len(families) != 1 {
		return Family{}, fmt.Errorf("response contains multiple Family")
	}
	return families[0], nil

}
