package main

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

func Is_inside(ip_addr string) bool {
	str_slice := strings.Split(ip_addr, ".")
	int_slice := make([]int, len(str_slice), cap(str_slice))

	fmt.Println(len(str_slice), cap(str_slice))

	if len(str_slice) != 4 {
		return false
	}

	for i, v := range str_slice {
		int_slice[i], _ = strconv.Atoi(v)
	}

	switch int_slice[0] {
	case 192:
		if int_slice[1] == 168 {
			return true
		}
	case 172:
		if int_slice[1] >= 16 || int_slice[1] <= 32 {
			return true
		}
	case 100, 10:
		return true
	}

	return false
}

func Get_IPv4_MAC_Ifname(interfaceName string) (net.HardwareAddr, net.IP) {
	// 获取接口信息
	ifi, err := net.InterfaceByName(interfaceName)
	if err != nil {
		// fmt.Println("Error:", err)
		return nil, nil
	}

	// 获取接口的 MAC 地址
	macAddr := ifi.HardwareAddr
	// fmt.Println("MAC 地址:", macAddr)

	var ipAddr net.IP

	// 获取接口的 IPv4 地址
	addrs, err := ifi.Addrs()
	if err != nil {
		// fmt.Println("Error:", err)
		return nil, nil
	}
	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				// fmt.Println("IPv4 地址:", ipNet.IP)
				ipAddr = ipNet.IP
				break
			}
		}
	}

	return macAddr, ipAddr
}
