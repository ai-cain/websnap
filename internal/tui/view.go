package tui

import (
	"fmt"
	"strings"

	"github.com/ai-cain/websnap/internal/domain"
	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		m.renderHeader(),
		m.renderBody(),
		m.renderFooter(),
	)

	return m.styles.frame.Render(content)
}

func (m Model) renderHeader() string {
	title := lipgloss.JoinHorizontal(
		lipgloss.Center,
		m.styles.badge.Render("websnap"),
		" ",
		m.styles.title.Render("Interactive Capture Studio"),
	)

	ruleWidth := max(18, min(72, contentWidth(m.width)-2))
	subtitle := m.styles.muted.Render("Choose an app group first, then drill into windows and tabs only when needed")

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		subtitle,
		m.styles.muted.Render(strings.Repeat("─", ruleWidth)),
	)
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
		return lipgloss.JoinVertical(
			lipgloss.Left,
			m.renderSectionHeader("No open capture targets found", "Open a folder, app, or browser window first, then press R to refresh."),
		)
	}

	blocks := make([]string, 0, len(m.groups))
	for i, item := range m.groups {
		blocks = append(blocks, m.renderChoiceCard(item.title, item.detail, i == m.groupIndex, m.groupCardWidth()))
	}

	body := lipgloss.JoinVertical(lipgloss.Left, blocks...)
	if m.groupGridColumns() > 1 {
		body = renderGrid(blocks, m.groupGridColumns())
	}

	parts := []string{
		m.renderSectionHeader("Select the app group", "Example: Antigravity, Chrome, Explorer. Enter opens that group and shows its windows."),
		body,
		"",
		m.styles.accent.Render("Selected: " + m.currentGroupTitle()),
	}

	if !m.hasBrowserGroup() {
		parts = append(parts,
			"",
			m.renderBrowserExtensionHint(),
		)
	}

	return lipgloss.JoinVertical(lipgloss.Left, parts...)
}

func (m Model) renderTargetSelection() string {
	if len(m.targets) == 0 {
		return lipgloss.JoinVertical(
			lipgloss.Left,
			m.renderSectionHeader("This group has no visible windows", "Go back and choose another group, or press R to refresh from the main screen."),
		)
	}

	blocks := make([]string, 0, len(m.targets))
	cardWidth := m.selectionCardWidth()
	for i, item := range m.targets {
		blocks = append(blocks, m.renderChoiceCard(item.title, item.detail, i == m.targetIndex, cardWidth))
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		m.renderSectionHeader("Choose the window inside the selected group", m.selectedGroup.title+" • "+m.selectedGroup.detail),
		lipgloss.JoinVertical(lipgloss.Left, blocks...),
	)
}

func (m Model) renderTabSelection() string {
	blocks := make([]string, 0, len(m.tabs))
	cardWidth := m.selectionCardWidth()
	for i, tab := range m.tabs {
		state := "open tab"
		if tab.Selected {
			state = "currently selected"
		}

		blocks = append(blocks, m.renderChoiceCard(tab.Title, state, i == m.tabIndex, cardWidth))
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		m.renderSectionHeader("Choose the already-open browser tab", m.selectedTarget.Title),
		lipgloss.JoinVertical(lipgloss.Left, blocks...),
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

	inputCardStyle := m.styles.choiceCardFocused
	inputCard := inputCardStyle.Render(lipgloss.JoinVertical(
		lipgloss.Left,
		m.styles.label.Render("Output path"),
		m.liveOut.View(),
	))

	return lipgloss.JoinVertical(
		lipgloss.Left,
		m.renderSectionHeader("Capture the current live state", "Review the selected target and confirm the output path."),
		lipgloss.JoinVertical(lipgloss.Left, summary...),
		"",
		inputCard,
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

func (m Model) renderBrowserExtensionHint() string {
	innerWidth := max(24, contentWidth(m.width)-m.styles.panel.GetHorizontalFrameSize())
	panelStyle := m.styles.panel.Width(innerWidth)
	bodyStyle := lipgloss.NewStyle().Foreground(m.styles.muted.GetForeground()).Width(innerWidth)

	return panelStyle.Render(lipgloss.JoinVertical(
		lipgloss.Left,
		m.styles.label.Render("Browser targets are extension-backed now"),
		bodyStyle.Render("If you expected Chrome or Edge here: load extensions/chromium-websnap, keep at least one http(s) tab open, click the extension once, then press R."),
	))
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
	}

	switch m.screen {
	case screenTargetSelect, screenTabSelect, screenLiveOptions:
		parts = append(parts,
			m.styles.shortcut.Render("Esc"),
			m.styles.muted.Render("back"),
		)
	default:
		parts = append(parts,
			m.styles.shortcut.Render("Esc"),
			m.styles.muted.Render("quit"),
		)
	}

	parts = append(parts,
		m.styles.shortcut.Render("Ctrl+C"),
		m.styles.muted.Render("quit"),
	)

	return m.styles.footer.Render(strings.Join(parts, "   "))
}

func (m Model) renderSectionHeader(title, subtitle string) string {
	return lipgloss.JoinVertical(
		lipgloss.Left,
		m.styles.label.Render(title),
		m.styles.muted.Render(subtitle),
		"",
	)
}

func contentWidth(total int) int {
	switch {
	case total >= 160:
		return 148
	case total >= 128:
		return total - 8
	case total >= 92:
		return total - 6
	case total >= 74:
		return total - 4
	default:
		return 70
	}
}

func (m Model) groupCardWidth() int {
	columns := m.groupGridColumns()
	if columns <= 1 {
		return max(32, contentWidth(m.width)-2)
	}

	return max(24, (contentWidth(m.width)-2)/columns)
}

func (m Model) selectionCardWidth() int {
	return max(32, contentWidth(m.width)-2)
}

func (m Model) renderChoiceCard(title, detail string, focused bool, width int) string {
	style := m.styles.choiceCard
	titleStyle := m.styles.label
	prefix := "  "
	if focused {
		style = m.styles.choiceCardFocused
		titleStyle = m.styles.accent.Bold(true)
		prefix = "▶ "
	}

	innerWidth := max(12, width-style.GetHorizontalFrameSize())
	content := m.renderChoiceLine(prefix+title, detail, titleStyle, innerWidth)
	return style.Render(padVisibleWidth(content, innerWidth))
}

func (m Model) renderChoiceLine(title, detail string, titleStyle lipgloss.Style, width int) string {
	if width <= 0 {
		return lipgloss.JoinHorizontal(
			lipgloss.Left,
			titleStyle.Render(title),
			m.styles.muted.Render("  "+detail),
		)
	}

	rightText := trimText(detail, min(20, max(10, width/3)))
	rightWidth := lipgloss.Width(rightText)
	leftMax := max(10, width-rightWidth-2)
	leftText := trimText(title, leftMax)
	leftWidth := lipgloss.Width(leftText)
	gap := width - leftWidth - rightWidth
	if gap < 2 {
		gap = 2
	}

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		titleStyle.Render(leftText),
		strings.Repeat(" ", gap),
		m.styles.muted.Render(rightText),
	)
}

func trimText(text string, limit int) string {
	runes := []rune(text)
	if len(runes) <= limit {
		return text
	}
	if limit <= 1 {
		return string(runes[:max(0, limit)])
	}
	return string(runes[:limit-1]) + "…"
}

func padVisibleWidth(text string, width int) string {
	diff := width - lipgloss.Width(text)
	if diff <= 0 {
		return text
	}

	return text + strings.Repeat(" ", diff)
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

		cells := make([]string, 0, end-start+columns-1)
		for i, item := range items[start:end] {
			if i > 0 {
				cells = append(cells, "  ")
			}
			cells = append(cells, item)
		}

		rows = append(rows, lipgloss.JoinHorizontal(lipgloss.Top, cells...))
	}

	return lipgloss.JoinVertical(lipgloss.Left, rows...)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
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

func (m Model) hasBrowserGroup() bool {
	for _, group := range m.groups {
		if len(group.targets) == 0 {
			continue
		}

		if group.targets[0].Type == domain.LiveTargetBrowser {
			return true
		}
	}

	return false
}
