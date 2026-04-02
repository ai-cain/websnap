package tui

import (
	"context"

	"github.com/ai-cain/websnap/internal/domain"
	tea "github.com/charmbracelet/bubbletea"
)

type captureCompletedMsg struct {
	result domain.CaptureResult
	err    error
}

func submitCaptureCmd(runner ShotRunner, req domain.CaptureRequest) tea.Cmd {
	return func() tea.Msg {
		result, err := runner.Execute(context.Background(), req)
		return captureCompletedMsg{
			result: result,
			err:    err,
		}
	}
}
