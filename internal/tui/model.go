package tui

import (
	"path/filepath"
	"sort"
	"strconv"
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
	screenGroupSelect screen = iota
	screenTargetSelect
	screenTabSelect
	screenLiveOptions
)

type groupMenuItem struct {
	title   string
	detail  string
	targets []domain.LiveTarget
}

type targetMenuItem struct {
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

	groups           []groupMenuItem
	groupIndex       int
	selectedGroup    groupMenuItem
	hasSelectedGroup bool

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
		studio:     studio,
		liveOut:    liveOut,
		spinner:    spin,
		styles:     s,
		screen:     screenGroupSelect,
		phase:      phaseBusy,
		busyLabel:  "Discovering open apps, folders, and browser windows…",
		groupIndex: 0,
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

func (m *Model) enterGroupSelection() {
	m.screen = screenGroupSelect
	m.phase = phaseEditing
	m.blurLiveInput()
}

func (m *Model) enterTargetSelection(group groupMenuItem) {
	m.selectedGroup = group
	m.hasSelectedGroup = true
	m.targets = buildTargetMenuItems(group.targets)
	m.targetIndex = 0
	m.screen = screenTargetSelect
	m.phase = phaseEditing
	m.lastErr = nil
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

func buildGroupMenuItems(targets []domain.LiveTarget) []groupMenuItem {
	type grouped struct {
		appName string
		tType   domain.LiveTargetType
		targets []domain.LiveTarget
	}

	buckets := map[string]*grouped{}
	order := make([]string, 0)

	for _, target := range targets {
		appName := normalizedAppName(target.AppName)
		key := string(target.Type) + "|" + appName
		if _, ok := buckets[key]; !ok {
			buckets[key] = &grouped{
				appName: appName,
				tType:   target.Type,
				targets: make([]domain.LiveTarget, 0),
			}
			order = append(order, key)
		}

		buckets[key].targets = append(buckets[key].targets, target)
	}

	sort.Slice(order, func(i, j int) bool {
		left := buckets[order[i]]
		right := buckets[order[j]]
		if groupRank(left.tType) != groupRank(right.tType) {
			return groupRank(left.tType) < groupRank(right.tType)
		}
		return displayAppName(left.appName) < displayAppName(right.appName)
	})

	items := make([]groupMenuItem, 0, len(order))
	for _, key := range order {
		group := buckets[key]
		sort.Slice(group.targets, func(i, j int) bool {
			return group.targets[i].Title < group.targets[j].Title
		})

		items = append(items, groupMenuItem{
			title:   displayAppName(group.appName),
			detail:  describeGroup(group.tType, len(group.targets)),
			targets: group.targets,
		})
	}

	return items
}

func (m Model) groupGridColumns() int {
	if len(m.groups) <= 1 {
		return 1
	}

	available := contentWidth(m.width)
	if available >= 138 && len(m.groups) >= 3 {
		return 3
	}

	if available >= 88 {
		return 2
	}

	return 1
}

func buildTargetMenuItems(targets []domain.LiveTarget) []targetMenuItem {
	items := make([]targetMenuItem, 0, len(targets))
	for _, target := range targets {
		detailParts := []string{string(target.Type)}
		if target.Type == domain.LiveTargetBrowser && target.CanListTabs {
			detailParts = append(detailParts, "tabs available")
		}

		items = append(items, targetMenuItem{
			title:  target.Title,
			detail: strings.Join(compactNonEmpty(detailParts), " • "),
			target: target,
		})
	}

	return items
}

func normalizedAppName(appName string) string {
	appName = strings.TrimSpace(strings.ToLower(appName))
	if appName == "" {
		return "unknown"
	}

	return appName
}

func displayAppName(appName string) string {
	switch normalizedAppName(appName) {
	case "chrome":
		return "Chrome"
	case "msedge":
		return "Edge"
	case "explorer":
		return "Explorador"
	case "applicationframehost":
		return "Windows Host"
	case "systemsettings":
		return "Configuración"
	case "textinputhost":
		return "Entrada de Windows"
	case "unknown":
		return "Desconocido"
	default:
		return titleCase(appName)
	}
}

func titleCase(value string) string {
	if value == "" {
		return ""
	}

	parts := strings.FieldsFunc(value, func(r rune) bool {
		return r == '-' || r == '_' || unicode.IsSpace(r)
	})

	for i, part := range parts {
		runes := []rune(part)
		if len(runes) == 0 {
			continue
		}
		runes[0] = unicode.ToUpper(runes[0])
		for j := 1; j < len(runes); j++ {
			runes[j] = unicode.ToLower(runes[j])
		}
		parts[i] = string(runes)
	}

	return strings.Join(parts, " ")
}

func describeGroup(targetType domain.LiveTargetType, count int) string {
	itemLabel := "ventanas"
	if count == 1 {
		itemLabel = "ventana"
	}

	switch targetType {
	case domain.LiveTargetBrowser:
		return pluralize(count, itemLabel) + " • navegador"
	case domain.LiveTargetFolder:
		return pluralize(count, itemLabel) + " • carpetas"
	default:
		return pluralize(count, itemLabel) + " • app"
	}
}

func pluralize(count int, noun string) string {
	return strconv.Itoa(count) + " " + noun
}

func groupRank(targetType domain.LiveTargetType) int {
	switch targetType {
	case domain.LiveTargetBrowser:
		return 0
	case domain.LiveTargetFolder:
		return 1
	default:
		return 2
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
