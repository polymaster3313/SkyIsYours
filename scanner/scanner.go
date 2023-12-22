package main

import (
	"SkyIsYours/scanner/portdetect"
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"sync"
	"time"

	zmq "github.com/pebbe/zmq4"
)

func scanner(target string, rate int, push *zmq.Socket) {
	var wg sync.WaitGroup
	var list []string
	clist := make(chan string)
	sem := make(chan struct{}, rate)

	file, err := os.Open(target)
	if err != nil {
		log.Fatalln("Error opening file:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		sem <- struct{}{}

		if line == "" {
			<-sem
			continue
		}

		if !isPublicIPv4(line) {
			log.Printf("\033[31m%s is not a public ip\033[0m", line)
			<-sem
			continue
		}
		wg.Add(1)
		go checker(line, clist, &wg, sem)
	}

	go func() {
		wg.Wait()
		close(clist)
	}()

	for result := range clist {
		list = append(list, result)
	}

	if len(list) == 0 {
		log.Println("\033[31mno target to push to bruter\033[0m")
		return
	}

	log.Printf("pushing %d target/s to bruter\n", len(list))

	time.Sleep(time.Millisecond * 100)
	for _, value := range list {

		log.Printf("sending %s to bruter", value)
		push.Send(value, 0)
	}

	time.Sleep(time.Millisecond * 100)
}

func isPublicIPv4(address string) bool {
	ip := net.ParseIP(address)
	if ip == nil || ip.To4() == nil {
		return false
	}

	return !ip.IsLoopback() && !ip.IsLinkLocalMulticast() && !ip.IsLinkLocalUnicast() && !ip.IsMulticast() && !ip.IsUnspecified() && !ip.IsPrivate()
}

func checker(ip string, clist chan<- string, wg *sync.WaitGroup, sem <-chan struct{}) {
	defer func() {
		<-sem
		wg.Done()
	}()

	result := portdetect.PortDetect(ip)
	if result == "ssh" {
		clist <- ip
	}
}

func getDefaultValue(value, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}

func getConnection(ipc bool, host, port string) string {
	if !ipc {
		return fmt.Sprintf("tcp://%s:%s", host, port)
	}
	return "ipc:///tmp/Sky"
}

func main() {
	ipcPtr := flag.Bool("ipc", false, "use interprocess communication")
	filePtr := flag.String("target", "", "path to the target file")
	hostPtr := flag.String("host", "", "host IP")
	portPtr := flag.String("port", "", "port number")
	ratePtr := flag.String("rate", "", "rate of IP scanning")

	flag.Parse()

	file := *filePtr
	port := getDefaultValue(*portPtr, "5544")
	host := getDefaultValue(*hostPtr, "127.0.0.1")
	ipc := *ipcPtr
	rateStr := *ratePtr

	rateInt, err := strconv.Atoi(rateStr)
	if err != nil || rateStr == "" || rateInt < 0 {
		rateInt = 100
	}

	if file == "" {
		log.Fatalln("No target file specified (--target <path>)")
	}
	if _, err := os.Stat(file); err != nil {
		if os.IsNotExist(err) {
			log.Fatalf("The file %s does not exist.\n", file)
		} else {
			log.Fatalf("Error checking the file: %v\n", err)
		}
	}

	context, err := zmq.NewContext()
	if err != nil {
		log.Fatalln(err)
	}

	push, err := context.NewSocket(zmq.PUSH)

	if err != nil {
		log.Fatalln(err)
	}

	if err != nil {
		log.Fatalln(err)
	}

	connection := getConnection(ipc, host, port)
	err = push.Bind(connection)
	if err != nil {
		log.Fatalf("Error binding to %s: %v\n", connection, err)
	}

	log.Printf("Binded on %s\n", connection)
	scanner(file, rateInt, push)
}
