package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"

	"github.com/woshilapp/win-go-nat/globals"
)

// func SParseMAC(mac_addr string) []byte {
// 	smac := make([]byte, 6)
// 	mac, _ := net.ParseMAC(mac_addr)
// 	copy(smac[:], mac)
// 	return smac
// }

func SParseIP(ip_addr string) []byte {
	ip := net.ParseIP(ip_addr).To4()
	sip := make([]byte, 4)
	copy(sip[:], ip)
	return sip
}

func buildARPPacket(srcMAC net.HardwareAddr, srcIP, targetIP net.IP) []byte {
	ethLayer := &layers.Ethernet{
		SrcMAC:       srcMAC,
		DstMAC:       net.HardwareAddr{0xff, 0xff, 0xff, 0xff, 0xff, 0xff}, // 广播地址
		EthernetType: layers.EthernetTypeARP,
	}

	arpLayer := &layers.ARP{
		AddrType:          layers.LinkTypeEthernet,
		Protocol:          layers.EthernetTypeIPv4,
		HwAddressSize:     6,
		ProtAddressSize:   4,
		Operation:         layers.ARPRequest,
		SourceHwAddress:   srcMAC,
		SourceProtAddress: SParseIP(srcIP.String()),
		DstHwAddress:      net.HardwareAddr{0, 0, 0, 0, 0, 0}, // 目标 MAC 地址为 0
		DstProtAddress:    SParseIP(targetIP.String()),
	}

	buffer := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{}
	if err := gopacket.SerializeLayers(buffer, opts, ethLayer, arpLayer); err != nil {
		fmt.Println("Error serializing packet:", err)
		return nil
	}

	return buffer.Bytes()
}

func Arp_Handle(handle *pcap.Handle) {
	// 为ARP数据包设置过滤器
	bpf, err := pcap.NewBPF(layers.LinkTypeEthernet, 1600, "arp")
	if err != nil {
		log.Fatal(err)
	}

	// 开始捕获数据包
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		if bpf.Matches(packet.Metadata().CaptureInfo, packet.Data()) { //filter
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
}

func Get_Mac_Addr(IPAddr string) interface{} {
	if globals.Arp_Maps[IPAddr] == "" {
		//send arp packet
		i_pkg := buildARPPacket(globals.Ins_Mac, globals.Ins_IP, net.ParseIP(IPAddr))
		globals.Ins_Handle.WritePacketData(i_pkg)

		o_pkg := buildARPPacket(globals.Out_Mac, globals.Out_IP, net.ParseIP(IPAddr))
		globals.Out_Handle.WritePacketData(o_pkg)

		mc := make(chan int, 1)

		go func() {
			for {
				_, ok := globals.Arp_Maps[IPAddr]
				if ok {
					mc <- 1
					break
				}
				time.Sleep(time.Millisecond * 50)
			}
		}()

		select {
		case <-mc:
			return globals.Arp_Maps[IPAddr]
		case <-time.After(time.Second * 1):
			return nil
		}

	} else {
		return globals.Arp_Maps[IPAddr]
	}
}
