package main

import (
	"log"

	"github.com/google/gopacket/pcap"
	"github.com/woshilapp/win-go-nat/globals"
)

func Init() {
	globals.Ins_if = "virbr1"
	globals.Out_if = "enp2s0"

	var i_err error
	var o_err error

	globals.Ins_Handle, i_err = pcap.OpenLive(globals.Ins_if, 1600, true, pcap.BlockForever)
	globals.Out_Handle, o_err = pcap.OpenLive(globals.Out_if, 1600, true, pcap.BlockForever)

	if i_err != nil || o_err != nil {
		log.Fatal("Error when open handle")
	}

	globals.Ins_Mac, globals.Ins_IP = Get_IPv4_MAC_Ifname(globals.Ins_if)
	globals.Out_Mac, globals.Out_IP = Get_IPv4_MAC_Ifname(globals.Out_if)

	if globals.Ins_Mac == nil || globals.Out_Mac == nil || globals.Ins_IP == nil || globals.Out_IP == nil {
		log.Fatal("Error when get addresses")
	}
}
