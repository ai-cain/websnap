package cli

import (
	"context"
	"io"

	"github.com/ai-cain/websnap/internal/domain"
	"github.com/ai-cain/websnap/internal/tui"
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

type fakeStudio struct {
	targets []domain.LiveTarget
}

func (f *fakeStudio) CaptureURL(ctx context.Context, req domain.CaptureRequest) (domain.CaptureResult, error) {
	runner := &fakeShotRunner{}
	return runner.Execute(ctx, req)
}

func (f *fakeStudio) ListTargets(_ context.Context) ([]domain.LiveTarget, error) {
	return f.targets, nil
}

func (f *fakeStudio) ListTabs(_ context.Context, _ domain.LiveTarget) ([]domain.BrowserTab, error) {
	return nil, nil
}

func (f *fakeStudio) CaptureLive(_ context.Context, _ domain.LiveCaptureRequest) (domain.CaptureResult, error) {
	return domain.CaptureResult{}, nil
}

var _ tui.Studio = (*fakeStudio)(nil)

type fakeInteractiveUI struct {
	called bool
	err    error
}

func (f *fakeInteractiveUI) Run(_ io.Reader, _ io.Writer, _ tui.Studio) error {
	f.called = true
	return f.err
}
