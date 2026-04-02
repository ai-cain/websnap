package tui

import (
	"strconv"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case spinner.TickMsg:
		if m.phase != phaseCapturing {
			return m, nil
		}

		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case captureCompletedMsg:
		if msg.err != nil {
			m.phase = phaseEditing
			m.lastErr = msg.err
			return m, nil
		}

		m.setSuccess(msg.result)
		return m, nil
	case tea.KeyMsg:
		return m.handleKey(msg)
	default:
		return m.updateFocusedInput(msg)
	}
}

func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch m.phase {
	case phaseCapturing:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

		return m, nil
	case phaseSuccess:
		switch msg.String() {
		case "enter", "esc", "q", "ctrl+c":
			return m, tea.Quit
		}

		return m, nil
	}

	switch msg.String() {
	case "ctrl+c", "esc":
		return m, tea.Quit
	case "tab", "shift+tab", "up", "down":
		m.moveFocus(msg.String())
		return m, nil
	case "left", "h":
		if m.focus == fieldMode {
			m.mode = m.mode.previous()
			return m, nil
		}
	case "right", "l":
		if m.focus == fieldMode {
			m.mode = m.mode.next()
			return m, nil
		}
	case "enter":
		return m.submitOrAdvance()
	}

	return m.updateFocusedInput(msg)
}

func (m *Model) moveFocus(key string) {
	next := m.focus
	if key == "shift+tab" || key == "up" {
		next--
	} else {
		next++
	}

	next = m.normalizeField(next)
	for !m.isFocusableField(next) {
		if key == "shift+tab" || key == "up" {
			next = m.normalizeField(next - 1)
			continue
		}

		next = m.normalizeField(next + 1)
	}

	m.setFocus(next)
}

func (m Model) submitOrAdvance() (tea.Model, tea.Cmd) {
	if m.focus < fieldOut {
		next := m.focus + 1
		next = m.normalizeField(next)
		for !m.isFocusableField(next) {
			next = m.normalizeField(next + 1)
		}

		m.setFocus(next)
		return m, nil
	}

	req, err := m.buildRequest()
	if err != nil {
		m.lastErr = friendlyInputError(err)
		return m, nil
	}

	m.phase = phaseCapturing
	m.lastErr = nil
	return m, tea.Batch(m.spinner.Tick, submitCaptureCmd(m.runner, req))
}

func (m Model) updateFocusedInput(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.phase != phaseEditing {
		return m, nil
	}

	idx := inputIndex(m.focus)
	if idx < 0 {
		return m, nil
	}

	var cmd tea.Cmd
	m.inputs[idx], cmd = m.inputs[idx].Update(msg)
	return m, cmd
}

func (m Model) isFocusableField(field int) bool {
	if field != fieldSelector {
		return true
	}

	return m.mode == modeSelector
}

func (m Model) normalizeField(field int) int {
	if field < fieldURL {
		return fieldOut
	}

	if field > fieldOut {
		return fieldURL
	}

	return field
}

func friendlyInputError(err error) error {
	if err == nil {
		return nil
	}

	if _, convErr := strconv.ParseInt(err.Error(), 10, 64); convErr == nil {
		return err
	}

	return err
}
