package netlink

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

type NetlinkMessage struct {
	syscall.NetlinkMessage
}

type NetlinkSocket struct {
	fd       int
	protocol int
	groups   uint32
}

func NewNetlinkSocket(protocol int, groups uint32) (*NetlinkSocket, error) {
	// 创建 netlink socket
	fd, err := syscall.Socket(syscall.AF_NETLINK, syscall.SOCK_RAW, int(protocol))
	if err != nil {
		return nil, err
	}
	ns := &NetlinkSocket{
		fd:       fd,
		protocol: protocol,
		groups:   groups,
	}
	// 构造地址进行绑定
	sockaddr := &syscall.SockaddrNetlink{
		Family: syscall.AF_NETLINK,
		Pid:    uint32(os.Getpid()),
		Groups: groups,
	}
	// 绑定地址
	if err = syscall.Bind(fd, sockaddr); err != nil {
		syscall.Close(fd)
		return nil, err
	}
	return ns, nil
}

func (ns *NetlinkSocket) Close() {
	syscall.Close(ns.fd)
}

// 初始化一个 netlink 请求结构
func (ns *NetlinkSocket) NewNetlinkMessage(t uint16, f uint16) *NetlinkMessage {
	return &NetlinkMessage{
		syscall.NetlinkMessage{
			Header: syscall.NlMsghdr{
				Len: uint32(NlmsgAlign(syscall.SizeofNlMsghdr)),
				//Len:   uint32(syscall.SizeofNlMsghdr), // Length of message (including header)
				Type:  t,                   // Message content
				Flags: f,                   // Additional flags
				Seq:   uint32(0),           // Sequence number
				Pid:   uint32(os.Getpid()), // Sending process port ID
			},
		},
	}
}

// 按照 netlink 协议格式进行数据 append
// NetlinkMessageData 可能是数据，也可能是指令
func (req *NetlinkMessage) AddData(data []byte) {
	if data != nil {
		dataAlign := make([]byte, NlmsgAlign(len(data)))
		copy(dataAlign, data)
		req.Data = append(req.Data, dataAlign...)
		req.Header.Len += uint32(NlmsgAlign(len(data)))
	}
}

func (s *NetlinkSocket) Receive() ([]NetlinkMessage, error) {
	//s.mutex.Lock()
	//defer s.mutex.Unlock()

	// 获取运行时支持的 page 大小并创建 buffer
	rb := make([]byte, syscall.Getpagesize())
	// 从 socket 上获取数据保存到 buffer 中
	nr, _, err := syscall.Recvfrom(s.fd, rb, 0)
	if err != nil {
		return nil, err
	}

	if nr < syscall.NLMSG_HDRLEN {
		return nil, fmt.Errorf("got short response from netlink")
	}
	rb = rb[:nr]
	// 在 syscall/netlink_linux.go 中
	// ParseNetlinkMessage parses b as an array of netlink messages
	// and returns the slice containing the NetlinkMessage structures.
	//syscall.ParseNetlinkMessage()已做长度判断
	nlmsgs, err := syscall.ParseNetlinkMessage(rb)
	if err != nil {
		return nil, err
	}
	n := len(nlmsgs)
	ret := make([]NetlinkMessage, 0, n)
	for i := 0; i < len(nlmsgs); i++ {
		//判断Type标记
		if nlmsgs[i].Header.Type == syscall.NLMSG_ERROR {
			continue
		}
		ret = append(ret, NetlinkMessage{
			NetlinkMessage: nlmsgs[i],
		})
	}
	return ret, nil
	//return syscall.ParseNetlinkMessage(rb)

}

// 将 NetlinkMessage 中的全部内容序列化到一个大 buffer 中
func (req *NetlinkMessage) Serialize() []byte {
	//总长度
	length := NlmsgAlign(int(req.Header.Len))
	// 创建保存所有数据的大 buffer
	b := make([]byte, length)
	// 取 header 内容
	hdr := (*(*[syscall.SizeofNlMsghdr]byte)(unsafe.Pointer(req)))[:]
	next := NlmsgAlign(syscall.SizeofNlMsghdr)
	// 拷贝 header 到大 buffer
	copy(b[0:next], hdr)
	//拷贝数据
	copy(b[next:], req.Data)
	return b
}

func (ns *NetlinkSocket) Send(request *NetlinkMessage) error {
	//ns.mutex.Lock()
	//defer ns.mutex.Unlock()

	// 在 syscall/syscall_linux.go 中定义
	//type SockaddrNetlink struct {
	//    Family uint16     /* AF_NETLINK */
	//    Pad    uint16     /* zero       */
	//    Pid    uint32     /* 1. 为消息发送进程的 pid */
	//                      /* 2. 若希望内核处理或者进行消息多播，则设置为 0 */
	//                      /* 3. 否则设置为处理消息的线程组 id */
	//                      /* 4. 在一个进程的多个线程使用 netlink socket 的情况下，进程中的 Pid 字段可以设置为其他值 */
	//                      /*    因此，该字段实际上未必是进程 ID ，其只是用于区分不同接收者或发送者的一个标识 */
	//                      /*    用户可以根据自己的需要设置该字段 */
	//    Groups uint32     /* 1. 用于指定多播组 */
	//                      /* 2. 通过 bind() 函数能够把调用进程加入到该字段指定的多播组中 */
	//                      /* 3. 若设置为 0 ，则表示不加入任何多播组 */
	//    raw    RawSockaddrNetlink
	//}
	sockaddr := &syscall.SockaddrNetlink{
		Family: syscall.AF_NETLINK,
		Pid:    0, //For Linux Kernel
		Groups: 0, //Linux 驱动无netlink组
	}
	if err := syscall.Sendto(ns.fd, request.Serialize(), 0, sockaddr); err != nil {
		return err
	}
	return nil
}
