package cli

import (
	"context"
	"io"

	"github.com/ai-cain/websnap/internal/domain"
)

type fakeShotRunner struct {
	received domain.CaptureRequest
	result   domain.CaptureResult
	err      error
}

func (f *fakeShotRunner) Execute(_ context.Context, req domain.CaptureRequest) (domain.CaptureResult, error) {
	f.received = req
	return f.result, f.err
}

type fakeInteractiveUI struct {
	called bool
	err    error
}

func (f *fakeInteractiveUI) Run(_ io.Reader, _ io.Writer, _ ShotRunner) error {
	f.called = true
	return f.err
}
