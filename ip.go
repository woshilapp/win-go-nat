package main

import (
	"fmt"
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
