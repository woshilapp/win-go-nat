// wtf it's
// go is a good language
// but its format too strict
package main

import (
	"fmt"
	"log"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

func main() {
	// 打开网络接口
	handle, err := pcap.OpenLive("wlp3s0", 65535, true, pcap.BlockForever)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	// 设置过滤器，只抓取TCP协议的数据包
	filter := "tcp"
	err = handle.SetBPFFilter(filter)
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

		// 解析TCP层
		tcpLayer := packet.Layer(layers.LayerTypeTCP)

		srcaddr := packet.NetworkLayer().(*layers.IPv4).SrcIP
		srcport := packet.TransportLayer().(*layers.TCP).SrcPort

		dstaddr := packet.NetworkLayer().(*layers.IPv4).DstIP
		dstport := packet.TransportLayer().(*layers.TCP).DstPort

		mtext := srcaddr.String() + ":" + srcport.String() + " --> " + dstaddr.String() + ":" + dstport.String()

		if tcpLayer != nil {
			tcp, _ := tcpLayer.(*layers.TCP)
			if tcp.SYN && tcp.ACK {
				fmt.Println(mtext, "连接建立")
			} else if tcp.FIN && tcp.ACK {
				fmt.Println(mtext, "连接关闭")
			} else if tcp.RST {
				fmt.Println(mtext, "连接重置")
			} else if tcp.ACK {
				fmt.Println(mtext, "连接正常")
			}
		}
	}
}
