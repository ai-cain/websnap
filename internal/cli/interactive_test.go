package cli

import (
	"bytes"
	"context"
	"testing"

	"github.com/ai-cain/websnap/internal/domain"
)

func TestInteractiveCommandRun(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		input        string
		runner       *fakeShotRunner
		wantExitCode int
		wantStdout   string
		wantStderr   string
	}{
		{
			name:  "successful flow with defaults",
			input: "https://example.com\n\n\n\n",
			runner: &fakeShotRunner{
				result: domain.CaptureResult{Path: "C:/captures/default.png"},
			},
			wantExitCode: 0,
			wantStdout:   "Saved screenshot:\nC:/captures/default.png",
		},
		{
			name:  "retries invalid width",
			input: "https://example.com\nabc\n1200\n800\ncaptures/home.png\n",
			runner: &fakeShotRunner{
				result: domain.CaptureResult{Path: "C:/captures/home.png"},
			},
			wantExitCode: 0,
			wantStdout:   "Saved screenshot:\nC:/captures/home.png",
			wantStderr:   "error: width must be a positive integer",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var stdout bytes.Buffer
			var stderr bytes.Buffer

			cmd := NewInteractiveCommand(tt.runner, bytes.NewBufferString(tt.input), &stdout, &stderr)

			exitCode := cmd.Run()
			if exitCode != tt.wantExitCode {
				t.Fatalf("exitCode = %d, want %d", exitCode, tt.wantExitCode)
			}

			if got := stdout.String(); !bytes.Contains([]byte(got), []byte(tt.wantStdout)) {
				t.Fatalf("stdout = %q, want to contain %q", got, tt.wantStdout)
			}

			if tt.wantStderr != "" && !bytes.Contains(stderr.Bytes(), []byte(tt.wantStderr)) {
				t.Fatalf("stderr = %q, want to contain %q", stderr.String(), tt.wantStderr)
			}
		})
	}
}

func TestInteractiveCommandPassesCollectedInput(t *testing.T) {
	t.Parallel()

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	runner := &fakeShotRunner{
		result: domain.CaptureResult{Path: "C:/captures/home.png"},
	}

	cmd := NewInteractiveCommand(
		runner,
		bytes.NewBufferString("https://example.com\n1600\n1000\ncaptures/home.png\n"),
		&stdout,
		&stderr,
	)

	exitCode := cmd.Run()
	if exitCode != 0 {
		t.Fatalf("exitCode = %d, want 0", exitCode)
	}

	if runner.received.URL != "https://example.com" {
		t.Fatalf("URL = %q, want %q", runner.received.URL, "https://example.com")
	}

	if runner.received.Width != 1600 {
		t.Fatalf("Width = %d, want 1600", runner.received.Width)
	}

	if runner.received.Height != 1000 {
		t.Fatalf("Height = %d, want 1000", runner.received.Height)
	}

	if runner.received.Out != "captures/home.png" {
		t.Fatalf("Out = %q, want %q", runner.received.Out, "captures/home.png")
	}
}

type fakeShotRunner struct {
	received domain.CaptureRequest
	result   domain.CaptureResult
	err      error
}

func (f *fakeShotRunner) Execute(_ context.Context, req domain.CaptureRequest) (domain.CaptureResult, error) {
	f.received = req
	return f.result, f.err
}
