package main

import (
	"flag"
	"fmt"
	"limestone/routes/auth"
	"limestone/routes/channels"
	"limestone/routes/servers"
	"limestone/util"
	"log"
)

func main() {
	directory := flag.String("dir", "", "specify an absolute directory to download files to, defaults to ~/Downloads")
	flag.Parse()

	var config util.Config
	config, err := util.ReadFromCache()
	if err != nil {
		fmt.Println("** Failed to read from cache, maybe you've never logged in yet **")
		fmt.Println("** If not, delete ~/.config/limestone to regenerate cache next time you log in **")

		util.GetLoginDetails(&config)
	} else {
		fmt.Printf("Logging in as %s\n", config.Email)
		config.Cached = true
	}

	sesh := auth.NewSession(config.Email, config.Password, "Limestone")
	err = sesh.Login()
	if err != nil {
		log.Fatal("Failed to login.")
	}

	if !config.Cached {
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

	var albumUrl string
	fmt.Println("Input the album/track to download:")
	fmt.Scanln(&albumUrl)

	valid, err := util.IsUrlValid(albumUrl)
	if !valid || err != nil {
		sesh.Logout()
		log.Fatal("Invalid URL provided.")
	}

	id, err := channels.SendDownloadMessage(&sesh, albumUrl)
	if err != nil {
		sesh.Logout()
		log.Fatal(err)
	}

	message, err := channels.GetUploadMessage(&sesh, id)
	if err != nil {
		sesh.Logout()
		log.Fatal(err)
	}

	err = util.DownloadFileFromDescription(message.Embeds[0].Description, *directory)
	if err != nil {
		log.Fatal(err)
	}

	err = sesh.Logout()
	if err != nil {
		fmt.Println("Failed to log out this session, go to your Divolt settings and remove it.")
	}
}
