package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// Listener struct represents a TCP listener
type Listener struct {
	Port     string
	Status   string
	Listener net.Listener
	//Conns    []net.Conn // Track all connections associated with this listener if you want to do multiple
	//Conns net.Conn //Handles only 1 cnnection
	Next *Listener
}

// ListenerList represents a linked list of listeners
type ListenerList struct {
	Head *Listener
	Stop chan struct{}
}

// Session struct represents a TCP Session (each session can only have 1 net.Conn)
type Session struct {
	id     int
	Port   string
	Status string
	Conn   net.Conn
	Next   *Session
}

// ListenerList represents a linked list of listeners
type SessionList struct {
	Head *Session
	Stop chan struct{}
}

func main() {
	fmt.Println("CODENAME: SAMURAI")
	listenerList := NewListenerList()
	sessionList := NewSessionList()
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("DEMOLITRON >>> ")
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
		case "session":
			fmt.Println("To activate a session: session --id <sessionID>")
			fmt.Println("List active sessions: session --ls")
			fmt.Println("Close all sessions: session --close")
		case "session --ls":
			sessionList.displaySessions()
		case "session --close":
			sessionList.closeSessions()
			fmt.Println("All sessions closed. Za3altneh ...")
		default:
			// Check if the command matches "listen -p <port>"
			regexListen := regexp.MustCompile(`^listen -p \d+$`)
			matchListen := regexListen.FindString(command)
			regexSession := regexp.MustCompile(`^session --id \d+$`)
			matchSession := regexSession.FindString(command)
			switch {
			case matchListen != "":
				{
					port := strings.Split(command, " ")[2]
					listenerList.registerListener(port, sessionList)
				}
			case matchSession != "":
				{
					idStr := strings.Split(command, " ")[2]
					id, _ := strconv.Atoi(idStr)
					openSession(id, sessionList)
				}
			default:
				{
					fmt.Println("Invalid command. Use 'listen', 'session'")
				}
			}
		}
	}
}
