package divolt

import (
	"errors"
	"net/url"
	"strings"

	util "github.com/dxbednarczyk/limestone/config"
	"github.com/dxbednarczyk/limestone/download"
	"github.com/urfave/cli/v2"
)

func Download(ctx *cli.Context, config util.Config) error {
	validated, err := validateURL(ctx.Args().First())
	if err != nil {
		return errors.New("invalid url provided")
	}

	session := NewSession(config.Email, config.Password)
	err = session.Login()

	defer session.Logout()

	if err != nil {
		return errors.New("failed to login")
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

	err = download.DownloadFromMessage(ctx, message.Embeds[0].Description)
	if err != nil {
		return errors.New("failed to download bot output")
	}

	return nil
}

func validateURL(u string) (string, error) {
	urls := []string{
		"qobuz.com",
		"deezer.com",
		"tidal.com",
		"soundcloud.com",
		"open.spotify.com",
		"music.youtube.com",
		"music.apple.com",
	}

	var contains bool

	for _, p := range urls {
		if strings.Contains(u, p) {
			contains = true
			break
		}
	}

	if !contains {
		return "", errors.New("url does not contain one of the valid sources")
	}

	// remove invalid query at end of some urls, especially deezer
	parsed, err := url.Parse(u)
	if err != nil {
		return "", err
	}

	queries := parsed.Query()

	queries.Del("deferredFl")
	queries.Del("utm_campaign")
	queries.Del("utm_source")
	queries.Del("utm_medium")

	parsed.RawQuery = queries.Encode()

	return parsed.String(), nil
}
