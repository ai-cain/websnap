package chromedp

import (
	"testing"

	"github.com/ai-cain/websnap/internal/domain"
)

func TestCaptureFailureMessage(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		req  domain.CaptureRequest
		want string
	}{
		{
			name: "selector capture",
			req:  domain.CaptureRequest{Selector: "#app"},
			want: "failed to capture the requested selector",
		},
		{
			name: "full-page capture",
			req:  domain.CaptureRequest{FullPage: true},
			want: "failed to capture the full page",
		},
		{
			name: "default viewport capture",
			req:  domain.CaptureRequest{},
			want: "failed to render page in a headless browser",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := captureFailureMessage(tt.req)
			if got != tt.want {
				t.Fatalf("captureFailureMessage() = %q, want %q", got, tt.want)
			}
		})
	}
}
