package main

import (
	"fmt"
	// "net"

	"time"

	"github.com/woshilapp/win-go-nat/globals"
	// "github.com/woshilapp/win-go-nat/tests"
)

func main() {
	Init()

	fmt.Println("started")
	fmt.Println(globals.Ins_IP, globals.Ins_Mac)
	fmt.Println(globals.Out_IP, globals.Out_Mac)

	// tests.Timeout_ip()

	// tests.Start()
	// tests.Icmpid_launch_catch()

	// fmt.Println(Is_inside("8.8.8.8"))
	go Arp_Handle(globals.Ins_Handle)
	go Arp_Handle(globals.Out_Handle)
	// Route()

	for {
		ma := Get_Mac_Addr("192.168.10.1")
		fmt.Printf("%T %s\n", ma, ma)
		time.Sleep(time.Second * 5)
	}
}
