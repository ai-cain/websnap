package cli

import (
	"context"
	"io"

	"github.com/ai-cain/websnap/internal/domain"
	"github.com/ai-cain/websnap/internal/tui"
)

type tuiBridge struct {
	runner tui.Runner
}

func newTUIBridge() tuiBridge {
	return tuiBridge{runner: tui.NewRunner()}
}

func (b tuiBridge) Run(input io.Reader, output io.Writer, runner ShotRunner) error {
	return b.runner.Run(input, output, tuiShotRunner{runner: runner})
}

type tuiShotRunner struct {
	runner ShotRunner
}

func (r tuiShotRunner) Execute(ctx context.Context, req domain.CaptureRequest) (domain.CaptureResult, error) {
	return r.runner.Execute(ctx, req)
}
