package windows

import (
	"context"
	"encoding/base64"
	"errors"
	"testing"

	"github.com/ai-cain/websnap/internal/domain"
)

func TestDesktopListTargets(t *testing.T) {
	t.Parallel()

	executor := &fakeExecutor{
		outputs: []executorResult{
			{
				output: []byte(`[
					{"windowHandle":131584,"title":"WhatsApp - Google Chrome","appName":"chrome","type":"browser","canListTabs":true},
					{"windowHandle":1312510,"title":"portfolio - Explorador de archivos","appName":"explorer","type":"folder","canListTabs":false}
				]`),
			},
		},
	}

	desktop := newDesktop(executor)

	targets, err := desktop.ListTargets(context.Background())
	if err != nil {
		t.Fatalf("ListTargets() error = %v", err)
	}

	if len(targets) != 2 {
		t.Fatalf("len(targets) = %d, want 2", len(targets))
	}

	if targets[0].Type != domain.LiveTargetBrowser || !targets[0].CanListTabs {
		t.Fatalf("targets[0] = %#v, want browser with tabs", targets[0])
	}

	if targets[1].Type != domain.LiveTargetFolder {
		t.Fatalf("targets[1] = %#v, want folder target", targets[1])
	}
}

func TestDesktopListTabs(t *testing.T) {
	t.Parallel()

	executor := &fakeExecutor{
		outputs: []executorResult{
			{
				output: []byte(`[
					{"index":0,"title":"WhatsApp","selected":true},
					{"index":1,"title":"YouTube","selected":false}
				]`),
			},
		},
	}

	desktop := newDesktop(executor)

	tabs, err := desktop.ListTabs(context.Background(), domain.LiveTarget{
		WindowHandle: 131584,
		Title:        "WhatsApp - Google Chrome",
		AppName:      "chrome",
		Type:         domain.LiveTargetBrowser,
		CanListTabs:  true,
	})
	if err != nil {
		t.Fatalf("ListTabs() error = %v", err)
	}

	if len(tabs) != 2 {
		t.Fatalf("len(tabs) = %d, want 2", len(tabs))
	}

	if !tabs[0].Selected || tabs[1].Selected {
		t.Fatalf("tabs = %#v, want first selected only", tabs)
	}
}

func TestDesktopCapture(t *testing.T) {
	t.Parallel()

	executor := &fakeExecutor{
		outputs: []executorResult{
			{
				output: []byte(`{"width":1550,"height":830,"pngBase64":"` + base64.StdEncoding.EncodeToString([]byte("png")) + `"}`),
			},
		},
	}

	desktop := newDesktop(executor)

	result, err := desktop.Capture(context.Background(), domain.LiveCaptureRequest{
		Target: domain.LiveTarget{
			WindowHandle: 131584,
			Title:        "WhatsApp - Google Chrome",
			AppName:      "chrome",
			Type:         domain.LiveTargetBrowser,
			CanListTabs:  true,
		},
		TabIndex: 0,
	})
	if err != nil {
		t.Fatalf("Capture() error = %v", err)
	}

	if string(result.PNG) != "png" {
		t.Fatalf("result.PNG = %q, want %q", string(result.PNG), "png")
	}

	if result.Width != 1550 || result.Height != 830 {
		t.Fatalf("dimensions = %dx%d, want 1550x830", result.Width, result.Height)
	}
}

func TestDesktopPropagatesExecutorError(t *testing.T) {
	t.Parallel()

	executor := &fakeExecutor{
		outputs: []executorResult{
			{err: errors.New("powershell failed")},
		},
	}

	desktop := newDesktop(executor)

	_, err := desktop.ListTargets(context.Background())
	if err == nil {
		t.Fatal("ListTargets() should return an error")
	}
}

func TestDesktopRejectsInvalidJSONOutput(t *testing.T) {
	t.Parallel()

	executor := &fakeExecutor{
		outputs: []executorResult{
			{output: []byte("#< CLIXML")},
		},
	}

	desktop := newDesktop(executor)

	_, err := desktop.ListTargets(context.Background())
	if err == nil {
		t.Fatal("ListTargets() should reject invalid json output")
	}
}

type fakeExecutor struct {
	outputs []executorResult
	index   int
}

type executorResult struct {
	output []byte
	err    error
}

func (f *fakeExecutor) Run(_ context.Context, _ string) ([]byte, error) {
	if f.index >= len(f.outputs) {
		return nil, errors.New("unexpected executor call")
	}

	current := f.outputs[f.index]
	f.index++
	return current.output, current.err
}
