package download

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dchest/uniuri"
	"github.com/urfave/cli/v2"
)

func DownloadFromMessage(ctx *cli.Context, description string) error {
	fmt.Println("Downloading...")

	splitDesc := strings.Split(description, "\n")
	url := strings.TrimSpace(splitDesc[len(splitDesc)-1])

	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	path := GetDownloadPath(ctx)
	filename := fmt.Sprintf("%s/limestone-%s.zip", path, uniuri.New())

	return DownloadWithProgressBar(resp, path, filename)
}

func DownloadFromWeb(ctx *cli.Context, trackID int, performerName, name string) error {
	fmt.Printf("Downloading %s - %s...\n", performerName, name)

	resp, err := http.Get(fmt.Sprintf("https://slavart-api.gamesdrive.net/api/download/track?id=%d", trackID))
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	path := GetDownloadPath(ctx)
	filename := fmt.Sprintf("%s/%s - %s.flac", path, performerName, name)

	return DownloadWithProgressBar(resp, path, filename)
}

func DownloadWithProgressBar(resp *http.Response, path, absoluteFilename string) error {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}

	dest, err := os.Create(absoluteFilename)
	if err != nil {
		return err
	}

	defer dest.Close()

	var p *tea.Program

	pw := &progressWriter{
		total:       int(resp.ContentLength),
		destination: dest,
		reader:      resp.Body,
		onProgress: func(ratio float64) {
			p.Send(ratio)
		},
	}

	m := model{
		writer:   pw,
		progress: progress.New(progress.WithDefaultGradient()),
	}

	p = tea.NewProgram(m)

	go pw.Start(p)

	if _, err := p.Run(); err != nil {
		return err
	}

	fmt.Println("Downloaded to ", absoluteFilename)

	return nil
}

func GetDownloadPath(ctx *cli.Context) string {
	path := ctx.Path("dir")

	if path == "" {
		home, err := os.UserHomeDir()
		if err == nil {
			path = home + "/Downloads"
		}
	}

	if path == "" {
		wd, err := os.Getwd()
		if err == nil {
			path = wd
		}
	}

	return path
}
