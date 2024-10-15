package main

import (
	"Demolitron/cmd"
	"fmt"
	"net/http"
	"os"
	"strings"
)

var notifURL string = ""

func main() {
	//Designed to work with ntfy service for push notifications to the phone :D
	if len(os.Args) == 2 {
		notificationURI := os.Args[1]
		//fmt.Println("[+]Notifications enabled! Check your device")
		notifURL = "https://ntfy.sh/" + notificationURI
		req, _ := http.NewRequest("POST", notifURL,
			strings.NewReader("Demolitron is armed with notifications enabled !"))
		req.Header.Set("Title", "Samurai !")
		req.Header.Set("Tags", "warning,skull")
		http.DefaultClient.Do(req)
		fmt.Println("[!]Demolitron is armed with notifications enabled !")
	}
	//Keep in mind exported functions should be uppercase ya 7mar
	cmd.InitServer(notifURL)
}
