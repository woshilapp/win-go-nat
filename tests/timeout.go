package tests

import (
	"fmt"
	"time"
)

func Timeout_ip() {
	istr := ""
	c1 := make(chan int8, 1)
	go func() {
		fmt.Scanf("%s", &istr)
		c1 <- 1
	}()

	select {
	case <-c1:
		fmt.Println(istr)
	case <-time.After(time.Second * 5):
		fmt.Println("timeout")
	}
}
