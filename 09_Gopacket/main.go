package main

import (
	"fmt"
	"log"
	"main/utils"
	"math/rand"
	"net"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

var (
	//snapshot,每个数据包最大长度
	snapshot_len int32 = 65535
	//promiscuous,混杂模式,即是否接受目的地址不为本机的包
	promiscuous bool = false

	//timeout,抓到包返回的超时
	timeout time.Duration = 1 * time.Second
	//handle  *pcap.Handle
	localIP net.IP = net.IPv4(192, 168, 1, 6)
)

func main() {
	//随机数种子
	rand.Seed(time.Now().UnixNano())
	// 找到网卡设备，后续发包时使用
	//handle, err := pcap.OpenLive("ens33", snapshot_len, promiscuous, timeout)
	handle, err := pcap.OpenLive("ens33", snapshot_len, promiscuous, pcap.BlockForever)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	// Construct all the network layers we need.
	eth := layers.Ethernet{
		SrcMAC:       net.HardwareAddr{0x00, 0x0c, 0x29, 0xc3, 0x73, 0xd9},
		DstMAC:       net.HardwareAddr{0xa8, 0x5e, 0x45, 0xb4, 0x00, 0x9c},
		EthernetType: layers.EthernetTypeIPv4,
	}
	ip4 := layers.IPv4{
		SrcIP:    localIP,
		DstIP:    net.IPv4(192, 168, 1, 3),
		Version:  4,
		TTL:      64,
		Protocol: layers.IPProtocolTCP,

		Id: uint16(rand.Intn(0xffff)),
	}
	tcp := layers.TCP{
		SrcPort: layers.TCPPort(54321),
		DstPort: layers.TCPPort(1027),
		SYN:     true,

		Window: uint16(rand.Intn(0xffff)),
		Seq:    utils.IPToUint32(localIP),
	}
	fmt.Println(tcp.Seq)
	tcp.SetNetworkLayerForChecksum(&ip4)
	opts := gopacket.SerializeOptions{
		FixLengths:       true,
		ComputeChecksums: true,
	}

	buff := gopacket.NewSerializeBuffer()
	if err := gopacket.SerializeLayers(buff, opts, &eth, &ip4, &tcp); err != nil {
		return
	}
	handle.WritePacketData(buff.Bytes())

	for {
		data, _, err := handle.ReadPacketData()
		if err != nil {
			log.Printf("error reading packet: %v", err)
			continue
		}
		packet := gopacket.NewPacket(data, layers.LayerTypeEthernet, gopacket.NoCopy)

		ipLayer := packet.Layer(layers.LayerTypeIPv4)
		if ipLayer == nil {
			fmt.Println("ipLayer is nil")
		}
		ip, ok := ipLayer.(*layers.IPv4)
		if !ok {
			fmt.Println("parse ip false")
		}
		fmt.Println("ip addr is : ", ip.DstIP)

		tcpLayer := packet.Layer(layers.LayerTypeTCP)
		if tcpLayer == nil {
			fmt.Println("tcpLayer is nil")
		}
		tcp, ok := tcpLayer.(*layers.TCP)
		if !ok {
			fmt.Println("parse tcp false")
		}
		fmt.Println("tcp source: ", tcp.SrcPort)
		fmt.Println("tcp dest: ", tcp.DstPort)
		fmt.Println("tcp SYN: ", tcp.SYN)
		fmt.Println("tcp ACK: ", tcp.ACK)
		fmt.Println("tcp RST: ", tcp.RST)
		fmt.Println("tcp FIN: ", tcp.FIN)

		fmt.Println("tcp Ack seq: ", tcp.Ack)

		break
	}
}
