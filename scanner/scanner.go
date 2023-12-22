package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

	zmq "github.com/pebbe/zmq4"
)

func scanner(target string, push *zmq.Socket) {
	file, err := os.Open(target)
	if err != nil {
		log.Fatalln("Error opening file:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line)

	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}
}

func checker(ip string) {

}
func main() {
	fmt.Println("Scanner started")

	ipcPtr := flag.Bool("ipc", false, "interprocess communication")
	filePtr := flag.String("target", "", "Path to the target")
	hostPtr := flag.String("host", "", "Host ip")
	portPtr := flag.String("port", "", "Port number")

	file := *filePtr
	port := *portPtr
	host := *hostPtr
	ipc := *ipcPtr

	if !ipc {
		if port == "" {
			port = "1000"
		}

		if host == "" {
			host = "*"
		}
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

	if !ipc {
		err = push.Bind(fmt.Sprintf("tcp://%s:%s", host, port))

		if err != nil {
			log.Fatalln(err)
		}

		log.Printf("Binded on tcp://%s:%s\n", host, port)
		scanner(file, push)
	}

	err = push.Bind("ipc:///tmp/Sky")

	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Binded on ipc:///tmp/Sky")
	scanner(file, push)

}
