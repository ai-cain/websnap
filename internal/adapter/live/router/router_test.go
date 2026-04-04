package router

import (
	"context"
	"testing"

	"github.com/ai-cain/websnap/internal/domain"
)

func TestCatalogPrefersWebBrowserTargetsWhenAvailable(t *testing.T) {
	t.Parallel()

	catalog := NewCatalog(
		fakeCatalog{
			targets: []domain.LiveTarget{
				{WindowHandle: 1, Title: "Chrome", AppName: "chrome", Type: domain.LiveTargetBrowser},
				{WindowHandle: 2, Title: "Portfolio", AppName: "explorer", Type: domain.LiveTargetFolder},
			},
		},
		fakeCatalog{
			targets: []domain.LiveTarget{
				{
					WindowHandle:    11,
					BrowserWindowID: 11,
					Title:           "X",
					AppName:         "chrome",
					Type:            domain.LiveTargetBrowser,
					CanListTabs:     true,
					Provider:        domain.LiveTargetProviderBrowserExtension,
				},
			},
		},
	)

	targets, err := catalog.ListTargets(context.Background())
	if err != nil {
		t.Fatalf("ListTargets() error = %v", err)
	}

	if len(targets) != 2 {
		t.Fatalf("len(targets) = %d, want 2", len(targets))
	}

	if targets[0].Type != domain.LiveTargetFolder {
		t.Fatalf("targets[0] = %#v, want non-browser desktop target first", targets[0])
	}

	if targets[1].Provider != domain.LiveTargetProviderBrowserExtension {
		t.Fatalf("targets[1] = %#v, want extension-backed browser target", targets[1])
	}
}

func TestCatalogFallsBackToDesktopTargetsWhenNoWebTargets(t *testing.T) {
	t.Parallel()

	catalog := NewCatalog(
		fakeCatalog{
			targets: []domain.LiveTarget{
				{WindowHandle: 1, Title: "Chrome", AppName: "chrome", Type: domain.LiveTargetBrowser},
			},
		},
		fakeCatalog{},
	)

	targets, err := catalog.ListTargets(context.Background())
	if err != nil {
		t.Fatalf("ListTargets() error = %v", err)
	}

	if len(targets) != 1 || targets[0].Provider != "" {
		t.Fatalf("targets = %#v, want desktop browser fallback", targets)
	}
}

func TestCatalogRoutesTabsByProvider(t *testing.T) {
	t.Parallel()

	web := fakeCatalog{
		tabs: []domain.BrowserTab{{ID: 9, Title: "X"}},
	}
	desktop := fakeCatalog{
		tabs: []domain.BrowserTab{{Index: 1, Title: "Desktop"}},
	}

	catalog := NewCatalog(desktop, web)

	tabs, err := catalog.ListTabs(context.Background(), domain.LiveTarget{
		Provider:        domain.LiveTargetProviderBrowserExtension,
		AppName:         "chrome",
		Type:            domain.LiveTargetBrowser,
		BrowserWindowID: 7,
		Title:           "X",
	})
	if err != nil {
		t.Fatalf("ListTabs() error = %v", err)
	}

	if len(tabs) != 1 || tabs[0].ID != 9 {
		t.Fatalf("tabs = %#v, want web tabs", tabs)
	}
}

func TestCapturerRoutesBrowserExtensionRequestsToWebCapturer(t *testing.T) {
	t.Parallel()

	web := &fakeCapturer{
		image: domain.LiveCaptureImage{PNG: []byte("web"), Width: 1, Height: 1},
	}
	desktop := &fakeCapturer{
		image: domain.LiveCaptureImage{PNG: []byte("desktop"), Width: 1, Height: 1},
	}

	capturer := NewCapturer(desktop, web)

	image, err := capturer.Capture(context.Background(), domain.LiveCaptureRequest{
		Target: domain.LiveTarget{
			Provider:        domain.LiveTargetProviderBrowserExtension,
			BrowserWindowID: 5,
			WindowHandle:    5,
			Title:           "X",
			AppName:         "chrome",
			Type:            domain.LiveTargetBrowser,
		},
		TabID: 81,
	})
	if err != nil {
		t.Fatalf("Capture() error = %v", err)
	}

	if string(image.PNG) != "web" {
		t.Fatalf("image = %#v, want web capturer result", image)
	}
}

type fakeCatalog struct {
	targets []domain.LiveTarget
	tabs    []domain.BrowserTab
	err     error
}

func (f fakeCatalog) ListTargets(_ context.Context) ([]domain.LiveTarget, error) {
	return f.targets, f.err
}

func (f fakeCatalog) ListTabs(_ context.Context, _ domain.LiveTarget) ([]domain.BrowserTab, error) {
	return f.tabs, f.err
}

type fakeCapturer struct {
	image domain.LiveCaptureImage
	err   error
}

func (f *fakeCapturer) Capture(_ context.Context, _ domain.LiveCaptureRequest) (domain.LiveCaptureImage, error) {
	return f.image, f.err
}
