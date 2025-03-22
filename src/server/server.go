package server

import (
	"context"
)

type HTTPServer interface {
	ListenAndServe() error
	Shutdown(ctx context.Context) error
}
