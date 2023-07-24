package util

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
)

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

func ValidateURL(u string) (string, error) {
	urls := []string{
		"qobuz.com",
		"deezer.com",
		"tidal.com",
		"soundcloud.com",
		"open.spotify.com",
		"music.youtube.com",
		"music.apple.com",
	}

	var contains bool

	for _, p := range urls {
		if strings.Contains(u, p) {
			contains = true
			break
		}
	}

	if !contains {
		return "", errors.New("url does not contain one of the valid sources")
	}

	// remove invalid query at end of some urls, especially deezer
	parsed, err := url.Parse(u)
	if err != nil {
		return "", err
	}

	queries := parsed.Query()

	queries.Del("deferredFl")
	queries.Del("utm_campaign")
	queries.Del("utm_source")
	queries.Del("utm_medium")

	parsed.RawQuery = queries.Encode()

	return parsed.String(), nil
}
