package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
)

func generateImplant(c2Ip string, c2Port string) {
	fmt.Println("\n[!]Generating implant for " + c2Ip + ":" + c2Port)
	unifiedAddress := "c2Address := \"" + c2Ip + ":" + c2Port + "\""
	//fmt.Println("[!]TCP Address: " + unifiedAddress) //Debug Purposes
	content, err := os.ReadFile("cmd/Bushido/main.go")
	if err != nil {
		fmt.Println("[-]Error Reading the cmd/Bushido/main.go")
		return
	}
	strContent := string(content)
	addressRegex := regexp.MustCompile(`c2Address := "(?:\d{1,3}\.){3}\d{1,3}:\d+"`)
	newContent := addressRegex.ReplaceAllString(strContent, unifiedAddress)
	err = os.WriteFile("cmd/Bushido/main.go", []byte(newContent), 0777)
	if err != nil {
		fmt.Println("[-]Error writing to file:", err)
		return
	}
	//Not too elegant and relies on the presence of a script that does it. Will do for now.
	build := exec.Command("./scripts/winBuild.sh")
	build.Run()
	fmt.Println("[+]Implant written to /Bushido/client.exe")
	fmt.Println()
}

func generateImplantDebug(c2Ip string, c2Port string) {
	fmt.Println("\n[!]Generating implant for " + c2Ip + ":" + c2Port)
	unifiedAddress := "c2Address := \"" + c2Ip + ":" + c2Port + "\""
	//fmt.Println("[!]TCP Address: " + unifiedAddress) //Debug Purposes
	content, err := os.ReadFile("cmd/Bushido/main.go")
	if err != nil {
		fmt.Println("[-]Error Reading the cmd/Bushido/main.go")
		return
	}
	strContent := string(content)
	addressRegex := regexp.MustCompile(`c2Address := "(?:\d{1,3}\.){3}\d{1,3}:\d+"`)
	newContent := addressRegex.ReplaceAllString(strContent, unifiedAddress)
	err = os.WriteFile("cmd/Bushido/main.go", []byte(newContent), 0777)
	if err != nil {
		fmt.Println("[-]Error writing to file:", err)
		return
	}
	//Not too elegant and relies on the presence of a script that does it. Will do for now.
	build := exec.Command("./scripts/winBuildDebug.sh")
	build.Run()
	fmt.Println("[+]Implant written to /Bushido/clientDebug.exe")
	fmt.Println()
}
