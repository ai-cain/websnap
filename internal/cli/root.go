package cli

import (
	"fmt"
	"io"
	"strings"
)

type App struct {
	shotRunner ShotRunner
	stdin      io.Reader
	stdout     io.Writer
	stderr     io.Writer
	ui         interactiveUI
}

func NewApp(shotRunner ShotRunner, stdout, stderr io.Writer) App {
	return NewAppWithInput(shotRunner, nil, stdout, stderr)
}

func NewAppWithInput(shotRunner ShotRunner, stdin io.Reader, stdout, stderr io.Writer) App {
	return newAppWithDeps(shotRunner, stdin, stdout, stderr, newTUIBridge())
}

func newAppWithDeps(
	shotRunner ShotRunner,
	stdin io.Reader,
	stdout, stderr io.Writer,
	ui interactiveUI,
) App {
	if stdin == nil {
		stdin = strings.NewReader("")
	}

	if stdout == nil {
		stdout = io.Discard
	}

	if stderr == nil {
		stderr = io.Discard
	}

	return App{
		shotRunner: shotRunner,
		stdin:      stdin,
		stdout:     stdout,
		stderr:     stderr,
		ui:         ui,
	}
}

func (a App) Run(args []string) int {
	if len(args) == 0 {
		cmd := newInteractiveCommand(a.shotRunner, a.stdin, a.stdout, a.stderr, a.ui)
		return cmd.Run()
	}

	switch strings.ToLower(args[0]) {
	case "help", "-h", "--help":
		a.printUsage()
		return 0
	case "interactive":
		cmd := newInteractiveCommand(a.shotRunner, a.stdin, a.stdout, a.stderr, a.ui)
		return cmd.Run()
	case "shot":
		cmd := NewShotCommand(a.shotRunner, a.stdout, a.stderr)
		return cmd.Run(args[1:])
	default:
		fmt.Fprintf(a.stderr, "error: unknown command %q\n\n", args[0])
		a.printUsage()
		return 1
	}
}

func (a App) printUsage() {
	fmt.Fprintln(a.stderr, "websnap captures web UI screenshots from the terminal.")
	fmt.Fprintln(a.stderr)
	fmt.Fprintln(a.stderr, "Usage:")
	fmt.Fprintln(a.stderr, "  websnap")
	fmt.Fprintln(a.stderr, "  websnap interactive")
	fmt.Fprintln(a.stderr, "  websnap shot <url> [--width <px>] [--height <px>] [--out <path>]")
}
