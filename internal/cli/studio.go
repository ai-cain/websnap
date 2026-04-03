package cli

import (
	"context"

	"github.com/ai-cain/websnap/internal/domain"
	apperrors "github.com/ai-cain/websnap/internal/support/errors"
	"github.com/ai-cain/websnap/internal/tui"
)

type liveTargetCatalog interface {
	ListTargets(ctx context.Context) ([]domain.LiveTarget, error)
	ListTabs(ctx context.Context, target domain.LiveTarget) ([]domain.BrowserTab, error)
}

type liveCaptureRunner interface {
	Execute(ctx context.Context, req domain.LiveCaptureRequest) (domain.CaptureResult, error)
}

type interactiveStudio struct {
	shotRunner  ShotRunner
	targets     liveTargetCatalog
	liveCapture liveCaptureRunner
}

func NewInteractiveStudio(shotRunner ShotRunner, targets liveTargetCatalog, liveCapture liveCaptureRunner) tui.Studio {
	return interactiveStudio{
		shotRunner:  shotRunner,
		targets:     targets,
		liveCapture: liveCapture,
	}
}

func (s interactiveStudio) CaptureURL(ctx context.Context, req domain.CaptureRequest) (domain.CaptureResult, error) {
	if s.shotRunner == nil {
		return domain.CaptureResult{}, apperrors.New(apperrors.CodeInvalidArgument, "shot runner is not configured")
	}

	return s.shotRunner.Execute(ctx, req)
}

func (s interactiveStudio) ListTargets(ctx context.Context) ([]domain.LiveTarget, error) {
	if s.targets == nil {
		return nil, apperrors.New(apperrors.CodeInvalidArgument, "live target catalog is not configured")
	}

	return s.targets.ListTargets(ctx)
}

func (s interactiveStudio) ListTabs(ctx context.Context, target domain.LiveTarget) ([]domain.BrowserTab, error) {
	if s.targets == nil {
		return nil, apperrors.New(apperrors.CodeInvalidArgument, "live target catalog is not configured")
	}

	return s.targets.ListTabs(ctx, target)
}

func (s interactiveStudio) CaptureLive(ctx context.Context, req domain.LiveCaptureRequest) (domain.CaptureResult, error) {
	if s.liveCapture == nil {
		return domain.CaptureResult{}, apperrors.New(apperrors.CodeInvalidArgument, "live capture runner is not configured")
	}

	return s.liveCapture.Execute(ctx, req)
}

