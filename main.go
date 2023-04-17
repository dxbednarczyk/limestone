package main

import (
	"errors"
	"fmt"
	"limestone/divolt"
	"limestone/util"
	"os"
	"syscall"

	"github.com/urfave/cli/v2"
	"golang.org/x/term"
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
			{
				Name: "login",
				UsageText: "limestone login <email>",
				Action: func(ctx *cli.Context) error {
					email := ctx.Args().First()
					if email == "" {
						fmt.Fprintln(ctx.App.ErrWriter, "No email specified")
						return errors.New("no email specified")
					}

					fmt.Printf("Enter the password for %s: ", email)
					bp, err := term.ReadPassword(int(syscall.Stdin))
					if err != nil {
						return err
					}
					password := string(bp)

					sesh := divolt.NewSession(email, password, "login test")
					err = sesh.Login()
					if err != nil {
						return err
					}

					fmt.Println("\nLogin test successful.")

					err = sesh.Logout()
					if err != nil {
						fmt.Fprintln(ctx.App.ErrWriter, "Failed to logout")
						return err
					}

					config.Email = email
					config.Password = password
					err = util.CacheLoginDetails(config)
					if err != nil {
						return err
					}

					return nil
				},
			},
		},
	}

	app.Run(os.Args)
}
