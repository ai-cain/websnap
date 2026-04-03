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

type LiveTarget struct {
	WindowHandle int64
	Title        string
	AppName      string
	Type         LiveTargetType
	CanListTabs  bool
}

type BrowserTab struct {
	Index    int
	Title    string
	Selected bool
}

type LiveCaptureRequest struct {
	Target   LiveTarget
	TabIndex int
	Out      string
}

type LiveCaptureImage struct {
	PNG    []byte
	Width  int64
	Height int64
}

func (r LiveCaptureRequest) Validate() error {
	if r.Target.WindowHandle <= 0 {
		return apperrors.New(apperrors.CodeInvalidArgument, "window handle is required")
	}

	if strings.TrimSpace(r.Target.Title) == "" {
		return apperrors.New(apperrors.CodeInvalidArgument, "target title is required")
	}

	if r.TabIndex < -1 {
		return apperrors.New(apperrors.CodeInvalidArgument, "tab index must be -1 or greater")
	}

	if r.Target.Type != LiveTargetBrowser && r.TabIndex >= 0 {
		return apperrors.New(apperrors.CodeInvalidArgument, "tab selection is only valid for browser targets")
	}

	if strings.TrimSpace(r.Out) != "" {
		ext := strings.ToLower(filepath.Ext(r.Out))
		if ext != "" && ext != ".png" {
			return apperrors.New(apperrors.CodeInvalidArgument, "output path must use the .png extension")
		}
	}

	return nil
}
