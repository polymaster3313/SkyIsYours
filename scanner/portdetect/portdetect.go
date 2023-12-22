package portdetect

import (
	"fmt"
	"log"
	"net"
	"time"
)

func PortDetect(ip string) string {
	address := fmt.Sprintf("%s:%s", ip, "22")
	conn, err := net.DialTimeout("tcp", address, time.Second)
	if err != nil {
		log.Printf("%s no discovery of open ssh port\n", ip)
		return ""
	}
	defer conn.Close()
	return "ssh"
}
