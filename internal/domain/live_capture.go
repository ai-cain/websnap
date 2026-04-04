package domain

import (
	"path/filepath"
	"strings"

	apperrors "github.com/ai-cain/websnap/internal/support/errors"
)

type LiveTargetType string

const (
	LiveTargetApp     LiveTargetType = "app"
	LiveTargetBrowser LiveTargetType = "browser"
	LiveTargetFolder  LiveTargetType = "folder"
)

type LiveTargetProvider string

const (
	LiveTargetProviderDesktop          LiveTargetProvider = "desktop"
	LiveTargetProviderBrowserExtension LiveTargetProvider = "browser-extension"
)

type LiveTarget struct {
	WindowHandle    int64
	Title           string
	AppName         string
	Type            LiveTargetType
	CanListTabs     bool
	Provider        LiveTargetProvider
	BrowserWindowID int
}

type BrowserTab struct {
	Index    int
	ID       int
	WindowID int
	URL      string
	Title    string
	Selected bool
}

type LiveCaptureRequest struct {
	Target   LiveTarget
	TabIndex int
	TabID    int
	Out      string
}

type LiveCaptureImage struct {
	PNG    []byte
	Width  int64
	Height int64
}

func (r LiveCaptureRequest) Validate() error {
	if strings.TrimSpace(r.Target.Title) == "" {
		return apperrors.New(apperrors.CodeInvalidArgument, "target title is required")
	}

	provider := r.Target.Provider
	if provider == "" {
		provider = LiveTargetProviderDesktop
	}

	switch provider {
	case LiveTargetProviderDesktop:
		if r.Target.WindowHandle <= 0 {
			return apperrors.New(apperrors.CodeInvalidArgument, "window handle is required")
		}
	case LiveTargetProviderBrowserExtension:
		if r.Target.Type != LiveTargetBrowser {
			return apperrors.New(apperrors.CodeInvalidArgument, "browser extension capture is only valid for browser targets")
		}
		if r.Target.BrowserWindowID <= 0 && r.Target.WindowHandle <= 0 {
			return apperrors.New(apperrors.CodeInvalidArgument, "browser window id is required")
		}
	default:
		return apperrors.New(apperrors.CodeInvalidArgument, "target provider is invalid")
	}

	if r.TabIndex < -1 {
		return apperrors.New(apperrors.CodeInvalidArgument, "tab index must be -1 or greater")
	}

	if r.Target.Type != LiveTargetBrowser && r.TabIndex >= 0 {
		return apperrors.New(apperrors.CodeInvalidArgument, "tab selection is only valid for browser targets")
	}

	if r.Target.Type != LiveTargetBrowser && r.TabID > 0 {
		return apperrors.New(apperrors.CodeInvalidArgument, "tab id is only valid for browser targets")
	}

	if strings.TrimSpace(r.Out) != "" {
		ext := strings.ToLower(filepath.Ext(r.Out))
		if ext != "" && ext != ".png" {
			return apperrors.New(apperrors.CodeInvalidArgument, "output path must use the .png extension")
		}
	}

	return nil
}
