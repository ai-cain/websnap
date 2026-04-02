package tui

import (
	"context"
	"testing"

	"github.com/ai-cain/websnap/internal/domain"
	tea "github.com/charmbracelet/bubbletea"
)

func TestModelTabMovesFocus(t *testing.T) {
	t.Parallel()

	model := NewModel(&fakeShotRunner{})
	next, _ := model.Update(tea.KeyMsg{Type: tea.KeyTab})
	got := next.(Model)

	if got.focus != fieldWidth {
		t.Fatalf("focus = %d, want %d", got.focus, fieldWidth)
	}
}

func TestModelCanCycleCaptureMode(t *testing.T) {
	t.Parallel()

	model := NewModel(&fakeShotRunner{})
	model.setFocus(fieldMode)

	next, _ := model.Update(tea.KeyMsg{Type: tea.KeyRight})
	selector := next.(Model)
	if selector.mode != modeSelector {
		t.Fatalf("mode = %v, want %v", selector.mode, modeSelector)
	}

	next, _ = selector.Update(tea.KeyMsg{Type: tea.KeyRight})
	fullPage := next.(Model)
	if fullPage.mode != modeFullPage {
		t.Fatalf("mode = %v, want %v", fullPage.mode, modeFullPage)
	}
}

func TestModelEnterOnLastFieldTransitionsToSuccess(t *testing.T) {
	t.Parallel()

	runner := &fakeShotRunner{
		result: domain.CaptureResult{Path: "C:/captures/home.png", Width: 1440, Height: 900},
	}

	model := NewModel(runner)
	model.inputs[inputIndex(fieldURL)].SetValue("https://example.com")
	model.inputs[inputIndex(fieldWidth)].SetValue("1440")
	model.inputs[inputIndex(fieldHeight)].SetValue("900")
	model.inputs[inputIndex(fieldOut)].SetValue("captures/home.png")
	model.setFocus(fieldOut)

	next, _ := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	capturing := next.(Model)

	if capturing.phase != phaseCapturing {
		t.Fatalf("phase = %d, want %d", capturing.phase, phaseCapturing)
	}

	msg := submitCaptureCmd(runner, domain.CaptureRequest{
		URL:    "https://example.com",
		Width:  1440,
		Height: 900,
		Out:    "captures/home.png",
	})()

	final, _ := capturing.Update(msg)
	success := final.(Model)

	if success.phase != phaseSuccess {
		t.Fatalf("phase = %d, want %d", success.phase, phaseSuccess)
	}

	if success.lastPath != "C:/captures/home.png" {
		t.Fatalf("lastPath = %q, want %q", success.lastPath, "C:/captures/home.png")
	}
}

type fakeShotRunner struct {
	result domain.CaptureResult
	err    error
}

func (f *fakeShotRunner) Execute(_ context.Context, _ domain.CaptureRequest) (domain.CaptureResult, error) {
	return f.result, f.err
}
