package main

import (
	"errors"
	"log"
	"os"

	"github.com/dxbednarczyk/limestone/divolt"
	"github.com/dxbednarczyk/limestone/util"

	"github.com/urfave/cli/v2"
)

func divoltDownload(ctx *cli.Context, config util.Config) error {
	valid := util.IsURLValid(ctx.Args().First())
	if !valid {
		return errors.New("invalid url provided")
	}

	sesh := divolt.NewSession(config.Email, config.Password, "Limestone")
	err := sesh.Login()

	defer sesh.Logout()

	if err != nil {
		return errors.New("failed to login")
	}

	if !config.Cached {
		err = util.CacheLoginDetails(config)
		if err != nil {
			log.Printf("Failed to cache login details: %s\n", err)
		}
	}

	err = divolt.CheckServerStatus(&sesh)
	if err != nil {
		return errors.New("invalid server status")
	}

	id, err := divolt.SendDownloadMessage(&sesh, ctx.Args().First())
	if err != nil {
		return errors.New("failed to send download request")
	}

	message, err := divolt.GetUploadMessage(ctx, &sesh, id)
	if err != nil {
		return errors.New("failed to get upload response")
	}

	path := ctx.Path("dir")
	if path == "" {
		path, err = os.Getwd()
		if err != nil {
			return errors.New("failed to get working directory")
		}
	}

	err = util.DownloadFromMessage(ctx, message.Embeds[0].Description, path)
	if err != nil {
		return errors.New("failed to download bot output")
	}

	return nil
}
