package main

import (
	"bufio"
	"fmt"
	"net"
)

func NewListenerList() *ListenerList {
	return &ListenerList{
		Stop: make(chan struct{}),
	}
}

// handleClient handles the client connection
func handleClient(conn net.Conn) {
	defer conn.Close()

	// Read authentication message
	auth := make([]byte, 32)
	n, err := conn.Read(auth)
	if err != nil {
		fmt.Println("Error reading from connection:", err)
		return
	}
	authString := string(auth[:n])

	// Perform authentication
	if authString != "i_L0V_y0U_Ju5t1n_P3t3R\n" {
		fmt.Println("Authentication failed")
		return
	}

	fmt.Println("Agent authenticated successfully!")
	// Example: Echo back any message received
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		message := scanner.Text()
		fmt.Println("Received message:", message)
		_, err := fmt.Fprintln(conn, "Echo:", message)
		if err != nil {
			fmt.Println("Error writing to connection:", err)
			return
		}
	}
}

// registerListener registers a new listener
func (ll *ListenerList) registerListener(port string) {
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

	go func(stop <-chan struct{}) {
		for {
			select {
			case <-stop:
				return // Stop accepting new connections and exit
			default:
				// Accept connection
				conn, err := listener.Accept()
				if err != nil {
					// Handle error
					continue
				}

				// Handle the connection
				go handleClient(conn)
			}
		}
	}(ll.Stop) // Pass the stop channel to the goroutine
}

// displayListeners displays the active listeners
func (ll *ListenerList) displayListeners() {
	fmt.Println("Active Listeners:")
	current := ll.Head
	for current != nil {
		fmt.Println("Port:", current.Port, "- Status:", current.Status)
		current = current.Next
	}
}

// closeListeners closes all active listeners and associated connections
func (ll *ListenerList) closeListeners() {
	// Close the stop channel to signal stop to all goroutines
	close(ll.Stop)

	current := ll.Head
	for current != nil {
		fmt.Println("Closing listener on port:", current.Port)

		// Close all associated connections first
		for _, conn := range current.Conns {
			conn.Close()
		}
		current.Conns = nil // Clear the connections list

		// Close the listener
		current.Listener.Close()

		current = current.Next
	}

	ll.Head = nil // Reset the listener list
}
