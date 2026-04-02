package chromedp

import (
	"context"
	"os"
	"testing"

	"github.com/ai-cain/websnap/internal/domain"
)

func TestBrowserCaptureScreenshotIntegration(t *testing.T) {
	t.Parallel()

	targetURL := os.Getenv("WEBSNAP_E2E_URL")
	if targetURL == "" {
		t.Skip("set WEBSNAP_E2E_URL to run browser integration test")
	}

	browser := New()
	req := domain.CaptureRequest{
		URL:    targetURL,
		Width:  1280,
		Height: 720,
	}

	data, err := browser.CaptureScreenshot(context.Background(), req)
	if err != nil {
		t.Fatalf("CaptureScreenshot() error = %v", err)
	}

	if len(data) == 0 {
		t.Fatal("CaptureScreenshot() returned empty data")
	}
}

func TestBrowserCaptureSelectorScreenshotIntegration(t *testing.T) {
	t.Parallel()

	targetURL := os.Getenv("WEBSNAP_E2E_URL")
	selector := os.Getenv("WEBSNAP_E2E_SELECTOR")
	if targetURL == "" || selector == "" {
		t.Skip("set WEBSNAP_E2E_URL and WEBSNAP_E2E_SELECTOR to run selector integration test")
	}

	browser := New()
	req := domain.CaptureRequest{
		URL:      targetURL,
		Width:    1280,
		Height:   720,
		Selector: selector,
	}

	data, err := browser.CaptureScreenshot(context.Background(), req)
	if err != nil {
		t.Fatalf("CaptureScreenshot() error = %v", err)
	}

	if len(data) == 0 {
		t.Fatal("CaptureScreenshot() returned empty data")
	}
}
