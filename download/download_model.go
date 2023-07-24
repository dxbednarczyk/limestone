package download

import (
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
)

const maxWidth = 80

type model struct {
	writer   *progressWriter
	progress progress.Model
	err      error
}

func downloadComplete() tea.Cmd {
	return tea.Tick(time.Millisecond*750, func(_ time.Time) tea.Msg {
		return nil
	})
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.progress.Width = msg.Width - 4

		if m.progress.Width > maxWidth {
			m.progress.Width = maxWidth - 4
		}

		return m, nil

	case progressErr:
		m.err = msg.err

		return m, tea.Quit

	case float64:
		var cmds []tea.Cmd

		if msg >= 1.0 {
			cmds = append(cmds, tea.Batch(downloadComplete(), tea.Quit))
		}

		cmds = append(cmds, m.progress.SetPercent(msg))

		return m, tea.Batch(cmds...)

	// FrameMsg is sent when the progress bar wants to animate itself
	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)

		m.progress = progressModel.(progress.Model)

		return m, cmd

	default:
		return m, nil
	}
}

func (m model) View() string {
	if m.err != nil {
		return "Error downloading: " + m.err.Error() + "\n"
	}

	return "\n  " + m.progress.View() + "\n"
}
