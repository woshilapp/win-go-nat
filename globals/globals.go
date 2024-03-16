package globals

import (
	"net"

	"github.com/google/gopacket/pcap"
)

var (
	Ins_if = ""
	Out_if = ""

	Ins_Handle *pcap.Handle
	Out_Handle *pcap.Handle

	Ins_Mac net.HardwareAddr
	Out_Mac net.HardwareAddr
	Ins_IP  net.IP
	Out_IP  net.IP

	Arp_Maps = make(map[string]string)
)
