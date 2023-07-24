package divolt

import (
	"errors"
	"fmt"

	"github.com/dxbednarczyk/limestone/download"
	"github.com/dxbednarczyk/limestone/util"
	"github.com/urfave/cli/v2"
)

func Download(ctx *cli.Context, config util.Config) error {
	validated, err := util.ValidateURL(ctx.Args().First())
	if err != nil {
		return errors.New("invalid url provided")
	}

	session := NewSession(config.Email, config.Password, "Limestone")
	err = session.Login()

	defer session.Logout()

	if err != nil {
		return errors.New("failed to login")
	}

	if !config.Cached {
		err = util.CacheLoginDetails(config)
		if err != nil {
			fmt.Printf("Failed to cache login details: %s\n", err)
		}
	}

	err = CheckServerStatus(&session)
	if err != nil {
		return errors.New("invalid server status")
	}

	id, err := SendDownloadMessage(&session, validated, ctx.Uint("quality"))
	if err != nil {
		return errors.New("failed to send download request")
	}

	message, err := GetUploadMessage(ctx, &session, id)
	if err != nil {
		return errors.New("failed to get upload response")
	}

	path, err := download.GetDownloadPath(ctx)
	if err != nil {
		return err
	}

	err = download.DownloadFromMessage(ctx, message.Embeds[0].Description, path)
	if err != nil {
		return errors.New("failed to download bot output")
	}

	return nil
}
