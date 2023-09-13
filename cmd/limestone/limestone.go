package main

import (
	"log/slog"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:    "limestone",
		Version: "0.4.1",
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
			&login,
			&logout,
			&divoltdl,
			&webdl,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
