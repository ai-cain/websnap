package orchestrator

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/ai-cain/websnap/internal/domain"
	"github.com/ai-cain/websnap/internal/port"
	apperrors "github.com/ai-cain/websnap/internal/support/errors"
)

type CaptureScreenshot struct {
	browser port.Browser
	writer  port.Writer
	baseDir string
	now     func() time.Time
}

func NewCaptureScreenshot(browser port.Browser, writer port.Writer, baseDir string, now func() time.Time) *CaptureScreenshot {
	if now == nil {
		now = time.Now
	}

	return &CaptureScreenshot{
		browser: browser,
		writer:  writer,
		baseDir: baseDir,
		now:     now,
	}
}

func (oc *CaptureScreenshot) Execute(ctx context.Context, req domain.CaptureRequest) (domain.CaptureResult, error) {
	if oc == nil {
		return domain.CaptureResult{}, apperrors.New(apperrors.CodeInvalidArgument, "capture orchestrator is not configured")
	}

	if oc.browser == nil {
		return domain.CaptureResult{}, apperrors.New(apperrors.CodeInvalidArgument, "browser dependency is required")
	}

	if oc.writer == nil {
		return domain.CaptureResult{}, apperrors.New(apperrors.CodeInvalidArgument, "writer dependency is required")
	}

	if strings.TrimSpace(oc.baseDir) == "" {
		return domain.CaptureResult{}, apperrors.New(apperrors.CodeInvalidArgument, "base directory is required")
	}

	if err := req.Validate(); err != nil {
		return domain.CaptureResult{}, err
	}

	outputPath, err := resolveOutputPath(oc.baseDir, req.Out, oc.now)
	if err != nil {
		return domain.CaptureResult{}, err
	}

	payload, err := oc.browser.CaptureScreenshot(ctx, req)
	if err != nil {
		return domain.CaptureResult{}, apperrors.Wrap(apperrors.CodeCaptureFailed, "failed to capture screenshot", err)
	}

	if len(payload) == 0 {
		return domain.CaptureResult{}, apperrors.New(apperrors.CodeCaptureFailed, "browser returned an empty screenshot")
	}

	if err := oc.writer.Save(ctx, outputPath, payload); err != nil {
		return domain.CaptureResult{}, apperrors.Wrap(apperrors.CodeWriteFailed, "failed to write screenshot to disk", err)
	}

	return domain.CaptureResult{
		Path:   outputPath,
		Width:  req.Width,
		Height: req.Height,
	}, nil
}

func resolveOutputPath(baseDir, requested string, now func() time.Time) (string, error) {
	if strings.TrimSpace(baseDir) == "" {
		return "", apperrors.New(apperrors.CodeInvalidArgument, "base directory is required")
	}

	if strings.TrimSpace(requested) != "" {
		path := requested
		if filepath.Ext(path) == "" {
			path += ".png"
		}

		if filepath.IsAbs(path) {
			return filepath.Clean(path), nil
		}

		return filepath.Join(baseDir, path), nil
	}

	filename := fmt.Sprintf("screenshot-%s.png", now().UTC().Format("20060102-150405"))
	return filepath.Join(baseDir, "media", "img", filename), nil
}
