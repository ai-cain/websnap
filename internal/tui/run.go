package tui

import (
	"io"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type Runner struct{}

func NewRunner() Runner {
	return Runner{}
}

func (r Runner) Run(input io.Reader, output io.Writer, studio Studio) error {
	if input == nil {
		input = strings.NewReader("")
	}

	if output == nil {
		output = io.Discard
	}

	_, err := tea.NewProgram(
		NewModel(studio),
		tea.WithInput(input),
		tea.WithOutput(output),
	).Run()

	return err
}

