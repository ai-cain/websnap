package main

import (
	"fmt"
	"os"
	"time"

	chromedpadapter "github.com/ai-cain/websnap/internal/adapter/browser/chromedp"
	extensionbridge "github.com/ai-cain/websnap/internal/adapter/browser/extensionbridge"
	liverouter "github.com/ai-cain/websnap/internal/adapter/live/router"
	livewindowsadapter "github.com/ai-cain/websnap/internal/adapter/live/windows"
	filesystemwriter "github.com/ai-cain/websnap/internal/adapter/writer/filesystem"
	"github.com/ai-cain/websnap/internal/cli"
	"github.com/ai-cain/websnap/internal/orchestrator"
)

func main() {
	initConsoleUTF8()

	workingDir, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: unable to determine working directory: %v\n", err)
		os.Exit(1)
	}

	writer := filesystemwriter.New()
	liveDesktop := livewindowsadapter.New()
	webBridge := extensionbridge.New("")
	if err := webBridge.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "warning: browser extension bridge unavailable: %v\n", err)
	}
	defer func() {
		_ = webBridge.Close(nil)
	}()

	runner := orchestrator.NewCaptureScreenshot(
		chromedpadapter.New(),
		writer,
		workingDir,
		time.Now,
	)

	liveRunner := orchestrator.NewCaptureLiveTarget(
		liverouter.NewCapturer(liveDesktop, webBridge),
		writer,
		workingDir,
		time.Now,
	)

	studio := cli.NewInteractiveStudio(liverouter.NewCatalog(liveDesktop, webBridge), liveRunner)

	app := cli.NewAppWithInput(runner, studio, os.Stdin, os.Stdout, os.Stderr)
	os.Exit(app.Run(os.Args[1:]))
}
