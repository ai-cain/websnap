package orchestrator

import (
	"context"
	"strings"
	"time"

	"github.com/ai-cain/websnap/internal/domain"
	"github.com/ai-cain/websnap/internal/port"
	apperrors "github.com/ai-cain/websnap/internal/support/errors"
)

type CaptureLiveTarget struct {
	capturer port.LiveCapturer
	writer   port.Writer
	baseDir  string
	now      func() time.Time
}

func NewCaptureLiveTarget(capturer port.LiveCapturer, writer port.Writer, baseDir string, now func() time.Time) *CaptureLiveTarget {
	if now == nil {
		now = time.Now
	}

	return &CaptureLiveTarget{
		capturer: capturer,
		writer:   writer,
		baseDir:  baseDir,
		now:      now,
	}
}

func (oc *CaptureLiveTarget) Execute(ctx context.Context, req domain.LiveCaptureRequest) (domain.CaptureResult, error) {
	if oc == nil {
		return domain.CaptureResult{}, apperrors.New(apperrors.CodeInvalidArgument, "live capture orchestrator is not configured")
	}

	if oc.capturer == nil {
		return domain.CaptureResult{}, apperrors.New(apperrors.CodeInvalidArgument, "live capturer dependency is required")
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

	payload, err := oc.capturer.Capture(ctx, req)
	if err != nil {
		return domain.CaptureResult{}, apperrors.Wrap(apperrors.CodeCaptureFailed, "failed to capture selected target", err)
	}

	if len(payload.PNG) == 0 {
		return domain.CaptureResult{}, apperrors.New(apperrors.CodeCaptureFailed, "live capturer returned an empty screenshot")
	}

	if err := oc.writer.Save(ctx, outputPath, payload.PNG); err != nil {
		return domain.CaptureResult{}, apperrors.Wrap(apperrors.CodeWriteFailed, "failed to write screenshot to disk", err)
	}

	return domain.CaptureResult{
		Path:   outputPath,
		Width:  payload.Width,
		Height: payload.Height,
	}, nil
}
