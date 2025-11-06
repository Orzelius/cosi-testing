package backend

import "context"

type Backend interface {
	Init(ctx context.Context) error
	Apply(ctx context.Context, data []byte, dryRun bool) error
	Diff(ctx context.Context, data []byte) error
}
