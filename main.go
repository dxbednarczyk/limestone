package main

import (
	"fmt"
	"limestone/routes/auth/session"
	"limestone/routes/channels"
	"limestone/routes/servers"
	"limestone/util"
	"log"
	"time"
)

func main() {
	var email string
	fmt.Println("Enter your Divolt account's email address:")
	fmt.Scanln(&email)

	var password string
	fmt.Println("Enter your Divolt account's password:")
	fmt.Scanln(&password)

	sesh := session.NewSession(email, password, "Limestone")
	err := sesh.Login()
	if err != nil {
		log.Fatal("Failed to login.")
	}

	err = servers.CheckServerStatus(&sesh)
	if err != nil {
		sesh.Logout()
		log.Fatal(err)
	}

	var albumUrl string
	fmt.Println("Input the album/track to download:")
	fmt.Scanln(&albumUrl)

	valid, err := util.IsUrlValid(albumUrl)
	if !valid || err != nil {
		sesh.Logout()
		log.Fatal("Invalid URL provided.")
	}

	err = channels.SendDownloadMessage(&sesh, albumUrl)
	if err != nil {
		sesh.Logout()
		log.Fatal(err)
	}

	fmt.Println("Waiting a bit...")
	time.Sleep(3 * time.Second)

	message, err := channels.GetUploadMessage(&sesh)
	if err != nil {
		sesh.Logout()
		log.Fatal(err)
	}

	err = sesh.Logout()
	if err != nil {
		log.Println("Failed to log out this session, go to your Divolt settings and remove it.")
	}

	path, err := util.DownloadFileFromDescription(message.Embeds[0].Description)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Downloaded to %s.", path)
}
