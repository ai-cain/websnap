package extensionbridge

import (
	"context"
	"encoding/base64"
	"image"
	"image/color"
	"image/png"
	"os"
	"strings"
	"testing"

	"github.com/ai-cain/websnap/internal/domain"
)

func TestBuildTargetsFromSnapshot(t *testing.T) {
	t.Parallel()

	targets := buildTargetsFromSnapshot(browserSnapshot{
		Browser: "chrome",
		Windows: []browserWindowPayload{
			{
				WindowID: 7,
				AppName:  "chrome",
				Title:    "X - Chrome",
				Tabs: []browserTabPayload{
					{TabID: 91, Index: 0, Title: "X", Active: true},
				},
			},
		},
	})

	if len(targets) != 1 {
		t.Fatalf("len(targets) = %d, want 1", len(targets))
	}

	if targets[0].Provider != domain.LiveTargetProviderBrowserExtension || targets[0].BrowserWindowID != 7 {
		t.Fatalf("targets[0] = %#v, want extension-backed browser target", targets[0])
	}
}

func TestBridgeListTabsFromSnapshot(t *testing.T) {
	t.Parallel()

	bridge := &Bridge{
		clients: map[string]*client{
			"chrome": {
				browser: "chrome",
				snapshot: browserSnapshot{
					Browser: "chrome",
					Windows: []browserWindowPayload{
						{
							WindowID: 11,
							AppName:  "chrome",
							Tabs: []browserTabPayload{
								{TabID: 101, Index: 0, Title: "X", URL: "https://x.com", Active: true},
							},
						},
					},
				},
			},
		},
	}

	tabs, err := bridge.ListTabs(context.Background(), domain.LiveTarget{
		Provider:        domain.LiveTargetProviderBrowserExtension,
		AppName:         "chrome",
		Type:            domain.LiveTargetBrowser,
		BrowserWindowID: 11,
		Title:           "X",
	})
	if err != nil {
		t.Fatalf("ListTabs() error = %v", err)
	}

	if len(tabs) != 1 || tabs[0].ID != 101 || !tabs[0].Selected {
		t.Fatalf("tabs = %#v, want active extension tab", tabs)
	}
}

func TestDecodePNGDataURL(t *testing.T) {
	t.Parallel()

	dataURL := buildPNGDataURL(t, 3, 2)
	pngBytes, width, height, err := decodePNGDataURL(dataURL)
	if err != nil {
		t.Fatalf("decodePNGDataURL() error = %v", err)
	}

	if len(pngBytes) == 0 || width != 3 || height != 2 {
		t.Fatalf("decodePNGDataURL() = len %d, %dx%d, want png bytes and 3x2", len(pngBytes), width, height)
	}
}

func TestNormalizedBrowserName(t *testing.T) {
	t.Parallel()

	tests := map[string]string{
		"chrome":         "chrome",
		"Google-Chrome":  "chrome",
		"edge":           "edge",
		"microsoft-edge": "edge",
	}

	for input, want := range tests {
		if got := normalizedBrowserName(input); got != want {
			t.Fatalf("normalizedBrowserName(%q) = %q, want %q", input, got, want)
		}
	}
}

func TestBuildTargetsFromSnapshotFallsBackToActiveTabTitle(t *testing.T) {
	t.Parallel()

	targets := buildTargetsFromSnapshot(browserSnapshot{
		Browser: "chrome",
		Windows: []browserWindowPayload{
			{
				WindowID: 1,
				Tabs: []browserTabPayload{
					{TabID: 1, Index: 0, Title: "Claude", Active: true},
				},
			},
		},
	})

	if len(targets) != 1 || !strings.Contains(targets[0].Title, "Claude") {
		t.Fatalf("targets = %#v, want title from active tab", targets)
	}
}

func buildPNGDataURL(t *testing.T, width, height int) string {
	t.Helper()

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	img.Set(0, 0, color.RGBA{R: 255, A: 255})

	file, err := os.CreateTemp(t.TempDir(), "*.png")
	if err != nil {
		t.Fatalf("CreateTemp() error = %v", err)
	}
	defer file.Close()

	if err := png.Encode(file, img); err != nil {
		t.Fatalf("png.Encode() error = %v", err)
	}

	bytes, err := os.ReadFile(file.Name())
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}

	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(bytes)
}
