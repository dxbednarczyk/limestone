package web

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/charmbracelet/bubbles/list"
	"github.com/schollz/closestmatch"
	"github.com/urfave/cli/v2"
	"golang.org/x/exp/maps"
)

type searchResponse struct {
	Tracks struct {
		Items []Track `json:"items"`
	} `json:"tracks"`
}

func GetTrack(ctx *cli.Context) (*Track, error) {
	query := ctx.Args().First()

	slog.Info("Getting results for query " + query)

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

	return filterTrack(searchData.Tracks.Items, query, ctx.Bool("closest"))
}

func filterTrack(responseItems []Track, query string, closest bool) (*Track, error) {
	if closest {
		return getClosestMatch(responseItems, query)
	}

	// this is extremely stupid.
	items := make([]list.Item, len(responseItems))
	for i := range responseItems {
		items[i] = responseItems[i]
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

func getClosestMatch(responseItems []Track, query string) (*Track, error) {
	tracks := map[string]*Track{}

	for i := range responseItems {
		track := &responseItems[i]

		desc := fmt.Sprintf("%s - %s", track.Performer.Name, track.Name)

		tracks[desc] = track
	}

	bagSizes := []int{2, 3}

	match := closestmatch.New(maps.Keys(tracks), bagSizes)

	closestKey := match.Closest(query)

	return tracks[closestKey], nil
}
