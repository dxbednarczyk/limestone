package download

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/uuid"
	"github.com/urfave/cli/v2"
)

func DownloadFromMessage(ctx *cli.Context, description string, path string) error {
	splitDesc := strings.Split(description, "\n")
	url := strings.TrimSpace(splitDesc[len(splitDesc)-1])

	err := os.Mkdir(path, os.ModePerm)
	if !os.IsExist(err) {
		return err
	}

	fmt.Println("Downloading...")

	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return errors.New("status not ok")
	}

	filename := fmt.Sprintf("%s/%s.zip", path, uuid.NewString())

	return DownloadWithProgressBar(resp, filename)
}

func DownloadWithProgressBar(resp *http.Response, filename string) error {
	dest, err := os.Create(filename)
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

	fmt.Println("Downloaded to ", filename)

	return nil
}

func GetDownloadPath(ctx *cli.Context) (string, error) {
	var path string
	var err error

	path = ctx.Path("dir")
	if path == "" {
		path, err = os.Getwd()
		if err != nil {
			return "", errors.New("failed to get working directory")
		}
	}

	return path, err
}
