package domain

import "testing"

func TestCaptureRequestValidate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		req     CaptureRequest
		wantErr bool
	}{
		{
			name: "valid request",
			req: CaptureRequest{
				URL:      "https://example.com",
				Width:    1440,
				Height:   900,
				Selector: "#app",
				Out:      "captures/home.png",
			},
			wantErr: false,
		},
		{
			name: "missing URL",
			req: CaptureRequest{
				Width:  1440,
				Height: 900,
			},
			wantErr: true,
		},
		{
			name: "URL without scheme",
			req: CaptureRequest{
				URL:    "example.com",
				Width:  1440,
				Height: 900,
			},
			wantErr: true,
		},
		{
			name: "zero width",
			req: CaptureRequest{
				URL:    "https://example.com",
				Width:  0,
				Height: 900,
			},
			wantErr: true,
		},
		{
			name: "wrong output extension",
			req: CaptureRequest{
				URL:    "https://example.com",
				Width:  1440,
				Height: 900,
				Out:    "captures/home.jpg",
			},
			wantErr: true,
		},
		{
			name: "selector cannot be whitespace only",
			req: CaptureRequest{
				URL:      "https://example.com",
				Width:    1440,
				Height:   900,
				Selector: "   ",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.req.Validate()
			if (err != nil) != tt.wantErr {
				t.Fatalf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
