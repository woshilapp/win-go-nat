package main

import (
	"fmt"
	"log"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

func Arp_Start() {
	// 打开网络设备以进行数据包捕获
	handle, err := pcap.OpenLive("enp2s0", 1600, true, pcap.BlockForever)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	// 为ARP数据包设置过滤器
	err = handle.SetBPFFilter("arp")
	if err != nil {
		log.Fatal(err)
	}

	// 初始化ARP缓存
	arpCache := make(map[string]string)

	// 开始捕获数据包
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		// 检查数据包是否为ARP数据包
		arpLayer := packet.Layer(layers.LayerTypeARP)
		arpPacket, _ := arpLayer.(*layers.ARP)

		// 从ARP数据包中提取IP和MAC地址
		ip := fmt.Sprintf("%d.%d.%d.%d", arpPacket.SourceProtAddress[0], arpPacket.SourceProtAddress[1],
			arpPacket.SourceProtAddress[2], arpPacket.SourceProtAddress[3])
		mac := fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x", arpPacket.SourceHwAddress[0], arpPacket.SourceHwAddress[1],
			arpPacket.SourceHwAddress[2], arpPacket.SourceHwAddress[3], arpPacket.SourceHwAddress[4], arpPacket.SourceHwAddress[5])

		// 使用新条目更新或刷新ARP缓存
		arpCache[ip] = mac

		// 打印ARP缓存
		fmt.Println("ARP 缓存:")
		for ip, mac := range arpCache {
			fmt.Printf("IP: %s, MAC: %s\n", ip, mac)
		}
		fmt.Println("--------------")
	}
}
