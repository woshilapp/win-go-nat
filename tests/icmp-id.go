package tests

import (
	"fmt"
	"log"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

func Icmpid_launch_catch() {
	// 打开网络接口进行数据包捕获
	handle, err := pcap.OpenLive("enp2s0", 65536, true, pcap.BlockForever)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	// 设置过滤器，只捕获 ICMP 数据包
	err = handle.SetBPFFilter("icmp")
	if err != nil {
		log.Fatal(err)
	}

	// 使用 gopacket PacketSource 进行数据包解析
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	// 循环读取数据包
	for packet := range packetSource.Packets() {
		// 解析 ICMP 数据包
		icmpLayer := packet.Layer(layers.LayerTypeICMPv4)
		if icmpLayer != nil {
			icmp, _ := icmpLayer.(*layers.ICMPv4)

			// 检查 ICMP 类型是否为 Echo 请求
			if icmp.TypeCode.Type() == layers.ICMPv4TypeEchoRequest || icmp.TypeCode.Type() == layers.ICMPv4TypeEchoReply {
				// 获取查询标识符（Query Identifier）
				fmt.Println(packet)
				fmt.Println("Query ID:", icmp.Id)
			}
		}
	}
}
