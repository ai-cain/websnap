package cli

import (
	"fmt"
	"io"
	"strings"
)

type interactiveUI interface {
	Run(input io.Reader, output io.Writer, runner ShotRunner) error
}

type InteractiveCommand struct {
	runner ShotRunner
	input  io.Reader
	stdout io.Writer
	stderr io.Writer
	ui     interactiveUI
}

func NewInteractiveCommand(runner ShotRunner, input io.Reader, stdout, stderr io.Writer) InteractiveCommand {
	return newInteractiveCommand(runner, input, stdout, stderr, newTUIBridge())
}

func newInteractiveCommand(
	runner ShotRunner,
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
		runner: runner,
		input:  input,
		stdout: stdout,
		stderr: stderr,
		ui:     ui,
	}
}

func (c InteractiveCommand) Run() int {
	if c.runner == nil {
		fmt.Fprintln(c.stderr, "error: interactive runner is not configured")
		return 1
	}

	if c.ui == nil {
		fmt.Fprintln(c.stderr, "error: interactive ui is not configured")
		return 1
	}

	if err := c.ui.Run(c.input, c.stdout, c.runner); err != nil {
		renderError(c.stderr, err)
		return 1
	}

	return 0
}
