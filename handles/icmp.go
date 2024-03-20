package handles

import (
	// "fmt"
	"log"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/woshilapp/win-go-nat/globals"
)

type ICMP_handle struct {
	Src_IP string
	Id     int32
	Close  bool
}

func (icmp *ICMP_handle) Get_src() string {
	return icmp.Src_IP
}

func (icmp *ICMP_handle) Start_Handle() {
	bpf, err := pcap.NewBPF(layers.LinkTypeEthernet, 1600, "icmp")
	if err != nil {
		log.Fatal(err)
	}

	packetSource := gopacket.NewPacketSource(globals.Ins_Handle, globals.Ins_Handle.LinkType())
	for packet := range packetSource.Packets() {
		if bpf.Matches(packet.Metadata().CaptureInfo, packet.Data()) {
			//do sth.
		}
	}
}
