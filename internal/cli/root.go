package cli

import (
	"fmt"
	"io"
	"strings"
)

type App struct {
	shotRunner ShotRunner
	stdout     io.Writer
	stderr     io.Writer
}

func NewApp(shotRunner ShotRunner, stdout, stderr io.Writer) App {
	if stdout == nil {
		stdout = io.Discard
	}

	if stderr == nil {
		stderr = io.Discard
	}

	return App{
		shotRunner: shotRunner,
		stdout:     stdout,
		stderr:     stderr,
	}
}

func (a App) Run(args []string) int {
	if len(args) == 0 {
		a.printUsage()
		return 1
	}

	switch strings.ToLower(args[0]) {
	case "help", "-h", "--help":
		a.printUsage()
		return 0
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
	fmt.Fprintln(a.stderr, "  websnap shot <url> [--width <px>] [--height <px>] [--out <path>]")
}
