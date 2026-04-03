package cli

import (
	"fmt"
	"io"
	"strings"

	"github.com/ai-cain/websnap/internal/tui"
)

type interactiveUI interface {
	Run(input io.Reader, output io.Writer, studio tui.Studio) error
}

type InteractiveCommand struct {
	studio tui.Studio
	input  io.Reader
	stdout io.Writer
	stderr io.Writer
	ui     interactiveUI
}

func NewInteractiveCommand(studio tui.Studio, input io.Reader, stdout, stderr io.Writer) InteractiveCommand {
	return newInteractiveCommand(studio, input, stdout, stderr, newTUIBridge())
}

func newInteractiveCommand(
	studio tui.Studio,
	input io.Reader,
	stdout, stderr io.Writer,
	ui interactiveUI,
) InteractiveCommand {
	if input == nil {
		input = strings.NewReader("")
	}

	if stdout == nil {
		stdout = io.Discard
	}

	if stderr == nil {
		stderr = io.Discard
	}

	return InteractiveCommand{
		studio: studio,
		input:  input,
		stdout: stdout,
		stderr: stderr,
		ui:     ui,
	}
}

func (c InteractiveCommand) Run() int {
	if c.studio == nil {
		fmt.Fprintln(c.stderr, "error: interactive studio is not configured")
		return 1
	}

	if c.ui == nil {
		fmt.Fprintln(c.stderr, "error: interactive ui is not configured")
		return 1
	}

	if err := c.ui.Run(c.input, c.stdout, c.studio); err != nil {
		renderError(c.stderr, err)
		return 1
	}

	return 0
}

