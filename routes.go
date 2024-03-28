package main

import (
	"bufio"
	"fmt"
	"net"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"github.com/woshilapp/win-go-nat/globals"
)

func hexToIPv4(hexStr string) string {
	ip := make([]string, 4)
	for i := 0; i < 4; i++ {
		hex := hexStr[i*2 : (i+1)*2]
		dec, err := strconv.ParseInt(hex, 16, 64)
		if err != nil {
			return ""
		}
		ip[i] = strconv.Itoa(int(dec))
	}
	for i, j := 0, 3; i < j; i, j = i+1, j-1 {
		ip[i], ip[j] = ip[j], ip[i]
	}
	return strings.Join(ip, ".")
}

func Get_Route() {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("cat", "/proc/net/route")
	case "windows":
		cmd = exec.Command("netstat", "-rn")
	case "darwin":
		cmd = exec.Command("netstat", "-rn")
	default:
		fmt.Println("Unsupported operating system")
		return
	}

	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Failed to execute command:", err)
		return
	}

	scanner := bufio.NewScanner(strings.NewReader(string(output)))

	switch runtime.GOOS {
	case "linux":
		// Skip the header line
		scanner.Scan()
		for scanner.Scan() {
			fields := strings.Fields(scanner.Text())
			if len(fields) < 11 {
				continue
			}

			route := globals.RouteItem{
				Destination: hexToIPv4(fields[1]),
				Gateway:     hexToIPv4(fields[2]),
				Mask:        hexToIPv4(fields[7]),
				Ifname:      fields[0],
			}

			if route.Ifname == globals.Out_if {
				globals.Routes = append(globals.Routes, route)
			}

		}
	case "darwin":
		// Skip the header lines
		scanner.Scan()
		scanner.Scan()
		for scanner.Scan() {
			fields := strings.Fields(scanner.Text())
			if len(fields) < 4 {
				continue
			}

			route := globals.RouteItem{
				Destination: fields[0],
				Gateway:     fields[2],
				Mask:        fields[3],
				Ifname:      fields[len(fields)-1],
			}

			globals.Routes = append(globals.Routes, route)
		}
	case "windows":
		// package main

		// import (
		// 	"fmt"
		// 	"github.com/StackExchange/wmi"
		// )

		// type Win32_IP4RouteTable struct {
		// 	Destination string
		// 	InterfaceIndex int
		// 	NextHop string
		// 	Metric1 uint
		// }

		// func main() {
		// 	var routes []Win32_IP4RouteTable
		// 	query := "SELECT Destination, InterfaceIndex, NextHop, Metric1 FROM Win32_IP4RouteTable"
		// 	err := wmi.Query(query, &routes)
		// 	if err != nil {
		// 		fmt.Println("Error querying Win32_IP4RouteTable:", err)
		// 		return
		// 	}

		// 	fmt.Println("Windows Route Table:")
		// 	for _, route := range routes {
		// 		fmt.Printf("Destination: %s, InterfaceIndex: %d, NextHop: %s, Metric1: %d\n", route.Destination, route.InterfaceIndex, route.NextHop, route.Metric1)
		// 	}
		// }
		// 	}

		// 	if err := scanner.Err(); err != nil {
		// 		fmt.Println("Error reading command output:", err)
		// 		return
		// 	}

		// 	for _, route := range globals.Routes {
		// 		fmt.Printf("Destination: %s, Gateway: %s, Mask: %s, Interface: %s\n", route.Destination, route.Gateway, route.Mask, route.Ifname)
		// 	}
		// }
	}
}

func CheckRouteTable(ip string) string {
	ipAddr := net.ParseIP(ip)
	if ipAddr == nil {
		return ip // 无效的IP地址，直接返回
	}

	for _, entry := range globals.Routes {
		_, ipNet, _ := net.ParseCIDR(entry.Destination + "/" + entry.Mask)
		if ipNet.Contains(ipAddr) {
			return entry.Gateway // IP地址匹配路由表项，返回下一跳地址
		}
	}

	return ip // 未匹配任何路由表项，返回原始IP地址
}
