package cli

import (
	"io"

	"github.com/ai-cain/websnap/internal/tui"
)

type tuiBridge struct {
	runner tui.Runner
}

func newTUIBridge() tuiBridge {
	return tuiBridge{runner: tui.NewRunner()}
}

func (b tuiBridge) Run(input io.Reader, output io.Writer, studio tui.Studio) error {
	return b.runner.Run(input, output, studio)
}

