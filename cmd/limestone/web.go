package main

import (
	"errors"

	"github.com/dxbednarczyk/limestone/internal/download"
	"github.com/dxbednarczyk/limestone/internal/web"
	"github.com/urfave/cli/v2"
)

var webdl = cli.Command{
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
	Action: func(ctx *cli.Context) error {
		if ctx.Args().First() == "" {
			return errors.New("you must provide a query")
		}

		track, err := web.Query(ctx)
		if err != nil {
			return err
		}

		if track == nil {
			return errors.New("no response or result from download request")
		}

		err = download.FromWeb(ctx, track.ID, track.Performer.Name, track.Name)
		if err != nil {
			return err
		}

		return nil
	},
}
