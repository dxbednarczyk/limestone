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
		Version: "0.3.0",
		Authors: []*cli.Author{
			{
				Name:  "Damian Bednarczyk",
				Email: "me@dxbednarczyk.com",
			},
		},
		Usage:     "Unofficial Slav Art CLI",
		UsageText: "limestone [divolt | web] [... args] <url>",
		Flags: []cli.Flag{
			&cli.PathFlag{Name: "dir", Usage: "directory to save downloaded music to"},
		},
		Commands: []*cli.Command{
			{
				Name: "divolt",
				UsageText: `limestone divolt [... args] <url>
		
				You can download individual tracks or full albums using Divolt.`,
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "email"},
					&cli.StringFlag{Name: "pass"},
				},
				Action: func(ctx *cli.Context) error {
					err := config.GetLoginDetails()
					if err != nil && !os.IsNotExist(err) {
						return err
					}

					if !config.Cached || os.IsNotExist(err) {
						return errors.New("please run `limestone login` before downloading. (no login details cached)")
					}

					err = divoltDownload(ctx, config)
					if err != nil {
						return err
					}

					return nil
				},
			},
			{
				Name: "web",
				UsageText: `limestone web <query>
				
				You can only download individual tracks from Qobuz using the web download method.`,
				Before: func(ctx *cli.Context) error {
					if ctx.Args().First() == "" {
						return errors.New("you must provide a query")
					}

					return nil
				},
				Action: func(ctx *cli.Context) error {
					err := webDownload(ctx)
					if err != nil {
						return err
					}

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
					passwordBytes, err := term.ReadPassword(int(syscall.Stdin))
					if err != nil {
						return err
					}

					fmt.Print("\nLogging in... ")

					session := divolt.NewSession(email, string(passwordBytes), "login test")
					err = session.Login()
					if err != nil {
						return err
					}

					fmt.Println("login successful.")

					config.Email = email
					config.Password = string(passwordBytes)
					err = util.CacheLoginDetails(config)
					if err != nil {
						return err
					}

					fmt.Println("Login details cached.")

					session.Logout()

					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}
