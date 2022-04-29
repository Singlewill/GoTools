package netlink

import (
	"syscall"
	"unsafe"
)

//报头长度
//nlmsg报头长度
var NlmsgHeaderLen = NlmsgAlign(syscall.SizeofNlMsghdr)

//genl报头长度
var GenlmsgHeadrLen = NlmsgAlign(int(unsafe.Sizeof(GenlMsgHdr{})))

//attr报头长度
var NlaHeaderLen = NlaAlign(syscall.NLA_HDRLEN)

//对齐
//NLMSG对齐, netlink header和genl header都要用这个对齐
const nlmsgAlignTo = syscall.NLMSG_ALIGNTO

func NlmsgAlign(len int) int {
	return ((len) + nlmsgAlignTo - 1) & ^(nlmsgAlignTo - 1)
}

//ATTR对齐, attr用这个对齐
const nlaAlignTo = syscall.NLA_ALIGNTO

func NlaAlign(len int) int {
	return ((len) + nlaAlignTo - 1) & ^(nlaAlignTo - 1)
}
