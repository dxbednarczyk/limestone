package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

type helper struct {
	HTTPClient *http.Client
}

type response struct {
	Albums album_parent `json:"albums"`
}

type album_parent struct {
	Items []album `json:"items"`
}

type album struct {
	Artist artist `json:"artist"`
	Title  string `json:"title"`
	ID     string `json:"id"`
}

type artist struct {
	Name string `json:"name"`
}

func (h helper) albumsByName(name string) ([]album, error) {
	url := fmt.Sprintf("https://slavart.gamesdrive.net/api/search?q=%s", name)
	req, err := h.HTTPClient.Get(url)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	var resp response
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}

	return resp.Albums.Items, nil
}

func (h helper) downloadAlbum(album teaAlbum, path string) error {
	url := fmt.Sprintf("https://slavart-api.gamesdrive.net/api/download/album?id=%s", album.id)
	req, err := h.HTTPClient.Get(url)
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(path, 0755); err != nil {
		return err
	}

	filename := fmt.Sprintf("%s/%s.%s.zip", path, album.title, album.id)
	if _, err := os.NewFile(0, filename).Write(body); err != nil {
		return err
	}

	return nil
}
