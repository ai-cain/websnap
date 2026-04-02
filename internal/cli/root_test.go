package cli

import (
	"bytes"
	"testing"

	"github.com/ai-cain/websnap/internal/domain"
)

func TestAppRunUnknownCommand(t *testing.T) {
	t.Parallel()

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	app := NewApp(nil, &stdout, &stderr)

	exitCode := app.Run([]string{"unknown"})
	if exitCode != 1 {
		t.Fatalf("exitCode = %d, want 1", exitCode)
	}

	if got := stderr.String(); got == "" {
		t.Fatal("stderr should contain usage output")
	}
}

func TestAppRunWithoutArgsStartsInteractiveMode(t *testing.T) {
	t.Parallel()

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	runner := &fakeShotRunner{
		result: domain.CaptureResult{Path: "C:/captures/interactive.png"},
	}

	app := NewAppWithInput(
		runner,
		bytes.NewBufferString("https://example.com\n\n\n\n"),
		&stdout,
		&stderr,
	)

	exitCode := app.Run(nil)
	if exitCode != 0 {
		t.Fatalf("exitCode = %d, want 0", exitCode)
	}

	if got := stdout.String(); !bytes.Contains([]byte(got), []byte("websnap interactive mode")) {
		t.Fatalf("stdout = %q, want interactive banner", got)
	}
}
