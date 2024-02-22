package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

func A_Handler(handle *pcap.Handle, t_handle *pcap.Handle) {
	// 设置过滤器，只抓取TCP协议的数据包
	filter := "icmp"
	err := handle.SetBPFFilter(filter)
	if err != nil {
		log.Fatal(err)
	}

	// 开始抓包
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		switch packet.NetworkLayer().(type) {
		case (*layers.IPv6):
			{
				fmt.Println("skipped ipv6 packet") //skip ipv6
				continue
			}

		default:
			{
				//nothing
			}
		}

		srcaddr := packet.NetworkLayer().(*layers.IPv4).SrcIP

		dstaddr := packet.NetworkLayer().(*layers.IPv4).DstIP

		if srcaddr.Equal(net.ParseIP("192.168.100.2")) && dstaddr.Equal(net.ParseIP("192.168.113.80")) {
			newpacket := gopacket.NewPacket(packet.Data(), layers.LayerTypeEthernet, gopacket.Default)

			//change mac addr
			linkLayer := newpacket.LinkLayer()
			if linkLayer == nil {
				fmt.Println("No link layer found in the packet")
			}

			ethLayer := linkLayer.(*layers.Ethernet)

			// srcMAC := ethLayer.SrcMAC
			// dstMAC := ethLayer.DstMAC

			//wlp3s0 if
			ethLayer.SrcMAC, _ = net.ParseMAC("18:47:3d:f2:b4:99") //local
			ethLayer.DstMAC, _ = net.ParseMAC("aa:4b:bd:7a:07:bd") //host

			// 修改IP地址
			networkLayer := newpacket.NetworkLayer()
			if networkLayer == nil {
				fmt.Println("No network layer found in the packet")
			}

			ipLayer, _ := networkLayer.(*layers.IPv4)
			if ipLayer != nil {
				ipLayer.SrcIP = net.IP{192, 168, 113, 219} // 修改为你想要的源IP地址
				// ipLayer.DstIP = net.IP{8, 8, 8, 8}     // 修改为你想要的目标IP地址
			} else {
				fmt.Println("No IPv4 layer found in the packet")
			}

			// 构建新的Packet对象

			fmt.Println(newpacket.Data())

			// 序列化修改后的数据包
			buf := gopacket.NewSerializeBuffer()
			opts := gopacket.SerializeOptions{
				ComputeChecksums: true,
				FixLengths:       true,
			}

			var serializableLayers []gopacket.SerializableLayer
			for _, layer := range newpacket.Layers() {
				if serializableLayer, ok := layer.(gopacket.SerializableLayer); ok {
					serializableLayers = append(serializableLayers, serializableLayer)
				}
			}

			err := gopacket.SerializeLayers(buf, opts, serializableLayers...)

			if err != nil {
				fmt.Println("Error serializing packet:", err)
				return
			}

			// 发送修改后的数据包
			err = t_handle.WritePacketData(buf.Bytes())
			if err != nil {
				fmt.Println("ERR:", err)
			}

			fmt.Println(newpacket)
		}

		// mtext := srcaddr.String() + ":" + srcport.String() + " --> " + dstaddr.String() + ":" + dstport.String()

	}
}

func B_Handler(handle *pcap.Handle, t_handle *pcap.Handle) {
	// 设置过滤器，只抓取TCP协议的数据包
	filter := "icmp"
	err := handle.SetBPFFilter(filter)
	if err != nil {
		log.Fatal(err)
	}

	// 开始抓包
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		switch packet.NetworkLayer().(type) {
		case (*layers.IPv6):
			{
				fmt.Println("skipped ipv6 packet") //skip ipv6
				continue
			}

		default:
			{
				//nothing
			}
		}

		srcaddr := packet.NetworkLayer().(*layers.IPv4).SrcIP

		dstaddr := packet.NetworkLayer().(*layers.IPv4).DstIP

		if srcaddr.Equal(net.ParseIP("192.168.113.80")) && dstaddr.Equal(net.ParseIP("192.168.113.219")) {
			newpacket := gopacket.NewPacket(packet.Data(), layers.LayerTypeEthernet, gopacket.Default)

			//change mac addr
			linkLayer := newpacket.LinkLayer()
			if linkLayer == nil {
				fmt.Println("No link layer found in the packet")
			}

			ethLayer := linkLayer.(*layers.Ethernet)

			// srcMAC := ethLayer.SrcMAC
			// dstMAC := ethLayer.DstMAC

			//virbr1 if
			ethLayer.SrcMAC, _ = net.ParseMAC("52:54:00:75:ba:00") //local
			ethLayer.DstMAC, _ = net.ParseMAC("52:54:00:68:25:38") //host

			// 修改IP地址
			networkLayer := newpacket.NetworkLayer()
			if networkLayer == nil {
				fmt.Println("No network layer found in the packet")
			}

			ipLayer, _ := networkLayer.(*layers.IPv4)
			if ipLayer != nil {
				// ipLayer.SrcIP = net.IP{192, 168, 184, 208} // 修改为你想要的源IP地址
				ipLayer.DstIP = net.IP{192, 168, 100, 2} // 修改为你想要的目标IP地址
			} else {
				fmt.Println("No IPv4 layer found in the packet")
			}

			// 构建新的Packet对象

			fmt.Println(newpacket.Data())

			// 序列化修改后的数据包
			buf := gopacket.NewSerializeBuffer()
			opts := gopacket.SerializeOptions{
				ComputeChecksums: true,
				FixLengths:       true,
			}

			var serializableLayers []gopacket.SerializableLayer
			for _, layer := range newpacket.Layers() {
				if serializableLayer, ok := layer.(gopacket.SerializableLayer); ok {
					serializableLayers = append(serializableLayers, serializableLayer)
				}
			}

			err := gopacket.SerializeLayers(buf, opts, serializableLayers...)

			if err != nil {
				fmt.Println("Error serializing packet:", err)
				return
			}

			// 发送修改后的数据包
			err = t_handle.WritePacketData(buf.Bytes())
			if err != nil {
				fmt.Println("ERR:", err)
			}

			fmt.Println(newpacket)
		}

		// mtext := srcaddr.String() + ":" + srcport.String() + " --> " + dstaddr.String() + ":" + dstport.String()

	}
}

func main() {
	// 打开网络接口
	a_handle, err := pcap.OpenLive("virbr1", 65535, true, pcap.BlockForever)
	if err != nil {
		log.Fatal(err)
	}
	defer a_handle.Close()

	b_handle, err := pcap.OpenLive("wlp3s0", 65535, true, pcap.BlockForever)
	if err != nil {
		log.Fatal(err)
	}
	defer b_handle.Close()

	go A_Handler(a_handle, b_handle)
	go B_Handler(b_handle, a_handle)

	for {
		time.Sleep(time.Second * 30)
	}
}
