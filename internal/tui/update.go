package tui

import (
	"strconv"

	"github.com/ai-cain/websnap/internal/domain"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case spinner.TickMsg:
		if m.phase != phaseBusy {
			return m, nil
		}

		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case targetsLoadedMsg:
		m.enterTargetSelection()
		if msg.err != nil {
			m.lastErr = msg.err
			m.targets = []targetMenuItem{newURLMenuItem()}
			m.targetIndex = 0
			return m, nil
		}

		m.targets = []targetMenuItem{newURLMenuItem()}
		for _, target := range msg.targets {
			m.targets = append(m.targets, newLiveTargetMenuItem(target))
		}
		m.targetIndex = 0
		return m, nil
	case tabsLoadedMsg:
		m.phase = phaseEditing
		if msg.err != nil {
			m.lastErr = msg.err
			m.enterTargetSelection()
			return m, nil
		}

		m.selectedTarget = msg.target
		m.hasSelectedTarget = true
		m.tabs = msg.tabs
		if len(msg.tabs) <= 1 {
			m.hasSelectedTab = len(msg.tabs) == 1
			if m.hasSelectedTab {
				m.selectedTab = msg.tabs[0]
			}
			m.enterLiveOptions()
			return m, textinput.Blink
		}

		m.hasSelectedTab = false
		m.enterTabSelection(msg.tabs)
		return m, nil
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
	case phaseBusy:
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

	switch m.screen {
	case screenTargetSelect:
		return m.handleTargetSelectionKey(msg)
	case screenTabSelect:
		return m.handleTabSelectionKey(msg)
	case screenLiveOptions:
		return m.handleLiveOptionsKey(msg)
	case screenURLForm:
		return m.handleURLFormKey(msg)
	default:
		return m, nil
	}
}

func (m Model) handleTargetSelectionKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "esc":
		return m, tea.Quit
	case "r":
		return m.reloadTargets()
	case "up", "k", "shift+tab":
		m.moveTargetSelection(-1)
		return m, nil
	case "down", "j", "tab":
		m.moveTargetSelection(1)
		return m, nil
	case "enter":
		return m.selectCurrentTarget()
	default:
		return m, nil
	}
}

func (m Model) handleTabSelectionKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "esc":
		m.enterTargetSelection()
		return m, nil
	case "up", "k", "shift+tab":
		m.moveTabSelection(-1)
		return m, nil
	case "down", "j", "tab":
		m.moveTabSelection(1)
		return m, nil
	case "enter":
		if len(m.tabs) == 0 {
			return m, nil
		}
		m.selectedTab = m.tabs[m.tabIndex]
		m.hasSelectedTab = true
		m.enterLiveOptions()
		return m, textinput.Blink
	default:
		return m, nil
	}
}

func (m Model) handleLiveOptionsKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "esc":
		if len(m.tabs) > 1 {
			m.enterTabSelection(m.tabs)
			return m, nil
		}
		m.enterTargetSelection()
		return m, nil
	case "enter":
		req, err := m.buildLiveRequest()
		if err != nil {
			m.lastErr = friendlyInputError(err)
			return m, nil
		}

		return m.submitLiveCapture(req)
	default:
		return m.updateFocusedInput(msg)
	}
}

func (m Model) handleURLFormKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		return m, tea.Quit
	case "esc":
		m.enterTargetSelection()
		return m, nil
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
		return m.submitOrAdvanceURLForm()
	}

	return m.updateFocusedInput(msg)
}

func (m *Model) moveTargetSelection(delta int) {
	if len(m.targets) == 0 {
		return
	}

	m.targetIndex += delta
	if m.targetIndex < 0 {
		m.targetIndex = len(m.targets) - 1
	}
	if m.targetIndex >= len(m.targets) {
		m.targetIndex = 0
	}
}

func (m *Model) moveTabSelection(delta int) {
	if len(m.tabs) == 0 {
		return
	}

	m.tabIndex += delta
	if m.tabIndex < 0 {
		m.tabIndex = len(m.tabs) - 1
	}
	if m.tabIndex >= len(m.tabs) {
		m.tabIndex = 0
	}
}

func (m Model) selectCurrentTarget() (tea.Model, tea.Cmd) {
	if len(m.targets) == 0 {
		return m, nil
	}

	selected := m.targets[m.targetIndex]
	m.lastErr = nil
	m.tabs = nil
	m.tabIndex = 0
	m.hasSelectedTab = false

	switch selected.kind {
	case menuItemURL:
		m.enterURLForm()
		return m, textinput.Blink
	case menuItemLiveTarget:
		m.selectedTarget = selected.target
		m.hasSelectedTarget = true
		if selected.target.Type == domain.LiveTargetBrowser && selected.target.CanListTabs {
			return m.loadTabsForSelectedTarget()
		}

		m.enterLiveOptions()
		return m, textinput.Blink
	default:
		return m, nil
	}
}

func (m Model) reloadTargets() (tea.Model, tea.Cmd) {
	m.startBusy("Refreshing open apps, folders, and browser windows…")
	return m, loadTargetsCmd(m.studio)
}

func (m Model) loadTabsForSelectedTarget() (tea.Model, tea.Cmd) {
	m.startBusy("Inspecting open browser tabs…")
	return m, loadTabsCmd(m.studio, m.selectedTarget)
}

func (m Model) submitLiveCapture(req domain.LiveCaptureRequest) (tea.Model, tea.Cmd) {
	m.startBusy("Capturing the currently selected live target…")
	return m, submitLiveCaptureCmd(m.studio, req)
}

func (m Model) submitOrAdvanceURLForm() (tea.Model, tea.Cmd) {
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

	m.startBusy("Capturing a fresh reproducible URL snapshot…")
	return m, submitURLCaptureCmd(m.studio, req)
}

func (m Model) updateFocusedInput(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.phase != phaseEditing {
		return m, nil
	}

	switch m.screen {
	case screenLiveOptions:
		var cmd tea.Cmd
		m.liveOut, cmd = m.liveOut.Update(msg)
		return m, cmd
	case screenURLForm:
		idx := inputIndex(m.focus)
		if idx < 0 {
			return m, nil
		}

		var cmd tea.Cmd
		m.urlInputs[idx], cmd = m.urlInputs[idx].Update(msg)
		return m, cmd
	default:
		return m, nil
	}
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

