package main

import (
	"errors"
	"fmt"
	"limestone/util"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	var config util.Config

	app := &cli.App{
		Name:    "limestone",
		Version: "0.2.0",
		Authors: []*cli.Author{
			{
				Name:  "Damian Bednarczyk",
				Email: "me@dxbednarczyk.com",
			},
		},
		Usage:     "Unofficial Slav Art CLI",
		UsageText: "limestone [divolt | web] [... args] <url>",
		Flags: []cli.Flag{
			&cli.PathFlag{Name: "dir"},
		},
		Commands: []*cli.Command{
			{
				Name:      "divolt",
				UsageText: "limestone divolt [... args] <url>",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "email"},
					&cli.StringFlag{Name: "pass"},
				},
				Before: func(ctx *cli.Context) error {
					err := config.GetLoginDetails()
					if err != nil && !os.IsNotExist(err) {
						fmt.Fprintln(ctx.App.ErrWriter, err.Error())
						return err
					}

					if config.Cached {
						return nil
					}

					if ctx.String("email") == "" || ctx.String("pass") == "" {
						fmt.Fprintln(ctx.App.ErrWriter, "No email or password specified")
						return errors.New("no email or password specified")
					}

					config.Email = ctx.String("email")
					config.Password = ctx.String("pass")

					return nil
				},
				Action: func(ctx *cli.Context) error {
					err := divoltDownload(ctx, config)
					if err != nil {
						fmt.Fprintln(ctx.App.ErrWriter, err.Error())
						return err
					}

					return nil
				},
			},
			{
				Name:      "web",
				UsageText: "limestone web <url>",
				Action: func(ctx *cli.Context) error {
					err := webDownload(ctx)
					if err != nil {
						fmt.Fprintln(ctx.App.ErrWriter, err.Error())
						return err
					}

					return nil
				},
			},
		},
	}

	app.Run(os.Args)
}
