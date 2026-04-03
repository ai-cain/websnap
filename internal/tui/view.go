package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var fieldLabels = []string{
	"Target URL",
	"Viewport width",
	"Viewport height",
	"Capture mode",
	"Selector",
	"Output path",
}

func (m Model) View() string {
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		m.renderHeader(),
		m.renderBody(),
		m.renderFooter(),
	)

	return m.styles.frame.Width(contentWidth(m.width)).Render(content)
}

func (m Model) renderHeader() string {
	title := lipgloss.JoinHorizontal(
		lipgloss.Center,
		m.styles.badge.Render("websnap"),
		" ",
		m.styles.title.Render("Interactive Capture Studio"),
	)

	subtitle := m.styles.muted.Render("Choose an already-open target or fall back to a fresh reproducible URL capture")
	return m.styles.panel.Render(lipgloss.JoinVertical(lipgloss.Left, title, subtitle))
}

func (m Model) renderBody() string {
	switch m.phase {
	case phaseBusy:
		return m.styles.panel.Render(m.spinner.View() + " " + m.busyLabel)
	case phaseSuccess:
		return m.styles.successPanel.Render(
			lipgloss.JoinVertical(
				lipgloss.Left,
				m.styles.successTitle.Render("Screenshot saved"),
				m.styles.text.Render(m.lastPath),
				m.styles.muted.Render(fmt.Sprintf("%dx%d • press Enter to close", m.lastWidth, m.lastHeight)),
			),
		)
	default:
		return lipgloss.JoinVertical(
			lipgloss.Left,
			m.renderEditingScreen(),
			m.renderStatus(),
		)
	}
}

func (m Model) renderEditingScreen() string {
	switch m.screen {
	case screenTabSelect:
		return m.renderTabSelection()
	case screenLiveOptions:
		return m.renderLiveOptions()
	case screenURLForm:
		return m.renderURLForm()
	default:
		return m.renderTargetSelection()
	}
}

func (m Model) renderTargetSelection() string {
	blocks := make([]string, 0, len(m.targets))
	for i, item := range m.targets {
		style := m.styles.field
		prefix := "  "
		if i == m.targetIndex {
			style = m.styles.fieldFocused
			prefix = "> "
		}

		content := lipgloss.JoinVertical(
			lipgloss.Left,
			m.styles.label.Render(prefix+item.title),
			m.styles.muted.Render(item.detail),
		)
		blocks = append(blocks, style.Render(content))
	}

	return m.styles.panel.Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			m.styles.label.Render("Select what you want to capture"),
			m.styles.muted.Render("Open folders, apps, and browser windows appear here. Browser tabs are suggested next when available."),
			lipgloss.JoinVertical(lipgloss.Left, blocks...),
		),
	)
}

func (m Model) renderTabSelection() string {
	blocks := make([]string, 0, len(m.tabs))
	for i, tab := range m.tabs {
		style := m.styles.field
		prefix := "  "
		if i == m.tabIndex {
			style = m.styles.fieldFocused
			prefix = "> "
		}

		state := "open tab"
		if tab.Selected {
			state = "currently selected"
		}

		content := lipgloss.JoinVertical(
			lipgloss.Left,
			m.styles.label.Render(prefix+tab.Title),
			m.styles.muted.Render(state),
		)
		blocks = append(blocks, style.Render(content))
	}

	return m.styles.panel.Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			m.styles.label.Render("Choose the already-open browser tab"),
			m.styles.muted.Render(m.selectedTarget.Title),
			lipgloss.JoinVertical(lipgloss.Left, blocks...),
		),
	)
}

func (m Model) renderLiveOptions() string {
	summary := []string{
		m.styles.label.Render("Selected live target"),
		m.styles.text.Render(m.selectedTarget.Title),
		m.styles.muted.Render(strings.Join(compactNonEmpty([]string{m.selectedTarget.AppName, string(m.selectedTarget.Type)}), " • ")),
	}

	if m.hasSelectedTab {
		summary = append(summary,
			m.styles.label.Render("Selected browser tab"),
			m.styles.text.Render(m.selectedTab.Title),
		)
	}

	outputStyle := m.styles.fieldFocused
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.JoinVertical(lipgloss.Left, summary...),
		"",
		outputStyle.Render(lipgloss.JoinVertical(
			lipgloss.Left,
			m.styles.label.Render("Output path"),
			m.liveOut.View(),
		)),
	)

	return m.styles.panel.Render(content)
}

func (m Model) renderURLForm() string {
	blocks := make([]string, 0, len(fieldLabels))
	for field := range fieldLabels {
		if field == fieldMode {
			blocks = append(blocks, m.renderModeField(field))
			continue
		}

		idx := inputIndex(field)
		content := lipgloss.JoinVertical(
			lipgloss.Left,
			m.styles.label.Render(fieldLabels[field]),
			m.renderFieldValue(field, idx),
		)

		style := m.styles.field
		if field == m.focus {
			style = m.styles.fieldFocused
		}

		blocks = append(blocks, style.Render(content))
	}

	return m.styles.panel.Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			m.styles.label.Render("Fresh reproducible URL capture"),
			m.styles.muted.Render("Use this when you want websnap to open a new clean page and capture it headlessly."),
			lipgloss.JoinVertical(lipgloss.Left, blocks...),
		),
	)
}

func (m Model) renderStatus() string {
	instructions := m.instructionsForCurrentScreen()
	parts := []string{m.styles.panel.Render(m.styles.muted.Render(instructions))}

	if m.lastErr != nil {
		parts = append(parts, m.styles.errorPanel.Render(
			lipgloss.JoinVertical(
				lipgloss.Left,
				m.styles.errorTitle.Render("Capture error"),
				m.styles.text.Render(m.lastErr.Error()),
			),
		))
	}

	return lipgloss.JoinVertical(lipgloss.Left, parts...)
}

func (m Model) instructionsForCurrentScreen() string {
	switch m.screen {
	case screenTabSelect:
		return "?/? choose tab • Enter confirm • Esc back • Ctrl+C quit"
	case screenLiveOptions:
		return "Type output path • Enter capture current state • Esc back • Ctrl+C quit"
	case screenURLForm:
		return "Tab move • ?/? change mode • Enter continue/capture • Esc back • Ctrl+C quit"
	default:
		return "?/? choose target • Enter inspect/capture • R reload • Ctrl+C quit"
	}
}

func (m Model) renderFooter() string {
	parts := []string{
		m.styles.shortcut.Render("Enter"),
		m.styles.muted.Render("confirm"),
		m.styles.shortcut.Render("Tab"),
		m.styles.muted.Render("next item/field"),
		m.styles.shortcut.Render("Esc"),
		m.styles.muted.Render("back"),
		m.styles.shortcut.Render("Ctrl+C"),
		m.styles.muted.Render("quit"),
	}

	return m.styles.footer.Render(strings.Join(parts, "   "))
}

func (m Model) renderModeField(field int) string {
	options := make([]string, 0, len(captureModes))
	for _, option := range captureModes {
		style := m.styles.muted
		if option == m.mode {
			style = m.styles.accent.Bold(true)
		}

		options = append(options, style.Render(option.label()))
	}

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		m.styles.label.Render(fieldLabels[field]),
		strings.Join(options, "   "),
	)

	style := m.styles.field
	if field == m.focus {
		style = m.styles.fieldFocused
	}

	return style.Render(content)
}

func (m Model) renderFieldValue(field, idx int) string {
	if field == fieldSelector && m.mode != modeSelector {
		return m.styles.muted.Render("Used only in selector mode")
	}

	return m.urlInputs[idx].View()
}

func contentWidth(total int) int {
	switch {
	case total >= 92:
		return 88
	case total >= 74:
		return total - 4
	default:
		return 70
	}
}

