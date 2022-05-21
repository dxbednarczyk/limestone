package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	a := textinput.NewModel()
	a.Focus()

	s := spinner.NewModel()
	s.Spinner = spinner.Dot

	initialModel := model{
		albumInput:   a,
		spinner:      s,
		pathInput:    textinput.NewModel(),
		currentState: 0,
		helper:       &helper{HTTPClient: http.DefaultClient},
	}

	err := tea.NewProgram(initialModel, tea.WithAltScreen()).Start()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
