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

func startHeartbeat(conn net.Conn, interval time.Duration, done chan struct{}) {
	buffer := make([]byte, 1024)

	for {
		// Send ping
		_, err := conn.Write([]byte("ping\n"))
		if err != nil {
			fmt.Println("[-] Failed to send ping:", err)
			conn.Close()
			close(done)
			return
		}

		// Set read deadline
		conn.SetReadDeadline(time.Now().Add(interval))

		// Read response
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("[-] No pong received (timeout or error):", err)
			conn.Close()
			close(done)
			return
		}

		msg := strings.TrimSpace(string(buffer[:n]))
		if msg != "pong" {
			fmt.Println("[-] Invalid response to ping:", msg)
			conn.Close()
			close(done)
			return
		}

		//fmt.Println("[+] Pong received")
		time.Sleep(interval) // wait before next ping
	}
}

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
	//Ran into an issue where sometimes this function starts before we reccieve a response due to slow internet. This makes sure we populate list with the data before starting heartbeats
	//Need a better fix. It still happens but less frequently
	// for {
	// 	if hostname != "" {
	// 		break
	// 	}
	// }
	//A better fix for the race condition me thinks
	done := make(chan struct{})
	go func() {
		time.Sleep(5 * time.Second)
		startHeartbeat(*conn, 120*time.Second, done)
	}()

	for {
		select {
		case <-session.StopChan:
			fmt.Println("[!] Session requested to stop")
			sl.closeSession(currentSessionID)
			return

		case <-done:
			fmt.Println("[!] Heartbeat failed, closing session")
			sl.closeSession(currentSessionID)
			return

		default:
			time.Sleep(2 * time.Second) // Keep loop light, real work is async
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
