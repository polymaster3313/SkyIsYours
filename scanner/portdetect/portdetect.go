package portdetect

import (
	"fmt"
	"net"
	"time"
)

func PortDetect(ip string) bool {
	address := fmt.Sprintf("%s:%s", ip, "22")
	conn, err := net.DialTimeout("tcp", address, time.Second)
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}
