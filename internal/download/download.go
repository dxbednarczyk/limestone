package download

import (
	"container/list"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dchest/uniuri"
	"github.com/urfave/cli/v2"
)

var queue = list.New()

type Download struct {
	resp 	*http.Response
	dest 	*os.File
}

func FromMessage(ctx *cli.Context, message string) error {
	slog.Info("Downloading from message")

	splitMessage := strings.Split(message, "\n")

	last := splitMessage[len(splitMessage) - 1]

	url := strings.TrimSpace(last)

	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	path := GetDownloadPath(ctx)
	filename := fmt.Sprintf("%s/limestone-%s.zip", path, uniuri.New())
	
	dest, err := createOutputFile(path, filename)
	if err != nil {
		return err
	}

	defer dest.Close()

	download := Download {
		resp,
		dest,
	}

	return addToQueue(download)
}

func FromWeb(ctx *cli.Context, trackID int, performerName, trackName string) error {
	slog.Info("Downloading from web")

	resp, err := http.Get(fmt.Sprintf("https://slavart-api.gamesdrive.net/api/download/track?id=%d", trackID))
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	path := GetDownloadPath(ctx)
	filename := fmt.Sprintf("%s/%s - %s.flac", path, performerName, trackName)

	dest, err := createOutputFile(path, filename)
	if err != nil {
		return err
	}

	defer dest.Close()

	download := Download{
		resp,
		dest,
	}

	return addToQueue(download)
}

func createOutputFile(path, filename string) (*os.File, error) {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return nil, err
	}

	return os.Create(filename)
}

func addToQueue(download Download) error {
	var prog *tea.Program

	writer := &progressWriter{
		total:       int(download.resp.ContentLength),
		destination: download.dest,
		reader:      download.resp.Body,
		onProgress: func(ratio float64) {
			prog.Send(ratio)
		},
	}

	m := model{
		writer:   writer,
		progress: progress.New(progress.WithDefaultGradient()),
	}

	queue.PushBack(m)

	return nil
}

func FlushQueue() error {
	for e := queue.Front(); e != nil; e = e.Next() {
		model := e.Value.(model)

		prog := tea.NewProgram(model)

		go model.writer.Start(prog)

		if _, err := prog.Run(); err != nil {
			return err
		}

		slog.Info("Downloaded to " + model.writer.destination.Name())
	}

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
