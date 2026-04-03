package tui

import (
	"context"

	"github.com/ai-cain/websnap/internal/domain"
	tea "github.com/charmbracelet/bubbletea"
)

type targetsLoadedMsg struct {
	targets []domain.LiveTarget
	err     error
}

type tabsLoadedMsg struct {
	target domain.LiveTarget
	tabs   []domain.BrowserTab
	err    error
}

type captureCompletedMsg struct {
	result domain.CaptureResult
	err    error
}

func loadTargetsCmd(studio Studio) tea.Cmd {
	return func() tea.Msg {
		targets, err := studio.ListTargets(context.Background())
		return targetsLoadedMsg{targets: targets, err: err}
	}
}

func loadTabsCmd(studio Studio, target domain.LiveTarget) tea.Cmd {
	return func() tea.Msg {
		tabs, err := studio.ListTabs(context.Background(), target)
		return tabsLoadedMsg{target: target, tabs: tabs, err: err}
	}
}

func submitLiveCaptureCmd(studio Studio, req domain.LiveCaptureRequest) tea.Cmd {
	return func() tea.Msg {
		result, err := studio.CaptureLive(context.Background(), req)
		return captureCompletedMsg{result: result, err: err}
	}
}
