package main

import (
	"errors"
	"fmt"
	"limestone/divolt"
	"limestone/util"
	"os"

	"github.com/urfave/cli/v2"
)

func divoltDownload(ctx *cli.Context, config util.Config) error {
	valid := util.IsUrlValid(ctx.Args().First())
	if !valid {
		return errors.New("invalid url provided")
	}

	sesh := divolt.NewSession(config.Email, config.Password, "Limestone")
	err := sesh.Login()
	if err != nil {
		return errors.New("failed to login")
	}

	if !config.Cached {
		err = util.CacheLoginDetails(config)
		if err != nil {
			fmt.Fprintf(ctx.App.ErrWriter, "Failed to cache login details - %s\n", err)
		}
	}

	err = divolt.CheckServerStatus(&sesh)
	if err != nil {
		sesh.Logout()
		return errors.New("invalid server status")
	}

	id, err := divolt.SendDownloadMessage(&sesh, ctx.Args().First())
	if err != nil {
		sesh.Logout()
		return errors.New("failed to send download request")
	}

	message, err := divolt.GetUploadMessage(&sesh, id)
	if err != nil {
		sesh.Logout()
		return errors.New("failed to get upload response")
	}

	path := ctx.Path("dir")
	if path == "" {
		path, err = os.Getwd()
		if err != nil {
			return errors.New("failed to get working directory")
		}
	}

	err = util.DownloadFromMessage(message.Embeds[0].Description, path)
	if err != nil {
		return errors.New("failed to download bot output")
	}

	err = sesh.Logout()
	if err != nil {
		fmt.Fprintln(ctx.App.ErrWriter, "Failed to log out this session")
	}

	return nil
}
