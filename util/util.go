package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/cheggaaa/pb"
	"github.com/google/uuid"
)

type Config struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Cached   bool
}

var homeDir, _ = os.UserHomeDir()
var SessionToken string

func IsUrlValid(url string) (bool, error) {
	urls := []string{
		"qobuz.com",
		"deezer.com",
		"listen.tidal.com",
		"tidal.com/browse",
		"music.youtube.com",
		"soundcloud.com",
		"music.apple.com",
		"open.spotify.com",
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

func AuthenticatedRequest(method string, path string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(
		method,
		fmt.Sprintf("https://api.divolt.xyz/%s", path),
		body,
	)
	if err != nil {
		return &http.Request{}, err
	}

	req.Header.Add("x-session-token", SessionToken)

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

func DownloadFileFromDescription(description string, path string) error {
	splitDesc := strings.Split(description, "\n")
	idx := len(splitDesc)
	url := strings.TrimSpace(splitDesc[idx-1])

	var err error
	if path == "" {
		path = fmt.Sprintf("%s/Downloads", homeDir)
	}

	err = os.Mkdir(path, os.ModePerm)
	if !os.IsExist(err) {
		return err
	}

	fmt.Println("Downloading...")

	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode > http.StatusPartialContent {
		return errors.New("status not ok")
	}

	newUuid := uuid.NewString()
	filename := fmt.Sprintf("%s/%s.zip", path, newUuid)

	dest, err := os.Create(filename)
	if err != nil {
		return err
	}

	defer dest.Close()

	// create bar
	contentLength, err := strconv.Atoi(resp.Header.Get("Content-Length"))
	if err != nil {
		return err
	}

	bar := pb.New(contentLength).SetUnits(pb.U_BYTES).SetRefreshRate(time.Millisecond * 10)
	bar.ShowSpeed = true
	bar.Start()

	// create proxy reader
	reader := bar.NewProxyReader(resp.Body)

	// and copy from reader
	io.Copy(dest, reader)
	bar.Finish()

	fmt.Printf("Downloaded to %s.\n", filename)
	return nil
}

func CacheLoginDetails(config Config) error {
	path := fmt.Sprintf("%s/.config/limestone", homeDir)
	file_path := fmt.Sprintf("%s/config.json", path)

	_, err := os.Stat(file_path)
	if !os.IsNotExist(err) {
		return err
	}

	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}

	dest, err := os.Create(file_path)
	if err != nil {
		return err
	}

	defer dest.Close()

	encoded, err := json.Marshal(config)
	if err != nil {
		return err
	}

	b, err := dest.Write(encoded)
	if err != nil {
		return err
	}
	dest.Sync()

	fmt.Printf("Wrote %d bytes to config.json\n", b)

	return nil
}

func ReadFromCache() (Config, error) {
	path := fmt.Sprintf("%s/.config/limestone/config.json", homeDir)
	content, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	var config Config
	err = json.Unmarshal(content, &config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}

func GetLoginDetails(config *Config) {
	var email string
	fmt.Println("Enter your Divolt account's email address:")
	fmt.Scanln(&email)

	var password string
	fmt.Println("Enter your Divolt account's password:")
	fmt.Scanln(&password)

	config.Email = email
	config.Password = password
}
