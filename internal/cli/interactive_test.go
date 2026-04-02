package cli

import (
	"bytes"
	"errors"
	"testing"
)

func TestInteractiveCommandRun(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		uiErr        error
		wantExitCode int
		wantStderr   string
	}{
		{name: "successful tui execution", wantExitCode: 0},
		{name: "tui returns error", uiErr: errors.New("boom"), wantExitCode: 1, wantStderr: "error: boom"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var stdout bytes.Buffer
			var stderr bytes.Buffer
			ui := &fakeInteractiveUI{err: tt.uiErr}

			cmd := newInteractiveCommand(&fakeShotRunner{}, bytes.NewBufferString(""), &stdout, &stderr, ui)

			exitCode := cmd.Run()
			if exitCode != tt.wantExitCode {
				t.Fatalf("exitCode = %d, want %d", exitCode, tt.wantExitCode)
			}

			if tt.wantStderr != "" && !bytes.Contains(stderr.Bytes(), []byte(tt.wantStderr)) {
				t.Fatalf("stderr = %q, want to contain %q", stderr.String(), tt.wantStderr)
			}

			if !ui.called {
				t.Fatal("interactive UI should be called")
			}
		})
	}
}
