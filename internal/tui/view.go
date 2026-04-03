package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

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

	subtitle := m.styles.muted.Render("Choose an app group first, then drill into windows and tabs only when needed")
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
	case screenTargetSelect:
		return m.renderTargetSelection()
	case screenTabSelect:
		return m.renderTabSelection()
	case screenLiveOptions:
		return m.renderLiveOptions()
	default:
		return m.renderGroupSelection()
	}
}

func (m Model) renderGroupSelection() string {
	if len(m.groups) == 0 {
		return m.styles.panel.Render(
			lipgloss.JoinVertical(
				lipgloss.Left,
				m.styles.label.Render("No open capture targets found"),
				m.styles.muted.Render("Open a folder, app, or browser window first, then press R to refresh."),
			),
		)
	}

	blocks := make([]string, 0, len(m.groups))
	for i, item := range m.groups {
		style := m.styles.field
		titleStyle := m.styles.label
		prefix := "  "
		if i == m.groupIndex {
			style = m.styles.fieldFocused.
				Background(lipgloss.Color("#083344")).
				Padding(0, 1)
			titleStyle = m.styles.accent.Bold(true)
			prefix = "▶ "
		}

		content := lipgloss.JoinVertical(
			lipgloss.Left,
			titleStyle.Render(prefix + item.title),
			m.styles.muted.Render(item.detail),
		)
		blocks = append(blocks, style.Width(m.groupCardWidth()).Render(content))
	}

	body := lipgloss.JoinVertical(lipgloss.Left, blocks...)
	if m.groupGridColumns() > 1 {
		body = renderGrid(blocks, m.groupGridColumns())
	}

	return m.styles.panel.Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			m.styles.label.Render("Select the app group"),
			m.styles.muted.Render("Example: Antigravity, Chrome, Explorer. Enter opens that group and shows its windows."),
			body,
			"",
			m.styles.accent.Render("Selected: "+m.currentGroupTitle()),
		),
	)
}

func (m Model) renderTargetSelection() string {
	if len(m.targets) == 0 {
		return m.styles.panel.Render(
			lipgloss.JoinVertical(
				lipgloss.Left,
				m.styles.label.Render("This group has no visible windows"),
				m.styles.muted.Render("Go back and choose another group, or press R to refresh from the main screen."),
			),
		)
	}

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
			m.styles.label.Render("Choose the window inside the selected group"),
			m.styles.muted.Render(m.selectedGroup.title+" • "+m.selectedGroup.detail),
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
		m.styles.label.Render("Selected group"),
		m.styles.text.Render(m.selectedGroup.title),
		m.styles.muted.Render(m.selectedGroup.detail),
		"",
		m.styles.label.Render("Selected window"),
		m.styles.text.Render(m.selectedTarget.Title),
		m.styles.muted.Render(strings.Join(compactNonEmpty([]string{displayAppName(m.selectedTarget.AppName), string(m.selectedTarget.Type)}), " • ")),
	}

	if m.hasSelectedTab {
		summary = append(summary,
			"",
			m.styles.label.Render("Selected browser tab"),
			m.styles.text.Render(m.selectedTab.Title),
		)
	}

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.JoinVertical(lipgloss.Left, summary...),
		"",
		m.styles.fieldFocused.Render(lipgloss.JoinVertical(
			lipgloss.Left,
			m.styles.label.Render("Output path"),
			m.liveOut.View(),
		)),
	)

	return m.styles.panel.Render(content)
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
	case screenTargetSelect:
		return "↑/↓ choose window • Enter open/capture • Esc back to groups • Ctrl+C quit"
	case screenTabSelect:
		return "↑/↓ choose tab • Enter confirm • Esc back to windows • Ctrl+C quit"
	case screenLiveOptions:
		return "Type output path • Enter capture current state • Esc back • Ctrl+C quit"
	default:
		return "↑/↓/←/→ choose group • Enter open group • R reload • Ctrl+C quit"
	}
}

func (m Model) renderFooter() string {
	parts := []string{
		m.styles.shortcut.Render("Enter"),
		m.styles.muted.Render("confirm"),
		m.styles.shortcut.Render("Tab"),
		m.styles.muted.Render("next item"),
		m.styles.shortcut.Render("Esc"),
		m.styles.muted.Render("back"),
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

func (m Model) groupCardWidth() int {
	if m.groupGridColumns() <= 1 {
		return max(40, contentWidth(m.width)-8)
	}

	return max(26, (contentWidth(m.width)-10)/2)
}

func renderGrid(items []string, columns int) string {
	if columns <= 1 {
		return lipgloss.JoinVertical(lipgloss.Left, items...)
	}

	rows := make([]string, 0, (len(items)+columns-1)/columns)
	for start := 0; start < len(items); start += columns {
		end := start + columns
		if end > len(items) {
			end = len(items)
		}
		rows = append(rows, lipgloss.JoinHorizontal(lipgloss.Top, items[start:end]...))
	}

	return lipgloss.JoinVertical(lipgloss.Left, rows...)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (m Model) currentGroupTitle() string {
	if len(m.groups) == 0 || m.groupIndex < 0 || m.groupIndex >= len(m.groups) {
		return "none"
	}

	return m.groups[m.groupIndex].title
}
