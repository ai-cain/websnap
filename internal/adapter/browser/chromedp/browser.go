package chromedp

import (
	"context"

	"github.com/ai-cain/websnap/internal/domain"
	apperrors "github.com/ai-cain/websnap/internal/support/errors"
	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/page"
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

	var screenshot []byte

	tasks := chromedp.Tasks{
		emulation.SetDeviceMetricsOverride(req.Width, req.Height, 1, false),
		chromedp.Navigate(req.URL),
		chromedp.WaitReady("body", chromedp.ByQuery),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var err error
			screenshot, err = page.CaptureScreenshot().
				WithFormat(page.CaptureScreenshotFormatPng).
				Do(ctx)
			return err
		}),
	}

	if err := chromedp.Run(taskCtx, tasks); err != nil {
		return nil, apperrors.Wrap(apperrors.CodeBrowserFailed, "failed to render page in a headless browser", err)
	}

	if len(screenshot) == 0 {
		return nil, apperrors.New(apperrors.CodeCaptureFailed, "browser returned an empty screenshot")
	}

	return screenshot, nil
}
