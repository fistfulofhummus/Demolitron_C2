package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"regexp"
	"strings"
)

// Always define your structs before functions
type Listener struct {
	//cName    string
	port     string
	status   string
	listener net.Listener
	conn     net.Conn
	next     *Listener
}
type listenerList struct {
	head *Listener
}

type clientList struct {
	head *Client
}
type Client struct {
	name string
	//IP    string
	ID    int
	state string
	conn  net.Conn
	next  *Client
}

func handleListen(command string, lList *listenerList, cList *clientList) {
	command = strings.TrimSuffix(command, "\n")
	regexListen := regexp.MustCompile(`listen -p \d+`)
	matchListen := regexListen.FindString(command)
	if matchListen != "" {
		command = strings.Split(command, " -p ")[1]
		status, tcpListen := makeListener(command)
		if !status {
			fmt.Println("Could not create listener")
			return
		}
		defer tcpListen.Close()
		lList.registerListener(command, tcpListen)
		cConn := make(chan net.Conn, 1)
		//go listen(tcpListen) //Sexiest Shit Ever
		go listen2(tcpListen, cConn, lList)
		go cList.registerClient2(cConn)
	}
}
func handleListen2(command string, lList *listenerList, cList *clientList) {
	command = strings.TrimSuffix(command, "\n")
	regexListen := regexp.MustCompile(`listen -p \d+`)
	matchListen := regexListen.FindString(command)
	if matchListen != "" {
		command = strings.Split(command, " -p ")[1]
		command = ":" + command
		command = strings.TrimSuffix(command, "\n")
		tcpAddr, err := net.ResolveTCPAddr("tcp4", command)
		if err != nil {
			fmt.Println("Couldnt resolve tcp address")
			os.Exit(0)
		}
		tcpListener, err := net.ListenTCP("tcp", tcpAddr)
		if err != nil {
			fmt.Println("Something went wrong")
		}
		lList.registerListener(command, tcpListener)
		defer tcpListener.Close()
		fmt.Println("Listening on Port " + command)
		for {
			conn, err := tcpListener.Accept()
			if err != nil {
				fmt.Println("Could not accept connection")
				continue
			}
			go handleClient(conn, lList, command)
		}
	}

}

func main() {
	fmt.Println("Splinter's Cell")
	lList := listenerList{}
	cList := clientList{}
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("SPLINTER >>> ")
		command, err := reader.ReadString('\n') //Returns up to AND including \n
		if err != nil {
			fmt.Println("Error reading input:", err)
			continue
		}

		//Below the madness begins
		//handleListen(command, &lList, &cList)
		handleListen2(command, &lList, &cList)

		command = strings.TrimSpace(command)
		switch command {
		case "listen":
			fmt.Println("Create a listener with: listen -p <port>")
			fmt.Println("List active listeners : listen --ls")
			fmt.Println("Remove a listener with: listen --kill")
		case "listen --ls":
			lList.displayListeners()
		case "listen --kill": //Now it works ?!
			{
				fmt.Print("Specify a listener to kill: ")
				command, err = reader.ReadString('\n')
				if err != nil {
					fmt.Println("Error reading input:", err)
					continue
				}
				lList.delListener(command)
			}
		case "session ls":
			cList.displayClient()
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
