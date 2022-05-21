package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	albumInput textinput.Model
	spinner    spinner.Model
	albumList  list.Model
	pathInput  textinput.Model

	helper *helper

	w            int
	h            int
	currentState int

	err error
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			switch m.currentState {
			case 0:
				query := strings.TrimSpace(m.albumInput.Value())
				if query != "" {
					m.currentState++

					return m, tea.Batch(
						spinner.Tick,
						m.getAlbumList(query),
					)
				}
			case 2:
				m.currentState++
				m.pathInput.Focus()
				return m, nil
			case 3:
				path := strings.TrimSpace(m.pathInput.Value())
				if path != "" {
					selected_album := m.albumList.SelectedItem().(teaAlbum)
					m.currentState++

					return m, tea.Batch(
						spinner.Tick,
						m.downloadAlbum(selected_album, path),
					)
				}
			}
		}
	case GotAlbumList:
		if err := msg.err; err != nil {
			m.err = err
			return m, nil
		}

		m.albumList = msg.list
		m.currentState++
	case DownloadedAlbum:
		if err := msg.err; err != nil {
			m.err = err
			return m, nil
		}

		m.currentState++
	case tea.WindowSizeMsg:
		m.setHW(msg.Width, msg.Height)
	}

	var cmd tea.Cmd
	switch m.currentState {
	case 0:
		m.albumInput, cmd = m.albumInput.Update(msg)
		return m, cmd
	case 1, 4:
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case 2:
		m.albumList, cmd = m.albumList.Update(msg)
		return m, cmd
	case 3:
		m.pathInput, cmd = m.pathInput.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m model) View() string {
	if m.err != nil {
		return fmt.Sprintf("Something went wrong: %s\n", m.err)
	}

	switch m.currentState {
	case 0:
		return fmt.Sprintf("Enter an album or artist to search for:\n%s", m.albumInput.View())
	case 1:
		return fmt.Sprintf("%s Fetching albums...", m.spinner.View())
	case 2:
		return m.albumList.View()
	case 3:
		return fmt.Sprintf("Enter a path to download to:\n%s", m.pathInput.View())
	case 4:
		return fmt.Sprintf("%s Downloading album (this may take a while, hold on!)...", m.spinner.View())
	case 5:
		return "Done! You may now exit this program.\n\nPress ctrl+c to exit..."
	}

	return "Press ctrl+c to exit..."
}
