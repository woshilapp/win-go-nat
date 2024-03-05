package main

import (
	"fmt"
	"log"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"

	"github.com/woshilapp/win-go-nat/globals"
)

func Arp_Handle(if_name string) {
	// 打开网络设备以进行数据包捕获
	handle, err := pcap.OpenLive(if_name, 1600, true, pcap.BlockForever)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	// 为ARP数据包设置过滤器
	err = handle.SetBPFFilter("arp")
	if err != nil {
		log.Fatal(err)
	}

	// 开始捕获数据包
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		arpLayer := packet.Layer(layers.LayerTypeARP)
		arpPacket, _ := arpLayer.(*layers.ARP)

		// 从ARP数据包中提取IP和MAC地址
		ip := fmt.Sprintf("%d.%d.%d.%d", arpPacket.SourceProtAddress[0], arpPacket.SourceProtAddress[1],
			arpPacket.SourceProtAddress[2], arpPacket.SourceProtAddress[3])
		mac := fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x", arpPacket.SourceHwAddress[0], arpPacket.SourceHwAddress[1],
			arpPacket.SourceHwAddress[2], arpPacket.SourceHwAddress[3], arpPacket.SourceHwAddress[4], arpPacket.SourceHwAddress[5])

		// 使用新条目更新或刷新ARP缓存
		globals.Arp_Maps[ip] = mac

		// 打印ARP缓存
		fmt.Println("ARP 缓存:")
		for ip, mac := range globals.Arp_Maps {
			fmt.Printf("IP: %s, MAC: %s\n", ip, mac)
		}
		fmt.Println("--------------")
	}
}
