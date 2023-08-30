package main

import (
	"errors"
	"fmt"
	"os"
	"syscall"

	"github.com/dxbednarczyk/limestone/divolt"
	"github.com/dxbednarczyk/limestone/util"
	"github.com/dxbednarczyk/limestone/web"
	"github.com/urfave/cli/v2"
	"golang.org/x/term"
)

//nolint:funlen
func main() {
	app := &cli.App{
		Name:    "limestone",
		Version: "0.3.3",
		Authors: []*cli.Author{
			{
				Name:  "Damian Bednarczyk",
				Email: "me@dxbednarczyk.com",
			},
		},
		Usage: "Unofficial Slav Art CLI",
		UsageText: `limestone [divolt | web] [... args] <url>
See the FAQ at https://rentry.org/slavart`,
		Flags: []cli.Flag{
			&cli.PathFlag{Name: "dir", Usage: "directory to save downloaded music to"},
		},
		Commands: []*cli.Command{
			{
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
					var config util.Config
					err := config.GetLoginDetails()

					if os.IsNotExist(err) {
						return errors.New("please authenticate using `limestone login`")
					}

					if err != nil {
						return err
					}

					err = divolt.Download(ctx, config)
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
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "closest",
						Aliases: []string{"c"},
						Usage:   "download the closest match to the query",
					},
				},
				Before: func(ctx *cli.Context) error {
					if ctx.Args().First() == "" {
						return errors.New("you must provide a query")
					}

					return nil
				},
				Action: func(ctx *cli.Context) error {
					fmt.Printf(`Getting results for query "%s"...%s`, ctx.Args().First(), "\n")

					track, err := web.Query(ctx)
					if err != nil {
						fmt.Println()
						return err
					}

					if track == nil {
						return errors.New("no response or result from download request")
					}

					fmt.Printf("Downloading %s - %s...\n", track.Performer.Name, track.Name)

					err = web.Download(ctx, track)
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

					session := divolt.NewSession(email, string(passwordBytes))
					err = session.Login()
					if err != nil {
						fmt.Println()
						return err
					}

					fmt.Println("login successful.")

					var config util.Config

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
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
