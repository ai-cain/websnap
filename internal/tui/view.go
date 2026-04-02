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
	blocks := make([]string, 0, len(m.inputs))
	for i := range m.inputs {
		field := lipgloss.JoinVertical(
			lipgloss.Left,
			m.styles.label.Render(fieldLabels[i]),
			m.inputs[i].View(),
		)

		style := m.styles.field
		if i == m.focus {
			style = m.styles.fieldFocused
		}

		blocks = append(blocks, style.Render(field))
	}

	return m.styles.panel.Render(lipgloss.JoinVertical(lipgloss.Left, blocks...))
}

func (m Model) renderStatus() string {
	if m.lastErr == nil {
		return m.styles.panel.Render(
			m.styles.muted.Render("Tab / Shift+Tab to move • Enter to continue • Ctrl+C to quit"),
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
		m.styles.shortcut.Render("Ctrl+C"),
		m.styles.muted.Render("quit"),
	}

	return m.styles.footer.Render(strings.Join(parts, "   "))
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
