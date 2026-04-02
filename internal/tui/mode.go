package tui

type captureMode int

const (
	modeViewport captureMode = iota
	modeSelector
	modeFullPage
)

var captureModes = []captureMode{
	modeViewport,
	modeSelector,
	modeFullPage,
}

func (m captureMode) label() string {
	switch m {
	case modeSelector:
		return "selector"
	case modeFullPage:
		return "full-page"
	default:
		return "viewport"
	}
}

func (m captureMode) next() captureMode {
	index := 0
	for i, candidate := range captureModes {
		if candidate == m {
			index = i
			break
		}
	}

	return captureModes[(index+1)%len(captureModes)]
}

func (m captureMode) previous() captureMode {
	index := 0
	for i, candidate := range captureModes {
		if candidate == m {
			index = i
			break
		}
	}

	if index == 0 {
		return captureModes[len(captureModes)-1]
	}

	return captureModes[index-1]
}
