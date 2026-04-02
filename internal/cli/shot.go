package cli

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"strings"

	"github.com/ai-cain/websnap/internal/domain"
	apperrors "github.com/ai-cain/websnap/internal/support/errors"
)

type ShotRunner interface {
	Execute(ctx context.Context, req domain.CaptureRequest) (domain.CaptureResult, error)
}

type ShotCommand struct {
	runner ShotRunner
	stdout io.Writer
	stderr io.Writer
}

func NewShotCommand(runner ShotRunner, stdout, stderr io.Writer) ShotCommand {
	if stdout == nil {
		stdout = io.Discard
	}

	if stderr == nil {
		stderr = io.Discard
	}

	return ShotCommand{
		runner: runner,
		stdout: stdout,
		stderr: stderr,
	}
}

func (c ShotCommand) Run(args []string) int {
	if c.runner == nil {
		fmt.Fprintln(c.stderr, "error: shot runner is not configured")
		return 1
	}

	if len(args) == 0 {
		fmt.Fprintln(c.stderr, "error: shot requires exactly one URL argument")
		c.printUsage()
		return 1
	}

	if strings.HasPrefix(args[0], "-") {
		fmt.Fprintln(c.stderr, "error: shot expects the URL before any flags")
		c.printUsage()
		return 1
	}

	targetURL := args[0]

	fs := flag.NewFlagSet("shot", flag.ContinueOnError)
	fs.SetOutput(c.stderr)
	fs.Usage = func() {
		c.printUsage()
	}

	width := fs.Int64("width", 1440, "viewport width in pixels")
	height := fs.Int64("height", 900, "viewport height in pixels")
	out := fs.String("out", "", "output path for the PNG file")

	if err := fs.Parse(args[1:]); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return 0
		}

		return 1
	}

	if len(fs.Args()) != 0 {
		fmt.Fprintln(c.stderr, "error: shot received unexpected extra positional arguments")
		c.printUsage()
		return 1
	}

	req := domain.CaptureRequest{
		URL:    targetURL,
		Width:  *width,
		Height: *height,
		Out:    *out,
	}

	result, err := c.runner.Execute(context.Background(), req)
	if err != nil {
		renderError(c.stderr, err)
		return 1
	}

	fmt.Fprintln(c.stdout, result.Path)
	return 0
}

func (c ShotCommand) printUsage() {
	fmt.Fprintln(c.stderr, "Usage:")
	fmt.Fprintln(c.stderr, "  websnap shot <url> [--width <px>] [--height <px>] [--out <path>]")
}

func renderError(w io.Writer, err error) {
	code := apperrors.CodeOf(err)
	if code == apperrors.CodeUnknown {
		fmt.Fprintf(w, "error: %v\n", err)
		return
	}

	fmt.Fprintf(w, "error [%s]: %v\n", code, err)
}
