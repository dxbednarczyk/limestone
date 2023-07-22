package main

import (
	"fmt"
	"math"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/dxbednarczyk/limestone/util"
	"github.com/urfave/cli/v2"
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

func webDownload(ctx *cli.Context) error {
	escaped := url.QueryEscape(ctx.Args().First())

	fmt.Printf(`Getting results for query "%s"...`, escaped)

	resp, err := http.Get("https://slavart.gamesdrive.net/api/search?q=" + escaped)
	if err != nil {
		return err
	}

	var searchData searchResponse

	err = util.UnmarshalResponseBody(resp, &searchData)
	if err != nil {
		return err
	}

	resp.Body.Close()

	// this is extremely stupid.
	items := make([]list.Item, len(searchData.Tracks.Items))
	for i := range searchData.Tracks.Items {
		items[i] = searchData.Tracks.Items[i]
	}

	choice, err := trackModel(items)
	if err != nil {
		return err
	}

	resp, err = http.Get(fmt.Sprintf("https://slavart-api.gamesdrive.net/api/download/track?id=%d", choice.ID))
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	path, err := util.GetDownloadPath(ctx)
	if err != nil {
		return err
	}

	filename := fmt.Sprintf("%s/%s - %s.flac", path, choice.Performer.Name, choice.Name)

	err = util.DownloadWithProgressBar(resp, filename)

	return err
}
