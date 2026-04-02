package tui

import (
	"context"
	"io"
	"strings"

	"github.com/ai-cain/websnap/internal/domain"
	tea "github.com/charmbracelet/bubbletea"
)

type ShotRunner interface {
	Execute(ctx context.Context, req domain.CaptureRequest) (domain.CaptureResult, error)
}

type Runner struct{}

func NewRunner() Runner {
	return Runner{}
}

func (r Runner) Run(input io.Reader, output io.Writer, runner ShotRunner) error {
	if input == nil {
		input = strings.NewReader("")
	}

	if output == nil {
		output = io.Discard
	}

	_, err := tea.NewProgram(
		NewModel(runner),
		tea.WithInput(input),
		tea.WithOutput(output),
	).Run()

	return err
}
