package cli

import (
	"context"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/ai-cain/websnap/internal/domain"
)

type InteractiveCommand struct {
	runner   ShotRunner
	prompter linePrompter
	stdout   io.Writer
	stderr   io.Writer
}

func NewInteractiveCommand(runner ShotRunner, input io.Reader, stdout, stderr io.Writer) InteractiveCommand {
	return InteractiveCommand{
		runner:   runner,
		prompter: newLinePrompter(input, stdout),
		stdout:   stdout,
		stderr:   stderr,
	}
}

func (c InteractiveCommand) Run() int {
	if c.runner == nil {
		fmt.Fprintln(c.stderr, "error: interactive runner is not configured")
		return 1
	}

	fmt.Fprintln(c.stdout, "websnap interactive mode")
	fmt.Fprintln(c.stdout, "Press Enter to accept defaults.")
	fmt.Fprintln(c.stdout)

	urlValue, err := c.promptRequired("URL", "")
	if err != nil {
		renderError(c.stderr, err)
		return 1
	}

	widthValue, err := c.promptInt64("Width", 1440)
	if err != nil {
		renderError(c.stderr, err)
		return 1
	}

	heightValue, err := c.promptInt64("Height", 900)
	if err != nil {
		renderError(c.stderr, err)
		return 1
	}

	outValue, err := c.prompter.Ask("Output path (optional): ")
	if err != nil {
		renderError(c.stderr, err)
		return 1
	}

	result, err := c.runner.Execute(context.Background(), domain.CaptureRequest{
		URL:    urlValue,
		Width:  widthValue,
		Height: heightValue,
		Out:    strings.TrimSpace(outValue),
	})
	if err != nil {
		renderError(c.stderr, err)
		return 1
	}

	fmt.Fprintln(c.stdout)
	fmt.Fprintln(c.stdout, "Saved screenshot:")
	fmt.Fprintln(c.stdout, result.Path)
	return 0
}

func (c InteractiveCommand) promptRequired(label, fallback string) (string, error) {
	for {
		value, err := c.prompter.Ask(fmt.Sprintf("%s: ", label))
		if err != nil {
			return "", err
		}

		value = strings.TrimSpace(value)
		if value == "" && fallback != "" {
			return fallback, nil
		}

		if value != "" {
			return value, nil
		}

		fmt.Fprintf(c.stderr, "error: %s is required\n", strings.ToLower(label))
	}
}

func (c InteractiveCommand) promptInt64(label string, fallback int64) (int64, error) {
	for {
		value, err := c.prompter.Ask(fmt.Sprintf("%s [%d]: ", label, fallback))
		if err != nil {
			return 0, err
		}

		value = strings.TrimSpace(value)
		if value == "" {
			return fallback, nil
		}

		parsed, err := strconv.ParseInt(value, 10, 64)
		if err == nil && parsed > 0 {
			return parsed, nil
		}

		fmt.Fprintf(c.stderr, "error: %s must be a positive integer\n", strings.ToLower(label))
	}
}
