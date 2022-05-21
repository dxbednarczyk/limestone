package main

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type GotAlbumList struct {
	err  error
	list list.Model
}

type DownloadedAlbum struct {
	err error
}

type teaAlbum struct {
	title, desc, id string
}

func (i teaAlbum) Title() string       { return i.title }
func (i teaAlbum) Description() string { return i.desc }
func (i teaAlbum) FilterValue() string { return i.title }

func (m *model) setHW(w, h int) {
	m.w, m.h = w, h
}

func (m *model) getAlbumList(name string) tea.Cmd {
	return func() tea.Msg {
		albums, err := m.helper.albumsByName(name)
		if err != nil {
			return GotAlbumList{err: err}
		}

		items := make([]list.Item, 0, len(albums))
		for _, album := range albums {
			ta := teaAlbum{
				title: album.Title,
				desc:  album.Artist.Name,
				id:    album.ID,
			}

			items = append(items, ta)
		}

		albumList := list.New(items, list.NewDefaultDelegate(), 0, 0)
		albumList.Title = "Results"
		albumList.SetSize(m.h, m.h)

		return GotAlbumList{list: albumList}
	}
}

func (m *model) downloadAlbum(album teaAlbum, path string) tea.Cmd {
	return func() tea.Msg {
		err := m.helper.downloadAlbum(album, path)
		if err != nil {
			return DownloadedAlbum{err: err}
		}

		return DownloadedAlbum{}
	}
}
