package tui

import "github.com/charmbracelet/lipgloss"

type styles struct {
	frame             lipgloss.Style
	panel             lipgloss.Style
	field             lipgloss.Style
	fieldFocused      lipgloss.Style
	choiceCard        lipgloss.Style
	choiceCardFocused lipgloss.Style
	badge             lipgloss.Style
	title             lipgloss.Style
	label             lipgloss.Style
	text              lipgloss.Style
	muted             lipgloss.Style
	accent            lipgloss.Style
	focusedPrompt     lipgloss.Style
	footer            lipgloss.Style
	shortcut          lipgloss.Style
	errorPanel        lipgloss.Style
	errorTitle        lipgloss.Style
	successPanel      lipgloss.Style
	successTitle      lipgloss.Style
}

func newStyles() styles {
	basePanel := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#334155")).
		Background(lipgloss.Color("#020617")).
		Padding(1, 2)

	choiceCard := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#334155")).
		Background(lipgloss.Color("#0F172A")).
		Padding(0, 1).
		MarginRight(1).
		MarginBottom(0)

	choiceCardFocused := choiceCard.
		BorderForeground(lipgloss.Color("#22D3EE")).
		Background(lipgloss.Color("#0C4A6E"))

	return styles{
		frame:             lipgloss.NewStyle().Background(lipgloss.Color("#020617")).Padding(1, 2),
		panel:             basePanel,
		field:             lipgloss.NewStyle().BorderLeft(true).BorderForeground(lipgloss.Color("#1E293B")).Background(lipgloss.Color("#020617")).PaddingLeft(1),
		fieldFocused:      lipgloss.NewStyle().BorderLeft(true).BorderForeground(lipgloss.Color("#06B6D4")).Background(lipgloss.Color("#082F49")).PaddingLeft(1),
		choiceCard:        choiceCard,
		choiceCardFocused: choiceCardFocused,
		badge:             lipgloss.NewStyle().Foreground(lipgloss.Color("#0F172A")).Background(lipgloss.Color("#22D3EE")).Bold(true).Padding(0, 1),
		title:             lipgloss.NewStyle().Foreground(lipgloss.Color("#E2E8F0")).Bold(true),
		label:             lipgloss.NewStyle().Foreground(lipgloss.Color("#F8FAFC")).Bold(true),
		text:              lipgloss.NewStyle().Foreground(lipgloss.Color("#E2E8F0")),
		muted:             lipgloss.NewStyle().Foreground(lipgloss.Color("#94A3B8")),
		accent:            lipgloss.NewStyle().Foreground(lipgloss.Color("#22D3EE")),
		focusedPrompt:     lipgloss.NewStyle().Foreground(lipgloss.Color("#22D3EE")).Bold(true),
		footer:            lipgloss.NewStyle().Foreground(lipgloss.Color("#CBD5E1")),
		shortcut:          lipgloss.NewStyle().Foreground(lipgloss.Color("#C084FC")).Bold(true),
		errorPanel:        basePanel.BorderForeground(lipgloss.Color("#EF4444")),
		errorTitle:        lipgloss.NewStyle().Foreground(lipgloss.Color("#FCA5A5")).Bold(true),
		successPanel:      basePanel.BorderForeground(lipgloss.Color("#22C55E")),
		successTitle:      lipgloss.NewStyle().Foreground(lipgloss.Color("#86EFAC")).Bold(true),
	}
}
