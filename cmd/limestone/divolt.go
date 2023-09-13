package main

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strings"
	"syscall"

	"github.com/dxbednarczyk/limestone/internal/config"
	"github.com/dxbednarczyk/limestone/internal/divolt"
	"github.com/dxbednarczyk/limestone/internal/download"
	"github.com/urfave/cli/v2"
	"golang.org/x/term"
)

var allowedUrls = []string{
	"qobuz.com",
	"deezer.com",
	"tidal.com",
	"soundcloud.com",
	"open.spotify.com",
	"music.youtube.com",
	"music.apple.com",
}

var divoltdl = cli.Command{
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

		formatted, err := formatURL(ctx.Args().First())
		if err != nil {
			return errors.New("invalid url provided")
		}

		session := divolt.NewSession(&cfg)

		if cfg.Auth.Token != "" {
			slog.Info("Trying existing session token")

			resp, err := session.AuthenticatedRequest(
				divolt.RequestInfo{
					Method: http.MethodGet,
					Path:   "/users/@me",
				},
			)
			if err != nil || 400 <= resp.StatusCode {
				return errors.New("token is invalid, must re-authenticate")
			}

			slog.Info("Token is valid")
		}

		err = divolt.CheckServerStatus(&session)
		if err != nil {
			return err
		}

		err = divolt.SendDownloadMessage(&session, formatted, ctx.Uint("quality"))
		if err != nil {
			return err
		}

		message, err := divolt.GetUploadMessage(&session)
		if err != nil {
			return errors.New("failed to get upload response")
		}

		err = download.FromMessage(ctx, message.Embeds[0].Description)
		if err != nil {
			return errors.New("failed to download bot output")
		}

		return nil
	},
}

var login = cli.Command{
	Name:      "login",
	UsageText: "limestone login <email>",
	Action: func(ctx *cli.Context) error {
		email := ctx.Args().First()
		if email == "" {
			return errors.New("no email specified")
		}

		fmt.Printf("Enter the password for %s: ", email)
		passwordBytes, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return err
		}

		slog.Info("Logging in")

		cfg := config.Config{
			Email:    email,
			Password: string(passwordBytes),
		}

		session := divolt.NewSession(&cfg)
		err = session.Login()
		if err != nil {
			return err
		}

		slog.Info("Login successful")

		err = cfg.CacheLoginDetails()
		if err != nil {
			return err
		}

		slog.Info("Login details cached")

		return nil
	},
}

var logout = cli.Command{
	Name:      "logout",
	UsageText: "limestone logout",
	Action: func(ctx *cli.Context) error {
		slog.Info("Logging out")

		cfg, err := config.GetLoginDetails()
		if err != nil {
			return err
		}

		// naming seems counterintuitive, but we need
		// to authenticate before we can de-authenticate
		session := divolt.NewSession(&cfg)
		err = session.Logout()
		if err != nil {
			return err
		}

		err = config.RemoveConfigDetails()
		if err != nil {
			return err
		}

		slog.Info("Logged out successfully")

		return nil
	},
}

func formatURL(unformatted string) (string, error) {
	var contains bool

	for _, p := range allowedUrls {
		if strings.Contains(unformatted, p) {
			contains = true
			break
		}
	}

	if !contains {
		return "", errors.New("url does not contain one of the valid sources")
	}

	// remove invalid query at end of some urls, especially deezer
	parsed, err := url.Parse(unformatted)
	if err != nil {
		return "", err
	}

	queries := parsed.Query()

	queries.Del("deferredFl")
	queries.Del("utm_campaign")
	queries.Del("utm_source")
	queries.Del("utm_medium")

	return parsed.String(), nil
}
