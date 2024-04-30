package main

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
)

func generateImplant(c2Ip string, c2Port string) {
	fmt.Println("Generating implant for " + c2Ip + ":" + c2Port)
	unifiedAddress := c2Ip + ":" + c2Port
	fmt.Println("TCP Address: " + unifiedAddress)
	content, err := os.ReadFile("Bushido/main.go")
	if err != nil {
		fmt.Println("Error Reading the Bushido/main.go")
		return
	}
	strContent := string(content)
	addressRegex := regexp.MustCompile(`(?:\d{1,3}\.){3}\d{1,3}:\d+`)
	newContent := addressRegex.ReplaceAllString(strContent, unifiedAddress)
	err = os.WriteFile("Bushido/main.go", []byte(newContent), 0777)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}
	//Not too elegant and relies on the presence of a script that does it. Will do for now.
	build := exec.Command("./scripts/winBuild.sh")
	build.Run()
	fmt.Println("Implant written to /Bushido/client.exe")
	fmt.Println()
}
