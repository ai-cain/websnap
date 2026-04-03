package domain

import "testing"

func TestLiveCaptureRequestValidate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		req     LiveCaptureRequest
		wantErr bool
	}{
		{
			name: "valid browser tab capture",
			req: LiveCaptureRequest{
				Target: LiveTarget{
					WindowHandle: 101,
					Title:        "WhatsApp - Google Chrome",
					AppName:      "chrome",
					Type:         LiveTargetBrowser,
					CanListTabs:  true,
				},
				TabIndex: 1,
				Out:      "captures/live-browser.png",
			},
		},
		{
			name: "valid generic app capture without tab",
			req: LiveCaptureRequest{
				Target: LiveTarget{
					WindowHandle: 202,
					Title:        "portfolio - Explorador de archivos",
					AppName:      "explorer",
					Type:         LiveTargetFolder,
				},
				TabIndex: -1,
			},
		},
		{
			name: "rejects missing handle",
			req: LiveCaptureRequest{
				Target: LiveTarget{
					Title:   "Chrome",
					AppName: "chrome",
					Type:    LiveTargetBrowser,
				},
				TabIndex: -1,
			},
			wantErr: true,
		},
		{
			name: "rejects tab index on non browser target",
			req: LiveCaptureRequest{
				Target: LiveTarget{
					WindowHandle: 303,
					Title:        "portfolio - Explorador de archivos",
					AppName:      "explorer",
					Type:         LiveTargetFolder,
				},
				TabIndex: 0,
			},
			wantErr: true,
		},
		{
			name: "rejects non png output",
			req: LiveCaptureRequest{
				Target: LiveTarget{
					WindowHandle: 404,
					Title:        "Chrome",
					AppName:      "chrome",
					Type:         LiveTargetBrowser,
				},
				TabIndex: -1,
				Out:      "captures/live-browser.jpg",
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
