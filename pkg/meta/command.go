package meta

import "context"

type CommandFunc func(ctx context.Context) error
