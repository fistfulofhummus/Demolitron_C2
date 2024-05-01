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

// handleClient handles the client connection //Rewrite this. I dont know why it has to be a goroutine.
func handleClient(ll *ListenerList, conn net.Conn, port string /*, listener net.Listener*/, sl *SessionList) {
	if !authSession(&conn) {
		return
	}
	fmt.Println("Agent authenticated successfully!")
	fmt.Println("Session Created")
	//conn.Write([]byte("SessionOpen\n")) Will Use This late to get hostinfo and initial config
	sl.registerSession(port, conn)
	//listener.Close() //TCP connections once open do not require a listener. I am doing 1 listener and session per port for now.
	ll.updateListenerStatus(port, "SESSION")
	//Checks if the session is still alive. Rewrite this in the sessions and not in the listeners section
	for {
		//alive := sha256.Sum256([]byte("Areyoualive?!"))
		_, err := conn.Write([]byte("AreYouAlive\n"))
		if err != nil {
			conn.Close()
			fmt.Println("The Client closed the remote connection.")
			conn = nil
			return
		}
		time.Sleep(600 * time.Second)
	}
}

// registerListener registers a new listener
func (ll *ListenerList) registerListener(port string, sl *SessionList) {
	addr := ":" + port

	// Resolve TCP address
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		fmt.Println("Could not resolve TCP address:", err)
		return
	}

	// Listen on the specified port
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		fmt.Println("Could not listen on port:", err)
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

	fmt.Println("Listening on port:", addr)

	// go func(stop <-chan struct{}) {
	// 	for {
	// 		select {
	// 		case <-stop:
	// 			listener.Close() //New
	// 			return           // Stop accepting new connections and exit
	// 		default:
	// 			// Accept connection
	// 			conn, err := listener.Accept()
	// 			if err != nil {
	// 				// Handle errors later
	// 				continue
	// 			}
	// 			// Handle the connection
	// 			handleClient(ll, conn, port /*, listener*/, sl)
	// 		}
	// 	}
	// }(ll.Stop) // Doesnt help much. You can just terminate without getting into this headache.

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
			handleClient(ll, conn, port /*, listener*/, sl)
		}
	}()
}

// displayListeners displays the active listeners
func (ll *ListenerList) displayListeners() {
	fmt.Println("\nActive Listeners:")
	current := ll.Head
	for current != nil {
		fmt.Println("Port:", current.Port, "- Status:", current.Status)
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
		fmt.Println("Closing listener on port:", current.Port)

		current.Listener.Close()

		current = current.Next
	}
	fmt.Println()
	ll.Head = nil // Reset the listener list
}
