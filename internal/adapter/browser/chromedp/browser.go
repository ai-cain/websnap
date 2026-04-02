package chromedp

import (
	"context"
	"strings"

	"github.com/ai-cain/websnap/internal/domain"
	apperrors "github.com/ai-cain/websnap/internal/support/errors"
	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/chromedp"
)

type Browser struct{}

func New() *Browser {
	return &Browser{}
}

func (b *Browser) CaptureScreenshot(ctx context.Context, req domain.CaptureRequest) ([]byte, error) {
	allocOptions := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Headless,
		chromedp.DisableGPU,
		chromedp.Flag("hide-scrollbars", true),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(ctx, allocOptions...)
	defer cancel()

	taskCtx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	tasks := chromedp.Tasks{
		emulation.SetDeviceMetricsOverride(req.Width, req.Height, 1, false),
		chromedp.Navigate(req.URL),
		chromedp.WaitReady("body", chromedp.ByQuery),
	}

	screenshot, captureAction := buildCaptureAction(req)
	tasks = append(tasks, captureAction)

	if err := chromedp.Run(taskCtx, tasks); err != nil {
		return nil, apperrors.Wrap(apperrors.CodeBrowserFailed, "failed to render page in a headless browser", err)
	}

	if len(*screenshot) == 0 {
		return nil, apperrors.New(apperrors.CodeCaptureFailed, "browser returned an empty screenshot")
	}

	return *screenshot, nil
}

func buildCaptureAction(req domain.CaptureRequest) (*[]byte, chromedp.Action) {
	screenshot := new([]byte)

	if strings.TrimSpace(req.Selector) != "" {
		return screenshot, chromedp.Screenshot(req.Selector, screenshot, chromedp.ByQuery)
	}

	if req.FullPage {
		return screenshot, chromedp.FullScreenshot(screenshot, 100)
	}

	return screenshot, chromedp.CaptureScreenshot(screenshot)
}
