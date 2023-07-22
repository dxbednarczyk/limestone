package main

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type model struct {
	choices list.Model
	choice  track
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			item, ok := m.choices.SelectedItem().(track)
			if ok {
				m.choice = item
			}

			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.choices.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.choices, cmd = m.choices.Update(msg)

	return m, cmd
}

func (m model) View() string {
	return docStyle.Render(m.choices.View())
}

func trackModel(tracks []list.Item) (track, error) {
	initialModel := model{choices: list.New(tracks, list.NewDefaultDelegate(), 0, 0)}

	prog := tea.NewProgram(initialModel, tea.WithAltScreen())

	returnedModel, err := prog.Run()

	return returnedModel.(model).choice, err
}
