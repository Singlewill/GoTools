package netlink

import (
	"encoding/binary"
	"fmt"
	"syscall"
	"unsafe"
)

type GenlMsgHdr struct {
	Command uint8
	Version uint8
	Reverse uint16
}

// linux/netlink.h
type NetlinkGenericMessage struct {
	Header GenlMsgHdr
	Data   []byte
}

func (genlReq NetlinkGenericMessage) Serial() []byte {
	b := make([]byte, NlmsgAlign(syscall.NLA_HDRLEN))
	b[0] = genlReq.Header.Command
	b[1] = genlReq.Header.Version
	//b[2]和b[3]是pading
	return append(b, genlReq.Data...)

}

// ParseNetlinkMessage parses b as an array of netlink messages and
// returns the slice containing the NetlinkMessage structures.
func ParseGenlMessage(b []byte) (*NetlinkGenericMessage, error) {
	if len(b) < syscall.NLMSG_HDRLEN {
		return &NetlinkGenericMessage{}, fmt.Errorf("netlinkGenericMessage : wrong format")
	}
	h, dbuf, err := GenlMessageHeaderAndData(b)
	if err != nil {
		return &NetlinkGenericMessage{}, err
	}
	genlmsg := &NetlinkGenericMessage{Header: *h, Data: dbuf}
	return genlmsg, nil
}

func GenlMessageHeaderAndData(b []byte) (*GenlMsgHdr, []byte, error) {
	if len(b) < GenlmsgHeadrLen {
		return nil, nil, fmt.Errorf("genlMsgHdr : wrong format")
	}
	h := (*GenlMsgHdr)(unsafe.Pointer(&b[0]))
	return h, b[GenlmsgHeadrLen:], nil
}

func (genlReq *NetlinkGenericMessage) AddAttr(data *Attribute) error {
	if data == nil {
		return fmt.Errorf("attribute : Data wrong")
	}
	//attrLen := NlaAlign(syscall.NLA_HDRLEN) + NlaAlign(len(data.Data))
	attrLen := NlaAlign(int(data.Len))
	b := make([]byte, attrLen)
	binary.LittleEndian.PutUint16(b[0:2], data.Len)
	binary.LittleEndian.PutUint16(b[2:4], data.Type)
	copy(b[syscall.NLA_HDRLEN:], data.Data)
	genlReq.Data = append(genlReq.Data, b...)
	return nil
}

func (ns *NetlinkSocket) GenlRecv() ([]*NetlinkGenericMessage, error) {
	nlmsgs, err := ns.Receive()
	if err != nil {
		return make([]*NetlinkGenericMessage, 0), err

	}
	ret := make([]*NetlinkGenericMessage, 0, len(nlmsgs))
	for _, nlmsg := range nlmsgs {
		genlmsg, err := ParseGenlMessage(nlmsg.Data)
		if err != nil {
			return make([]*NetlinkGenericMessage, 0), err
		}
		ret = append(ret, genlmsg)
	}

	return ret, nil

}
