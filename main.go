package main

import (
	"flag"
	"fmt"
	"limestone/routes/auth"
	"limestone/routes/channels"
	"limestone/routes/servers"
	"limestone/util"
	"log"
	"os"
)

func main() {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Failed to find user home directory")
	}

	var path string
	if home != "" {
		path = home + "/Downloads"
	}

	directory := flag.String("d", path, "specify an absolute directory to download files to")
	flag.Parse()

	if *directory == "" {
		log.Fatal("No download path specified.")
	}

	var url string
	fmt.Println("Input the album/track to download:")
	fmt.Scanln(&url)

	valid := util.IsUrlValid(url)
	if !valid {
		log.Fatal("Invalid URL provided.")
	}

	var config util.Config
	err = config.GetLoginDetails()
	if err != nil {
		log.Fatal("Failed to get login details.")
	}

	sesh := auth.NewSession(config.Email, config.Password, "Limestone")
	err = sesh.Login()
	if err != nil {
		log.Fatal("Failed to login.")
	}

	if !config.Cached && home != "" {
		config.HomeDir = home
		err = util.CacheLoginDetails(config)
		if err != nil {
			fmt.Println("Failed to cache login details, you will need to input them again next time.")
		}
	}

	err = servers.CheckServerStatus(&sesh)
	if err != nil {
		sesh.Logout()
		log.Fatal(err)
	}

	id, err := channels.SendDownloadMessage(&sesh, url)
	if err != nil {
		sesh.Logout()
		log.Fatal(err)
	}

	message, err := channels.GetUploadMessage(&sesh, id)
	if err != nil {
		sesh.Logout()
		log.Fatal(err)
	}

	err = util.DownloadFromMessage(message.Embeds[0].Description, *directory)
	if err != nil {
		log.Fatal(err)
	}

	err = sesh.Logout()
	if err != nil {
		fmt.Println("Failed to log out this session, go to your Divolt settings and remove it.")
	}
}
