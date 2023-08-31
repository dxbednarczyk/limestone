package main

import (
	"errors"
	"fmt"
	"os"
	"syscall"

	"github.com/dxbednarczyk/limestone/config"
	"github.com/dxbednarczyk/limestone/divolt"
	"github.com/dxbednarczyk/limestone/download"
	"github.com/dxbednarczyk/limestone/web"
	"github.com/urfave/cli/v2"
	"golang.org/x/term"
)

//nolint:funlen
func main() {
	app := &cli.App{
		Name:    "limestone",
		Version: "0.4.0",
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
					cfg, err := config.GetLoginDetails()

					if os.IsNotExist(err) {
						return errors.New("please authenticate using `limestone login`")
					}

					if err != nil {
						return err
					}

					err = divolt.Download(ctx, &cfg)
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
					fmt.Printf("Getting results for query '%s'...\n", ctx.Args().First())

					track, err := web.Query(ctx)
					if err != nil {
						fmt.Println()
						return err
					}

					if track == nil {
						return errors.New("no response or result from download request")
					}

					err = download.DownloadFromWeb(ctx, track.ID, track.Performer.Name, track.Name)
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

					cfg := config.Config{
						Email:    email,
						Password: string(passwordBytes),
					}

					session := divolt.NewSession(&cfg)
					err = session.Login()
					if err != nil {
						fmt.Println()
						return err
					}

					fmt.Println("login successful.")

					err = config.CacheLoginDetails(cfg)
					if err != nil {
						return err
					}

					fmt.Println("Login details cached.")

					return nil
				},
			},
			{
				Name:      "logout",
				UsageText: "limestone logout",
				Action: func(ctx *cli.Context) error {
					fmt.Print("Logging out... ")

					cfg, err := config.GetLoginDetails()
					if err != nil {
						return err
					}

					// naming seems counterintuitive, but we obviously need
					// to authenticate before we can deauthenticate
					session := divolt.NewSession(&cfg)
					err = session.Logout()
					if err != nil {
						return err
					}

					err = config.RemoveConfigDetails()
					if err != nil {
						return err
					}

					fmt.Println("logged out successfully.")

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
