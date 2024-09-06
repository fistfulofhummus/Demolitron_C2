package main

import (
	//"crypto/sha256"
	"fmt"
	"strconv"
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
	fmt.Println("\n\n[+]Agent authenticated successfully !") //\n\n added at first for clarity
	fmt.Print("[+]Session Created")                          //2 New lines exist after this fuck me
	hostname, user := hostinfo(conn)                         //Get some hostinfo instantly without much headache
	currentSessionID := sl.registerSession(port, hostname, user, *conn)
	fmt.Println("[!]Sesssion ID: " + strconv.Itoa(currentSessionID))
	ll.updateListenerStatus(port, "SESSION")
	//Checks if the session is still alive. Rewrite this in the sessions and not in the listeners section
	for {
		//alive := sha256.Sum256([]byte("Areyoualive?!"))
		if !authSession(conn) {
			fmt.Println("[!]Cleaning Up")     //Need a way of exiting the go routine. Will impliment something cooler later.
			sl.closeSession(currentSessionID) //We can use the updateSession functions instead if we want to keep a record of the dead sessions. There is a status field in the session struct after all
			return
		}
		time.Sleep(30 * time.Second)
	}
}

// Registers a new listener
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

// Displays the active listeners
func (ll *ListenerList) displayListeners() {
	fmt.Println("\n[!]Active Listeners:")
	current := ll.Head
	for current != nil {
		fmt.Println("[+]Port:", current.Port, "- Status:", current.Status)
		current = current.Next
	}
	fmt.Println()
}

// Updates the status of the listener. Dont need it but I ll keep it
func (ll *ListenerList) updateListenerStatus(targetPort string, status string /*, conn net.Conn*/) { //This is useless only 1 place uses it
	current := ll.Head
	for current.Port != targetPort {
		current = current.Next
	}
	current.Status = status
}

// closeListeners closes all active listeners
func (ll *ListenerList) closeListeners() {
	fmt.Println()
	current := ll.Head
	for current != nil {
		fmt.Println("[!]Closing listener on port:", current.Port)
		current.Listener.Close()
		current = current.Next
	}
	ll.Head = nil // Reset the listener list
	fmt.Println("[+]All listeners closed")
	fmt.Println()
}
