package port

import (
	"context"

	"github.com/ai-cain/websnap/internal/domain"
)

type Browser interface {
	CaptureScreenshot(ctx context.Context, req domain.CaptureRequest) ([]byte, error)
}
