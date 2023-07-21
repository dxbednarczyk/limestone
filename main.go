package main

import (
	"errors"
	"fmt"
	"os"
	"syscall"

	"github.com/dxbednarczyk/limestone/divolt"
	"github.com/dxbednarczyk/limestone/util"
	"github.com/urfave/cli/v2"
	"golang.org/x/term"
)

//nolint:funlen
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
						return err
					}

					if !config.Cached || os.IsNotExist(err) {
						return errors.New("please run `limestone login` before downloading. (no login details cached)")
					}

					return nil
				},
				Action: func(ctx *cli.Context) error {
					err := divoltDownload(ctx, config)
					if err != nil {
						return err
					}

					return nil
				},
			},
			{
				Name:      "web",
				UsageText: "limestone web <url>",
				Action: func(ctx *cli.Context) error {
					fmt.Println("web is unimplemented.")
					return nil
				},
			},
			{
				Name:      "login",
				UsageText: "limestone login <email>",
				Action: func(ctx *cli.Context) error {
					email := ctx.Args().First()
					if email == "" {
						return errors.New("no email specified")
					}

					fmt.Printf("Enter the password for %s: ", email)
					bp, err := term.ReadPassword(int(syscall.Stdin))
					if err != nil {
						return err
					}

					fmt.Print("\nLogging in... ")

					sesh := divolt.NewSession(email, string(bp), "login test")
					err = sesh.Login()
					if err != nil {
						return err
					}

					fmt.Println("login successful.")

					config.Email = email
					config.Password = string(bp)
					err = util.CacheLoginDetails(config)
					if err != nil {
						return err
					}

					fmt.Println("Login details cached.")

					sesh.Logout()

					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
