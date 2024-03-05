package main

import (
	"bufio"
	"fmt"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

type routeItem struct {
	Destination string
	Gateway     string
	Mask        string
	Ifname      string
}

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

func Route() {
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

	var routes []routeItem
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

			route := routeItem{
				Destination: hexToIPv4(fields[1]),
				Gateway:     hexToIPv4(fields[2]),
				Mask:        hexToIPv4(fields[7]),
				Ifname:      fields[0],
			}

			routes = append(routes, route)
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

			route := routeItem{
				Destination: fields[0],
				Gateway:     fields[2],
				Mask:        fields[3],
				Ifname:      fields[len(fields)-1],
			}

			routes = append(routes, route)
		}
	case "windows":
		//fk windows
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading command output:", err)
		return
	}

	for _, route := range routes {
		fmt.Printf("Destination: %s, Gateway: %s, Mask: %s, Interface: %s\n", route.Destination, route.Gateway, route.Mask, route.Ifname)
	}
}
