package cli

import (
	"bytes"
	"testing"
)

func TestAppRunUnknownCommand(t *testing.T) {
	t.Parallel()

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	app := newAppWithDeps(nil, nil, &stdout, &stderr, &fakeInteractiveUI{})

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
	ui := &fakeInteractiveUI{}

	app := newAppWithDeps(
		&fakeShotRunner{},
		bytes.NewBufferString(""),
		&stdout,
		&stderr,
		ui,
	)

	exitCode := app.Run(nil)
	if exitCode != 0 {
		t.Fatalf("exitCode = %d, want 0", exitCode)
	}

	if !ui.called {
		t.Fatal("interactive UI should be called")
	}
}
