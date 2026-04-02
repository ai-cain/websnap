package tui

import (
	"strconv"
	"strings"

	"github.com/ai-cain/websnap/internal/domain"
	apperrors "github.com/ai-cain/websnap/internal/support/errors"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	fieldURL = iota
	fieldWidth
	fieldHeight
	fieldMode
	fieldSelector
	fieldOut
)

type phase int

const (
	phaseEditing phase = iota
	phaseCapturing
	phaseSuccess
)

type Model struct {
	runner     ShotRunner
	inputs     []textinput.Model
	focus      int
	mode       captureMode
	phase      phase
	spinner    spinner.Model
	styles     styles
	width      int
	height     int
	lastErr    error
	lastPath   string
	lastWidth  int64
	lastHeight int64
}

func NewModel(runner ShotRunner) Model {
	s := newStyles()
	inputs := []textinput.Model{
		newInput("https://example.com", "", s),
		newInput("", "1440", s),
		newInput("", "900", s),
		newInput(".hero", "", s),
		newInput("./captures/home.png", "", s),
	}

	spin := spinner.New()
	spin.Spinner = spinner.Dot
	spin.Style = s.accent

	model := Model{
		runner:  runner,
		inputs:  inputs,
		mode:    modeViewport,
		spinner: spin,
		styles:  s,
	}

	model.setFocus(fieldURL)
	return model
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) buildRequest() (domain.CaptureRequest, error) {
	width, err := strconv.ParseInt(m.inputs[inputIndex(fieldWidth)].Value(), 10, 64)
	if err != nil {
		return domain.CaptureRequest{}, err
	}

	height, err := strconv.ParseInt(m.inputs[inputIndex(fieldHeight)].Value(), 10, 64)
	if err != nil {
		return domain.CaptureRequest{}, err
	}

	req := domain.CaptureRequest{
		URL:    m.inputs[inputIndex(fieldURL)].Value(),
		Width:  width,
		Height: height,
		Out:    m.inputs[inputIndex(fieldOut)].Value(),
	}

	switch m.mode {
	case modeSelector:
		req.Selector = m.inputs[inputIndex(fieldSelector)].Value()
		if strings.TrimSpace(req.Selector) == "" {
			return domain.CaptureRequest{}, apperrors.New(apperrors.CodeInvalidArgument, "selector is required in selector mode")
		}
	case modeFullPage:
		req.FullPage = true
	}

	return req, req.Validate()
}

func (m *Model) setFocus(field int) {
	m.focus = field

	for current := fieldURL; current <= fieldOut; current++ {
		idx := inputIndex(current)
		if idx < 0 {
			continue
		}

		if current == field {
			m.inputs[idx].Focus()
			m.inputs[idx].PromptStyle = m.styles.focusedPrompt
			continue
		}

		m.inputs[idx].Blur()
		m.inputs[idx].PromptStyle = m.styles.muted
	}
}

func (m *Model) setSuccess(result domain.CaptureResult) {
	m.phase = phaseSuccess
	m.lastErr = nil
	m.lastPath = result.Path
	m.lastWidth = result.Width
	m.lastHeight = result.Height
}

func newInput(placeholder, value string, s styles) textinput.Model {
	input := textinput.New()
	input.Placeholder = placeholder
	input.SetValue(value)
	input.Prompt = "› "
	input.CharLimit = 256
	input.Cursor.Style = s.accent
	input.TextStyle = s.text
	input.PlaceholderStyle = s.muted
	return input
}

func inputIndex(field int) int {
	switch field {
	case fieldURL:
		return 0
	case fieldWidth:
		return 1
	case fieldHeight:
		return 2
	case fieldSelector:
		return 3
	case fieldOut:
		return 4
	default:
		return -1
	}
}
