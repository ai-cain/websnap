package filesystem

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestWriterSaveCreatesParentDirectories(t *testing.T) {
	t.Parallel()

	baseDir := t.TempDir()
	target := filepath.Join(baseDir, "media", "img", "capture.png")
	content := []byte("png-bytes")

	writer := New()

	if err := writer.Save(context.Background(), target, content); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	got, err := os.ReadFile(target)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}

	if string(got) != string(content) {
		t.Fatalf("content = %q, want %q", string(got), string(content))
	}
}
