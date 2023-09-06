package divolt

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/dxbednarczyk/limestone/config"
	"github.com/dxbednarczyk/limestone/download"
	"github.com/urfave/cli/v2"
)

var Divolt = cli.Command{
	Name: "divolt",
	UsageText: `limestone divolt [... args] <url>

You can download individual tracks or full albums using Divolt.`,
	Flags: []cli.Flag{
		&cli.UintFlag{
			Name:        "quality",
			Aliases:     []string{"q"},
			Usage:       "specify a number in the range 0-4",
			Value:       999,
			DefaultText: "highest available",
		},
	},
	Action: func(ctx *cli.Context) error {
		cfg, err := config.GetLoginDetails()

		if os.IsNotExist(err) {
			return errors.New("please authenticate using `limestone login`")
		}

		if err != nil {
			return err
		}

		validated, err := validateURL(ctx.Args().First())
		if err != nil {
			return errors.New("invalid url provided")
		}

		session := NewSession(&cfg)

		if cfg.Auth.Token != "" {
			fmt.Print("Trying existing session token... ")

			resp, err := session.AuthenticatedRequest(
				requestInfo{
					method: http.MethodGet,
					path:   "/users/@me",
				},
			)
			if err != nil || 400 <= resp.StatusCode {
				fmt.Println("token is invalid.")
				fmt.Println("Creating new session...")

				err := session.Login()
				if err != nil {
					return errors.New("failed to login")
				}

				config.CacheLoginDetails(session.Config)
			} else {
				fmt.Println("token is valid.")
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

		err = download.DownloadFromMessage(ctx, message.Embeds[0].Description)
		if err != nil {
			return errors.New("failed to download bot output")
		}

		return nil
	},
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
