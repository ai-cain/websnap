package domain

import (
	"net/url"
	"path/filepath"
	"strings"

	apperrors "github.com/ai-cain/websnap/internal/support/errors"
)

type CaptureRequest struct {
	URL    string
	Width  int64
	Height int64
	Out    string
}

func (r CaptureRequest) Validate() error {
	if strings.TrimSpace(r.URL) == "" {
		return apperrors.New(apperrors.CodeInvalidArgument, "url is required")
	}

	parsedURL, err := url.ParseRequestURI(r.URL)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		return apperrors.New(apperrors.CodeInvalidArgument, "url must be absolute and include scheme and host")
	}

	if r.Width <= 0 {
		return apperrors.New(apperrors.CodeInvalidArgument, "width must be greater than zero")
	}

	if r.Height <= 0 {
		return apperrors.New(apperrors.CodeInvalidArgument, "height must be greater than zero")
	}

	if r.Out != "" {
		ext := strings.ToLower(filepath.Ext(r.Out))
		if ext != "" && ext != ".png" {
			return apperrors.New(apperrors.CodeInvalidArgument, "output path must use the .png extension")
		}
	}

	return nil
}
