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

	subtitle := m.styles.muted.Render("Elegant terminal screenshots with Go + chromedp")
	return m.styles.panel.Render(lipgloss.JoinVertical(lipgloss.Left, title, subtitle))
}

func (m Model) renderBody() string {
	switch m.phase {
	case phaseCapturing:
		return m.styles.panel.Render(m.spinner.View() + " Capturing screenshot…")
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
			m.renderForm(),
			m.renderStatus(),
		)
	}
}

func (m Model) renderForm() string {
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

	return m.styles.panel.Render(lipgloss.JoinVertical(lipgloss.Left, blocks...))
}

func (m Model) renderStatus() string {
	if m.lastErr == nil {
		return m.styles.panel.Render(
			m.styles.muted.Render("Tab to move • ←/→ change mode • Enter to continue • Ctrl+C to quit"),
		)
	}

	return m.styles.errorPanel.Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			m.styles.errorTitle.Render("Capture error"),
			m.styles.text.Render(m.lastErr.Error()),
		),
	)
}

func (m Model) renderFooter() string {
	parts := []string{
		m.styles.shortcut.Render("Enter"),
		m.styles.muted.Render("confirm"),
		m.styles.shortcut.Render("Tab"),
		m.styles.muted.Render("next field"),
		m.styles.shortcut.Render("←/→"),
		m.styles.muted.Render("change mode"),
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

	return m.inputs[idx].View()
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
