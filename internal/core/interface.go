package core

import (
	"context"
	"sync"
)

type App interface {
	GracefulShutdownApp(ctx context.Context, cancel context.CancelFunc)
	Run(ctx context.Context, cond *sync.Cond, canWork *bool)
}
