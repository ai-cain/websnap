package tui

import (
	"path/filepath"
	"strings"
	"unicode"

	"github.com/ai-cain/websnap/internal/domain"
	apperrors "github.com/ai-cain/websnap/internal/support/errors"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type phase int

const (
	phaseEditing phase = iota
	phaseBusy
	phaseSuccess
)

type screen int

const (
	screenTargetSelect screen = iota
	screenTabSelect
	screenLiveOptions
)

type menuItemKind int

const (
	menuItemLiveTarget menuItemKind = iota
)

type targetMenuItem struct {
	kind   menuItemKind
	title  string
	detail string
	target domain.LiveTarget
}

type Model struct {
	studio Studio

	liveOut textinput.Model
	phase   phase
	screen  screen

	spinner   spinner.Model
	styles    styles
	width     int
	height    int
	busyLabel string

	lastErr    error
	lastPath   string
	lastWidth  int64
	lastHeight int64

	targets           []targetMenuItem
	targetIndex       int
	tabs              []domain.BrowserTab
	tabIndex          int
	selectedTarget    domain.LiveTarget
	hasSelectedTarget bool
	selectedTab       domain.BrowserTab
	hasSelectedTab    bool
}

func NewModel(studio Studio) Model {
	s := newStyles()
	liveOut := newInput("./captures/live-target.png", "", s)

	spin := spinner.New()
	spin.Spinner = spinner.Dot
	spin.Style = s.accent

	model := Model{
		studio:      studio,
		liveOut:     liveOut,
		spinner:     spin,
		styles:      s,
		screen:      screenTargetSelect,
		phase:       phaseBusy,
		busyLabel:   "Discovering open apps, folders, and browser windows…",
		targetIndex: 0,
	}

	model.blurLiveInput()
	return model
}

func (m Model) Init() tea.Cmd {
	return loadTargetsCmd(m.studio)
}

func (m Model) buildLiveRequest() (domain.LiveCaptureRequest, error) {
	if !m.hasSelectedTarget {
		return domain.LiveCaptureRequest{}, apperrors.New(apperrors.CodeInvalidArgument, "live target is required")
	}

	req := domain.LiveCaptureRequest{
		Target:   m.selectedTarget,
		TabIndex: -1,
		Out:      m.liveOut.Value(),
	}

	if m.hasSelectedTab {
		req.TabIndex = m.selectedTab.Index
	}

	return req, req.Validate()
}

func (m *Model) focusLiveInput() {
	m.liveOut.Focus()
	m.liveOut.PromptStyle = m.styles.focusedPrompt
	if strings.TrimSpace(m.liveOut.Value()) == "" {
		m.liveOut.SetValue(suggestLiveOutputPath(m.selectedTarget, m.selectedTab, m.hasSelectedTab))
	}
	m.liveOut.CursorEnd()
}

func (m *Model) blurLiveInput() {
	m.liveOut.Blur()
	m.liveOut.PromptStyle = m.styles.muted
	m.liveOut.CursorEnd()
}

func (m *Model) setSuccess(result domain.CaptureResult) {
	m.phase = phaseSuccess
	m.lastErr = nil
	m.lastPath = result.Path
	m.lastWidth = result.Width
	m.lastHeight = result.Height
}

func (m *Model) startBusy(label string) {
	m.phase = phaseBusy
	m.busyLabel = label
	m.lastErr = nil
}

func (m *Model) enterTargetSelection() {
	m.screen = screenTargetSelect
	m.phase = phaseEditing
	m.blurLiveInput()
}

func (m *Model) enterLiveOptions() {
	m.screen = screenLiveOptions
	m.phase = phaseEditing
	m.lastErr = nil
	m.focusLiveInput()
}

func (m *Model) enterTabSelection(tabs []domain.BrowserTab) {
	m.screen = screenTabSelect
	m.phase = phaseEditing
	m.tabs = tabs
	m.tabIndex = selectedTabIndex(tabs)
	m.lastErr = nil
	m.blurLiveInput()
}

func newInput(placeholder, value string, s styles) textinput.Model {
	input := textinput.New()
	input.Placeholder = placeholder
	input.SetValue(value)
	input.Prompt = "> "
	input.CharLimit = 256
	input.Cursor.Style = s.accent
	input.TextStyle = s.text
	input.PlaceholderStyle = s.muted
	return input
}

func newLiveTargetMenuItem(target domain.LiveTarget) targetMenuItem {
	label := target.Title
	if strings.TrimSpace(label) == "" {
		label = target.AppName
	}

	detailParts := []string{strings.TrimSpace(target.AppName), string(target.Type)}
	if target.Type == domain.LiveTargetBrowser && target.CanListTabs {
		detailParts = append(detailParts, "tabs available")
	}

	return targetMenuItem{
		kind:   menuItemLiveTarget,
		title:  label,
		detail: strings.Join(compactNonEmpty(detailParts), " • "),
		target: target,
	}
}

func selectedTabIndex(tabs []domain.BrowserTab) int {
	for i, tab := range tabs {
		if tab.Selected {
			return i
		}
	}

	return 0
}

func suggestLiveOutputPath(target domain.LiveTarget, tab domain.BrowserTab, hasTab bool) string {
	name := target.Title
	if hasTab && strings.TrimSpace(tab.Title) != "" {
		name = tab.Title
	}

	name = strings.TrimSpace(name)
	if name == "" {
		name = target.AppName
	}

	name = sanitizeFileName(name)
	if name == "" {
		name = "live-target"
	}

	return filepath.Join("captures", name+".png")
}

func sanitizeFileName(value string) string {
	value = strings.TrimSpace(strings.ToLower(value))
	if value == "" {
		return ""
	}

	var builder strings.Builder
	lastDash := false
	for _, r := range value {
		switch {
		case unicode.IsLetter(r), unicode.IsDigit(r):
			builder.WriteRune(r)
			lastDash = false
		case r == '.', r == '-', r == '_':
			builder.WriteRune(r)
			lastDash = false
		default:
			if !lastDash {
				builder.WriteRune('-')
				lastDash = true
			}
		}
	}

	return strings.Trim(builder.String(), "-._")
}

func compactNonEmpty(values []string) []string {
	compact := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		compact = append(compact, value)
	}

	return compact
}
