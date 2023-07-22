package main

import (
	"errors"
	"log"

	"github.com/dxbednarczyk/limestone/divolt"
	"github.com/dxbednarczyk/limestone/util"
	"github.com/urfave/cli/v2"
)

func divoltDownload(ctx *cli.Context, config util.Config) error {
	valid := util.IsURLValid(ctx.Args().First())
	if !valid {
		return errors.New("invalid url provided")
	}

	session := divolt.NewSession(config.Email, config.Password, "Limestone")
	err := session.Login()

	defer session.Logout()

	if err != nil {
		return errors.New("failed to login")
	}

	if !config.Cached {
		err = util.CacheLoginDetails(config)
		if err != nil {
			log.Printf("Failed to cache login details: %s\n", err)
		}
	}

	err = divolt.CheckServerStatus(&session)
	if err != nil {
		return errors.New("invalid server status")
	}

	id, err := divolt.SendDownloadMessage(&session, ctx.Args().First())
	if err != nil {
		return errors.New("failed to send download request")
	}

	message, err := divolt.GetUploadMessage(ctx, &session, id)
	if err != nil {
		return errors.New("failed to get upload response")
	}

	path, err := util.GetDownloadPath(ctx)
	if err != nil {
		return err
	}

	err = util.DownloadFromMessage(ctx, message.Embeds[0].Description, path)
	if err != nil {
		return errors.New("failed to download bot output")
	}

	return nil
}
