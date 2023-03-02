package main

import (
	"errors"
	"fmt"
	"limestone/routes/auth"
	"limestone/routes/channels"
	"limestone/routes/servers"
	"limestone/util"
	"os"

	"github.com/urfave/cli/v2"
)

func divoltDownload(ctx *cli.Context, config util.Config) error {
	var err error

	path := ctx.Path("dir")
	if path == "" {
		path, err = os.Getwd()
		if err != nil {
			return errors.New("failed to get working directory")
		}
	}

	valid := util.IsUrlValid(ctx.Args().First())
	if !valid {
		return errors.New("invalid url provided")
	}

	sesh := auth.NewSession(config.Email, config.Password, "Limestone")
	err = sesh.Login()
	if err != nil {
		return errors.New("failed to login")
	}

	if !config.Cached {
		dir, err := os.UserConfigDir()
		if err != nil {
			return errors.New("failed to get user config directory")
		}

		config.ConfigDir = dir

		err = util.CacheLoginDetails(config)
		if err != nil {
			fmt.Fprintln(ctx.App.ErrWriter, "failed to cache login details")
		}
	}

	err = servers.CheckServerStatus(&sesh)
	if err != nil {
		sesh.Logout()
		return errors.New("invalid server status")
	}

	id, err := channels.SendDownloadMessage(&sesh, ctx.Args().First())
	if err != nil {
		sesh.Logout()
		return errors.New("failed to send download request")
	}

	message, err := channels.GetUploadMessage(&sesh, id)
	if err != nil {
		sesh.Logout()
		return errors.New("failed to get upload response")
	}

	err = util.DownloadFromMessage(message.Embeds[0].Description, path)
	if err != nil {
		return errors.New("failed to download bot output")
	}

	err = sesh.Logout()
	if err != nil {
		fmt.Fprintln(ctx.App.ErrWriter, "failed to log out this session")
	}

	return nil
}

func webDownload(ctx *cli.Context) error {
	fmt.Fprintln(ctx.App.ErrWriter, "unimplemented")
	return errors.New("unnimplemented")
}
