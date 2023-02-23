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
	HomeDir  string
}

var SessionToken string

func IsUrlValid(url string) bool {
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

	return contains
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

func DownloadFromMessage(description string, path string) error {
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

	uuid := uuid.NewString()
	filename := fmt.Sprintf("%s/%s.zip", path, uuid)

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

	// and copy from reader
	io.Copy(dest, reader)
	bar.Finish()

	fmt.Printf("Downloaded to %s.\n", filename)
	return nil
}

func CacheLoginDetails(config Config) error {
	path := fmt.Sprintf("%s/.config/limestone", config.HomeDir)
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

func readFromCache(config *Config) error {
	path := fmt.Sprintf("%s/.config/limestone/config.json", config.HomeDir)
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	err = json.Unmarshal(content, &config)
	if err != nil {
		return err
	}

	if config.Email == "" || config.Password == "" {
		return errors.New("empty email or password")
	}

	return nil
}

func (config *Config) GetLoginDetails() error {
	err := readFromCache(config)
	if err != nil {
		return err
	}

	if !(config.Email == "" || config.Password == "") {
		config.Cached = true
		return nil
	}

	fmt.Println("** Failed to read from cache, maybe you've never logged in yet **")
	fmt.Println("** If not, delete ~/.config/limestone to regenerate cache next time you log in **")

	var email string
	fmt.Println("Enter your Divolt account's email address:")
	fmt.Scanln(&email)

	var password string
	fmt.Println("Enter your Divolt account's password:")
	fmt.Scanln(&password)

	config.Email = email
	config.Password = password

	return nil
}
