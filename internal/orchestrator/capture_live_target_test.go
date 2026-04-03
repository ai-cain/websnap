package orchestrator

import (
	"context"
	"errors"
	"path/filepath"
	"testing"
	"time"

	"github.com/ai-cain/websnap/internal/domain"
)

func TestCaptureLiveTargetExecute(t *testing.T) {
	t.Parallel()

	fixedTime := time.Date(2026, time.April, 3, 19, 15, 10, 0, time.UTC)

	tests := []struct {
		name      string
		req       domain.LiveCaptureRequest
		baseDir   string
		capturer  *fakeLiveCapturer
		writer    *fakeWriter
		wantPath  string
		wantError bool
	}{
		{
			name: "writes live capture to default media directory",
			req: domain.LiveCaptureRequest{
				Target: domain.LiveTarget{
					WindowHandle: 1001,
					Title:        "WhatsApp - Google Chrome",
					AppName:      "chrome",
					Type:         domain.LiveTargetBrowser,
				},
				TabIndex: -1,
			},
			baseDir: "C:/workspace/websnap",
			capturer: &fakeLiveCapturer{
				payload: domain.LiveCaptureImage{PNG: []byte("png"), Width: 1550, Height: 830},
			},
			writer:   &fakeWriter{},
			wantPath: filepath.Join("C:/workspace/websnap", "media", "img", "screenshot-20260403-191510.png"),
		},
		{
			name: "appends png extension for requested live output",
			req: domain.LiveCaptureRequest{
				Target: domain.LiveTarget{
					WindowHandle: 2002,
					Title:        "portfolio - Explorador de archivos",
					AppName:      "explorer",
					Type:         domain.LiveTargetFolder,
				},
				TabIndex: -1,
				Out:      "captures/live-window",
			},
			baseDir: "C:/workspace/websnap",
			capturer: &fakeLiveCapturer{
				payload: domain.LiveCaptureImage{PNG: []byte("png"), Width: 1280, Height: 720},
			},
			writer:   &fakeWriter{},
			wantPath: filepath.Join("C:/workspace/websnap", "captures", "live-window.png"),
		},
		{
			name: "returns writer error",
			req: domain.LiveCaptureRequest{
				Target: domain.LiveTarget{
					WindowHandle: 3003,
					Title:        "WhatsApp - Google Chrome",
					AppName:      "chrome",
					Type:         domain.LiveTargetBrowser,
				},
				TabIndex: -1,
			},
			baseDir: "C:/workspace/websnap",
			capturer: &fakeLiveCapturer{
				payload: domain.LiveCaptureImage{PNG: []byte("png"), Width: 1550, Height: 830},
			},
			writer:    &fakeWriter{err: errors.New("disk full")},
			wantError: true,
		},
		{
			name: "rejects invalid tab selection for folder",
			req: domain.LiveCaptureRequest{
				Target: domain.LiveTarget{
					WindowHandle: 4004,
					Title:        "portfolio - Explorador de archivos",
					AppName:      "explorer",
					Type:         domain.LiveTargetFolder,
				},
				TabIndex: 0,
			},
			baseDir:   "C:/workspace/websnap",
			capturer:  &fakeLiveCapturer{},
			writer:    &fakeWriter{},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			oc := NewCaptureLiveTarget(tt.capturer, tt.writer, tt.baseDir, func() time.Time {
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

			if result.Width != tt.capturer.payload.Width || result.Height != tt.capturer.payload.Height {
				t.Fatalf("result dimensions = %dx%d, want %dx%d", result.Width, result.Height, tt.capturer.payload.Width, tt.capturer.payload.Height)
			}

			if tt.writer.savedPath != tt.wantPath {
				t.Fatalf("writer.savedPath = %q, want %q", tt.writer.savedPath, tt.wantPath)
			}
		})
	}
}

type fakeLiveCapturer struct {
	payload domain.LiveCaptureImage
	err     error
}

func (f *fakeLiveCapturer) Capture(_ context.Context, _ domain.LiveCaptureRequest) (domain.LiveCaptureImage, error) {
	return f.payload, f.err
}
