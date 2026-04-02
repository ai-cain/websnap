package orchestrator

import (
	"context"
	"errors"
	"path/filepath"
	"testing"
	"time"

	"github.com/ai-cain/websnap/internal/domain"
)

func TestCaptureScreenshotExecute(t *testing.T) {
	t.Parallel()

	fixedTime := time.Date(2026, time.April, 2, 18, 30, 45, 0, time.UTC)

	tests := []struct {
		name      string
		req       domain.CaptureRequest
		baseDir   string
		browser   *fakeBrowser
		writer    *fakeWriter
		wantPath  string
		wantError bool
	}{
		{
			name: "writes to default media directory",
			req: domain.CaptureRequest{
				URL:    "https://example.com",
				Width:  1440,
				Height: 900,
			},
			baseDir:  "C:/workspace/websnap",
			browser:  &fakeBrowser{data: []byte("png")},
			writer:   &fakeWriter{},
			wantPath: filepath.Join("C:/workspace/websnap", "media", "img", "screenshot-20260402-183045.png"),
		},
		{
			name: "appends png extension when output has no extension",
			req: domain.CaptureRequest{
				URL:    "https://example.com",
				Width:  1440,
				Height: 900,
				Out:    "captures/home",
			},
			baseDir:  "C:/workspace/websnap",
			browser:  &fakeBrowser{data: []byte("png")},
			writer:   &fakeWriter{},
			wantPath: filepath.Join("C:/workspace/websnap", "captures", "home.png"),
		},
		{
			name: "returns writer error",
			req: domain.CaptureRequest{
				URL:    "https://example.com",
				Width:  1440,
				Height: 900,
			},
			baseDir:   "C:/workspace/websnap",
			browser:   &fakeBrowser{data: []byte("png")},
			writer:    &fakeWriter{err: errors.New("disk full")},
			wantError: true,
		},
		{
			name: "rejects incompatible selector and full-page",
			req: domain.CaptureRequest{
				URL:      "https://example.com",
				Width:    1440,
				Height:   900,
				Selector: "#app",
				FullPage: true,
			},
			baseDir:   "C:/workspace/websnap",
			browser:   &fakeBrowser{data: []byte("png")},
			writer:    &fakeWriter{},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			oc := NewCaptureScreenshot(tt.browser, tt.writer, tt.baseDir, func() time.Time {
				return fixedTime
			})

			result, err := oc.Execute(context.Background(), tt.req)
			if (err != nil) != tt.wantError {
				t.Fatalf("Execute() error = %v, wantError %v", err, tt.wantError)
			}

			if tt.wantError {
				return
			}

			if result.Path != tt.wantPath {
				t.Fatalf("result.Path = %q, want %q", result.Path, tt.wantPath)
			}

			if tt.writer.savedPath != tt.wantPath {
				t.Fatalf("writer.savedPath = %q, want %q", tt.writer.savedPath, tt.wantPath)
			}

			if string(tt.writer.savedData) != "png" {
				t.Fatalf("savedData = %q, want %q", string(tt.writer.savedData), "png")
			}
		})
	}
}

type fakeBrowser struct {
	data []byte
	err  error
}

func (f *fakeBrowser) CaptureScreenshot(_ context.Context, _ domain.CaptureRequest) ([]byte, error) {
	return f.data, f.err
}

type fakeWriter struct {
	savedPath string
	savedData []byte
	err       error
}

func (f *fakeWriter) Save(_ context.Context, path string, data []byte) error {
	f.savedPath = path
	f.savedData = data
	return f.err
}
