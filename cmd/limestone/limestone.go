package main

import (
	"io"
	"log/slog"
	"os"

	"github.com/dxbednarczyk/limestone/internal/download"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:    "limestone",
		Version: "0.5.0",
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
			&cli.BoolFlag{Name: "log", Usage: "log errors and download information"},
		},
		Before: func(ctx *cli.Context) error {
			var handler io.Writer

			if ctx.Bool("log") {
				handler = os.Stderr
			} else {
				handler = io.Discard
			}

			slog.SetDefault(slog.New(slog.NewTextHandler(handler, nil)))
			slog.Info("Logging is enabled")

			return nil
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

	err = download.FlushQueue()
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
