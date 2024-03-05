package main

import (
	"github.com/woshilapp/win-go-nat/globals"
)

func Init() {
	globals.Ins_if = "virbr1"
	globals.Out_if = "enp1s0"
}
