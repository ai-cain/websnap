package port

import "context"

type Writer interface {
	Save(ctx context.Context, path string, data []byte) error
}
