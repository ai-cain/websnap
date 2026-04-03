package tui

import (
	"context"

	"github.com/ai-cain/websnap/internal/domain"
)

type Studio interface {
	CaptureURL(ctx context.Context, req domain.CaptureRequest) (domain.CaptureResult, error)
	ListTargets(ctx context.Context) ([]domain.LiveTarget, error)
	ListTabs(ctx context.Context, target domain.LiveTarget) ([]domain.BrowserTab, error)
	CaptureLive(ctx context.Context, req domain.LiveCaptureRequest) (domain.CaptureResult, error)
}

