package port

import (
	"context"

	"github.com/ai-cain/websnap/internal/domain"
)

type LiveTargetCatalog interface {
	ListTargets(ctx context.Context) ([]domain.LiveTarget, error)
	ListTabs(ctx context.Context, target domain.LiveTarget) ([]domain.BrowserTab, error)
}

type LiveCapturer interface {
	Capture(ctx context.Context, req domain.LiveCaptureRequest) (domain.LiveCaptureImage, error)
}
