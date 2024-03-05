package main

import (
	"fmt"
	"time"

	"github.com/woshilapp/win-go-nat/globals"
	// "github.com/woshilapp/win-go-nat/tests"
)

func main() {
	Init()

	fmt.Println("started")

	// tests.Start()
	// tests.Icmpid_launch_catch()

	// fmt.Println(Is_inside("8.8.8.8"))
	go Arp_Handle(globals.Ins_if)
	go Arp_Handle(globals.Out_if)
	// Route()

	for {
		time.Sleep(time.Second * 100)
	}
}
