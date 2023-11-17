package worker

import (
	"context"
	"sync"
)

type WorkerApp interface {
	GracefulShutdownWorker(ctx context.Context)
	Run(ctx context.Context, cond *sync.Cond, canWork *bool)
}
