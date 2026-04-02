package tui

import (
	"strconv"

	"github.com/ai-cain/websnap/internal/domain"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	fieldURL = iota
	fieldWidth
	fieldHeight
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
		newInput("./captures/home.png", "", s),
	}

	spin := spinner.New()
	spin.Spinner = spinner.Dot
	spin.Style = s.accent

	model := Model{
		runner:  runner,
		inputs:  inputs,
		spinner: spin,
		styles:  s,
	}

	model.setFocus(0)
	return model
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) buildRequest() (domain.CaptureRequest, error) {
	width, err := strconv.ParseInt(m.inputs[fieldWidth].Value(), 10, 64)
	if err != nil {
		return domain.CaptureRequest{}, err
	}

	height, err := strconv.ParseInt(m.inputs[fieldHeight].Value(), 10, 64)
	if err != nil {
		return domain.CaptureRequest{}, err
	}

	req := domain.CaptureRequest{
		URL:    m.inputs[fieldURL].Value(),
		Width:  width,
		Height: height,
		Out:    m.inputs[fieldOut].Value(),
	}

	return req, req.Validate()
}

func (m *Model) setFocus(index int) {
	m.focus = index

	for i := range m.inputs {
		if i == index {
			m.inputs[i].Focus()
			m.inputs[i].PromptStyle = m.styles.focusedPrompt
			continue
		}

		m.inputs[i].Blur()
		m.inputs[i].PromptStyle = m.styles.muted
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
