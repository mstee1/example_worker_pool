package worker

import (
	"context"

	"github.com/mstee1/example_worker_pool/internal/config"
)

func (w *workerApp) sendErr(ctx context.Context, err error) {
	select {
	case <-ctx.Done():
	default:
		w.chans.errCh <- config.ErrCh{
			Name: w.cfgw.worker.Name,
			Log:  w.log,
			Err:  err,
		}
	}
}
