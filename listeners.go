package main

import (
	//"crypto/sha256"
	"fmt"
	"time"

	//"hash/sha256"

	"net"
)

func NewListenerList() *ListenerList {
	return &ListenerList{
		Stop: make(chan struct{}),
	}
}

// handleClient handles the client connection
func handleClient(ll *ListenerList, conn *net.Conn, port string, sl *SessionList) {
	if !authSession(conn) {
		return
	}
	fmt.Println("[+]Agent authenticated successfully!")
	fmt.Println("[+]Session Created")
	//conn.Write([]byte("SessionOpen\n")) Will Use This late to get hostinfo and initial config
	sl.registerSession(port, *conn)
	ll.updateListenerStatus(port, "SESSION")
	//Checks if the session is still alive. Rewrite this in the sessions and not in the listeners section
	for {
		//alive := sha256.Sum256([]byte("Areyoualive?!"))
		authSession(conn)
		time.Sleep(180 * time.Second)
	}
}

// registerListener registers a new listener
func (ll *ListenerList) registerListener(port string, sl *SessionList) {
	addr := ":" + port

	// Resolve TCP address
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		fmt.Println("[-]Could not resolve TCP address:", err)
		return
	}

	// Listen on the specified port
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		fmt.Println("[-]Could not listen on port:", err)
		return
	}

	// Create a new Listener struct
	newListener := &Listener{
		Port:     port,
		Status:   "Listening",
		Listener: listener,
	}

	// Add to the head of the listener list
	newListener.Next = ll.Head
	ll.Head = newListener

	fmt.Println()
	fmt.Println("[+]Listening on port:", addr)
	fmt.Println()

	//Simple and elegant
	go func() {
		for {
			// Accept connection
			conn, err := listener.Accept()
			if err != nil {
				// Handle errors later
				continue
			}
			// Handle the connection
			listener.Close() //TCP connections once open do not require a listener. I just close it. I am doing 1 listener and session per port for now.
			handleClient(ll, &conn, port, sl)
		}
	}()
}

// displayListeners displays the active listeners
func (ll *ListenerList) displayListeners() {
	fmt.Println("\n[!]Active Listeners:")
	current := ll.Head
	for current != nil {
		fmt.Println("[+]Port:", current.Port, "- Status:", current.Status)
		current = current.Next
	}
}

func (ll *ListenerList) updateListenerStatus(targetPort string, status string /*, conn net.Conn*/) { //This is useless only 1 place uses it
	current := ll.Head
	for current.Port != targetPort {
		current = current.Next
	}
	current.Status = status
}

// closeListeners closes all active listeners
func (ll *ListenerList) closeListeners() {
	// Close the stop channel to signal stop to all goroutines. Overkill but okay.
	//close(ll.Stop)
	fmt.Println()
	current := ll.Head
	for current != nil {
		fmt.Println("[!]Closing listener on port:", current.Port)

		current.Listener.Close()

		current = current.Next
	}
	fmt.Println()
	ll.Head = nil // Reset the listener list
}
