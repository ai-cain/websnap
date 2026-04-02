package main

import (
	"fmt"
	"os"
	"time"

	chromedpadapter "github.com/ai-cain/websnap/internal/adapter/browser/chromedp"
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

	runner := orchestrator.NewCaptureScreenshot(
		chromedpadapter.New(),
		filesystemwriter.New(),
		workingDir,
		time.Now,
	)

	app := cli.NewAppWithInput(runner, os.Stdin, os.Stdout, os.Stderr)
	os.Exit(app.Run(os.Args[1:]))
}
