package cli

import (
	"bytes"
	"context"
	"testing"

	"github.com/ai-cain/websnap/internal/domain"
)

func TestShotCommandRun(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		args         []string
		runner       *fakeShotRunner
		wantExitCode int
		wantStdout   string
		wantStderr   string
	}{
		{
			name:         "missing URL argument",
			args:         []string{},
			runner:       &fakeShotRunner{},
			wantExitCode: 1,
			wantStderr:   "error: shot requires exactly one URL argument",
		},
		{
			name:         "rejects flags before URL",
			args:         []string{"--width", "1200"},
			runner:       &fakeShotRunner{},
			wantExitCode: 1,
			wantStderr:   "error: shot expects the URL before any flags",
		},
		{
			name: "successful execution",
			args: []string{"https://example.com", "--width", "1200", "--height", "800"},
			runner: &fakeShotRunner{
				result: domain.CaptureResult{Path: "C:/captures/home.png"},
			},
			wantExitCode: 0,
			wantStdout:   "C:/captures/home.png\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var stdout bytes.Buffer
			var stderr bytes.Buffer

			cmd := NewShotCommand(tt.runner, &stdout, &stderr)

			exitCode := cmd.Run(tt.args)
			if exitCode != tt.wantExitCode {
				t.Fatalf("exitCode = %d, want %d", exitCode, tt.wantExitCode)
			}

			if got := stdout.String(); got != tt.wantStdout {
				t.Fatalf("stdout = %q, want %q", got, tt.wantStdout)
			}

			if tt.wantStderr != "" && !bytes.Contains(stderr.Bytes(), []byte(tt.wantStderr)) {
				t.Fatalf("stderr = %q, want to contain %q", stderr.String(), tt.wantStderr)
			}

			if tt.wantExitCode == 0 {
				if tt.runner.received.URL != "https://example.com" {
					t.Fatalf("runner received URL = %q, want %q", tt.runner.received.URL, "https://example.com")
				}

				if tt.runner.received.Width != 1200 {
					t.Fatalf("runner received width = %d, want 1200", tt.runner.received.Width)
				}

				if tt.runner.received.Height != 800 {
					t.Fatalf("runner received height = %d, want 800", tt.runner.received.Height)
				}
			}
		})
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
