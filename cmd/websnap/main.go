package main

import (
	"fmt"
	"os"
	"time"

	chromedpadapter "github.com/ai-cain/websnap/internal/adapter/browser/chromedp"
	livewindowsadapter "github.com/ai-cain/websnap/internal/adapter/live/windows"
	filesystemwriter "github.com/ai-cain/websnap/internal/adapter/writer/filesystem"
	"github.com/ai-cain/websnap/internal/cli"
	"github.com/ai-cain/websnap/internal/orchestrator"
)

func main() {
	workingDir, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: unable to determine working directory: %v\n", err)
		os.Exit(1)
	}

	writer := filesystemwriter.New()
	liveDesktop := livewindowsadapter.New()

	runner := orchestrator.NewCaptureScreenshot(
		chromedpadapter.New(),
		writer,
		workingDir,
		time.Now,
	)

	liveRunner := orchestrator.NewCaptureLiveTarget(
		liveDesktop,
		writer,
		workingDir,
		time.Now,
	)

	studio := cli.NewInteractiveStudio(liveDesktop, liveRunner)

	app := cli.NewAppWithInput(runner, studio, os.Stdin, os.Stdout, os.Stderr)
	os.Exit(app.Run(os.Args[1:]))
}
