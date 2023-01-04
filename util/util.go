package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/cheggaaa/pb"
	"github.com/google/uuid"
)

func IsUrlValid(url string) (bool, error) {
	urls := []string{
		"www.qobuz.com/",
		"play.qobuz.com/",
		"www.deezer.com/",
		"listen.tidal.com/",
		"tidal.com/browse/",
		"music.youtube.com/",
		"soundcloud.com/",
		"music.apple.com/",
		"open.spotify.com/",
	}

	var contains bool
	for _, p := range urls {
		if strings.Contains(url, p) {
			contains = true
			break
		}
	}

	if !contains {
		return false, nil
	}

	resp, err := http.Get(url)
	if err != nil {
		return false, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, nil
	}

	return true, nil
}

func RequestWithSessionToken(method string, path string, body io.Reader, token string) (*http.Request, error) {
	req, err := http.NewRequest(
		method,
		"https://api.divolt.xyz/"+path,
		body,
	)
	if err != nil {
		return &http.Request{}, err
	}

	req.Header.Add("x-session-token", token)

	return req, nil
}

func UnmarshalResponseBody[T any](resp *http.Response, to *T) error {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, to)
	if err != nil {
		return err
	}

	return nil
}

func DownloadFileFromDescription(description string) (string, error) {
	splitDesc := strings.Split(description, "\n")
	idx := len(splitDesc)
	url := splitDesc[idx-1]
	url = strings.TrimSpace(url)

	path, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	path += "/Downloads"
	_ = os.Mkdir(path, os.ModePerm)

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("status not ok")
	}

	contentLength, _ := strconv.Atoi(resp.Header.Get("Content-Length"))
	newUuid := uuid.NewString()

	filename := fmt.Sprintf("%s/%s.zip", path, newUuid)
	dest, err := os.Create(filename)
	if err != nil {
		return "", err
	}

	defer dest.Close()

	// create bar
	bar := pb.New(contentLength).SetUnits(pb.U_BYTES).SetRefreshRate(time.Millisecond * 10)
	bar.ShowSpeed = true
	bar.Start()

	// create proxy reader
	reader := bar.NewProxyReader(resp.Body)

	// and copy from reader
	io.Copy(dest, reader)
	bar.Finish()

	return filename, nil
}
