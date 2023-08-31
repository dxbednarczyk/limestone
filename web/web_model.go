package web

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	docStyle = lipgloss.NewStyle().Margin(1, 2)
	cmd      tea.Cmd
)

type model struct {
	choices list.Model
	choice  Track
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
			item, ok := m.choices.SelectedItem().(Track)
			if ok {
				m.choice = item
			}

			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.choices.SetSize(msg.Width-h, msg.Height-v)
	}

	m.choices, cmd = m.choices.Update(msg)

	return m, cmd
}

func (m model) View() string {
	return docStyle.Render(m.choices.View())
}

func trackModel(tracks []list.Item) (Track, error) {
	initialModel := model{choices: list.New(tracks, list.NewDefaultDelegate(), 0, 0)}

	prog := tea.NewProgram(initialModel, tea.WithAltScreen())

	returnedModel, err := prog.Run()

	return returnedModel.(model).choice, err
}
