package backend

import "context"

type Backend interface {
	Init() error
	Apply(ctx context.Context, data []byte) error
	Diff(ctx context.Context, data []byte) error
}
