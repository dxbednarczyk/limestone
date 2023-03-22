package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/cheggaaa/pb"
	"github.com/google/uuid"
	"github.com/urfave/cli"
)

type Config struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Cached   bool
}

var linkRegex = regexp.MustCompile(`((http|https)://)(www.)?[a-zA-Z0-9@:%._\+~#?&//=]{2,256}\.[a-z]{2,6}\b([-a-zA-Z0-9@:%._\+~#?&//=]*)`)

func IsUrlValid(url string) bool {
	urls := []string{
		"qobuz",
		"deezer.com",
		"tidal",
		"soundcloud",
		"spotify",
	}

	var contains bool
	for _, p := range urls {
		if strings.Contains(url, p) {
			contains = true
			break
		}
	}

	return contains && linkRegex.MatchString(url)
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

func DownloadFromMessage(ctx *cli.Context, description string, path string) error {
	splitDesc := strings.Split(description, "\n")
	url := strings.TrimSpace(splitDesc[len(splitDesc)-1])

	err := os.Mkdir(path, os.ModePerm)
	if !os.IsExist(err) {
		return err
	}

	fmt.Fprintln(ctx.App.Writer, "Downloading...")

	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return errors.New("status not ok")
	}

	filename := fmt.Sprintf("%s/%s.zip", path, uuid.NewString())

	dest, err := os.Create(filename)
	if err != nil {
		return err
	}

	defer dest.Close()

	// create bar
	length, err := strconv.Atoi(resp.Header.Get("Content-Length"))
	if err != nil {
		return err
	}

	bar := pb.New(length).SetUnits(pb.U_BYTES).SetRefreshRate(time.Millisecond * 10)
	bar.ShowSpeed = true
	bar.Start()

	// create proxy reader
	reader := bar.NewProxyReader(resp.Body)

	// copy from reader
	io.Copy(dest, reader)
	bar.Finish()

	fmt.Fprintf(ctx.App.Writer, "Downloaded to %s.\n", filename)
	return nil
}

func CacheLoginDetails(config Config) error {
	dir, err := os.UserConfigDir()
	if err != nil {
		return errors.New("failed to get user config directory")
	}

	err = os.MkdirAll(dir+"/limestone", os.ModePerm)
	if err != nil {
		return err
	}

	file_path := dir + "/limestone/config.json"
	var dest *os.File

	_, err = os.Stat(file_path)
	if os.IsNotExist(err) {
		dest, err = os.Create(file_path)
		if err != nil {
			return err
		}

		defer dest.Close()
	} else {
		return err
	}

	encoded, err := json.Marshal(config)
	if err != nil {
		return err
	}

	_, err = dest.Write(encoded)
	if err != nil {
		return err
	}

	dest.Sync()

	return nil
}

func (config *Config) GetLoginDetails() error {
	dir, err := os.UserConfigDir()
	if err != nil {
		return err
	}

	path := dir + "/limestone/config.json"
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	err = json.Unmarshal(content, &config)
	if err != nil {
		return err
	}

	if !(config.Email == "" || config.Password == "") {
		config.Cached = true
	}

	return nil
}
