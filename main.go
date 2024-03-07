package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"regexp"
	"strings"
)

// Listener struct represents a TCP listener
type Listener struct {
	Port     string
	Status   string
	Listener net.Listener
	Conns    []net.Conn // Track all connections associated with this listener
	Next     *Listener
}

// ListenerList represents a linked list of listeners
type ListenerList struct {
	Head *Listener
	Stop chan struct{}
}

func main() {
	fmt.Println("Splinter's Cell")
	listenerList := NewListenerList()
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("SPLINTER >>> ")
		command, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			continue
		}

		command = strings.TrimSpace(command)

		switch command {
		case "listen":
			fmt.Println("Create a listener with: listen -p <port>")
			fmt.Println("List active listeners: listen --ls")
			fmt.Println("Close all listeners: listen --close")
		case "listen --ls":
			listenerList.displayListeners()
		case "listen --close":
			listenerList.closeListeners()
			fmt.Println("All listeners closed.")
		default:
			// Check if the command matches "listen -p <port>"
			regexListen := regexp.MustCompile(`^listen -p \d+$`)
			matchListen := regexListen.FindString(command)
			if matchListen != "" {
				port := strings.Split(command, " ")[2]
				listenerList.registerListener(port)
			} else {
				fmt.Println("Invalid command. Use 'listen', 'listen --ls', 'listen --close', or 'listen -p <port>'.")
			}
		}
	}
}
