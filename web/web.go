package web

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/charmbracelet/bubbles/list"
	"github.com/dxbednarczyk/limestone/download"
	"github.com/schollz/closestmatch"
	"github.com/urfave/cli/v2"
	"golang.org/x/exp/maps"
)

type searchResponse struct {
	Tracks struct {
		Items []Track `json:"items"`
	} `json:"tracks"`
}

var Web = cli.Command{
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
		track, err := Query(ctx)
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

func Query(ctx *cli.Context) (*Track, error) {
	query := ctx.Args().First()

	fmt.Printf("Getting results for query '%s'...\n", query)

	escaped := url.QueryEscape(query)

	resp, err := http.Get("https://slavart.gamesdrive.net/api/search?q=" + escaped)
	if err != nil {
		return nil, err
	}

	var searchData searchResponse

	err = json.NewDecoder(resp.Body).Decode(&searchData)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if ctx.Bool("closest") {
		tracks := make(map[string]*Track)

		for i := range searchData.Tracks.Items {
			track := &searchData.Tracks.Items[i]

			desc := fmt.Sprintf("%s - %s", track.Performer.Name, track.Name)

			tracks[desc] = track
		}

		bagSizes := []int{2, 3}
		cm := closestmatch.New(maps.Keys(tracks), bagSizes)

		closest := cm.Closest(query)

		return tracks[closest], nil
	}

	// this is extremely stupid.
	items := make([]list.Item, len(searchData.Tracks.Items))
	for i := range searchData.Tracks.Items {
		items[i] = searchData.Tracks.Items[i]
	}

	choice, err := trackModel(items)
	if err != nil {
		return nil, err
	}

	if choice.ID == 0 {
		return nil, errors.New("no choice selected")
	}

	return &choice, nil
}
