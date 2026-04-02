package filesystem

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	apperrors "github.com/ai-cain/websnap/internal/support/errors"
)

type Writer struct{}

func New() *Writer {
	return &Writer{}
}

func (w *Writer) Save(ctx context.Context, path string, data []byte) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	if strings.TrimSpace(path) == "" {
		return apperrors.New(apperrors.CodeInvalidArgument, "output path is required")
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return apperrors.Wrap(apperrors.CodeWriteFailed, "failed to create output directory", err)
	}

	if err := os.WriteFile(path, data, 0o644); err != nil {
		return apperrors.Wrap(apperrors.CodeWriteFailed, "failed to persist screenshot", err)
	}

	return nil
}
