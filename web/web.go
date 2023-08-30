package web

import (
	"errors"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/dxbednarczyk/limestone/util"
	"github.com/schollz/closestmatch"
	"github.com/urfave/cli/v2"
	"golang.org/x/exp/maps"
)

type searchResponse struct {
	Tracks struct {
		Items []track `json:"items"`
	} `json:"tracks"`
}

type track struct {
	Name      string `json:"title"`
	Performer struct {
		Name string `json:"name"`
	} `json:"performer"`
	Duration        int  `json:"duration"`
	ParentalWarning bool `json:"parental_warning"`
	ID              int  `json:"id"`
}

func (t track) Title() string {
	var sb strings.Builder

	if t.ParentalWarning {
		sb.WriteString("[E] ")
	}

	sb.WriteString(t.Name)

	return sb.String()
}

func (t track) Description() string { return fmt.Sprintf("%s | %s", t.Performer.Name, t.FormatTime()) }
func (t track) FilterValue() string { return t.Name }

func (t *track) FormatTime() string {
	duration := time.Duration(t.Duration) * time.Second

	minutes := math.Floor(duration.Minutes())
	seconds := math.Floor(duration.Seconds()) - (minutes * 60)

	if seconds == 0 {
		return fmt.Sprintf("%.0f minutes", minutes)
	}

	return fmt.Sprintf("%.0f minutes, %.0f seconds", minutes, seconds)
}

func Query(ctx *cli.Context) (*track, error) {
	query := ctx.Args().First()

	escaped := url.QueryEscape(query)

	resp, err := http.Get("https://slavart.gamesdrive.net/api/search?q=" + escaped)
	if err != nil {
		return nil, err
	}

	var searchData searchResponse

	err = util.UnmarshalResponseBody(resp, &searchData)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if ctx.Bool("closest") {
		tracks := make(map[string]*track)

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
