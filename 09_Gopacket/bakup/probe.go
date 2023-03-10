package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"math/rand"
	"reflect"
	"syscall"
	"time"
	"unsafe"
)

const ENDIAN = "LittleEndian"

type ipHeader struct {
	VersionIHL uint8
	TOS        uint8
	TotalLen   uint16
	ID         uint16
	FlagsFrag  uint16
	TTL        uint8
	Protocol   uint8
	IPChecksum uint16
	SrcAddr    uint32
	DstAddr    uint32
}
type tcpPsdHeader struct {
	SrcAddr  uint32
	DstAddr  uint32
	Reverse  uint8
	Protocol uint8
	TcpLen   uint16
}
type tcpHeader struct {
	SrcPort  uint16
	DstPort  uint16
	Sequence uint32
	AckNo    uint32
	//Offset = Header Lengh + Reverse + flags
	Offset        uint16
	Window        uint16
	TCPChecksum   uint16
	UrgentPointer uint16
}

func rawSend(descriptor int, sockaddr syscall.SockaddrInet4, payload []byte) {
	err := syscall.Sendto(descriptor, payload, 0, &sockaddr)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("send success")
}

func rawRecv(descriptor int) {
	buff := make([]byte, 1514)
	//var packet TCPIP
	fmt.Println("enter rawRecv")
	//size, addr, err := syscall.Recvfrom(descriptor, buff, syscall.MSG_NOSIGNAL)
	size, addr, err := syscall.Recvfrom(descriptor, buff, 0)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Recvfrom return %d\n", size)
	fmt.Printf("addr: %v\n", addr)
}

/*
	@buildPayload : 通过结构体构造byte.Buffer
	结构体成员变量只能是type (bool, int8, uint8, int16, float32, complex64, …)或对应array
*/
func buildPayload(o interface{}, order binary.ByteOrder) (*bytes.Buffer, error) {
	var t reflect.Type
	var v reflect.Value
	var err error

	payload := bytes.NewBuffer([]byte{})

	if reflect.ValueOf(o).Type().Kind() == reflect.Struct {
		t = reflect.TypeOf(o)
		v = reflect.ValueOf(o)
	} else if reflect.ValueOf(o).Type().Kind() == reflect.Ptr {
		t = reflect.TypeOf(o).Elem()
		v = reflect.ValueOf(o).Elem()
	} else {
		return payload, fmt.Errorf("wrong parameter")

	}

	for i := 0; i < t.NumField(); i++ {
		value := v.Field(i).Interface()
		err = binary.Write(payload, order, value)
		if err != nil {
			fmt.Println(err)
			payload.Reset()
			break
		}
	}
	return payload, err

}

/*
	@parsePayload : 通过[]byte解析struct
	结构体成员变量只能是type (bool, int8, uint8, int16, float32, complex64, …)或对应array
*/
func parsePayload(data []byte, o interface{}, order binary.ByteOrder) error {
	var t reflect.Type
	var v reflect.Value

	if reflect.ValueOf(o).Type().Kind() == reflect.Ptr {
		t = reflect.TypeOf(o).Elem()
		v = reflect.ValueOf(o).Elem()
	} else {
		return fmt.Errorf("wrong parameter")
	}
	if len(data) < int(unsafe.Sizeof(o)) {
		return fmt.Errorf("data length is not long")
	}

	index := 0
	for i := 0; i < t.NumField(); i++ {
		value := v.Field(i)
		switch value.Kind() {
		case reflect.Uint8:
			//uint8(binary.BigEndian.Uint8())
			value.SetUint(uint64(uint8(data[index])))
			index += 1
		case reflect.Uint16:
			//uint8(binary.BigEndian.Uint8())
			value.SetUint(uint64(binary.BigEndian.Uint16(data[index : index+2])))
			index += 2
		case reflect.Uint32:
			value.SetUint(uint64(binary.BigEndian.Uint32(data[index : index+4])))
			index += 4
		}
	}
	return nil

}

//网际校验和算法适用于IP、UDP、ICMP等协议的校验。
func CheckSum_generic(data []byte) uint16 {
	var (
		sum    uint32
		length int = len(data)
		index  int
	)
	//以每16位为单位进行求和，直到所有的字节全部求完或者只剩下一个8位字节（如果剩余一个8位字节说明字节数为奇数个）
	for length > 1 {
		sum += uint32(data[index])<<8 + uint32(data[index+1])
		index += 2
		length -= 2
	}
	//如果字节数为奇数个，要加上最后剩下的那个8位字节
	if length > 0 {
		sum += uint32(data[index])
	}
	//加上高16位进位的部分
	sum += (sum >> 16)
	//返回的时候求反
	return uint16(^sum)
}

func buildIP(a uint8, b uint8, c uint8, d uint8) uint32 {
	return (uint32(a) << 24) | (uint32(b) << 16) | (uint32(c) << 8) | (uint32(d))
}

func CheckFlags(offset uint16) bool {
	//syn, bit1
	fin := (offset >> 0) & 0x1
	syn := (offset >> 1) & 0x1
	rst := (offset >> 2) & 0x1
	ack := (offset >> 4) & 0x1

	fmt.Println(fin, syn, rst, ack)
	if syn == 1 && ack == 1 && rst == 0 && fin == 0 {
		return true
	} else {
		return false
	}

}

func main() {
	//初始化时间种子
	rand.Seed(time.Now().UnixNano()) // 纳秒时间戳

	//IP头赋值
	iph := ipHeader{
		//低4位，version =4,
		//高4位，首部长度=5*4=20字节
		VersionIHL: (uint8(4) << 4) | (uint8(5) & 0xf),
		TOS:        0x00,
		//总长度为ip头+tcp头
		TotalLen: 20 + 20,
		ID:       uint16(rand.Intn(math.MaxUint16)),
		//ID:        1000,
		FlagsFrag: 0,
		TTL:       64,
		//IPPROTO_TCP = 6
		Protocol:   6,
		IPChecksum: 0,
		SrcAddr:    buildIP(192, 168, 1, 14),
		DstAddr:    buildIP(192, 168, 1, 1),
	}
	//构造[]byte以计算IP头部校验和
	ipPayload, err := buildPayload(iph, binary.BigEndian)
	iph.IPChecksum = CheckSum_generic(ipPayload.Bytes())

	//tcp头部赋值
	tcph := tcpHeader{
		SrcPort: 1024,
		DstPort: 80,
		//序列号赋值为目的ip地址，便于接收确认
		Sequence: iph.DstAddr,
		AckNo:    0,
		// Header Lengh + Reverse + flags
		Offset: 0x5002,
		//Window:        uint16(rand.Intn(math.MaxUint16)),
		Window:        1001,
		TCPChecksum:   0,
		UrgentPointer: 0,
	}

	//伪tcp头部赋值，用以计算tcp校验和
	psd_tcph := tcpPsdHeader{
		SrcAddr:  iph.SrcAddr,
		DstAddr:  iph.DstAddr,
		Reverse:  0,
		Protocol: iph.Protocol,
		TcpLen:   uint16(unsafe.Sizeof(tcph)),
	}
	//将tcp头部和tcp伪头部合成[]byte计算校验和
	tcpPayload, err := buildPayload(tcph, binary.BigEndian)
	tcpPsdPayload, err := buildPayload(psd_tcph, binary.BigEndian)
	tcph.TCPChecksum = CheckSum_generic(append(tcpPsdPayload.Bytes(), tcpPayload.Bytes()...))

	//将已经填充校验和的IP头和TCP头合并成一个[]byte
	ipPayload_checksum, _ := buildPayload(iph, binary.BigEndian)
	tcpPayload_checksum, _ := buildPayload(tcph, binary.BigEndian)
	sendPayload := append(ipPayload_checksum.Bytes(), tcpPayload_checksum.Bytes()...)

	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_TCP)
	//fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_RAW)
	if err != nil {
		fmt.Println(err)
	}
	defer syscall.Close(fd)
	err = syscall.SetNonblock(fd, true)
	if err != nil {
		fmt.Println(err)
	}
	err = syscall.SetsockoptInt(fd, syscall.IPPROTO_IP, syscall.IP_HDRINCL, 1)
	if err != nil {
		fmt.Println(err)
	}
	//不绑定网卡
	//syscall.BindToDevice(fd, "ens33")

	var dest_slice = make([]byte, 4)
	var dest_arr [4]byte

	binary.LittleEndian.PutUint32(dest_slice, iph.DstAddr)
	copy(dest_arr[:], dest_slice[:4])

	addr := syscall.SockaddrInet4{
		Port: int(tcph.DstPort),
		Addr: dest_arr,
	}
	rawSend(fd, addr, sendPayload)

	//rawRecv(fd)
	buff := make([]byte, 15140)
	fmt.Println("enter rawRecv")
	for {
		size, _, _ := syscall.Recvfrom(fd, buff, syscall.MSG_NOSIGNAL)
		if err != nil {
			fmt.Println(err)
		}
		//fmt.Printf("Recvfrom return %d\n", size)
		//fmt.Printf("addr: %v\n", addr2)
		if size < int(unsafe.Sizeof(iph))+int(unsafe.Sizeof(tcph)) {
			continue
		}
		recvip := ipHeader{}
		recvtcp := tcpHeader{}
		err = parsePayload(buff, &recvip, binary.BigEndian)
		err = parsePayload(buff[20:], &recvtcp, binary.BigEndian)
		//目的地址非本机IP
		if recvip.DstAddr != iph.SrcAddr {
			continue
		}
		//非tcp
		if recvip.Protocol != syscall.IPPROTO_TCP {
			continue
		}
		//目的端口是否一致
		if recvtcp.DstPort != tcph.SrcPort {
			continue
		}
		if !CheckFlags(recvtcp.Offset) {
			continue
		}
		fmt.Println(recvip.SrcAddr, recvtcp.AckNo)
		fmt.Printf("ip1 is %x\n", recvip.DstAddr)
		fmt.Printf("ip2 is %x\n", recvip.SrcAddr)
		fmt.Printf("port1 is %x\n", recvtcp.DstPort)
		fmt.Printf("port2 is %x\n", recvtcp.SrcPort)
		break
	}
	/////////////////////////
	/*
		fp, err := os.Create("binfile")
		if err != nil {
			fmt.Println(err)
		}
		defer fp.Close()
		fp.Write(buff)
	*/
	////////////////////////////////////

	/*
		recvip := ipHeader{}
		recvtcp := tcpHeader{}
		err = parsePayload(buff, &recvip, binary.BigEndian)
		fmt.Println(err)
		err = parsePayload(buff[20:], &recvtcp, binary.BigEndian)
		fmt.Println(err)
	*/

}
