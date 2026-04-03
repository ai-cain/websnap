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
	baseBackground := lipgloss.Color("#14110F")
	panelBorder := lipgloss.Color("#3B312A")
	cardBackground := lipgloss.Color("#1B1714")
	focusedBackground := lipgloss.Color("#2A211B")
	accent := lipgloss.Color("#D4A574")
	titleText := lipgloss.Color("#F5EFE6")
	bodyText := lipgloss.Color("#E9DFD2")
	mutedText := lipgloss.Color("#B9A999")

	basePanel := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(panelBorder).
		Background(baseBackground).
		Padding(1, 2)

	choiceCard := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(panelBorder).
		Background(cardBackground).
		Padding(0, 0).
		MarginRight(1).
		MarginBottom(0)

	choiceCardFocused := choiceCard.
		BorderForeground(accent).
		Background(focusedBackground)

	return styles{
		frame:             lipgloss.NewStyle().Background(baseBackground).Padding(1, 2),
		panel:             basePanel,
		field:             lipgloss.NewStyle().BorderLeft(true).BorderForeground(panelBorder).Background(baseBackground).PaddingLeft(1),
		fieldFocused:      lipgloss.NewStyle().BorderLeft(true).BorderForeground(accent).Background(focusedBackground).PaddingLeft(1),
		choiceCard:        choiceCard,
		choiceCardFocused: choiceCardFocused,
		badge:             lipgloss.NewStyle().Foreground(baseBackground).Background(accent).Bold(true).Padding(0, 1),
		title:             lipgloss.NewStyle().Foreground(titleText).Bold(true),
		label:             lipgloss.NewStyle().Foreground(titleText).Bold(true),
		text:              lipgloss.NewStyle().Foreground(bodyText),
		muted:             lipgloss.NewStyle().Foreground(mutedText),
		accent:            lipgloss.NewStyle().Foreground(accent),
		focusedPrompt:     lipgloss.NewStyle().Foreground(accent).Bold(true),
		footer:            lipgloss.NewStyle().Foreground(bodyText),
		shortcut:          lipgloss.NewStyle().Foreground(accent).Bold(true),
		errorPanel:        basePanel.BorderForeground(lipgloss.Color("#E57373")),
		errorTitle:        lipgloss.NewStyle().Foreground(lipgloss.Color("#FCA5A5")).Bold(true),
		successPanel:      basePanel.BorderForeground(lipgloss.Color("#7FB069")),
		successTitle:      lipgloss.NewStyle().Foreground(lipgloss.Color("#B5D99C")).Bold(true),
	}
}
