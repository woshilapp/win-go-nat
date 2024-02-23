package main

import (
	"fmt"

	"github.com/woshilapp/win-go-nat/tests"
)

func main() {
	fmt.Println("started")

	// tests.Start()
	tests.Icmpid_launch_catch()
}
