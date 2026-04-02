package cli

import (
	"bytes"
	"testing"
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
